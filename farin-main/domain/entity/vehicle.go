package entity

type Vehicle struct {
	Base

	VehicleID                 int64  `gorm:"type:bigint;unique;not null"`
	CodeVehicle               string `gorm:"type:varchar(120);unique;not null"`
	VIN                       string `gorm:"type:varchar(120);unique;not null"`
	PlateLicense              string `gorm:"type:varchar(115);unique;not null"`
	TypeVehicle               string `gorm:"type:varchar(50)"`
	Brand                     string `gorm:"type:varchar(50)"`
	Model                     string `gorm:"type:varchar(50)"`
	Color                     string `gorm:"type:varchar(30)"`
	ManufactureOfYear         int    `gorm:"type:int"`
	KilometersInitial         int64  `gorm:"type:bigint"`
	ExpiryInsurancePartyThird int64  `gorm:"type:bigint"`
	ExpiryInsuranceBody       int64  `gorm:"type:bigint"`
	ImageDocumentVehicle      string `gorm:"type:varchar(256)"`
	ImageCardVehicle          string `gorm:"type:varchar(256)"`
	ThirdPartyInsuranceImage  string `gorm:"type:varchar(256)"`
	BodyInsuranceImage        string `gorm:"type:varchar(256)"`
	ContractorID              string `gorm:"type:uuid;not null;index"`
	Status                    string `gorm:"type:varchar(20)"`
	Description               string `gorm:"type:text"`
}
