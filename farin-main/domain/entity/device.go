package entity

type Device struct {
	Base
	DeviceID            int64  `gorm:"type:varchar(120);not null"`
	CodeDevice          string `gorm:"type:varchar(120);not null"`
	NumberSerial        string `gorm:"type:varchar(50);not null"`
	Model               string `gorm:"type:varchar(50)"`
	DateInstallation    int64  `gorm:"type:bigint"`
	DateExpiryWarranty  int64  `gorm:"type:bigint"`
	DateExpiryInsurance int64  `gorm:"type:bigint"`
	ClassDevice         string `gorm:"type:varchar(50)"`
	ImageContract       string `gorm:"type:varchar(256)"`
	ImageInsurance      string `gorm:"type:varchar(256)"`
	ContractorID        string `gorm:"type:uuid;not null;index"`
	VehicleID           string `gorm:"type:uuid;not null;index"`
	Description         string `gorm:"type:text"`
}
