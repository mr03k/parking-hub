package dto

// RingIDRequest represents the request structure for fetching a ring by ID.
type RingIDRequest struct {
	ID string `json:"id"`
}

// RingListResponse represents the response structure for listing rings.
type RingListResponse struct {
	Rings []RingResponse `json:"rings"`
}

// RingDetailResponse represents the response structure for ring details.
type RingDetailResponse struct {
	ID                   string  `json:"id"`
	RingName             string  `json:"ring_name"`
	RingCode             string  `json:"ring_code"`
	RingLength           float64 `json:"ring_length"`
	RingBoundary         string  `json:"ring_boundary"`
	ParkingSpots         int     `json:"parking_spots"`
	DisabledParkingSpots int     `json:"disabled_parking_spots"`
	TrafficSigns         int     `json:"traffic_signs"`
	DisabledTrafficSigns int     `json:"disabled_traffic_signs"`
	StartPoint           string  `json:"start_point"`
	BufferDistance       float64 `json:"buffer_distance"`
	Description          string  `json:"description"`
}

// RingResponse is a simplified representation of a ring in a list.
type RingResponse struct {
	ID       string `json:"id"`
	RingName string `json:"ring_name"`
	RingCode string `json:"ring_code"`
}
