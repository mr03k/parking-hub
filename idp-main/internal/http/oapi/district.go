package oapi

import (
	"application/internal/http/dto"
	"application/internal/http/response"
	"github.com/swaggest/openapi-go"
	"log/slog"
	"net/http"
)

// DistrictOpenApi represents the OpenAPI specification for district endpoints.
type DistrictOpenApi struct {
	logger *slog.Logger
}

// DistrictID represents the path parameter for district operations.
type DistrictID struct {
	ID string `json:"id" query:"id"`
}

// NewDistrictOpenApi creates a new instance of DistrictOpenApi.
func NewDistrictOpenApi(logger *slog.Logger) *DistrictOpenApi {
	return &DistrictOpenApi{
		logger: logger.With("layer", "DistrictOpenApi"),
	}
}

// ListDistricts describes the OpenAPI specification for listing districts.
func (s *DistrictOpenApi) ListDistricts(op openapi.OperationContext) {
	op.SetSummary("List Districts")
	op.SetDescription("Fetch a list of all district records with basic information.")

	op.AddRespStructure(new(dto.DistrictListResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

// GetDistrictDetail describes the OpenAPI specification for fetching a single district's details.
func (s *DistrictOpenApi) GetDistrictDetail(op openapi.OperationContext) {
	op.SetSummary("Get District Detail")
	op.SetDescription("Fetch detailed information for a specific district by its unique ID.")

	op.AddReqStructure(new(DistrictID), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(dto.DistrictDetailResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusNotFound))
}

// OpenApiSpec registers the district OpenAPI endpoints.
func (s *DistrictOpenApi) OpenApiSpec(api OAPI) {
	districtTags := WithTags("District")

	api.Register("GET", "/api/v1/districts", s.ListDistricts, districtTags)            // List Districts
	api.Register("GET", "/api/v1/districts/detail", s.GetDistrictDetail, districtTags) // Get District Detail
}
