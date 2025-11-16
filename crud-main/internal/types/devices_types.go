package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateDevicesRequest request params
type CreateDevicesRequest struct {
	CodeDevice          string    `json:"codeDevice" binding:""`
	NumberSerial        string    `json:"numberSerial" binding:""`
	Model               string    `json:"model" binding:""`
	DateInstallation    time.Time `json:"dateInstallation" binding:""`
	DateExpiryWarranty  time.Time `json:"dateExpiryWarranty" binding:""`
	DateExpiryInsurance time.Time `json:"dateExpiryInsurance" binding:""`
	ClassDevice         string    `json:"classDevice" binding:""`
	ImageContract       string    `json:"imageContract" binding:""`
	ImageInsurance      string    `json:"imageInsurance" binding:""`
	IDContractor        string    `json:"iDContractor" binding:""`
	Description         string    `json:"description" binding:""`
}

// UpdateDevicesByIDRequest request params
type UpdateDevicesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	CodeDevice          string    `json:"codeDevice" binding:""`
	NumberSerial        string    `json:"numberSerial" binding:""`
	Model               string    `json:"model" binding:""`
	DateInstallation    time.Time `json:"dateInstallation" binding:""`
	DateExpiryWarranty  time.Time `json:"dateExpiryWarranty" binding:""`
	DateExpiryInsurance time.Time `json:"dateExpiryInsurance" binding:""`
	ClassDevice         string    `json:"classDevice" binding:""`
	ImageContract       string    `json:"imageContract" binding:""`
	ImageInsurance      string    `json:"imageInsurance" binding:""`
	IDContractor        string    `json:"iDContractor" binding:""`
	Description         string    `json:"description" binding:""`
}

// DevicesObjDetail detail
type DevicesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	CodeDevice          string    `json:"codeDevice"`
	NumberSerial        string    `json:"numberSerial"`
	Model               string    `json:"model"`
	DateInstallation    time.Time `json:"dateInstallation"`
	DateExpiryWarranty  time.Time `json:"dateExpiryWarranty"`
	DateExpiryInsurance time.Time `json:"dateExpiryInsurance"`
	ClassDevice         string    `json:"classDevice"`
	ImageContract       string    `json:"imageContract"`
	ImageInsurance      string    `json:"imageInsurance"`
	IDContractor        string    `json:"iDContractor"`
	Description         string    `json:"description"`
	CreatedAt           int       `json:"createdAt"`
}

// CreateDevicesReply only for api docs
type CreateDevicesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteDevicesByIDReply only for api docs
type DeleteDevicesByIDReply struct {
	Result
}

// UpdateDevicesByIDReply only for api docs
type UpdateDevicesByIDReply struct {
	Result
}

// GetDevicesByIDReply only for api docs
type GetDevicesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Devices DevicesObjDetail `json:"devices"`
	} `json:"data"` // return data
}

// ListDevicessRequest request params
type ListDevicessRequest struct {
	query.Params
}

// ListDevicessReply only for api docs
type ListDevicessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Devicess []DevicesObjDetail `json:"devicess"`
	} `json:"data"` // return data
}
