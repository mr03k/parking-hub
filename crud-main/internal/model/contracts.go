package model

import (
	"time"
)

type Contracts struct {
	ID              uint64     `gorm:"column:id;type:uuid;primary_key" json:"id"`
	NumberContract  string     `gorm:"column:number_contract;type:varchar(50);NOT NULL" json:"numberContract"`
	DateContract    *time.Time `gorm:"column:date_contract;type:date;NOT NULL" json:"dateContract"`
	DateStart       *time.Time `gorm:"column:date_start;type:date;NOT NULL" json:"dateStart"`
	DateEnd         *time.Time `gorm:"column:date_end;type:date;NOT NULL" json:"dateEnd"`
	AmountContract  int64      `gorm:"column:amount_contract;type:int8;NOT NULL" json:"amountContract"`
	TypeContract    string     `gorm:"column:type_contract;type:varchar(50);NOT NULL" json:"typeContract"`
	IDContractor    string     `gorm:"column:id_contractor;type:uuid;NOT NULL" json:"iDContractor"`
	PeriodOperation int        `gorm:"column:period_operation;type:int4" json:"periodOperation"`
	PeriodEquipment int        `gorm:"column:period_equipment;type:int4" json:"periodEquipment"`
	Description     string     `gorm:"column:description;type:text" json:"description"`
	CreatedAt       int        `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
}
