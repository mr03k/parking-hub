package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateBaseRatesRequest request params
type CreateBaseRatesRequest struct {
	IDCategoryVehicle string `json:"iDCategoryVehicle" binding:""`
	FromMinutes       int    `json:"fromMinutes" binding:""`
	ToMinutes         int    `json:"toMinutes" binding:""`
	BaseRate          string `json:"baseRate" binding:""`
	Description       string `json:"description" binding:""`
}

// UpdateBaseRatesByIDRequest request params
type UpdateBaseRatesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	IDCategoryVehicle string `json:"iDCategoryVehicle" binding:""`
	FromMinutes       int    `json:"fromMinutes" binding:""`
	ToMinutes         int    `json:"toMinutes" binding:""`
	BaseRate          string `json:"baseRate" binding:""`
	Description       string `json:"description" binding:""`
}

// BaseRatesObjDetail detail
type BaseRatesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	IDCategoryVehicle string `json:"iDCategoryVehicle"`
	FromMinutes       int    `json:"fromMinutes"`
	ToMinutes         int    `json:"toMinutes"`
	BaseRate          string `json:"baseRate"`
	Description       string `json:"description"`
	CreatedAt         int    `json:"createdAt"`
}

// CreateBaseRatesReply only for api docs
type CreateBaseRatesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteBaseRatesByIDReply only for api docs
type DeleteBaseRatesByIDReply struct {
	Result
}

// UpdateBaseRatesByIDReply only for api docs
type UpdateBaseRatesByIDReply struct {
	Result
}

// GetBaseRatesByIDReply only for api docs
type GetBaseRatesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		BaseRates BaseRatesObjDetail `json:"baseRates"`
	} `json:"data"` // return data
}

// ListBaseRatessRequest request params
type ListBaseRatessRequest struct {
	query.Params
}

// ListBaseRatessReply only for api docs
type ListBaseRatessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		BaseRatess []BaseRatesObjDetail `json:"baseRatess"`
	} `json:"data"` // return data
}
