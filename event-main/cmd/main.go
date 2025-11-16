package main

import (
	"context"
	"errors"
	"flag"
	"git.abanppc.com/farin-project/event/config"
	"git.abanppc.com/farin-project/event/internal/healthz"
	mongoPKG "git.abanppc.com/farin-project/event/pkg/mongo"
	"git.abanppc.com/farin-project/event/pkg/rabbit"
	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

var ErrorRequestTimeout = errors.New("request take too long to respond")

func main() {
	cp := flag.String("config", "", "config file address")
	flag.Parse()
	configPath := *cp
	flag.Parse()
	ctx := context.Background()
	defer ctx.Done()

	// Initialize configuration
	cfg := initConfig(configPath)

	// Initialize tracing
	initTracing(ctx)

	// Initialize logger
	logger := initLogger(cfg)
	var rbtConfig config.RabbitMQ

	if err := cfg.Unmarshal("rabbitMQ", &rbtConfig); err != nil {
		log.Fatal("unmarshal rabbitMQ fail", err)
	}

	rbt := rabbit.NewRabbit(&rbtConfig)
	if err := rbt.Setup(); err != nil {
		log.Fatal("setup rabbitMQ fail", err)
	}

	var mongoConfig config.MongoDB

	if err := cfg.Unmarshal("mongoDB", &mongoConfig); err != nil {
		log.Fatal("unmarshal mongoDB fail ", err)
	}

	mongo := mongoPKG.NewMongo(&mongoConfig)
	if err := mongo.Setup(ctx); err != nil {
		log.Fatal("setup mongodb fail ", err)
	}

	go func() {
		hz := healthz.NewHealthz(logger, rbt, mongo)
		for {
			hz.HealthzReadiness()
			time.Sleep(5 * time.Second)

		}
	}()

	bootstrap, err := wireApp(logger, &cfg, rbt, &rbtConfig, mongo)
	if err != nil {
		log.Fatal("bootstrap app fail", err)
	}
	bootstrap.Boot()
	// Handle graceful shutdown
	handleGracefulShutdown(ctx, logger)
}

func initConfig(configPath string) config.Config {
	if configPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		configPath = path.Join(wd, "config.example.yaml")
	}

	config, err := config.NewKoanfConfig(config.WithYamlConfigPath(configPath))
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

func initLogger(c config.Config) *slog.Logger {
	var cfg config.LogingConfig
	if err := c.Unmarshal("", &cfg); err != nil {
		log.Fatal(err)
	}
	logger := initSlogLogger(cfg)
	logger.Info("logger started", "config", cfg)
	return logger
}

func handleGracefulShutdown(ctx context.Context, logger *slog.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	logger.Info("app stopping...")

	logger.Info("app stopped", "signal", sig)
}

func initSlogLogger(cfg config.LogingConfig) *slog.Logger {
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
