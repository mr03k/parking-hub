package oapi

import (
	"application/internal/http/dto"
	"application/internal/http/response"
	"github.com/swaggest/openapi-go"
	"log/slog"
	"net/http"
)

// WorkCalendarOpenApi represents the OpenAPI specification for work calendar endpoints.
type WorkCalendarOpenApi struct {
	logger *slog.Logger
}

// WorkCalendarID represents the path parameter for work calendar operations.
type WorkCalendarID struct {
	ID string `json:"id" path:"id"`
}

// NewWorkCalendarOpenApi creates a new instance of WorkCalendarOpenApi.
func NewWorkCalendarOpenApi(logger *slog.Logger) *WorkCalendarOpenApi {
	return &WorkCalendarOpenApi{
		logger: logger.With("layer", "WorkCalendarOpenApi"),
	}
}

// ListWorkCalendars describes the OpenAPI specification for listing work calendars.
func (s *WorkCalendarOpenApi) ListWorkCalendars(op openapi.OperationContext) {
	op.SetSummary("List Work Calendars")
	op.SetDescription("Fetch a list of all work calendar records with basic information.")

	op.AddRespStructure(new(dto.WorkCalendarListResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

// GetWorkCalendarDetail describes the OpenAPI specification for fetching a single work calendar's details.
func (s *WorkCalendarOpenApi) GetWorkCalendarDetail(op openapi.OperationContext) {
	op.SetSummary("Get Work Calendar Detail")
	op.SetDescription("Fetch detailed information for a specific work calendar by its unique ID.")

	op.AddReqStructure(new(WorkCalendarID), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(dto.WorkCalendarDetailResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusNotFound))
}

// OpenApiSpec registers the work calendar OpenAPI endpoints.
func (s *WorkCalendarOpenApi) OpenApiSpec(api OAPI) {
	workCalendarTags := WithTags("WorkCalendar")

	api.Register("GET", "/api/v1/work-calendars", s.ListWorkCalendars, workCalendarTags)          // List Work Calendars
	api.Register("GET", "/api/v1/work-calendars/{id}", s.GetWorkCalendarDetail, workCalendarTags) // Get Work Calendar Detail
}
