package handler

import (
	"log/slog"
	"net/http"

	"application/internal/datasource"
	"application/internal/http/dto"
	"application/internal/http/response"

	"github.com/swaggest/openapi-go"
)

// DriverHandler represents the HTTP handler for drivers.
type DriverHandler struct {
	logger *slog.Logger
	repo   *datasource.DriverRepository
}

// NewDriverHandler creates a new instance of DriverHandler.
func NewDriverHandler(logger *slog.Logger, repo *datasource.DriverRepository) *DriverHandler {
	return &DriverHandler{
		logger: logger.With("layer", "DriverHandler"),
		repo:   repo,
	}
}

// ListDrivers handles listing all drivers and responds with a DTO.
func (h *DriverHandler) ListDrivers(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "ListDrivers")
	w.Header().Set("Content-Type", "application/json")

	logger.Info("Fetching list of drivers")
	drivers := h.repo.GetDrivers()

	// Convert drivers to DTO
	driverResponses := make([]dto.DriverResponse, len(drivers))
	for i, driver := range drivers {
		driverResponses[i] = dto.DriverResponse{
			ID:               driver.ID,
			Address:          driver.Address,
			DriverType:       driver.DriverType,
			ShiftType:        driver.ShiftType,
			EmploymentStatus: driver.EmploymentStatus,
		}
	}

	response.Pure(w, http.StatusOK, dto.DriverListResponse{Drivers: driverResponses})
}

// GetDriverDetail handles fetching details of a single driver by ID and responds with a DTO.
func (h *DriverHandler) GetDriverDetail(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "GetDriverDetail")
	w.Header().Set("Content-Type", "application/json")

	// Parse the ID from the query parameters
	driverID := r.PathValue("id")
	if driverID == "" {
		logger.Error("missing driver ID in request")
		response.BadRequest(w, "missing driver ID")
		return
	}

	logger.Info("Fetching details for driver", "driverID", driverID)
	driver, err := h.repo.GetDriverByID(driverID)
	if err != nil {
		logger.Error("driver not found", "driverID", driverID)
		response.NotFound(w)
		return
	}

	// Convert driver to DTO
	driverResponse := dto.DriverDetailResponse{
		ID:                          driver.ID,
		Address:                     driver.Address,
		DriverType:                  driver.DriverType,
		ShiftType:                   driver.ShiftType,
		EmploymentStatus:            driver.EmploymentStatus,
		EmploymentStartDate:         driver.EmploymentStartDate,
		EmploymentEndDate:           driver.EmploymentEndDate,
		DriverPhotoURL:              driver.DriverPhotoURL,
		IDCardImageURL:              driver.IDCardImageURL,
		BirthCertificateImageURL:    driver.BirthCertificateImageURL,
		MilitaryServiceCardImageURL: driver.MilitaryServiceCardImageURL,
		HealthCertificateImageURL:   driver.HealthCertificateImageURL,
		CriminalRecordImageURL:      driver.CriminalRecordImageURL,
		Description:                 driver.Description,
	}

	// Respond with the driver details DTO
	response.Pure(w, http.StatusOK, driverResponse)
}

// RegisterMuxRouter registers the driver handler routes with an HTTP multiplexer.
func (h *DriverHandler) RegisterMuxRouter(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/drivers", h.ListDrivers)          // List Drivers
	mux.HandleFunc("/api/v1/drivers/{id}", h.GetDriverDetail) // Get Driver Detail
}

type DriverID struct {
	ID string `json:"id" path:"id"`
}

// List Drivers
func (s *DriverHandler) ListDriversOAPI(op openapi.OperationContext) {
	op.SetSummary("List Drivers")
	op.SetDescription("Fetch a list of all drivers with their basic information.")

	op.AddRespStructure(new(dto.DriverListResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

// Get Driver By ID
func (s *DriverHandler) GetDriverOAPI(op openapi.OperationContext) {
	op.SetSummary("Get Driver By ID")
	op.SetDescription("Fetch detailed information for a specific driver by their unique ID.")

	op.AddReqStructure(new(DriverID), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(dto.DriverDetailResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusNotFound))
}

// OpenApiSpec registers the Driver OpenAPI endpoints.
func (s *DriverHandler) OpenApiSpec(api OAPI) {
	driverTags := WithTags("Driver")

	api.Register("GET", "/api/v1/drivers", s.ListDriversOAPI, driverTags)
	api.Register("GET", "/api/v1/drivers/{id}", s.GetDriverOAPI, driverTags)
}
