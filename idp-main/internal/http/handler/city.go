package handler

import (
	"log/slog"
	"net/http"

	"application/internal/http/dto"
	"application/internal/http/response"

	"github.com/google/uuid"
	"github.com/swaggest/openapi-go"
)

type CityHandler struct {
	logger   *slog.Logger
	mockCity *dto.CityResponse
}

var _ Handler = (*CityHandler)(nil)

func NewCityHandler(logger *slog.Logger) *CityHandler {
	return &CityHandler{
		logger: logger.With("layer", "CityHandler"),
		mockCity: &dto.CityResponse{
			ID: uuid.NewString(),
			CityRequest: dto.CityRequest{
				Name:          "Jakarta",
				Code:          "JKT",
				ContryID:      "1",
				GeoBoundaries: "POINT(106.8272 -6.1751)",
			},
		},
	}
}

// list city
func (h *CityHandler) List(w http.ResponseWriter, r *http.Request) {
	list := dto.CityListResponse{
		Count:  1,
		Cities: []dto.CityResponse{*h.mockCity},
	}

	w.Header().Set("Content-Type", "application/json")
	response.Pure(w, http.StatusOK, list)
}

// get city by id
func (h *CityHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = r.PathValue("city_id")

	response.Pure(w, http.StatusOK, h.mockCity)
}

// register Mux for City
func (h *CityHandler) RegisterMuxRouter(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/cities", h.List)
	mux.HandleFunc("GET /api/v1/cities/{city_id}", h.Get)
}

// get city by id
func (s *CityHandler) GetCityOAPI(op openapi.OperationContext) {
	op.SetSummary("Get City By ID")
	op.AddReqStructure(new(CityID), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(dto.CityResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

// list city

func (s *CityHandler) ListCityOAPI(op openapi.OperationContext) {
	op.SetSummary("List City")
	op.AddRespStructure(new(dto.CityListResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

func (s *CityHandler) OpenApiSpec(api OAPI) {
	tags := WithTags("city")

	api.Register("GET", "/api/v1/cities", s.ListCityOAPI, tags)
	api.Register("GET", "/api/v1/cities/{city_id}", s.GetCityOAPI, tags)
}

type CityID struct {
	CityID string `json:"city_id" path:"city_id"`
}
