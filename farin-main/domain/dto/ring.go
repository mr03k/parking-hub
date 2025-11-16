package dto

import (
	"farin/domain/entity"
)

type RingRequest struct {
	ID       int64   `json:"id" binding:"required"`
	RingCode string  `json:"ringCode" binding:"required,min=1,max=100"`
	Length   float64 `json:"length" binding:"required,min=0"`
	RingName string  `json:"ringName" binding:"required,min=1,max=120"`
	Geom     string  `json:"geom" binding:"required"` // Assuming geometry is stored as a string (e.g., WKT format)
}

func (req *RingRequest) ToEntity() *entity.Ring {
	return &entity.Ring{
		ID:       req.ID,
		RingCode: req.RingCode,
		Length:   req.Length,
		RingName: req.RingName,
		Geom:     req.Geom, // Assuming geometry comes as a WKT string or similar
	}
}

type RingResponse struct {
	ID        int64   `json:"id"`
	RingCode  string  `json:"ringCode"`
	Length    float64 `json:"length"`
	RingName  string  `json:"ringName"`
	Geom      string  `json:"geom"`
	CreatedAt int64   `json:"createdAt"`
	UpdatedAt int64   `json:"updatedAt"`
}

func (resp *RingResponse) FromEntity(ring *entity.Ring) {
	resp.ID = ring.ID
	resp.RingCode = ring.RingCode
	resp.Length = ring.Length
	resp.RingName = ring.RingName
	resp.Geom = ring.Geom // Assuming the geometry is returned as the first element in the array
	resp.CreatedAt = ring.CreatedAt
	resp.UpdatedAt = ring.UpdatedAt
}

type RingListResponse struct {
	Rings []RingResponse `json:"rings"`
	Total int64          `json:"total"`
}
