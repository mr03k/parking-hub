package oapi

import (
	"log/slog"
	"net/http"

	"application/internal/http/dto"
	"application/internal/http/response"

	"github.com/swaggest/openapi-go"
)

type CityOpenApi struct {
	logger *slog.Logger
}

type CityID struct {
	CityID string `json:"city_id" path:"city_id"`
}

func NewCityOpenApi(logger *slog.Logger) *CityOpenApi {
	return &CityOpenApi{
		logger: logger.With("layer", "VehicleOpenApi"),
	}
}

// get city by id
func (s *CityOpenApi) GetCity(op openapi.OperationContext) {
	op.SetSummary("Get City By ID")
	op.AddReqStructure(new(CityID), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(dto.CityResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

// list city

func (s *CityOpenApi) ListCity(op openapi.OperationContext) {
	op.SetSummary("List City")
	op.AddRespStructure(new(dto.CityListResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

func (s *CityOpenApi) OpenApiSpec(api OAPI) {
	tags := WithTags("city")

	api.Register("GET", "/api/v1/cities", s.ListCity, tags)
	api.Register("GET", "/api/v1/cities/{city_id}", s.GetCity, tags)
}
