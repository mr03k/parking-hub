package oapi

import (
	"application/internal/http/dto"
	"application/internal/http/response"
	"log/slog"
	"net/http"

	"github.com/swaggest/openapi-go"
)

// DeviceOpenApi represents the OpenAPI specification for device endpoints.
type DeviceOpenApi struct {
	logger *slog.Logger
}

// DeviceID represents the path parameter for device operations.
type DeviceID struct {
	ID string `json:"id" path:"id"`
}

// NewDeviceOpenApi creates a new instance of DeviceOpenApi.
func NewDeviceOpenApi(logger *slog.Logger) *DeviceOpenApi {
	return &DeviceOpenApi{
		logger: logger.With("layer", "DeviceOpenApi"),
	}
}

// ListDevices describes the OpenAPI specification for listing devices.
func (s *DeviceOpenApi) ListDevices(op openapi.OperationContext) {
	op.SetSummary("List Devices")
	op.SetDescription("Fetch a list of all license plate reader devices with basic information.")

	op.AddRespStructure(new(dto.DeviceListResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

// GetDeviceDetail describes the OpenAPI specification for fetching a single device's details.
func (s *DeviceOpenApi) GetDeviceDetail(op openapi.OperationContext) {
	op.SetSummary("Get Device Detail")
	op.SetDescription("Fetch detailed information for a specific device by its unique ID.")

	op.AddReqStructure(new(DeviceID), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(dto.DeviceDetailResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusNotFound))
}

// OpenApiSpec registers the device OpenAPI endpoints.
func (s *DeviceOpenApi) OpenApiSpec(api OAPI) {
	deviceTags := WithTags("Device")

	api.Register("GET", "/api/v1/devices", s.ListDevices, deviceTags)          // List Devices
	api.Register("GET", "/api/v1/devices/{id}", s.GetDeviceDetail, deviceTags) // Get Device Detail
}
