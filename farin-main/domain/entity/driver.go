package entity

type DriverType string

const (
	DriverTypeOriginal DriverType = "original"
	DriverTypeReserved DriverType = "reserve"
)

type ShiftType string

const (
	ShiftTypeMorning   ShiftType = "morning"
	ShiftTypeAfterNoon ShiftType = "afternoon"
	ShiftTypeBoth      ShiftType = "both"
)

type Driver struct {
	Base
	ContractorID             *string    `gorm:"type:uuid;index"`
	UserID                   string     `gorm:"type:uuid;not null;index"`
	DriverType               DriverType `gorm:"type:varchar(10)"`
	ShiftType                string     `gorm:"type:varchar(10)"`
	EmploymentStatus         string     `gorm:"type:varchar(20)"`
	EmploymentStartDate      int64      `gorm:"not null"`
	EmploymentEndDate        *int64     `gorm:""`
	DriverPhoto              string     `gorm:"type:varchar(256)"`
	IDCardImage              string     `gorm:"type:varchar(256)"`
	BirthCertificateImage    string     `gorm:"type:varchar(256)"`
	MilitaryServiceCardImage string     `gorm:"type:varchar(256)"`
	HealthCertificateImage   string     `gorm:"type:varchar(256)"`
	CriminalRecordImage      string     `gorm:"type:varchar(256)"`
	Description              string     `gorm:"type:varchar(256)"`
}
