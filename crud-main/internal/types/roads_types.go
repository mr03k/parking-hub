package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateRoadsRequest request params
type CreateRoadsRequest struct {
	RoadName             string `json:"roadName" binding:""`
	CodeRoad             string `json:"codeRoad" binding:""`
	TypeRoad             string `json:"typeRoad" binding:""`
	GradeRoad            string `json:"gradeRoad" binding:""`
	LengthRoad           string `json:"lengthRoad" binding:""`
	WidthRoad            string `json:"widthRoad" binding:""`
	LimitSpeed           int    `json:"limitSpeed" binding:""`
	BoundaryRoad         string `json:"boundaryRoad" binding:""`
	SpotsParking         int    `json:"spotsParking" binding:""`
	SpotsParkingDisabled int    `json:"spotsParkingDisabled" binding:""`
	Description          string `json:"description" binding:""`
}

// UpdateRoadsByIDRequest request params
type UpdateRoadsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	RoadName             string `json:"roadName" binding:""`
	CodeRoad             string `json:"codeRoad" binding:""`
	TypeRoad             string `json:"typeRoad" binding:""`
	GradeRoad            string `json:"gradeRoad" binding:""`
	LengthRoad           string `json:"lengthRoad" binding:""`
	WidthRoad            string `json:"widthRoad" binding:""`
	LimitSpeed           int    `json:"limitSpeed" binding:""`
	BoundaryRoad         string `json:"boundaryRoad" binding:""`
	SpotsParking         int    `json:"spotsParking" binding:""`
	SpotsParkingDisabled int    `json:"spotsParkingDisabled" binding:""`
	Description          string `json:"description" binding:""`
}

// RoadsObjDetail detail
type RoadsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	RoadName             string `json:"roadName"`
	CodeRoad             string `json:"codeRoad"`
	TypeRoad             string `json:"typeRoad"`
	GradeRoad            string `json:"gradeRoad"`
	LengthRoad           string `json:"lengthRoad"`
	WidthRoad            string `json:"widthRoad"`
	LimitSpeed           int    `json:"limitSpeed"`
	BoundaryRoad         string `json:"boundaryRoad"`
	SpotsParking         int    `json:"spotsParking"`
	SpotsParkingDisabled int    `json:"spotsParkingDisabled"`
	CreatedAt            int    `json:"createdAt"`
	Description          string `json:"description"`
}

// CreateRoadsReply only for api docs
type CreateRoadsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteRoadsByIDReply only for api docs
type DeleteRoadsByIDReply struct {
	Result
}

// UpdateRoadsByIDReply only for api docs
type UpdateRoadsByIDReply struct {
	Result
}

// GetRoadsByIDReply only for api docs
type GetRoadsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Roads RoadsObjDetail `json:"roads"`
	} `json:"data"` // return data
}

// ListRoadssRequest request params
type ListRoadssRequest struct {
	query.Params
}

// ListRoadssReply only for api docs
type ListRoadssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Roadss []RoadsObjDetail `json:"roadss"`
	} `json:"data"` // return data
}
