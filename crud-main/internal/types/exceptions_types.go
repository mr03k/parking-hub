package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateExceptionsRequest request params
type CreateExceptionsRequest struct {
	CarLicensePlates        string    `json:"carLicensePlates" binding:""`
	MotorcycleLicensePlates string    `json:"motorcycleLicensePlates" binding:""`
	ExceptionMultiplier     string    `json:"exceptionMultiplier" binding:""`
	StartDate               time.Time `json:"startDate" binding:""`
	EndDate                 time.Time `json:"endDate" binding:""`
	Description             string    `json:"description" binding:""`
	NotificationNumber      string    `json:"notificationNumber" binding:""`
	NotificationDate        time.Time `json:"notificationDate" binding:""`
	DocumentImage           string    `json:"documentImage" binding:""`
	UserID                  string    `json:"userID" binding:""`
	VehicleType             string    `json:"vehicleType" binding:""`
}

// UpdateExceptionsByIDRequest request params
type UpdateExceptionsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	CarLicensePlates        string    `json:"carLicensePlates" binding:""`
	MotorcycleLicensePlates string    `json:"motorcycleLicensePlates" binding:""`
	ExceptionMultiplier     string    `json:"exceptionMultiplier" binding:""`
	StartDate               time.Time `json:"startDate" binding:""`
	EndDate                 time.Time `json:"endDate" binding:""`
	Description             string    `json:"description" binding:""`
	NotificationNumber      string    `json:"notificationNumber" binding:""`
	NotificationDate        time.Time `json:"notificationDate" binding:""`
	DocumentImage           string    `json:"documentImage" binding:""`
	UserID                  string    `json:"userID" binding:""`
	VehicleType             string    `json:"vehicleType" binding:""`
}

// ExceptionsObjDetail detail
type ExceptionsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	CarLicensePlates        string    `json:"carLicensePlates"`
	MotorcycleLicensePlates string    `json:"motorcycleLicensePlates"`
	ExceptionMultiplier     string    `json:"exceptionMultiplier"`
	StartDate               time.Time `json:"startDate"`
	EndDate                 time.Time `json:"endDate"`
	Description             string    `json:"description"`
	NotificationNumber      string    `json:"notificationNumber"`
	NotificationDate        time.Time `json:"notificationDate"`
	DocumentImage           string    `json:"documentImage"`
	UserID                  string    `json:"userID"`
	VehicleType             string    `json:"vehicleType"`
	CreatedAt               int       `json:"createdAt"`
}

// CreateExceptionsReply only for api docs
type CreateExceptionsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteExceptionsByIDReply only for api docs
type DeleteExceptionsByIDReply struct {
	Result
}

// UpdateExceptionsByIDReply only for api docs
type UpdateExceptionsByIDReply struct {
	Result
}

// GetExceptionsByIDReply only for api docs
type GetExceptionsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Exceptions ExceptionsObjDetail `json:"exceptions"`
	} `json:"data"` // return data
}

// ListExceptionssRequest request params
type ListExceptionssRequest struct {
	query.Params
}

// ListExceptionssReply only for api docs
type ListExceptionssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Exceptionss []ExceptionsObjDetail `json:"exceptionss"`
	} `json:"data"` // return data
}
