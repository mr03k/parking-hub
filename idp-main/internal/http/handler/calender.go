package handler

import (
	"log/slog"
	"net/http"

	"application/internal/http/dto"
	"application/internal/http/response"
	"application/internal/repo"

	"github.com/swaggest/openapi-go"
)

// WorkCalendarHandler represents the HTTP handler for work calendar records.
type WorkCalendarHandler struct {
	logger *slog.Logger
	repo   *repo.WorkCalendarRepository
}

// NewWorkCalendarHandler creates a new instance of WorkCalendarHandler.
func NewWorkCalendarHandler(logger *slog.Logger, repo *repo.WorkCalendarRepository) *WorkCalendarHandler {
	return &WorkCalendarHandler{
		logger: logger.With("layer", "WorkCalendarHandler"),
		repo:   repo,
	}
}

// ListWorkCalendars handles listing all work calendar records and responds with a DTO.
func (h *WorkCalendarHandler) ListWorkCalendars(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "ListWorkCalendars")

	logger.Info("Fetching list of work calendar records")
	workCalendars := h.repo.GetWorkCalendars()
	w.Header().Set("Content-Type", "application/json")

	// Convert work calendar records to DTO
	calendarResponses := make([]dto.WorkCalendarResponse, len(workCalendars))
	for i, calendar := range workCalendars {
		calendarResponses[i] = dto.WorkCalendarResponse{
			ID:          calendar.CalendarID.String(),
			ShamsiDate:  calendar.ShamsiDate,
			WorkDate:    calendar.WorkDate,
			Weekday:     calendar.Weekday,
			Year:        calendar.Year,
			IsHoliday:   calendar.IsHoliday,
			WorkShift:   calendar.WorkShift,
			Description: calendar.Description,
		}
	}

	response.Pure(w, http.StatusOK, dto.WorkCalendarListResponse{WorkCalendars: calendarResponses})
}

// GetWorkCalendarDetail handles fetching details of a single work calendar by ID and responds with a DTO.
func (h *WorkCalendarHandler) GetWorkCalendarDetail(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "GetWorkCalendarDetail")
	w.Header().Set("Content-Type", "application/json")

	// Parse the ID from the query parameters
	calendarID := r.PathValue("id")
	if calendarID == "" {
		logger.Error("missing calendar ID in request")
		response.BadRequest(w, "missing calendar ID")
		return
	}

	logger.Info("Fetching details for work calendar", "calendarID", calendarID)
	calendar, err := h.repo.GetWorkCalendarByID(calendarID)
	if err != nil {
		logger.Error("work calendar not found", "calendarID", calendarID)
		response.NotFound(w)
		return
	}

	// Convert work calendar to DTO
	calendarResponse := dto.WorkCalendarDetailResponse{
		ID:             calendar.CalendarID.String(),
		ContractID:     calendar.ContractID.String(),
		ShamsiDate:     calendar.ShamsiDate,
		WorkDate:       calendar.WorkDate,
		Weekday:        calendar.Weekday,
		Year:           calendar.Year,
		IsHoliday:      calendar.IsHoliday,
		WorkShift:      calendar.WorkShift,
		Description:    calendar.Description,
		WorkShiftStart: calendar.WorkShiftStart,
		WorkShiftEnd:   calendar.WorkShiftEnd,
	}

	response.Pure(w, http.StatusOK, calendarResponse)
}

// RegisterMuxRouter registers the work calendar handler routes with an HTTP multiplexer.
func (h *WorkCalendarHandler) RegisterMuxRouter(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/work-calendars", h.ListWorkCalendars)          // List Work Calendars
	mux.HandleFunc("/api/v1/work-calendars/{id}", h.GetWorkCalendarDetail) // Get Work Calendar Detail
}

// ListWorkCalendars describes the OpenAPI specification for listing work calendars.
func (s *WorkCalendarHandler) ListWorkCalendarsOAPI(op openapi.OperationContext) {
	op.SetSummary("List Work Calendars")
	op.SetDescription("Fetch a list of all work calendar records with basic information.")

	op.AddRespStructure(new(dto.WorkCalendarListResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

// GetWorkCalendarDetail describes the OpenAPI specification for fetching a single work calendar's details.
func (s *WorkCalendarHandler) GetWorkCalendarDetailOapi(op openapi.OperationContext) {
	op.SetSummary("Get Work Calendar Detail")
	op.SetDescription("Fetch detailed information for a specific work calendar by its unique ID.")

	op.AddReqStructure(new(WorkCalendarID), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(dto.WorkCalendarDetailResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusNotFound))
}

// OpenApiSpec registers the work calendar OpenAPI endpoints.
func (s *WorkCalendarHandler) OpenApiSpec(api OAPI) {
	workCalendarTags := WithTags("WorkCalendar")

	api.Register("GET", "/api/v1/work-calendars", s.ListWorkCalendarsOAPI, workCalendarTags)          // List Work Calendars
	api.Register("GET", "/api/v1/work-calendars/{id}", s.GetWorkCalendarDetailOapi, workCalendarTags) // Get Work Calendar Detail
}

// WorkCalendarID represents the path parameter for work calendar operations.
type WorkCalendarID struct {
	ID string `json:"id" path:"id"`
}
