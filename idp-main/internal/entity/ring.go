package entity

import (
	"time"
)

// Ring represents the data model for a ring.
type Ring struct {
	ID                   string    `json:"id" db:"ring_id"`                                    // Unique ID for the ring
	RingName             string    `json:"ring_name" db:"ring_name"`                           // Name or number of the ring
	RingCode             string    `json:"ring_code" db:"ring_code"`                           // Unique code for identifying the ring
	RingLength           float64   `json:"ring_length" db:"ring_length"`                       // Length of the ring in kilometers
	RingBoundary         string    `json:"ring_boundary" db:"ring_boundary"`                   // Geographic boundary of the ring as a polygon
	ParkingSpots         int       `json:"parking_spots" db:"parking_spots"`                   // Number of regular parking spots in the ring
	DisabledParkingSpots int       `json:"disabled_parking_spots" db:"disabled_parking_spots"` // Number of disabled parking spots in the ring
	TrafficSigns         int       `json:"traffic_signs" db:"traffic_signs"`                   // Number of regular traffic signs in the ring
	DisabledTrafficSigns int       `json:"disabled_traffic_signs" db:"disabled_traffic_signs"` // Number of disabled traffic signs in the ring
	StartPoint           string    `json:"start_point" db:"start_point"`                       // Start point of the ring as a geographic point
	BufferDistance       float64   `json:"buffer_distance" db:"buffer_distance"`               // Buffer distance related to the start point in meters
	Description          string    `json:"description" db:"description"`                       // Additional description
	CreatedAt            time.Time `json:"created_at" db:"created_at"`                         // Record creation timestamp
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`                         // Last updated timestamp
}
