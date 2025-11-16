package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"application/config"
	biz "application/internal/v1/biz/healthz"
	"application/internal/v1/http/response"
	"application/pkg/middlewares"
	"application/pkg/middlewares/httplogger"
	"application/pkg/middlewares/httprecovery"
	"application/pkg/utils"
	"go.opentelemetry.io/otel"
)

type HealthzHandler struct {
	logger     *slog.Logger
	uc         biz.HealthzUseCaseInterface
	httpConfig *config.HTTPServer
	lastReady  time.Time
}

var _ Handler = (*HealthzHandler)(nil)

func NewMuxHealthzHandler(uc biz.HealthzUseCaseInterface, logger *slog.Logger,
	httpConfig *config.HTTPServer,
) *HealthzHandler {
	return &HealthzHandler{
		logger:     logger.With("layer", "MuxHealthzService"),
		uc:         uc,
		httpConfig: httpConfig,
		lastReady:  time.Now(),
	}
}

// Healthz Liveness
func (s *HealthzHandler) HealthzLiveness(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := otel.Tracer("handler").Start(ctx, "rediness")
	defer span.End()
	logger := s.logger.With("method", "HealthzLiveness", "ctx", utils.GetLoggerContext(r.Context()))
	if s.lastReady.Compare(time.Now().Add(-(time.Minute * 5))) == -1 {
		response.Custom(w, http.StatusServiceUnavailable, nil, "")
		return
	}

	logger.Debug("Liveness")
	err := s.uc.Liveness(ctx)
	if err != nil {
		response.InternalError(w)
		return
	}

	response.Ok(w, nil, "ok")
}

// Healthz Readiness
func (s *HealthzHandler) HealthzReadiness(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := otel.Tracer("handler").Start(ctx, "rediness")
	defer span.End()
	logger := s.logger.With("method", "HealthzReadiness", "ctx", ctx)
	w.Header().Set("Content-Type", "application/json")

	err := s.uc.Readiness(ctx)
	if err != nil {
		response.InternalError(w)
		return
	}

	s.lastReady = time.Now()
	response.Ok(w, nil, "ok")
	logger.DebugContext(ctx, "HealthzReadiness", "url", r.Host, "status", http.StatusOK)
}

// panic
func (s *HealthzHandler) Panic(_ http.ResponseWriter, _ *http.Request) {
	panic("Panic for test")
}

func (s *HealthzHandler) LongRun(w http.ResponseWriter, r *http.Request) {
	// sleep 30 second
	timeString := r.PathValue("time")
	ctx := r.Context()
	logger := s.logger.With("method", "LongRun", "ctx", ctx)
	logger.Debug("LongRun", "time", timeString)

	// sleep to int
	duration, err := time.ParseDuration(timeString)
	if err != nil {
		logger.Error("LongRun", "err", err)
		response.InternalError(w)
		return
	}
	time.Sleep(duration)
	response.Ok(w, nil, "ok")
}

func (s *HealthzHandler) RegisterMuxRouter(mux FuncHandler) {
	recoverMiddleware, err := httprecovery.NewRecoveryMiddleware()
	if err != nil {
		panic(err)
	}

	loggerMiddleware, err := httplogger.NewLoggerMiddleware()
	if err != nil {
		panic(err)
	}
	loggerMiddlewareDebug, err := httplogger.NewLoggerMiddleware(httplogger.WithLevel(slog.LevelDebug))
	if err != nil {
		panic(err)
	}

	healthzMiddleware := []middlewares.Middleware{
		recoverMiddleware.RecoverMiddleware,
		httplogger.SetRequestContextLogger,
		loggerMiddlewareDebug.LoggerMiddleware,
	}

	otherMiddleware := []middlewares.Middleware{
		loggerMiddleware.LoggerMiddleware,
		recoverMiddleware.RecoverMiddleware,
		httplogger.SetRequestContextLogger,
	}
	mux.HandleFunc(fmt.Sprintf("GET /%s/healthz/liveness", s.httpConfig.BasePath),
		middlewares.MultipleMiddleware(s.HealthzLiveness, healthzMiddleware...))
	mux.HandleFunc(fmt.Sprintf("GET /%s/healthz/readiness", s.httpConfig.BasePath),
		middlewares.MultipleMiddleware(s.HealthzReadiness, healthzMiddleware...))
	mux.HandleFunc(fmt.Sprintf("GET /%s/healthz/panic", s.httpConfig.BasePath),
		middlewares.MultipleMiddleware(s.Panic, otherMiddleware...))
	mux.HandleFunc(fmt.Sprintf("GET /%s/healthz/sleep/{time}", s.httpConfig.BasePath),
		middlewares.MultipleMiddleware(s.LongRun, otherMiddleware...))
}
