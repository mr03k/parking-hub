package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateVehiclesRequest request params
type CreateVehiclesRequest struct {
	CodeVehicle               string    `json:"codeVehicle" binding:""`
	Vin                       string    `json:"vin" binding:""`
	PlateLicense              string    `json:"plateLicense" binding:""`
	TypeVehicle               string    `json:"typeVehicle" binding:""`
	Brand                     string    `json:"brand" binding:""`
	Model                     string    `json:"model" binding:""`
	Color                     string    `json:"color" binding:""`
	ManufactureOfYear         int       `json:"manufactureOfYear" binding:""`
	KilometersInitial         int64     `json:"kilometersInitial" binding:""`
	ExpiryInsurancePartyThird time.Time `json:"expiryInsurancePartyThird" binding:""`
	ExpiryInsuranceBody       time.Time `json:"expiryInsuranceBody" binding:""`
	ImageDocumentVehicle      string    `json:"imageDocumentVehicle" binding:""`
	ImageCardVehicle          string    `json:"imageCardVehicle" binding:""`
	ThirdPartyInsuranceImage  string    `json:"thirdPartyInsuranceImage" binding:""`
	BodyInsuranceImage        string    `json:"bodyInsuranceImage" binding:""`
	IDContractor              string    `json:"iDContractor" binding:""`
	Status                    string    `json:"status" binding:""`
	Description               string    `json:"description" binding:""`
}

// UpdateVehiclesByIDRequest request params
type UpdateVehiclesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	CodeVehicle               string    `json:"codeVehicle" binding:""`
	Vin                       string    `json:"vin" binding:""`
	PlateLicense              string    `json:"plateLicense" binding:""`
	TypeVehicle               string    `json:"typeVehicle" binding:""`
	Brand                     string    `json:"brand" binding:""`
	Model                     string    `json:"model" binding:""`
	Color                     string    `json:"color" binding:""`
	ManufactureOfYear         int       `json:"manufactureOfYear" binding:""`
	KilometersInitial         int64     `json:"kilometersInitial" binding:""`
	ExpiryInsurancePartyThird time.Time `json:"expiryInsurancePartyThird" binding:""`
	ExpiryInsuranceBody       time.Time `json:"expiryInsuranceBody" binding:""`
	ImageDocumentVehicle      string    `json:"imageDocumentVehicle" binding:""`
	ImageCardVehicle          string    `json:"imageCardVehicle" binding:""`
	ThirdPartyInsuranceImage  string    `json:"thirdPartyInsuranceImage" binding:""`
	BodyInsuranceImage        string    `json:"bodyInsuranceImage" binding:""`
	IDContractor              string    `json:"iDContractor" binding:""`
	Status                    string    `json:"status" binding:""`
	Description               string    `json:"description" binding:""`
}

// VehiclesObjDetail detail
type VehiclesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	CodeVehicle               string    `json:"codeVehicle"`
	Vin                       string    `json:"vin"`
	PlateLicense              string    `json:"plateLicense"`
	TypeVehicle               string    `json:"typeVehicle"`
	Brand                     string    `json:"brand"`
	Model                     string    `json:"model"`
	Color                     string    `json:"color"`
	ManufactureOfYear         int       `json:"manufactureOfYear"`
	KilometersInitial         int64     `json:"kilometersInitial"`
	ExpiryInsurancePartyThird time.Time `json:"expiryInsurancePartyThird"`
	ExpiryInsuranceBody       time.Time `json:"expiryInsuranceBody"`
	ImageDocumentVehicle      string    `json:"imageDocumentVehicle"`
	ImageCardVehicle          string    `json:"imageCardVehicle"`
	ThirdPartyInsuranceImage  string    `json:"thirdPartyInsuranceImage"`
	BodyInsuranceImage        string    `json:"bodyInsuranceImage"`
	IDContractor              string    `json:"iDContractor"`
	Status                    string    `json:"status"`
	Description               string    `json:"description"`
	CreatedAt                 int       `json:"createdAt"`
}

// CreateVehiclesReply only for api docs
type CreateVehiclesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteVehiclesByIDReply only for api docs
type DeleteVehiclesByIDReply struct {
	Result
}

// UpdateVehiclesByIDReply only for api docs
type UpdateVehiclesByIDReply struct {
	Result
}

// GetVehiclesByIDReply only for api docs
type GetVehiclesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Vehicles VehiclesObjDetail `json:"vehicles"`
	} `json:"data"` // return data
}

// ListVehiclessRequest request params
type ListVehiclessRequest struct {
	query.Params
}

// ListVehiclessReply only for api docs
type ListVehiclessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Vehicless []VehiclesObjDetail `json:"vehicless"`
	} `json:"data"` // return data
}
