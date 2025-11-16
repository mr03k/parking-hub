package entity

import (
	"time"
)

type Contract struct {
	Base
	ContractNumber  string    `gorm:"type:varchar(50);unique;not null"`
	ContractDate    time.Time `gorm:"not null"`
	StartDate       time.Time `gorm:"not null"`
	EndDate         time.Time `gorm:"not null"`
	ContractAmount  int64     `gorm:"not null"`
	ContractType    string    `gorm:"type:varchar(50);not null"`
	ContractorID    string    `gorm:"not null"`
	OperationPeriod int       `gorm:"not null"`
	EquipmentPeriod int       `gorm:"not null"`
	Description     string    `gorm:"type:text"`
}
