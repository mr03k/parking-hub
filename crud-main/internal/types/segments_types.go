package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateSegmentsRequest request params
type CreateSegmentsRequest struct {
	SegmentName          string `json:"segmentName" binding:""`
	SegmentCode          string `json:"segmentCode" binding:""`
	PartID               string `json:"partID" binding:""`
	RoadID               string `json:"roadID" binding:""`
	DistrictID           string `json:"districtID" binding:""`
	SegmentLength        string `json:"segmentLength" binding:""`
	SegmentBoundary      string `json:"segmentBoundary" binding:""`
	ParkingSpots         int    `json:"parkingSpots" binding:""`
	DisabledParkingSpots int    `json:"disabledParkingSpots" binding:""`
	Description          string `json:"description" binding:""`
}

// UpdateSegmentsByIDRequest request params
type UpdateSegmentsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	SegmentName          string `json:"segmentName" binding:""`
	SegmentCode          string `json:"segmentCode" binding:""`
	PartID               string `json:"partID" binding:""`
	RoadID               string `json:"roadID" binding:""`
	DistrictID           string `json:"districtID" binding:""`
	SegmentLength        string `json:"segmentLength" binding:""`
	SegmentBoundary      string `json:"segmentBoundary" binding:""`
	ParkingSpots         int    `json:"parkingSpots" binding:""`
	DisabledParkingSpots int    `json:"disabledParkingSpots" binding:""`
	Description          string `json:"description" binding:""`
}

// SegmentsObjDetail detail
type SegmentsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	SegmentName          string `json:"segmentName"`
	SegmentCode          string `json:"segmentCode"`
	PartID               string `json:"partID"`
	RoadID               string `json:"roadID"`
	DistrictID           string `json:"districtID"`
	SegmentLength        string `json:"segmentLength"`
	SegmentBoundary      string `json:"segmentBoundary"`
	ParkingSpots         int    `json:"parkingSpots"`
	DisabledParkingSpots int    `json:"disabledParkingSpots"`
	Description          string `json:"description"`
	CreatedAt            int    `json:"createdAt"`
}

// CreateSegmentsReply only for api docs
type CreateSegmentsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteSegmentsByIDReply only for api docs
type DeleteSegmentsByIDReply struct {
	Result
}

// UpdateSegmentsByIDReply only for api docs
type UpdateSegmentsByIDReply struct {
	Result
}

// GetSegmentsByIDReply only for api docs
type GetSegmentsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Segments SegmentsObjDetail `json:"segments"`
	} `json:"data"` // return data
}

// ListSegmentssRequest request params
type ListSegmentssRequest struct {
	query.Params
}

// ListSegmentssReply only for api docs
type ListSegmentssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Segmentss []SegmentsObjDetail `json:"segmentss"`
	} `json:"data"` // return data
}
