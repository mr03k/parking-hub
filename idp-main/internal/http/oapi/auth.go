package oapi

import (
	"log/slog"
	"net/http"

	"application/internal/http/dto"
	"application/internal/http/response"

	"github.com/swaggest/openapi-go"
)

type AuthOpenApi struct {
	logger *slog.Logger
}

type VehicalLogin struct {
	VehicalID string `path:"vehicle_id"`
	dto.DriverLogin
}

func NewAuthOpenApi(logger *slog.Logger) *AuthOpenApi {
	return &AuthOpenApi{
		logger: logger.With("layer", "AuthOpenApi"),
	}
}

// login
func (s *AuthOpenApi) Login(op openapi.OperationContext) {
	op.SetSummary("Login")
	op.AddReqStructure(new(VehicalLogin), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(dto.DriverLoginResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusBadRequest))
}

func (s *AuthOpenApi) OpenApiSpec(api OAPI) {
	userTags := WithTags("Chapar")
	api.Register(http.MethodPost, "/api/v1/vehicles/{vehicle_id}/login", s.Login, userTags)
}
