package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreatePeakHourMultipliersRequest request params
type CreatePeakHourMultipliersRequest struct {
	CodeTimePeak string    `json:"codeTimePeak" binding:""`
	Description  string    `json:"description" binding:""`
	Multiplier   string    `json:"multiplier" binding:""`
	Weekday      string    `json:"weekday" binding:""`
	TimeStart    string    `json:"timeStart" binding:""`
	TimeEnd      string    `json:"timeEnd" binding:""`
	FromValid    time.Time `json:"fromValid" binding:""`
	ToValid      time.Time `json:"toValid" binding:""`
	Flag         string    `json:"flag" binding:""`
}

// UpdatePeakHourMultipliersByIDRequest request params
type UpdatePeakHourMultipliersByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	CodeTimePeak string    `json:"codeTimePeak" binding:""`
	Description  string    `json:"description" binding:""`
	Multiplier   string    `json:"multiplier" binding:""`
	Weekday      string    `json:"weekday" binding:""`
	TimeStart    string    `json:"timeStart" binding:""`
	TimeEnd      string    `json:"timeEnd" binding:""`
	FromValid    time.Time `json:"fromValid" binding:""`
	ToValid      time.Time `json:"toValid" binding:""`
	Flag         string    `json:"flag" binding:""`
}

// PeakHourMultipliersObjDetail detail
type PeakHourMultipliersObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	CodeTimePeak string    `json:"codeTimePeak"`
	Description  string    `json:"description"`
	Multiplier   string    `json:"multiplier"`
	Weekday      string    `json:"weekday"`
	TimeStart    string    `json:"timeStart"`
	TimeEnd      string    `json:"timeEnd"`
	FromValid    time.Time `json:"fromValid"`
	ToValid      time.Time `json:"toValid"`
	Flag         string    `json:"flag"`
	CreatedAt    int       `json:"createdAt"`
}

// CreatePeakHourMultipliersReply only for api docs
type CreatePeakHourMultipliersReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeletePeakHourMultipliersByIDReply only for api docs
type DeletePeakHourMultipliersByIDReply struct {
	Result
}

// UpdatePeakHourMultipliersByIDReply only for api docs
type UpdatePeakHourMultipliersByIDReply struct {
	Result
}

// GetPeakHourMultipliersByIDReply only for api docs
type GetPeakHourMultipliersByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		PeakHourMultipliers PeakHourMultipliersObjDetail `json:"peakHourMultipliers"`
	} `json:"data"` // return data
}

// ListPeakHourMultiplierssRequest request params
type ListPeakHourMultiplierssRequest struct {
	query.Params
}

// ListPeakHourMultiplierssReply only for api docs
type ListPeakHourMultiplierssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		PeakHourMultiplierss []PeakHourMultipliersObjDetail `json:"peakHourMultiplierss"`
	} `json:"data"` // return data
}
