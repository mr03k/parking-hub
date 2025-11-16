package entity

import (
	"time"
)

// Driver represents the driver entity with UUIDs and camelCase JSON tags.
type Driver struct {
	ID                          string     `json:"id"`                                    // Unique ID for the driver
	Address                     string     `json:"address"`                               // Driver's residential address
	ContractorID                string     `json:"contractorId"`                          // Foreign key to the contractor table
	DriverType                  string     `json:"driverType"`                            // Type of driver (e.g., primary, reserve)
	ShiftType                   string     `json:"shiftType"`                             // Shift type (e.g., morning, evening, both)
	EmploymentStatus            string     `json:"employmentStatus"`                      // Employment status (e.g., active, inactive)
	EmploymentStartDate         *time.Time `json:"employmentStartDate"`                   // Employment start date
	EmploymentEndDate           *time.Time `json:"employmentEndDate,omitempty"`           // Employment end date (optional)
	DriverPhotoURL              string     `json:"driverPhotoUrl,omitempty"`              // URL of the driver's photo
	IDCardImageURL              string     `json:"idCardImageUrl,omitempty"`              // URL of the driver's national ID card image
	BirthCertificateImageURL    string     `json:"birthCertificateImageUrl,omitempty"`    // URL of the driver's birth certificate image
	MilitaryServiceCardImageURL string     `json:"militaryServiceCardImageUrl,omitempty"` // URL of the driver's military service card image
	HealthCertificateImageURL   string     `json:"healthCertificateImageUrl,omitempty"`   // URL of the driver's health certificate image
	CriminalRecordImageURL      string     `json:"criminalRecordImageUrl,omitempty"`      // URL of the driver's criminal record image
	Description                 string     `json:"description,omitempty"`                 // Additional description
	CreatedAt                   time.Time  `json:"createdAt"`                             // Timestamp of record creation
	UpdatedAt                   time.Time  `json:"updatedAt"`                             // Timestamp of last record update
}
