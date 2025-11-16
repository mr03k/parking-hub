package model

import (
	"time"
)

type PeakHourMultipliers struct {
	ID           uint64     `gorm:"column:id;type:uuid;primary_key" json:"id"`
	CodeTimePeak string     `gorm:"column:code_time_peak;type:varchar(20);NOT NULL" json:"codeTimePeak"`
	Description  string     `gorm:"column:description;type:varchar(100)" json:"description"`
	Multiplier   string     `gorm:"column:multiplier;type:numeric;NOT NULL" json:"multiplier"`
	Weekday      string     `gorm:"column:weekday;type:varchar(10)" json:"weekday"`
	TimeStart    string     `gorm:"column:time_start;type:varchar(20);NOT NULL" json:"timeStart"`
	TimeEnd      string     `gorm:"column:time_end;type:varchar(20);NOT NULL" json:"timeEnd"`
	FromValid    *time.Time `gorm:"column:from_valid;type:date;NOT NULL" json:"fromValid"`
	ToValid      *time.Time `gorm:"column:to_valid;type:date;NOT NULL" json:"toValid"`
	Flag         string     `gorm:"column:flag;type:varchar(20)" json:"flag"`
	CreatedAt    int        `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
}
