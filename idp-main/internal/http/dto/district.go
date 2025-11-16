package dto

import "time"

// DistrictResponse represents a summary of district information.
type DistrictResponse struct {
	ID           string  `json:"id"`
	DistrictName string  `json:"districtName"`
	DistrictCode string  `json:"districtCode"`
	CityID       string  `json:"cityId"`
	GeoBoundary  string  `json:"geoBoundary,omitempty"`
	Population   int64   `json:"population,omitempty"`
	Area         float64 `json:"area,omitempty"`
}

// DistrictDetailResponse represents detailed district information.
type DistrictDetailResponse struct {
	ID           string    `json:"id"`
	DistrictName string    `json:"districtName"`
	DistrictCode string    `json:"districtCode"`
	CityID       string    `json:"cityId"`
	GeoBoundary  string    `json:"geoBoundary,omitempty"`
	Population   int64     `json:"population,omitempty"`
	Area         float64   `json:"area,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
}

// DistrictListResponse represents the response for a list of districts.
type DistrictListResponse struct {
	Districts []DistrictResponse `json:"districts"`
}
