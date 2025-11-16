package dto

import (
	"farin/domain/entity"
)

type VehicleRequest struct {
	VehicleID                 int64  `json:"vehicleID" binding:"required,min=1,max=120"`
	CodeVehicle               string `json:"codeVehicle" binding:"required"`
	VIN                       string `json:"vin" binding:"required,min=1,max=120"`
	PlateLicense              string `json:"plateLicense" binding:"required,min=1,max=115"`
	TypeVehicle               string `json:"typeVehicle,omitempty"`
	Brand                     string `json:"brand,omitempty"`
	Model                     string `json:"model,omitempty"`
	Color                     string `json:"color,omitempty"`
	ManufactureOfYear         int    `json:"manufactureOfYear,omitempty"`
	KilometersInitial         int64  `json:"kilometersInitial,omitempty"`
	ExpiryInsurancePartyThird int64  `json:"expiryInsurancePartyThird,omitempty"`
	ExpiryInsuranceBody       int64  `json:"expiryInsuranceBody,omitempty"`
	Status                    string `json:"status,omitempty"`
	ContractorID              string `json:"contractorId" binding:"required,fkGorm=contractors"`
	Description               string `json:"description,omitempty"`
	ImageDocumentVehicle      []byte `json:"imageDocumentVehicle" binding:"fileData=image/jpeg&image/png;4096000"`
	ImageCardVehicle          []byte `json:"imageCardVehicle" binding:"fileData=image/jpeg&image/png;4096000"`
	ThirdPartyInsuranceImage  []byte `json:"thirdPartyInsuranceImage" binding:"fileData=image/jpeg&image/png;4096000"`
	BodyInsuranceImage        []byte `json:"bodyInsuranceImage" binding:"fileData=image/jpeg&image/png;4096000"`
}

func (req *VehicleRequest) ToEntity() *entity.Vehicle {
	return &entity.Vehicle{
		VehicleID:                 req.VehicleID,
		CodeVehicle:               req.CodeVehicle,
		VIN:                       req.VIN,
		PlateLicense:              req.PlateLicense,
		TypeVehicle:               req.TypeVehicle,
		Brand:                     req.Brand,
		Model:                     req.Model,
		Color:                     req.Color,
		ManufactureOfYear:         req.ManufactureOfYear,
		KilometersInitial:         req.KilometersInitial,
		ExpiryInsurancePartyThird: req.ExpiryInsurancePartyThird,
		ExpiryInsuranceBody:       req.ExpiryInsuranceBody,
		Status:                    req.Status,
		Description:               req.Description,
		ContractorID:              req.ContractorID,
	}
}

type VehicleResponse struct {
	ID                        string `json:"id"`
	VehicleID                 int64  `json:"vehicleID"`
	CodeVehicle               string `json:"codeVehicle"`
	VIN                       string `json:"vin"`
	PlateLicense              string `json:"plateLicense"`
	TypeVehicle               string `json:"typeVehicle"`
	Brand                     string `json:"brand"`
	Model                     string `json:"model"`
	Color                     string `json:"color"`
	ManufactureOfYear         int    `json:"manufactureOfYear"`
	KilometersInitial         int64  `json:"kilometersInitial"`
	ExpiryInsurancePartyThird int64  `json:"expiryInsurancePartyThird"`
	ExpiryInsuranceBody       int64  `json:"expiryInsuranceBody"`
	Status                    string `json:"status"`
	Description               string `json:"description,omitempty"`
	ContractorID              string `json:"contractorId"`
	CreatedAt                 int64  `json:"createdAt"`
	UpdatedAt                 int64  `json:"updatedAt"`
}

func (resp *VehicleResponse) FromEntity(vehicle *entity.Vehicle) {
	resp.ID = vehicle.ID
	resp.VehicleID = vehicle.VehicleID
	resp.CodeVehicle = vehicle.CodeVehicle
	resp.VIN = vehicle.VIN
	resp.PlateLicense = vehicle.PlateLicense
	resp.TypeVehicle = vehicle.TypeVehicle
	resp.Brand = vehicle.Brand
	resp.Model = vehicle.Model
	resp.Color = vehicle.Color
	resp.ManufactureOfYear = vehicle.ManufactureOfYear
	resp.KilometersInitial = vehicle.KilometersInitial
	resp.ExpiryInsurancePartyThird = vehicle.ExpiryInsurancePartyThird
	resp.ExpiryInsuranceBody = vehicle.ExpiryInsuranceBody
	resp.Status = vehicle.Status
	resp.Description = vehicle.Description
	resp.ContractorID = vehicle.ContractorID
	resp.CreatedAt = vehicle.CreatedAt
	resp.UpdatedAt = vehicle.UpdatedAt
}

type VehicleListResponse struct {
	Vehicles []VehicleResponse `json:"vehicles"`
	Total    int64             `json:"total"`
}
