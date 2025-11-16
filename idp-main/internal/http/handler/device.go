package handler

import (
	"log/slog"
	"net/http"

	"application/internal/http/dto"
	"application/internal/http/response"
	"application/internal/repo"

	"github.com/swaggest/openapi-go"
)

// DeviceHandler represents the HTTP handler for license plate reader devices.
type DeviceHandler struct {
	logger *slog.Logger
	repo   *repo.DeviceRepository
}

// NewDeviceHandler creates a new instance of DeviceHandler.
func NewDeviceHandler(logger *slog.Logger, repo *repo.DeviceRepository) *DeviceHandler {
	return &DeviceHandler{
		logger: logger.With("layer", "DeviceHandler"),
		repo:   repo,
	}
}

// ListDevices handles listing all devices and responds with a DTO.
func (h *DeviceHandler) ListDevices(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "ListDevices")
	w.Header().Set("Content-Type", "application/json")

	logger.Info("Fetching list of devices")
	devices := h.repo.GetDevices()

	// Convert devices to DTO
	deviceResponses := make([]dto.DeviceResponse, len(devices))
	for i, device := range devices {
		deviceResponses[i] = dto.DeviceResponse{
			ID:                  device.ID.String(),
			CodeDevice:          device.CodeDevice,
			NumberSerial:        device.NumberSerial,
			Model:               device.Model,
			DateInstallation:    device.DateInstallation,
			DateExpiryWarranty:  device.DateExpiryWarranty,
			DateExpiryInsurance: device.DateExpiryInsurance,
			ClassDevice:         device.ClassDevice,
			Status:              "Active", // Example status
		}
	}

	response.Pure(w, http.StatusOK, dto.DeviceListResponse{Devices: deviceResponses})
}

// GetDeviceDetail handles fetching details of a single device by ID and responds with a DTO.
func (h *DeviceHandler) GetDeviceDetail(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "GetDeviceDetail")
	w.Header().Set("Content-Type", "application/json")

	// Parse the ID from the query parameters
	deviceID := r.PathValue("id")
	if deviceID == "" {
		logger.Error("missing device ID in request")
		response.BadRequest(w, "missing device ID")
		return
	}

	logger.Info("Fetching details for device", "deviceID", deviceID)
	device, err := h.repo.GetDeviceByID(deviceID)
	if err != nil {
		logger.Error("device not found", "deviceID", deviceID)
		response.NotFound(w)
		return
	}

	// Convert device to DTO
	deviceResponse := dto.DeviceDetailResponse{
		ID:                  device.ID.String(),
		CodeDevice:          device.CodeDevice,
		NumberSerial:        device.NumberSerial,
		Model:               device.Model,
		DateInstallation:    device.DateInstallation,
		DateExpiryWarranty:  device.DateExpiryWarranty,
		DateExpiryInsurance: device.DateExpiryInsurance,
		ClassDevice:         device.ClassDevice,
		ImageContractURL:    device.ImageContractURL,
		ImageInsuranceURL:   device.ImageInsuranceURL,
		Description:         device.Description,
	}

	response.Pure(w, http.StatusOK, deviceResponse)
}

// RegisterMuxRouter registers the device handler routes with an HTTP multiplexer.
func (h *DeviceHandler) RegisterMuxRouter(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/devices", h.ListDevices)          // List Devices
	mux.HandleFunc("/api/v1/devices/{id}", h.GetDeviceDetail) // Get Device Detail
}

// DeviceID represents the path parameter for device operations.
type DeviceID struct {
	ID string `json:"id" path:"id"`
}

// ListDevices describes the OpenAPI specification for listing devices.
func (s *DeviceHandler) ListDevicesOAPI(op openapi.OperationContext) {
	op.SetSummary("List Devices")
	op.SetDescription("Fetch a list of all license plate reader devices with basic information.")

	op.AddRespStructure(new(dto.DeviceListResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

// GetDeviceDetail describes the OpenAPI specification for fetching a single device's details.
func (s *DeviceHandler) GetDeviceDetailOAPI(op openapi.OperationContext) {
	op.SetSummary("Get Device Detail")
	op.SetDescription("Fetch detailed information for a specific device by its unique ID.")

	op.AddReqStructure(new(DeviceID), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(dto.DeviceDetailResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusNotFound))
}

// OpenApiSpec registers the device OpenAPI endpoints.
func (s *DeviceHandler) OpenApiSpec(api OAPI) {
	deviceTags := WithTags("Device")

	api.Register("GET", "/api/v1/devices", s.ListDevicesOAPI, deviceTags)          // List Devices
	api.Register("GET", "/api/v1/devices/{id}", s.GetDeviceDetailOAPI, deviceTags) // Get Device Detail
}
