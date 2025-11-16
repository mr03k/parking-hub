package model

import (
	"time"
)

type Vehicles struct {
	ID                        uint64     `gorm:"column:id;type:uuid;primary_key" json:"id"`
	CodeVehicle               string     `gorm:"column:code_vehicle;type:varchar(20);NOT NULL" json:"codeVehicle"`
	Vin                       string     `gorm:"column:vin;type:varchar(20);NOT NULL" json:"vin"`
	PlateLicense              string     `gorm:"column:plate_license;type:varchar(15);NOT NULL" json:"plateLicense"`
	TypeVehicle               string     `gorm:"column:type_vehicle;type:varchar(50)" json:"typeVehicle"`
	Brand                     string     `gorm:"column:brand;type:varchar(50)" json:"brand"`
	Model                     string     `gorm:"column:model;type:varchar(50)" json:"model"`
	Color                     string     `gorm:"column:color;type:varchar(30)" json:"color"`
	ManufactureOfYear         int        `gorm:"column:manufacture_of_year;type:int4" json:"manufactureOfYear"`
	KilometersInitial         int64      `gorm:"column:kilometers_initial;type:int8" json:"kilometersInitial"`
	ExpiryInsurancePartyThird *time.Time `gorm:"column:expiry_insurance_party_third;type:date" json:"expiryInsurancePartyThird"`
	ExpiryInsuranceBody       *time.Time `gorm:"column:expiry_insurance_body;type:date" json:"expiryInsuranceBody"`
	ImageDocumentVehicle      string     `gorm:"column:image_document_vehicle;type:varchar(200)" json:"imageDocumentVehicle"`
	ImageCardVehicle          string     `gorm:"column:image_card_vehicle;type:varchar(200)" json:"imageCardVehicle"`
	ThirdPartyInsuranceImage  string     `gorm:"column:third_party_insurance_image;type:varchar(200)" json:"thirdPartyInsuranceImage"`
	BodyInsuranceImage        string     `gorm:"column:body_insurance_image;type:varchar(200)" json:"bodyInsuranceImage"`
	IDContractor              string     `gorm:"column:id_contractor;type:uuid" json:"iDContractor"`
	Status                    string     `gorm:"column:status;type:varchar(20)" json:"status"`
	Description               string     `gorm:"column:description;type:text" json:"description"`
	CreatedAt                 int        `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
}
