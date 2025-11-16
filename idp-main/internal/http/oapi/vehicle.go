package oapi

import (
	"log/slog"
	"net/http"

	"application/internal/http/dto"
	"application/internal/http/response"

	"github.com/swaggest/openapi-go"
)

type VehicleOpenApi struct {
	logger *slog.Logger
}

type VehicleID struct {
	VehicleID string `json:"vehicle_id" path:"vehicle_id"`
}

func NewVehicleOpenApi(logger *slog.Logger) *VehicleOpenApi {
	return &VehicleOpenApi{
		logger: logger.With("layer", "VehicleOpenApi"),
	}
}

// Create Vehicle
func (s *VehicleOpenApi) CreateVehicle(op openapi.OperationContext) {
	op.SetSummary("Create Vehicle")

	op.AddReqStructure(new(dto.VehicleRequest), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(dto.VehicleResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusBadRequest))
}

// List Vehicles
func (s *VehicleOpenApi) ListVehicles(op openapi.OperationContext) {
	op.SetSummary("List Vehicles")

	op.AddRespStructure(new(dto.VehicleListResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

// Delete Vehicle
func (s *VehicleOpenApi) DeleteVehicle(op openapi.OperationContext) {
	op.SetSummary("Delete Vehicle")
	op.AddReqStructure(new(VehicleID), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

// Get Vehicle By ID
func (s *VehicleOpenApi) GetVehicle(op openapi.OperationContext) {
	op.SetSummary("Get Vehicle By ID")
	op.AddReqStructure(new(VehicleID), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(dto.VehicleResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

func (s *VehicleOpenApi) OpenApiSpec(api OAPI) {
	vehicleTags := WithTags("Vehicle")
	api.Register("POST", "/api/v1/vehicles", s.CreateVehicle, vehicleTags)
	api.Register("GET", "/api/v1/vehicles", s.ListVehicles, vehicleTags)
	api.Register("DELETE", "/api/v1/vehicles/{vehicle_id}", s.DeleteVehicle, vehicleTags)
	api.Register("GET", "/api/v1/vehicles/{vehicle_id}", s.GetVehicle, vehicleTags)
}
