package oapi

import (
	"application/internal/http/dto"
	"application/internal/http/response"
	"log/slog"
	"net/http"

	"github.com/swaggest/openapi-go"
)

type RingOpenApi struct {
	logger *slog.Logger
}

type RingID struct {
	ID string `json:"id" path:"id"`
}

func NewRingOpenApi(logger *slog.Logger) *RingOpenApi {
	return &RingOpenApi{
		logger: logger.With("layer", "RingOpenApi"),
	}
}

// List Rings
func (s *RingOpenApi) ListRings(op openapi.OperationContext) {
	op.SetSummary("List Rings")

	op.AddRespStructure(new(dto.RingListResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

// Get Ring By ID
func (s *RingOpenApi) GetRing(op openapi.OperationContext) {
	op.SetSummary("Get Ring By ID")

	op.AddReqStructure(new(RingID), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(dto.RingResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

// OpenApiSpec registers the Ring OpenAPI endpoints.
func (s *RingOpenApi) OpenApiSpec(api OAPI) {
	ringTags := WithTags("Ring")

	api.Register("GET", "/api/v1/rings", s.ListRings, ringTags)
	api.Register("GET", "/api/v1/rings/{id}", s.GetRing, ringTags)
}
