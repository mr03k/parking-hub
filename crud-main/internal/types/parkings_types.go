package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateParkingsRequest request params
type CreateParkingsRequest struct {
	CodeParking        string `json:"codeParking" binding:""`
	IDSegment          string `json:"iDSegment" binding:""`
	TypeParking        string `json:"typeParking" binding:""`
	BoundaryParking    string `json:"boundaryParking" binding:""`
	Position           string `json:"position" binding:""`
	StatusAvailability string `json:"statusAvailability" binding:""`
	Description        string `json:"description" binding:""`
}

// UpdateParkingsByIDRequest request params
type UpdateParkingsByIDRequest struct {
	ID string `json:"id" binding:""` // uint64 id

	CodeParking        string `json:"codeParking" binding:""`
	IDSegment          string `json:"iDSegment" binding:""`
	TypeParking        string `json:"typeParking" binding:""`
	BoundaryParking    string `json:"boundaryParking" binding:""`
	Position           string `json:"position" binding:""`
	StatusAvailability string `json:"statusAvailability" binding:""`
	Description        string `json:"description" binding:""`
}

// ParkingsObjDetail detail
type ParkingsObjDetail struct {
	ID string `json:"id"` // convert to uint64 id

	CodeParking        string `json:"codeParking"`
	IDSegment          string `json:"iDSegment"`
	TypeParking        string `json:"typeParking"`
	BoundaryParking    string `json:"boundaryParking"`
	Position           string `json:"position"`
	StatusAvailability string `json:"statusAvailability"`
	Description        string `json:"description"`
	CreatedAt          int    `json:"createdAt"`
}

// CreateParkingsReply only for api docs
type CreateParkingsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID string `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteParkingsByIDReply only for api docs
type DeleteParkingsByIDReply struct {
	Result
}

// UpdateParkingsByIDReply only for api docs
type UpdateParkingsByIDReply struct {
	Result
}

// GetParkingsByIDReply only for api docs
type GetParkingsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Parkings ParkingsObjDetail `json:"parkings"`
	} `json:"data"` // return data
}

// ListParkingssRequest request params
type ListParkingssRequest struct {
	query.Params
}

// ListParkingssReply only for api docs
type ListParkingssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Parkingss []ParkingsObjDetail `json:"parkingss"`
	} `json:"data"` // return data
}
