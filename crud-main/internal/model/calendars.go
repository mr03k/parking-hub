package model

import (
	"time"
)

type Calendars struct {
	ID          uint64     `gorm:"column:id;type:uuid;primary_key" json:"id"`
	IDContract  string     `gorm:"column:id_contract;type:uuid;NOT NULL" json:"iDContract"`
	ShamsiDate  string     `gorm:"column:shamsi_date;type:varchar(10);NOT NULL" json:"shamsiDate"`
	WorkDate    *time.Time `gorm:"column:work_date;type:date;NOT NULL" json:"workDate"`
	Weekday     string     `gorm:"column:weekday;type:weekday;NOT NULL" json:"weekday"`
	Year        int        `gorm:"column:year;type:int4" json:"year"`
	HolidayIs   bool       `gorm:"column:holiday_is;type:bool;NOT NULL" json:"holidayIs"`
	ShiftWork   string     `gorm:"column:shift_work;type:shift_work;NOT NULL" json:"shiftWork"`
	Description string     `gorm:"column:description;type:varchar(255)" json:"description"`
	CreatedAt   int        `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
}
