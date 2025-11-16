package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateRingsRequest request params
type CreateRingsRequest struct {
	NameRing             string `json:"nameRing" binding:""`
	CodeRing             string `json:"codeRing" binding:""`
	LengthRing           string `json:"lengthRing" binding:""`
	BoundaryRing         string `json:"boundaryRing" binding:""`
	SpotsParking         int    `json:"spotsParking" binding:""`
	SpotsParkingDisabled int    `json:"spotsParkingDisabled" binding:""`
	SignsTraffic         int    `json:"signsTraffic" binding:""`
	SignsTrafficDisabled int    `json:"signsTrafficDisabled" binding:""`
	PointStart           string `json:"pointStart" binding:""`
	DistanceBuffer       string `json:"distanceBuffer" binding:""`
	Description          string `json:"description" binding:""`
}

// UpdateRingsByIDRequest request params
type UpdateRingsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	NameRing             string `json:"nameRing" binding:""`
	CodeRing             string `json:"codeRing" binding:""`
	LengthRing           string `json:"lengthRing" binding:""`
	BoundaryRing         string `json:"boundaryRing" binding:""`
	SpotsParking         int    `json:"spotsParking" binding:""`
	SpotsParkingDisabled int    `json:"spotsParkingDisabled" binding:""`
	SignsTraffic         int    `json:"signsTraffic" binding:""`
	SignsTrafficDisabled int    `json:"signsTrafficDisabled" binding:""`
	PointStart           string `json:"pointStart" binding:""`
	DistanceBuffer       string `json:"distanceBuffer" binding:""`
	Description          string `json:"description" binding:""`
}

// RingsObjDetail detail
type RingsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	NameRing             string `json:"nameRing"`
	CodeRing             string `json:"codeRing"`
	LengthRing           string `json:"lengthRing"`
	BoundaryRing         string `json:"boundaryRing"`
	SpotsParking         int    `json:"spotsParking"`
	SpotsParkingDisabled int    `json:"spotsParkingDisabled"`
	SignsTraffic         int    `json:"signsTraffic"`
	SignsTrafficDisabled int    `json:"signsTrafficDisabled"`
	PointStart           string `json:"pointStart"`
	DistanceBuffer       string `json:"distanceBuffer"`
	CreatedAt            int    `json:"createdAt"`
	Description          string `json:"description"`
}

// CreateRingsReply only for api docs
type CreateRingsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteRingsByIDReply only for api docs
type DeleteRingsByIDReply struct {
	Result
}

// UpdateRingsByIDReply only for api docs
type UpdateRingsByIDReply struct {
	Result
}

// GetRingsByIDReply only for api docs
type GetRingsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Rings RingsObjDetail `json:"rings"`
	} `json:"data"` // return data
}

// ListRingssRequest request params
type ListRingssRequest struct {
	query.Params
}

// ListRingssReply only for api docs
type ListRingssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Ringss []RingsObjDetail `json:"ringss"`
	} `json:"data"` // return data
}
