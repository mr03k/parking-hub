package entity

type DriverAssignment struct {
	Base
	DriverID    string `gorm:"type:uuid;not null;index" json:"driverId"`
	CodeVehicle string `gorm:"type:varchar(120);not null;index" json:"codeVehicle"`
	RingID      int64  `gorm:"type:uuid;not null;index" json:"ringId"`
	CalenderID  string `gorm:"type:uuid;not null;index" json:"calenderId"` // Changed json tag to match struct

	Driver   Driver   `gorm:"foreignKey:DriverID;references:id;constraint:OnDelete:CASCADE" json:"driver"`
	Vehicle  Vehicle  `gorm:"foreignKey:CodeVehicle;references:code_vehicle;constraint:OnDelete:CASCADE" json:"vehicle"`
	Ring     Ring     `gorm:"foreignKey:RingID;references:id;constraint:OnDelete:CASCADE" json:"ring"`
	Calender Calender `gorm:"foreignKey:CalenderID;references:id;constraint:OnDelete:CASCADE" json:"calender"` // Changed foreignKey and json tag
}
