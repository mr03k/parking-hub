package oapi

import (
	"log/slog"
	"net/http"

	"application/internal/http/response"

	"github.com/swaggest/openapi-go"
)

type HealthzOpenApi struct {
	logger *slog.Logger
	api    OAPI
}

func NewHealthzOpenApi(logger *slog.Logger, api OAPI) *HealthzOpenApi {
	o := &HealthzOpenApi{
		logger: logger.With("layer", "HealthzOpenApi"),
		api:    api,
	}
	o.OpenApiSpec()
	return o
}

// rediness
func (s *HealthzOpenApi) HealthzLiveness(op openapi.OperationContext) {
	op.SetSummary("Healthz Liveness")
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

// rediness
func (s *HealthzOpenApi) HealthzReadiness(op openapi.OperationContext) {
	op.SetSummary("Healthz Readiness")
	op.AddRespStructure(new(response.Response[string]))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

func (s *HealthzOpenApi) OpenApiSpec() {
	healthzTag := WithTags("Healthz(Internal)")

	s.api.Register(http.MethodGet, "/healthz/liveness", s.HealthzLiveness, healthzTag)
	s.api.Register(http.MethodGet, "/healthz/readiness", s.HealthzReadiness, healthzTag)
}
