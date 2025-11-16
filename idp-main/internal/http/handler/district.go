package handler

import (
	"log/slog"
	"net/http"

	"application/internal/http/dto"
	"application/internal/http/response"
	"application/internal/repo"

	"github.com/swaggest/openapi-go"
)

// DistrictHandler represents the HTTP handler for district records.
type DistrictHandler struct {
	logger *slog.Logger
	repo   *repo.DistrictRepository
}

// NewDistrictHandler creates a new instance of DistrictHandler.
func NewDistrictHandler(logger *slog.Logger, repo *repo.DistrictRepository) *DistrictHandler {
	return &DistrictHandler{
		logger: logger.With("layer", "DistrictHandler"),
		repo:   repo,
	}
}

// ListDistricts handles listing all district records and responds with a DTO.
func (h *DistrictHandler) ListDistricts(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "ListDistricts")

	logger.Info("Fetching list of district records")
	districts := h.repo.GetDistricts()
	w.Header().Set("Content-Type", "application/json")

	// Convert district records to DTO
	districtResponses := make([]dto.DistrictResponse, len(districts))
	for i, district := range districts {
		districtResponses[i] = dto.DistrictResponse{
			ID:           district.ID.String(),
			DistrictName: district.DistrictName,
			DistrictCode: district.DistrictCode,
			CityID:       district.CityID.String(),
			GeoBoundary:  district.GeoBoundary,
			Population:   district.Population,
			Area:         district.Area,
		}
	}

	response.Pure(w, http.StatusOK, dto.DistrictListResponse{Districts: districtResponses})
}

// GetDistrictDetail handles fetching details of a single district by ID and responds with a DTO.
func (h *DistrictHandler) GetDistrictDetail(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "GetDistrictDetail")
	w.Header().Set("Content-Type", "application/json")

	// Parse the ID from the query parameters
	districtID := r.URL.Query().Get("id")
	if districtID == "" {
		logger.Error("missing district ID in request")
		response.BadRequest(w, "missing district ID")
		return
	}

	logger.Info("Fetching details for district", "districtID", districtID)
	district, err := h.repo.GetDistrictByID(districtID)
	if err != nil {
		logger.Error("district not found", "districtID", districtID)
		response.NotFound(w)
		return
	}

	// Convert district to DTO
	districtResponse := dto.DistrictDetailResponse{
		ID:           district.ID.String(),
		DistrictName: district.DistrictName,
		DistrictCode: district.DistrictCode,
		CityID:       district.CityID.String(),
		GeoBoundary:  district.GeoBoundary,
		Population:   district.Population,
		Area:         district.Area,
		CreatedAt:    district.CreatedAt,
	}

	response.Pure(w, http.StatusOK, districtResponse)
}

// RegisterMuxRouter registers the district handler routes with an HTTP multiplexer.
func (h *DistrictHandler) RegisterMuxRouter(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/districts", h.ListDistricts)            // List Districts
	mux.HandleFunc("/api/v1/districts/detail", h.GetDistrictDetail) // Get District Detail
}

// DistrictID represents the path parameter for district operations.
type DistrictID struct {
	ID string `json:"id" query:"id"`
}

// ListDistricts describes the OpenAPI specification for listing districts.
func (s *DistrictHandler) ListDistrictsOAPI(op openapi.OperationContext) {
	op.SetSummary("List Districts")
	op.SetDescription("Fetch a list of all district records with basic information.")

	op.AddRespStructure(new(dto.DistrictListResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

// GetDistrictDetail describes the OpenAPI specification for fetching a single district's details.
func (s *DistrictHandler) GetDistrictDetailOAPI(op openapi.OperationContext) {
	op.SetSummary("Get District Detail")
	op.SetDescription("Fetch detailed information for a specific district by its unique ID.")

	op.AddReqStructure(new(DistrictID), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(dto.DistrictDetailResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusNotFound))
}

// OpenApiSpec registers the district OpenAPI endpoints.
func (s *DistrictHandler) OpenApiSpec(api OAPI) {
	districtTags := WithTags("District")

	api.Register("GET", "/api/v1/districts", s.ListDistrictsOAPI, districtTags)            // List Districts
	api.Register("GET", "/api/v1/districts/detail", s.GetDistrictDetailOAPI, districtTags) // Get District Detail
}
