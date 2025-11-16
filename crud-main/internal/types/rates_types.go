package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateRatesRequest request params
type CreateRatesRequest struct {
	Code                   string    `json:"code" binding:""`
	RoadCategoryID         string    `json:"roadCategoryID" binding:""`
	TimeCycleMinutes       int       `json:"timeCycleMinutes" binding:""`
	RateMultiplier         string    `json:"rateMultiplier" binding:""`
	PeakHourMultiplier     string    `json:"peakHourMultiplier" binding:""`
	GoodPercentage         int       `json:"goodPercentage" binding:""`
	NormalSettlementPeriod int       `json:"normalSettlementPeriod" binding:""`
	LatePenalty            string    `json:"latePenalty" binding:""`
	LatePenaltyMax         string    `json:"latePenaltyMax" binding:""`
	ValidFrom              time.Time `json:"validFrom" binding:""`
	ValidTo                time.Time `json:"validTo" binding:""`
	Description            string    `json:"description" binding:""`
	StartTime              string    `json:"startTime" binding:""`
	EndTime                string    `json:"endTime" binding:""`
	CityID                 string    `json:"cityID" binding:""`
	ApprovalNumber         string    `json:"approvalNumber" binding:""`
	ApprovalDate           time.Time `json:"approvalDate" binding:""`
	Year                   int       `json:"year" binding:""`
	BaseRateID             string    `json:"baseRateID" binding:""`
	ExceptionsID           string    `json:"exceptionsID" binding:""`
}

// UpdateRatesByIDRequest request params
type UpdateRatesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	Code                   string    `json:"code" binding:""`
	RoadCategoryID         string    `json:"roadCategoryID" binding:""`
	TimeCycleMinutes       int       `json:"timeCycleMinutes" binding:""`
	RateMultiplier         string    `json:"rateMultiplier" binding:""`
	PeakHourMultiplier     string    `json:"peakHourMultiplier" binding:""`
	GoodPercentage         int       `json:"goodPercentage" binding:""`
	NormalSettlementPeriod int       `json:"normalSettlementPeriod" binding:""`
	LatePenalty            string    `json:"latePenalty" binding:""`
	LatePenaltyMax         string    `json:"latePenaltyMax" binding:""`
	ValidFrom              time.Time `json:"validFrom" binding:""`
	ValidTo                time.Time `json:"validTo" binding:""`
	Description            string    `json:"description" binding:""`
	StartTime              string    `json:"startTime" binding:""`
	EndTime                string    `json:"endTime" binding:""`
	CityID                 string    `json:"cityID" binding:""`
	ApprovalNumber         string    `json:"approvalNumber" binding:""`
	ApprovalDate           time.Time `json:"approvalDate" binding:""`
	Year                   int       `json:"year" binding:""`
	BaseRateID             string    `json:"baseRateID" binding:""`
	ExceptionsID           string    `json:"exceptionsID" binding:""`
}

// RatesObjDetail detail
type RatesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	Code                   string    `json:"code"`
	RoadCategoryID         string    `json:"roadCategoryID"`
	TimeCycleMinutes       int       `json:"timeCycleMinutes"`
	RateMultiplier         string    `json:"rateMultiplier"`
	PeakHourMultiplier     string    `json:"peakHourMultiplier"`
	GoodPercentage         int       `json:"goodPercentage"`
	NormalSettlementPeriod int       `json:"normalSettlementPeriod"`
	LatePenalty            string    `json:"latePenalty"`
	LatePenaltyMax         string    `json:"latePenaltyMax"`
	ValidFrom              time.Time `json:"validFrom"`
	ValidTo                time.Time `json:"validTo"`
	Description            string    `json:"description"`
	StartTime              string    `json:"startTime"`
	EndTime                string    `json:"endTime"`
	CityID                 string    `json:"cityID"`
	ApprovalNumber         string    `json:"approvalNumber"`
	ApprovalDate           time.Time `json:"approvalDate"`
	Year                   int       `json:"year"`
	BaseRateID             string    `json:"baseRateID"`
	ExceptionsID           string    `json:"exceptionsID"`
	CreatedAt              int       `json:"createdAt"`
}

// CreateRatesReply only for api docs
type CreateRatesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteRatesByIDReply only for api docs
type DeleteRatesByIDReply struct {
	Result
}

// UpdateRatesByIDReply only for api docs
type UpdateRatesByIDReply struct {
	Result
}

// GetRatesByIDReply only for api docs
type GetRatesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Rates RatesObjDetail `json:"rates"`
	} `json:"data"` // return data
}

// ListRatessRequest request params
type ListRatessRequest struct {
	query.Params
}

// ListRatessReply only for api docs
type ListRatessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Ratess []RatesObjDetail `json:"ratess"`
	} `json:"data"` // return data
}
