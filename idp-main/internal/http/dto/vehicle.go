package dto

import (
	"application/internal/entity/device"
	"encoding/json"
	"net/http"
	"time"
)

// VehicleRequest represents the JSON structure for creating or updating a vehicle
type VehicleRequest struct {
	CodeVehicle               string `json:"code_vehicle"`
	VIN                       string `json:"vin"`
	PlateLicense              string `json:"plate_license"`
	TypeVehicle               string `json:"type_vehicle"`
	Brand                     string `json:"brand"`
	Model                     string `json:"model"`
	Color                     string `json:"color"`
	ManufactureOfYear         int    `json:"manufacture_of_year"`
	KilometersInitial         int64  `json:"kilometers_initial"`
	ExpiryInsurancePartyThird int64  `json:"expiry_insurance_party_third"`
	ExpiryInsuranceBody       int64  `json:"expiry_insurance_body"`
	ContractorID              int64  `json:"contractor_id"`
	Status                    string `json:"status"`
	Description               string `json:"description"`
}

// VehicleResponse represents the JSON structure for a vehicle in response
type VehicleResponse struct {
	ID                        int64  `json:"id"`
	CodeVehicle               string `json:"code_vehicle"`
	VIN                       string `json:"vin"`
	PlateLicense              string `json:"plate_license"`
	TypeVehicle               string `json:"type_vehicle"`
	Brand                     string `json:"brand"`
	Model                     string `json:"model"`
	Color                     string `json:"color"`
	ManufactureOfYear         int    `json:"manufacture_of_year"`
	KilometersInitial         int64  `json:"kilometers_initial"`
	ExpiryInsurancePartyThird int64  `json:"expiry_insurance_party_third"`
	ExpiryInsuranceBody       int64  `json:"expiry_insurance_body"`
	ContractorID              int64  `json:"contractor_id"`
	Status                    string `json:"status"`
	Description               string `json:"description"`
}

// VehicleListResponse represents the JSON structure for a list of vehicles
type VehicleListResponse struct {
	Count    int               `json:"count"`
	Vehicles []VehicleResponse `json:"vehicles"`
}

// ToEntity converts a VehicleRequest to a Vehicle entity
func (r *VehicleRequest) ToEntity() *device.Vehicle {
	return &device.Vehicle{
		CodeVehicle:               r.CodeVehicle,
		VIN:                       r.VIN,
		PlateLicense:              r.PlateLicense,
		TypeVehicle:               r.TypeVehicle,
		Brand:                     r.Brand,
		Model:                     r.Model,
		Color:                     r.Color,
		ManufactureOfYear:         r.ManufactureOfYear,
		KilometersInitial:         r.KilometersInitial,
		ExpiryInsurancePartyThird: time.Unix(r.ExpiryInsurancePartyThird, 0),
		ExpiryInsuranceBody:       time.Unix(r.ExpiryInsuranceBody, 0),
		ContractorID:              r.ContractorID,
		Status:                    r.Status,
		Description:               r.Description,
	}
}

// NewVehicleRequest creates a VehicleRequest from an HTTP request
func NewVehicleRequest(r *http.Request) (*VehicleRequest, error) {
	decoder := json.NewDecoder(r.Body)
	req := &VehicleRequest{}
	if err := decoder.Decode(req); err != nil {
		return nil, err
	}
	return req, nil
}

// NewVehicleResponse creates a VehicleResponse from a Vehicle entity
func NewVehicleResponse(v *device.Vehicle) *VehicleResponse {
	return &VehicleResponse{
		ID:                        v.ID,
		CodeVehicle:               v.CodeVehicle,
		VIN:                       v.VIN,
		PlateLicense:              v.PlateLicense,
		TypeVehicle:               v.TypeVehicle,
		Brand:                     v.Brand,
		Model:                     v.Model,
		Color:                     v.Color,
		ManufactureOfYear:         v.ManufactureOfYear,
		KilometersInitial:         v.KilometersInitial,
		ExpiryInsurancePartyThird: v.ExpiryInsurancePartyThird.Unix(),
		ExpiryInsuranceBody:       v.ExpiryInsuranceBody.Unix(),
		ContractorID:              v.ContractorID,
		Status:                    v.Status,
		Description:               v.Description,
	}
}

// NewVehicleListResponse creates a VehicleListResponse from a list of Vehicle entities
func NewVehicleListResponse(vehicles []device.Vehicle) *VehicleListResponse {
	vehicleList := make([]VehicleResponse, len(vehicles))
	for i, v := range vehicles {
		vehicleList[i] = *NewVehicleResponse(&v)
	}

	return &VehicleListResponse{
		Count:    len(vehicleList),
		Vehicles: vehicleList,
	}
}
