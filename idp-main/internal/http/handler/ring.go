package handler

import (
	"log/slog"
	"net/http"

	"application/internal/datasource"
	"application/internal/http/dto"
	"application/internal/http/response"

	"github.com/swaggest/openapi-go"
)

// RingHandler represents the HTTP handler for rings.
type RingHandler struct {
	logger *slog.Logger
	repo   *datasource.RingRepository
}

// NewRingHandler creates a new instance of RingHandler.
func NewRingHandler(logger *slog.Logger, repo *datasource.RingRepository) *RingHandler {
	return &RingHandler{
		logger: logger.With("layer", "RingHandler"),
		repo:   repo,
	}
}

// ListRings handles listing all rings and responds with a DTO.
func (h *RingHandler) ListRings(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "ListRings")
	w.Header().Set("Content-Type", "application/json")

	logger.Info("Fetching list of rings")
	rings := h.repo.GetRings()

	// Convert rings to DTO
	ringResponses := make([]dto.RingResponse, len(rings))
	for i, ring := range rings {
		ringResponses[i] = dto.RingResponse{
			ID:       ring.ID,
			RingName: ring.RingName,
			RingCode: ring.RingCode,
		}
	}

	response.Pure(w, http.StatusOK, dto.RingListResponse{Rings: ringResponses})
}

// GetRingDetail handles fetching details of a single ring by ID and responds with a DTO.
func (h *RingHandler) GetRingDetail(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "GetRingDetail")
	w.Header().Set("Content-Type", "application/json")

	// Parse the ID from the query parameters
	ringID := r.PathValue("id")
	if ringID == "" {
		logger.Error("missing ring ID in request")
		response.BadRequest(w, "missing ring ID")
		return
	}

	logger.Info("Fetching details for ring", "ringID", ringID)
	ring, err := h.repo.GetRingDetail(ringID)
	if err != nil {
		response.NotFound(w)
		return
	}

	// Convert ring to DTO
	ringResponse := dto.RingDetailResponse{
		ID:                   ring.ID,
		RingName:             ring.RingName,
		RingCode:             ring.RingCode,
		RingLength:           ring.RingLength,
		RingBoundary:         ring.RingBoundary,
		ParkingSpots:         ring.ParkingSpots,
		DisabledParkingSpots: ring.DisabledParkingSpots,
		TrafficSigns:         ring.TrafficSigns,
		DisabledTrafficSigns: ring.DisabledTrafficSigns,
		StartPoint:           ring.StartPoint,
		BufferDistance:       ring.BufferDistance,
		Description:          ring.Description,
	}

	// Respond with the ring details DTO
	response.Pure(w, http.StatusOK, ringResponse)
}

// RegisterMuxRouter registers the ring handler routes with an HTTP multiplexer.
func (h *RingHandler) RegisterMuxRouter(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/rings", h.ListRings)          // List Rings
	mux.HandleFunc("/api/v1/rings/{id}", h.GetRingDetail) // Get Ring Detail
}

type RingID struct {
	ID string `json:"id" path:"id"`
}

// List Rings
func (s *RingHandler) ListRingsOAPI(op openapi.OperationContext) {
	op.SetSummary("List Rings")

	op.AddRespStructure(new(dto.RingListResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

// Get Ring By ID
func (s *RingHandler) GetRingOAPI(op openapi.OperationContext) {
	op.SetSummary("Get Ring By ID")

	op.AddReqStructure(new(RingID), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(dto.RingResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

// OpenApiSpec registers the Ring OpenAPI endpoints.
func (s *RingHandler) OpenApiSpec(api OAPI) {
	ringTags := WithTags("Ring")

	api.Register("GET", "/api/v1/rings", s.ListRingsOAPI, ringTags)
	api.Register("GET", "/api/v1/rings/{id}", s.GetRingOAPI, ringTags)
}
