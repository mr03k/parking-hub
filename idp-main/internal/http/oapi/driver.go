package oapi

import (
	"application/internal/http/dto"
	"application/internal/http/response"
	"log/slog"
	"net/http"

	"github.com/swaggest/openapi-go"
)

type DriverOpenApi struct {
	logger *slog.Logger
}

type DriverID struct {
	ID string `json:"id" path:"id"`
}

func NewDriverOpenApi(logger *slog.Logger) *DriverOpenApi {
	return &DriverOpenApi{
		logger: logger.With("layer", "DriverOpenApi"),
	}
}

// List Drivers
func (s *DriverOpenApi) ListDrivers(op openapi.OperationContext) {
	op.SetSummary("List Drivers")
	op.SetDescription("Fetch a list of all drivers with their basic information.")

	op.AddRespStructure(new(dto.DriverListResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

// Get Driver By ID
func (s *DriverOpenApi) GetDriver(op openapi.OperationContext) {
	op.SetSummary("Get Driver By ID")
	op.SetDescription("Fetch detailed information for a specific driver by their unique ID.")

	op.AddReqStructure(new(DriverID), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(dto.DriverDetailResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusNotFound))
}

// OpenApiSpec registers the Driver OpenAPI endpoints.
func (s *DriverOpenApi) OpenApiSpec(api OAPI) {
	driverTags := WithTags("Driver")

	api.Register("GET", "/api/v1/drivers", s.ListDrivers, driverTags)
	api.Register("GET", "/api/v1/drivers/{id}", s.GetDriver, driverTags)
}
