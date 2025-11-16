package model

import (
	"time"
)

type Exceptions struct {
	ID                      uint64     `gorm:"column:id;type:uuid;primary_key" json:"id"`
	CarLicensePlates        string     `gorm:"column:car_license_plates;type:_text" json:"carLicensePlates"`
	MotorcycleLicensePlates string     `gorm:"column:motorcycle_license_plates;type:_text" json:"motorcycleLicensePlates"`
	ExceptionMultiplier     string     `gorm:"column:exception_multiplier;type:numeric;NOT NULL" json:"exceptionMultiplier"`
	StartDate               *time.Time `gorm:"column:start_date;type:date;NOT NULL" json:"startDate"`
	EndDate                 *time.Time `gorm:"column:end_date;type:date" json:"endDate"`
	Description             string     `gorm:"column:description;type:text" json:"description"`
	NotificationNumber      string     `gorm:"column:notification_number;type:varchar(20)" json:"notificationNumber"`
	NotificationDate        *time.Time `gorm:"column:notification_date;type:date" json:"notificationDate"`
	DocumentImage           string     `gorm:"column:document_image;type:varchar(200)" json:"documentImage"`
	UserID                  string     `gorm:"column:user_id;type:uuid;NOT NULL" json:"userID"`
	VehicleType             string     `gorm:"column:vehicle_type;type:varchar(20)" json:"vehicleType"`
	CreatedAt               int        `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
}
