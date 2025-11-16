package entity

type Weekday string

const (
	Saturday  Weekday = "Saturday"
	Sunday    Weekday = "Sunday"
	Monday    Weekday = "Monday"
	Tuesday   Weekday = "Tuesday"
	Wednesday Weekday = "Wednesday"
	Thursday  Weekday = "Thursday"
	Friday    Weekday = "Friday"
)

type WorkShift string

const (
	Morning   WorkShift = "Morning"
	Afternoon WorkShift = "Afternoon"
	Both      WorkShift = "Both"
)

type Calender struct {
	Base
	ContractID     string    `gorm:"type:uuid;index"`
	ShamsiDate     string    `gorm:"type:varchar(10);not null"`
	WorkDate       int64     `gorm:"type:date;not null"`
	Weekday        Weekday   `gorm:"type:enum('Saturday','Sunday','Monday','Tuesday','Wednesday','Thursday','Friday')"`
	Year           int       `gorm:"type:int;not null"`
	IsHoliday      bool      `gorm:"not null"`
	WorkShift      WorkShift `gorm:"type:enum('Morning','Afternoon','Both')"`
	Description    string    `gorm:"type:varchar(255)"`
	WorkShiftStart int64     `gorm:"type:time"`
	WorkShiftEnd   int64     `gorm:"type:time"`
}
