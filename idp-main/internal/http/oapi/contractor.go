package oapi

import (
	"log/slog"
	"net/http"

	"application/internal/http/dto"
	"application/internal/http/response"

	"github.com/swaggest/openapi-go"
)

type ContractorOpenApi struct {
	logger *slog.Logger
}

type ContractorID struct {
	ContractorID string `json:"contractor_id" path:"contractor_id"`
}

func NewContractorOpenApi(logger *slog.Logger) *ContractorOpenApi {
	return &ContractorOpenApi{
		logger: logger.With("layer", "VehicleOpenApi"),
	}
}

// Create Vehicle
func (s *ContractorOpenApi) CreateContractor(op openapi.OperationContext) {
	op.SetSummary("Create Vehicle")

	op.AddReqStructure(new(dto.ContractorRequest), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(dto.ContractorResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusBadRequest))
}

// List Vehicles
func (s *ContractorOpenApi) ListContractor(op openapi.OperationContext) {
	op.SetSummary("List Vehicles")

	op.AddRespStructure(new(dto.ConttactorListResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

// Delete Vehicle
func (s *ContractorOpenApi) DeleteContractor(op openapi.OperationContext) {
	op.SetSummary("Delete Vehicle")
	op.AddReqStructure(new(ContractorID), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

// Get Vehicle By ID
func (s *ContractorOpenApi) GetContractor(op openapi.OperationContext) {
	op.SetSummary("Get Vehicle By ID")
	op.AddReqStructure(new(ContractorID), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(dto.VehicleResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

func (s *ContractorOpenApi) OpenApiSpec(api OAPI) {
	vehicleTags := WithTags("Contractor")
	api.Register("POST", "/api/v1/contractor", s.CreateContractor, vehicleTags)
	api.Register("GET", "/api/v1/contractor", s.ListContractor, vehicleTags)
	api.Register("DELETE", "/api/v1/contractor/{contractor_id}", s.DeleteContractor, vehicleTags)
	api.Register("GET", "/api/v1/contractor/{contractor_id}", s.GetContractor, vehicleTags)
}
