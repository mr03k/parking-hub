package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateVehicleCategoriesRequest request params
type CreateVehicleCategoriesRequest struct {
	CodeCategory string `json:"codeCategory" binding:""`
	NameCategory string `json:"nameCategory" binding:""`
	Description  string `json:"description" binding:""`
}

// UpdateVehicleCategoriesByIDRequest request params
type UpdateVehicleCategoriesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	CodeCategory string `json:"codeCategory" binding:""`
	NameCategory string `json:"nameCategory" binding:""`
	Description  string `json:"description" binding:""`
}

// VehicleCategoriesObjDetail detail
type VehicleCategoriesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	CodeCategory string `json:"codeCategory"`
	NameCategory string `json:"nameCategory"`
	Description  string `json:"description"`
	CreatedAt    int    `json:"createdAt"`
}

// CreateVehicleCategoriesReply only for api docs
type CreateVehicleCategoriesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteVehicleCategoriesByIDReply only for api docs
type DeleteVehicleCategoriesByIDReply struct {
	Result
}

// UpdateVehicleCategoriesByIDReply only for api docs
type UpdateVehicleCategoriesByIDReply struct {
	Result
}

// GetVehicleCategoriesByIDReply only for api docs
type GetVehicleCategoriesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		VehicleCategories VehicleCategoriesObjDetail `json:"vehicleCategories"`
	} `json:"data"` // return data
}

// ListVehicleCategoriessRequest request params
type ListVehicleCategoriessRequest struct {
	query.Params
}

// ListVehicleCategoriessReply only for api docs
type ListVehicleCategoriessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		VehicleCategoriess []VehicleCategoriesObjDetail `json:"vehicleCategoriess"`
	} `json:"data"` // return data
}
