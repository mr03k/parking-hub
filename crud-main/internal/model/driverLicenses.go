package model

import (
	"time"
)

type DriverLicenses struct {
	ID            uint64     `gorm:"column:id;type:uuid;primary_key" json:"id"`
	IDDriver      string     `gorm:"column:id_driver;type:uuid;NOT NULL" json:"iDDriver"`
	LicenseNumber string     `gorm:"column:license_number;type:varchar(20);NOT NULL" json:"licenseNumber"`
	TypeLicense   string     `gorm:"column:type_license;type:varchar(50)" json:"typeLicense"`
	DateIssue     *time.Time `gorm:"column:date_issue;type:date" json:"dateIssue"`
	DateExpiry    *time.Time `gorm:"column:date_expiry;type:date" json:"dateExpiry"`
	ImageLicense  string     `gorm:"column:image_license;type:varchar(200)" json:"imageLicense"`
	Description   string     `gorm:"column:description;type:text" json:"description"`
	CreatedAt     int        `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
}
