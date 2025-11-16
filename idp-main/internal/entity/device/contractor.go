package device

import "github.com/google/uuid"

type Contractor struct {
	ID                    uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ContractorName        string    `gorm:"size:100;not null" json:"contractor_name"`
	CodeContractor        string    `gorm:"size:10;unique;not null" json:"code_contractor"`
	NumberRegistration    string    `gorm:"size:50" json:"number_registration"`
	PersonContact         string    `gorm:"size:100" json:"person_contact"`
	CEOName               string    `gorm:"size:100" json:"ceo_name"`
	SignatoriesAuthorized string    `json:"signatories_authorized"`
	PhoneNumber           string    `gorm:"size:15" json:"phone_number"`
	Email                 string    `gorm:"size:100" json:"email"`
	Address               string    `gorm:"size:255" json:"address"`
	TypeContract          string    `gorm:"size:50" json:"type_contract"`
	NumberAccountBank     string    `gorm:"size:30" json:"number_account_bank"`
	Description           string    `json:"description"`
}
