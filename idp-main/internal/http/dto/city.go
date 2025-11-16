package dto

import (
	"encoding/json"
	"net/http"
)

// CityRequest represents a request for city
type CityRequest struct {
	Name          string `json:"city_name"`
	Code          string `json:"city_code"`
	ContryID      string `json:"contry_id"`
	GeoBoundaries string `json:"geo_boundaries"`
}

// CityResponse represents a response for city
type CityResponse struct {
	// City ID
	ID string `json:"city_id"`
	// City Name
	CityRequest
}

// NewCityRequestFromRequest creates a new CityRequest from http request
func NewCityRequestFromRequest(r *http.Request) (*CityRequest, error) {
	req := new(CityRequest)
	err := json.NewDecoder(r.Body).Decode(req)
	return req, err
}

type CityListResponse struct {
	Cities []CityResponse `json:"cities"`
	Count  int            `json:"count"`
}
