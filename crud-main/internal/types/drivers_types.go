package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateDriversRequest request params
type CreateDriversRequest struct {
	FirstName                string    `json:"firstName" binding:""`
	NameLast                 string    `json:"nameLast" binding:""`
	Gender                   string    `json:"gender" binding:""`
	CodeDriver               string    `json:"codeDriver" binding:""`
	IDNational               string    `json:"iDNational" binding:""`
	CodePostal               string    `json:"codePostal" binding:""`
	NumberPhone              string    `json:"numberPhone" binding:""`
	NumberMobile             string    `json:"numberMobile" binding:""`
	Email                    string    `json:"email" binding:""`
	Address                  string    `json:"address" binding:""`
	IDContractor             string    `json:"iDContractor" binding:""`
	TypeDriver               string    `json:"typeDriver" binding:""`
	TypeShift                string    `json:"typeShift" binding:""`
	StatusEmployment         string    `json:"statusEmployment" binding:""`
	DateStartEmployment      time.Time `json:"dateStartEmployment" binding:""`
	DateEndEmployment        time.Time `json:"dateEndEmployment" binding:""`
	DriverPhoto              string    `json:"driverPhoto" binding:""`
	ImageCardID              string    `json:"imageCardID" binding:""`
	BirthCertificateImage    string    `json:"birthCertificateImage" binding:""`
	ImageCardServiceMilitary string    `json:"imageCardServiceMilitary" binding:""`
	ImageCertificateHealth   string    `json:"imageCertificateHealth" binding:""`
	ImageRecordCriminal      string    `json:"imageRecordCriminal" binding:""`
	Description              string    `json:"description" binding:""`
}

// UpdateDriversByIDRequest request params
type UpdateDriversByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	FirstName                string    `json:"firstName" binding:""`
	NameLast                 string    `json:"nameLast" binding:""`
	Gender                   string    `json:"gender" binding:""`
	CodeDriver               string    `json:"codeDriver" binding:""`
	IDNational               string    `json:"iDNational" binding:""`
	CodePostal               string    `json:"codePostal" binding:""`
	NumberPhone              string    `json:"numberPhone" binding:""`
	NumberMobile             string    `json:"numberMobile" binding:""`
	Email                    string    `json:"email" binding:""`
	Address                  string    `json:"address" binding:""`
	IDContractor             string    `json:"iDContractor" binding:""`
	TypeDriver               string    `json:"typeDriver" binding:""`
	TypeShift                string    `json:"typeShift" binding:""`
	StatusEmployment         string    `json:"statusEmployment" binding:""`
	DateStartEmployment      time.Time `json:"dateStartEmployment" binding:""`
	DateEndEmployment        time.Time `json:"dateEndEmployment" binding:""`
	DriverPhoto              string    `json:"driverPhoto" binding:""`
	ImageCardID              string    `json:"imageCardID" binding:""`
	BirthCertificateImage    string    `json:"birthCertificateImage" binding:""`
	ImageCardServiceMilitary string    `json:"imageCardServiceMilitary" binding:""`
	ImageCertificateHealth   string    `json:"imageCertificateHealth" binding:""`
	ImageRecordCriminal      string    `json:"imageRecordCriminal" binding:""`
	Description              string    `json:"description" binding:""`
}

// DriversObjDetail detail
type DriversObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	FirstName                string    `json:"firstName"`
	NameLast                 string    `json:"nameLast"`
	Gender                   string    `json:"gender"`
	CodeDriver               string    `json:"codeDriver"`
	IDNational               string    `json:"iDNational"`
	CodePostal               string    `json:"codePostal"`
	NumberPhone              string    `json:"numberPhone"`
	NumberMobile             string    `json:"numberMobile"`
	Email                    string    `json:"email"`
	Address                  string    `json:"address"`
	IDContractor             string    `json:"iDContractor"`
	TypeDriver               string    `json:"typeDriver"`
	TypeShift                string    `json:"typeShift"`
	StatusEmployment         string    `json:"statusEmployment"`
	DateStartEmployment      time.Time `json:"dateStartEmployment"`
	DateEndEmployment        time.Time `json:"dateEndEmployment"`
	DriverPhoto              string    `json:"driverPhoto"`
	ImageCardID              string    `json:"imageCardID"`
	BirthCertificateImage    string    `json:"birthCertificateImage"`
	ImageCardServiceMilitary string    `json:"imageCardServiceMilitary"`
	ImageCertificateHealth   string    `json:"imageCertificateHealth"`
	ImageRecordCriminal      string    `json:"imageRecordCriminal"`
	CreatedAt                int       `json:"createdAt"`
	Description              string    `json:"description"`
}

// CreateDriversReply only for api docs
type CreateDriversReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteDriversByIDReply only for api docs
type DeleteDriversByIDReply struct {
	Result
}

// UpdateDriversByIDReply only for api docs
type UpdateDriversByIDReply struct {
	Result
}

// GetDriversByIDReply only for api docs
type GetDriversByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Drivers DriversObjDetail `json:"drivers"`
	} `json:"data"` // return data
}

// ListDriverssRequest request params
type ListDriverssRequest struct {
	query.Params
}

// ListDriverssReply only for api docs
type ListDriverssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Driverss []DriversObjDetail `json:"driverss"`
	} `json:"data"` // return data
}
