package handler

import (
	"application/internal/biz/device"
	"application/internal/http/dto"
	"application/internal/http/response"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

var _ Handler = (*VehicleHandler)(nil)

type VehicleHandler struct {
	logger *slog.Logger
	uc     device.VehicleServiceInterface
}

// NewVehicleHandler creates a new VehicleHandler instance
func NewVehicleHandler(logger *slog.Logger, uc device.VehicleServiceInterface) *VehicleHandler {
	return &VehicleHandler{
		logger: logger.With("layer", "VehicleHandler"),
		uc:     uc,
	}
}

// CreateVehicle handles creating a new vehicle
func (h *VehicleHandler) CreateVehicle(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "CreateVehicle")
	ctx := r.Context()

	// Parse request data
	var req dto.VehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("error parsing request", "error", err)
		response.BadRequest(w, err.Error())
		return
	}

	// Create vehicle
	vehicle, err := h.uc.CreateVehicle(ctx, req.ToEntity())
	if err != nil {
		logger.Error("error creating vehicle", "error", err)
		response.Custom(w, http.StatusInternalServerError, nil, "internal server error")
		return
	}

	response.Pure(w, http.StatusCreated, dto.NewVehicleResponse(vehicle))
}

// GetVehicle handles fetching a single vehicle by ID
func (h *VehicleHandler) GetVehicle(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "GetVehicle")
	ctx := r.Context()

	// Extract vehicle ID from the URL
	vehicleID := r.PathValue("id")
	if vehicleID == "" {
		response.BadRequest(w, "vehicle ID is required")
		return
	}
	w.Header().Set("Content-Type", "application/json")

	// Retrieve vehicle
	vehicle, err := h.uc.GetVehicle(ctx, vehicleID)
	if err != nil {
		if errors.Is(err, device.ErrVehicleNotFound) {
			response.NotFound(w)
			return
		}
		logger.Error("error retrieving vehicle", "error", err)
		response.Custom(w, http.StatusInternalServerError, nil, "internal server error")
		return
	}

	response.Pure(w, http.StatusOK, dto.NewVehicleResponse(vehicle))
}

// ListVehicles handles listing all vehicles
func (h *VehicleHandler) ListVehicles(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "ListVehicles")
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	// List vehicles
	vehicles, err := h.uc.ListVehicles(ctx)
	if err != nil {
		logger.Error("error listing vehicles", "error", err)
		response.Custom(w, http.StatusInternalServerError, nil, "internal server error")
		return
	}

	response.Pure(w, http.StatusOK, dto.NewVehicleListResponse(vehicles))
}

// DeleteVehicle handles deleting a vehicle by ID
func (h *VehicleHandler) DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "DeleteVehicle")
	ctx := r.Context()

	// Extract vehicle ID from the URL
	vehicleID := r.PathValue("id")
	if vehicleID == "" {
		response.BadRequest(w, "vehicle ID is required")
		return
	}

	// Delete vehicle
	if err := h.uc.DeleteVehicle(ctx, vehicleID); err != nil {
		if errors.Is(err, device.ErrVehicleNotFound) {
			response.NotFound(w)
			return
		}
		logger.Error("error deleting vehicle", "error", err)
		response.Custom(w, http.StatusInternalServerError, nil, "internal server error")
		return
	}

	response.Pure(w, http.StatusOK, nil)
}

// RegisterVehicleHandler registers vehicle routes with an HTTP mux
func (h *VehicleHandler) RegisterMuxRouter(mux *http.ServeMux) {
	// Create
	mux.HandleFunc("POST /api/v1/vehicles", h.CreateVehicle)
	// Get
	mux.HandleFunc("/api/v1/vehicles/{id}", h.GetVehicle)
	// List
	mux.HandleFunc("/api/v1/vehicles", h.ListVehicles)
	// Delete
	mux.HandleFunc("DELETE /api/v1/vehicles/{id}", h.DeleteVehicle)
}
