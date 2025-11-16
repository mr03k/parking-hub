package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateDriverLicensesRequest request params
type CreateDriverLicensesRequest struct {
	IDDriver      string    `json:"iDDriver" binding:""`
	LicenseNumber string    `json:"licenseNumber" binding:""`
	TypeLicense   string    `json:"typeLicense" binding:""`
	DateIssue     time.Time `json:"dateIssue" binding:""`
	DateExpiry    time.Time `json:"dateExpiry" binding:""`
	ImageLicense  string    `json:"imageLicense" binding:""`
	Description   string    `json:"description" binding:""`
}

// UpdateDriverLicensesByIDRequest request params
type UpdateDriverLicensesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	IDDriver      string    `json:"iDDriver" binding:""`
	LicenseNumber string    `json:"licenseNumber" binding:""`
	TypeLicense   string    `json:"typeLicense" binding:""`
	DateIssue     time.Time `json:"dateIssue" binding:""`
	DateExpiry    time.Time `json:"dateExpiry" binding:""`
	ImageLicense  string    `json:"imageLicense" binding:""`
	Description   string    `json:"description" binding:""`
}

// DriverLicensesObjDetail detail
type DriverLicensesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	IDDriver      string    `json:"iDDriver"`
	LicenseNumber string    `json:"licenseNumber"`
	TypeLicense   string    `json:"typeLicense"`
	DateIssue     time.Time `json:"dateIssue"`
	DateExpiry    time.Time `json:"dateExpiry"`
	ImageLicense  string    `json:"imageLicense"`
	Description   string    `json:"description"`
	CreatedAt     int       `json:"createdAt"`
}

// CreateDriverLicensesReply only for api docs
type CreateDriverLicensesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteDriverLicensesByIDReply only for api docs
type DeleteDriverLicensesByIDReply struct {
	Result
}

// UpdateDriverLicensesByIDReply only for api docs
type UpdateDriverLicensesByIDReply struct {
	Result
}

// GetDriverLicensesByIDReply only for api docs
type GetDriverLicensesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		DriverLicenses DriverLicensesObjDetail `json:"driverLicenses"`
	} `json:"data"` // return data
}

// ListDriverLicensessRequest request params
type ListDriverLicensessRequest struct {
	query.Params
}

// ListDriverLicensessReply only for api docs
type ListDriverLicensessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		DriverLicensess []DriverLicensesObjDetail `json:"driverLicensess"`
	} `json:"data"` // return data
}
