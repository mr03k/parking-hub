package model

import (
	"time"
)

type Drivers struct {
	ID                       uint64     `gorm:"column:id;type:uuid;primary_key" json:"id"`
	FirstName                string     `gorm:"column:first_name;type:varchar(50);NOT NULL" json:"firstName"`
	NameLast                 string     `gorm:"column:name_last;type:varchar(50);NOT NULL" json:"nameLast"`
	Gender                   string     `gorm:"column:gender;type:varchar(10)" json:"gender"`
	CodeDriver               string     `gorm:"column:code_driver;type:varchar(10);NOT NULL" json:"codeDriver"`
	IDNational               string     `gorm:"column:id_national;type:varchar(20)" json:"iDNational"`
	CodePostal               string     `gorm:"column:code_postal;type:varchar(10)" json:"codePostal"`
	NumberPhone              string     `gorm:"column:number_phone;type:varchar(15)" json:"numberPhone"`
	NumberMobile             string     `gorm:"column:number_mobile;type:varchar(15)" json:"numberMobile"`
	Email                    string     `gorm:"column:email;type:varchar(100)" json:"email"`
	Address                  string     `gorm:"column:address;type:text" json:"address"`
	IDContractor             string     `gorm:"column:id_contractor;type:uuid" json:"iDContractor"`
	TypeDriver               string     `gorm:"column:type_driver;type:varchar(10)" json:"typeDriver"`
	TypeShift                string     `gorm:"column:type_shift;type:varchar(10)" json:"typeShift"`
	StatusEmployment         string     `gorm:"column:status_employment;type:varchar(20)" json:"statusEmployment"`
	DateStartEmployment      *time.Time `gorm:"column:date_start_employment;type:date" json:"dateStartEmployment"`
	DateEndEmployment        *time.Time `gorm:"column:date_end_employment;type:date" json:"dateEndEmployment"`
	DriverPhoto              string     `gorm:"column:driver_photo;type:varchar(200)" json:"driverPhoto"`
	ImageCardID              string     `gorm:"column:image_card_id;type:varchar(200)" json:"imageCardID"`
	BirthCertificateImage    string     `gorm:"column:birth_certificate_image;type:varchar(200)" json:"birthCertificateImage"`
	ImageCardServiceMilitary string     `gorm:"column:image_card_service_military;type:varchar(200)" json:"imageCardServiceMilitary"`
	ImageCertificateHealth   string     `gorm:"column:image_certificate_health;type:varchar(200)" json:"imageCertificateHealth"`
	ImageRecordCriminal      string     `gorm:"column:image_record_criminal;type:varchar(200)" json:"imageRecordCriminal"`
	CreatedAt                int        `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
	Description              string     `gorm:"column:description;type:text" json:"description"`
}
