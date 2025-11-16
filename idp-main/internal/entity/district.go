package entity

import (
	"time"

	"github.com/google/uuid"
)

// District represents the entity for a geographical district.
type District struct {
	ID           uuid.UUID `json:"id"`                    // Unique identifier for the district
	DistrictName string    `json:"districtName"`          // Name of the district
	DistrictCode string    `json:"districtCode"`          // Unique code for the district
	CityID       uuid.UUID `json:"cityId"`                // Foreign key to the cities table
	GeoBoundary  string    `json:"geoBoundary,omitempty"` // Geographic boundary of the district as a polygon
	Population   int64     `json:"population,omitempty"`  // Population of the district (optional)
	Area         float64   `json:"area,omitempty"`        // Area of the district in square kilometers (optional)
	CreatedAt    time.Time `json:"createdAt"`             // Timestamp when the record was created
}
