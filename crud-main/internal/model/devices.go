package model

import (
	"time"
)

type Devices struct {
	ID                  uint64     `gorm:"column:id;type:uuid;primary_key" json:"id"`
	CodeDevice          string     `gorm:"column:code_device;type:varchar(20);NOT NULL" json:"codeDevice"`
	NumberSerial        string     `gorm:"column:number_serial;type:varchar(50);NOT NULL" json:"numberSerial"`
	Model               string     `gorm:"column:model;type:varchar(50)" json:"model"`
	DateInstallation    *time.Time `gorm:"column:date_installation;type:date" json:"dateInstallation"`
	DateExpiryWarranty  *time.Time `gorm:"column:date_expiry_warranty;type:date" json:"dateExpiryWarranty"`
	DateExpiryInsurance *time.Time `gorm:"column:date_expiry_insurance;type:date" json:"dateExpiryInsurance"`
	ClassDevice         string     `gorm:"column:class_device;type:varchar(50)" json:"classDevice"`
	ImageContract       string     `gorm:"column:image_contract;type:varchar(200)" json:"imageContract"`
	ImageInsurance      string     `gorm:"column:image_insurance;type:varchar(200)" json:"imageInsurance"`
	IDContractor        string     `gorm:"column:id_contractor;type:uuid" json:"iDContractor"`
	Description         string     `gorm:"column:description;type:text" json:"description"`
	CreatedAt           int        `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
}
