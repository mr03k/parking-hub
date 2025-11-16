package model

import (
	"time"
)

type DeviceParts struct {
	ID                 uint64     `gorm:"column:id;type:uuid;primary_key" json:"id"`
	CodePart           string     `gorm:"column:code_part;type:varchar(20);NOT NULL" json:"codePart"`
	PartName           string     `gorm:"column:part_name;type:varchar(100)" json:"partName"`
	TypePart           string     `gorm:"column:type_part;type:varchar(50)" json:"typePart"`
	Brand              string     `gorm:"column:brand;type:varchar(50)" json:"brand"`
	Model              string     `gorm:"column:model;type:varchar(50)" json:"model"`
	NumberSerial       string     `gorm:"column:number_serial;type:varchar(50)" json:"numberSerial"`
	DateInstallation   *time.Time `gorm:"column:date_installation;type:date" json:"dateInstallation"`
	DateExpiryWarranty *time.Time `gorm:"column:date_expiry_warranty;type:date" json:"dateExpiryWarranty"`
	PeriodMaintenance  string     `gorm:"column:period_maintenance;type:varchar(50)" json:"periodMaintenance"`
	IDDevice           string     `gorm:"column:id_device;type:uuid" json:"iDDevice"`
	Description        string     `gorm:"column:description;type:text" json:"description"`
	CreatedAt          int        `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
}
