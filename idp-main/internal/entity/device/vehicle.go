package device

import (
	"time"
)

type Vehicle struct {
	ID                        int64
	CodeVehicle               string    `gorm:"size:20;unique;not null" json:"code_vehicle"`
	VIN                       string    `gorm:"size:20;unique;not null" json:"vin"`
	PlateLicense              string    `gorm:"size:15;unique;not null" json:"plate_license"`
	TypeVehicle               string    `gorm:"size:50" json:"type_vehicle"`
	Brand                     string    `gorm:"size:50" json:"brand"`
	Model                     string    `gorm:"size:50" json:"model"`
	Color                     string    `gorm:"size:30" json:"color"`
	ManufactureOfYear         int       `json:"manufacture_of_year"`
	KilometersInitial         int64     `json:"kilometers_initial"`
	ExpiryInsurancePartyThird time.Time `json:"expiry_insurance_party_third"`
	ExpiryInsuranceBody       time.Time `json:"expiry_insurance_body"`
	ImageDocumentVehicle      []byte    `json:"image_document_vehicle"`
	ImageCardVehicle          []byte    `json:"image_card_vehicle"`
	ThirdPartyInsuranceImage  []byte    `json:"third_party_insurance_image"`
	BodyInsuranceImage        []byte    `json:"body_insurance_image"`
	ContractorID              int64     `json:"id_contractor"` // Foreign key
	Status                    string    `gorm:"size:20" json:"status"`
	Description               string    `json:"description"`
}
