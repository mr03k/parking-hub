package device

import (
	"github.com/google/uuid"
	"time"
)

type LicensePlateReaderDevice struct {
	ID                  uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	CodeDevice          string    `gorm:"size:20;not null" json:"codeDevice"`   // Unique code for the device
	NumberSerial        string    `gorm:"size:50;not null" json:"numberSerial"` // Serial number of the device
	Model               string    `gorm:"size:50" json:"model"`                 // Model of the device
	DateInstallation    time.Time `json:"dateInstallation"`                     // Installation date of the device
	DateExpiryWarranty  time.Time `json:"dateExpiryWarranty"`                   // Warranty expiry date
	DateExpiryInsurance time.Time `json:"dateExpiryInsurance"`                  // Insurance expiry date
	ClassDevice         string    `gorm:"size:50" json:"classDevice"`           // Class of the device
	ImageContractURL    string    `json:"imageContractUrl,omitempty"`           // URL for the contract image
	ImageInsuranceURL   string    `json:"imageInsuranceUrl,omitempty"`          // URL for the insurance image
	ContractorID        uuid.UUID `json:"contractorId"`                         // Foreign key to the contractor
	Description         string    `json:"description,omitempty"`                // Additional description
	CreatedAt           time.Time `json:"createdAt"`                            // Timestamp of record creation
}
