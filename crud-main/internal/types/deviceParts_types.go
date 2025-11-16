package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateDevicePartsRequest request params
type CreateDevicePartsRequest struct {
	CodePart           string    `json:"codePart" binding:""`
	PartName           string    `json:"partName" binding:""`
	TypePart           string    `json:"typePart" binding:""`
	Brand              string    `json:"brand" binding:""`
	Model              string    `json:"model" binding:""`
	NumberSerial       string    `json:"numberSerial" binding:""`
	DateInstallation   time.Time `json:"dateInstallation" binding:""`
	DateExpiryWarranty time.Time `json:"dateExpiryWarranty" binding:""`
	PeriodMaintenance  string    `json:"periodMaintenance" binding:""`
	IDDevice           string    `json:"iDDevice" binding:""`
	Description        string    `json:"description" binding:""`
}

// UpdateDevicePartsByIDRequest request params
type UpdateDevicePartsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	CodePart           string    `json:"codePart" binding:""`
	PartName           string    `json:"partName" binding:""`
	TypePart           string    `json:"typePart" binding:""`
	Brand              string    `json:"brand" binding:""`
	Model              string    `json:"model" binding:""`
	NumberSerial       string    `json:"numberSerial" binding:""`
	DateInstallation   time.Time `json:"dateInstallation" binding:""`
	DateExpiryWarranty time.Time `json:"dateExpiryWarranty" binding:""`
	PeriodMaintenance  string    `json:"periodMaintenance" binding:""`
	IDDevice           string    `json:"iDDevice" binding:""`
	Description        string    `json:"description" binding:""`
}

// DevicePartsObjDetail detail
type DevicePartsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	CodePart           string    `json:"codePart"`
	PartName           string    `json:"partName"`
	TypePart           string    `json:"typePart"`
	Brand              string    `json:"brand"`
	Model              string    `json:"model"`
	NumberSerial       string    `json:"numberSerial"`
	DateInstallation   time.Time `json:"dateInstallation"`
	DateExpiryWarranty time.Time `json:"dateExpiryWarranty"`
	PeriodMaintenance  string    `json:"periodMaintenance"`
	IDDevice           string    `json:"iDDevice"`
	Description        string    `json:"description"`
	CreatedAt          int       `json:"createdAt"`
}

// CreateDevicePartsReply only for api docs
type CreateDevicePartsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteDevicePartsByIDReply only for api docs
type DeleteDevicePartsByIDReply struct {
	Result
}

// UpdateDevicePartsByIDReply only for api docs
type UpdateDevicePartsByIDReply struct {
	Result
}

// GetDevicePartsByIDReply only for api docs
type GetDevicePartsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		DeviceParts DevicePartsObjDetail `json:"deviceParts"`
	} `json:"data"` // return data
}

// ListDevicePartssRequest request params
type ListDevicePartssRequest struct {
	query.Params
}

// ListDevicePartssReply only for api docs
type ListDevicePartssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		DevicePartss []DevicePartsObjDetail `json:"devicePartss"`
	} `json:"data"` // return data
}
