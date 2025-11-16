package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreatePartsRequest request params
type CreatePartsRequest struct {
	PartName             string `json:"partName" binding:""`
	CodePart             string `json:"codePart" binding:""`
	IDRoad               string `json:"iDRoad" binding:""`
	LengthPart           string `json:"lengthPart" binding:""`
	BoundaryPart         string `json:"boundaryPart" binding:""`
	SpotsParking         int    `json:"spotsParking" binding:""`
	SpotsParkingDisabled int    `json:"spotsParkingDisabled" binding:""`
	Description          string `json:"description" binding:""`
}

// UpdatePartsByIDRequest request params
type UpdatePartsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	PartName             string `json:"partName" binding:""`
	CodePart             string `json:"codePart" binding:""`
	IDRoad               string `json:"iDRoad" binding:""`
	LengthPart           string `json:"lengthPart" binding:""`
	BoundaryPart         string `json:"boundaryPart" binding:""`
	SpotsParking         int    `json:"spotsParking" binding:""`
	SpotsParkingDisabled int    `json:"spotsParkingDisabled" binding:""`
	Description          string `json:"description" binding:""`
}

// PartsObjDetail detail
type PartsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	PartName             string `json:"partName"`
	CodePart             string `json:"codePart"`
	IDRoad               string `json:"iDRoad"`
	LengthPart           string `json:"lengthPart"`
	BoundaryPart         string `json:"boundaryPart"`
	SpotsParking         int    `json:"spotsParking"`
	SpotsParkingDisabled int    `json:"spotsParkingDisabled"`
	Description          string `json:"description"`
}

// CreatePartsReply only for api docs
type CreatePartsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeletePartsByIDReply only for api docs
type DeletePartsByIDReply struct {
	Result
}

// UpdatePartsByIDReply only for api docs
type UpdatePartsByIDReply struct {
	Result
}

// GetPartsByIDReply only for api docs
type GetPartsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Parts PartsObjDetail `json:"parts"`
	} `json:"data"` // return data
}

// ListPartssRequest request params
type ListPartssRequest struct {
	query.Params
}

// ListPartssReply only for api docs
type ListPartssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Partss []PartsObjDetail `json:"partss"`
	} `json:"data"` // return data
}
