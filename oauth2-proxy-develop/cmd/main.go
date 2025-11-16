package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	configPKG "application/config"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"golang.org/x/oauth2"
)

var ErrorRequestTimeout = errors.New("request take too long to respond")

func main() {
	cp := flag.String("config", "", "config file address")
	rp := flag.String("rules", "", "rules yaml file address")
	flag.Parse()
	configPath := *cp
	rulesPath := *rp
	flag.Parse()
	ctx := context.Background()
	defer ctx.Done()

	// Initialize configuration
	config := initConfig(configPath)

	// Initialize tracing
	initTracing(ctx)

	// Initialize logger
	logger := initLogger(config)

	// Initialize and start HTTP server
	httpServer := initHTTPServer(ctx, config, rulesPath, logger)

	// Handle graceful shutdown
	handleGracefulShutdown(ctx, httpServer, logger)
}

func initConfig(configPath string) configPKG.Config {
	if configPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		configPath = path.Join(wd, "config.example.yaml")
	}

	config, err := configPKG.NewKoanfConfig(configPKG.WithYamlConfigPath(configPath))
	if err != nil {
		panic(err)
	}

	return config
}

func initTracing(ctx context.Context) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		panic(err)
	}
	defer func() {
		err := exporter.Shutdown(ctx)
		if err != nil {
			panic(err)
		}
	}()

	r := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("myService"),
		semconv.ServiceVersionKey.String("1.0.0"),
		semconv.ServiceInstanceIDKey.String("abcdef12345"),
		semconv.ContainerName("myContainer"),
	)

	r2, err := resource.New(context.Background())
	if err != nil {
		panic(err)
	}

	resource, err := resource.Merge(r, r2)
	if err != nil {
		panic(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
	)

	p := b3.New()
	otel.SetTextMapPropagator(p)
	otel.SetTracerProvider(tp)
}

func initLogger(config configPKG.Config) *slog.Logger {
	var cfg configPKG.LogingConfig
	if err := config.Unmarshal("", &cfg); err != nil {
		log.Fatal(err)
	}
	logger := initSlogLogger(cfg)
	logger.Info("logger started", "config", cfg)
	return logger
}

func initHTTPServer(ctx context.Context, config configPKG.Config, rulesPath string, logger *slog.Logger) *http.Server {
	var server configPKG.Server
	if err := config.Unmarshal("server", &server); err != nil {
		log.Fatal(err)
	}

	httpConfig := server.HTTP
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if rulesPath == "" {
		rulesPath = path.Join(wd, "rules.yaml")
	}

	rulesKoanf := koanf.New("")
	if err := rulesKoanf.Load(file.Provider(rulesPath), yaml.Parser()); err != nil {
		logger.Error("failed to load routes file", "error", err)
		log.Fatal(err.Error())
	}

	var oauth configPKG.Oauth
	if err := config.Unmarshal("oauth", &oauth); err != nil {
		log.Fatal(err)
	}
	oidcConfig, verifier := setupOIDC(ctx, server, oauth)

	var securityConfig configPKG.Security
	if err := config.Unmarshal("security", &securityConfig); err != nil {
		log.Fatal(err)
	}

	engine, err := wireApp(ctx, config, logger, rulesKoanf,
		oidcConfig, verifier, &securityConfig, &httpConfig)
	if err != nil {
		logger.Error("failed to init app", "err", err)
		log.Fatal(err)
	}

	serviceAddr := fmt.Sprintf("%s:%d", httpConfig.Host, httpConfig.Port)
	logger.Info("running server at :", "server", serviceAddr, "emoji", "👷😳")
	httpServer := &http.Server{
		Addr:        serviceAddr,
		Handler:     engine,
		ReadTimeout: 3 * time.Second,
	}

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Error("failed to run app", "err", err)
			panic(err)
		}
	}()
	logger.Info(fmt.Sprintf("microservice started at %s", serviceAddr))

	return httpServer
}

func setupOIDC(ctx context.Context, serverConfig configPKG.Server, oauth configPKG.Oauth) (*oauth2.Config,
	*oidc.IDTokenVerifier,
) {
	provider, err := oidc.NewProvider(ctx, oauth.ProviderURL)
	if err != nil {
		log.Fatal(err)
	}
	oldConfig := &oidc.Config{
		ClientID: oauth.ClientID,
	}
	verifier := provider.Verifier(oldConfig)
	redirectPath, err := url.JoinPath(serverConfig.DomainURL, "oauth2/callback")
	if err != nil {
		log.Fatal(err)
	}
	config := oauth2.Config{
		ClientID:     oauth.ClientID,
		ClientSecret: oauth.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  redirectPath,
		Scopes: []string{
			oidc.ScopeOpenID, "profile", "email", oidc.ScopeOfflineAccess,
			"federated:id", "groups",
		},
	}
	return &config, verifier
}

func handleGracefulShutdown(ctx context.Context, httpServer *http.Server, logger *slog.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	logger.Info("app stopping...")

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("failed to shutdown app", "err", err)
		panic(err)
	}

	logger.Info("app stopped", "signal", sig)
}

func initSlogLogger(cfg configPKG.LogingConfig) *slog.Logger {
	slogHandlerOptions := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}

	level := cfg.Observability.Logging.Level

	switch level {
	case "debug":
		slogHandlerOptions.Level = slog.LevelDebug
	case "info":
		slogHandlerOptions.Level = slog.LevelInfo
	case "warn":
		slogHandlerOptions.Level = slog.LevelWarn
	case "error":
		slogHandlerOptions.Level = slog.LevelError
	default:
		slogHandlerOptions.Level = slog.LevelInfo
	}

	slogHandlers := []slog.Handler{}
	slogHandlers = append(slogHandlers, slog.NewJSONHandler(os.Stdout, slogHandlerOptions))

	if cfg.Observability.Logging.Logstash.Enabled {
		con, err := net.Dial("udp", cfg.Observability.Logging.Logstash.Address)
		if err != nil {
			panic(err)
		}
		slogHandlers = append(slogHandlers, slog.NewJSONHandler(con, slogHandlerOptions))
	}

	logger := slog.New(slogmulti.Fanout(slogHandlers...))
	slog.SetDefault(logger)

	return logger
}
