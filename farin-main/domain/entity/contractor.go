package entity

type Contractor struct {
	Base
	ContractorName        string `gorm:"type:varchar(100);not null" json:"contractorName"`
	CodeContractor        string `gorm:"type:varchar(10);unique;not null" json:"codeContractor"`
	NumberRegistration    string `gorm:"type:varchar(50)" json:"numberRegistration"`
	PersonContact         string `gorm:"type:varchar(100)" json:"personContact"`
	CEOName               string `gorm:"type:varchar(100)" json:"ceoName"`
	SignatoriesAuthorized string `gorm:"type:text" json:"signatoriesAuthorized"`
	PhoneNumber           string `gorm:"type:varchar(15)" json:"phoneNumber"`
	Email                 string `gorm:"type:varchar(100)" json:"email"`
	Address               string `gorm:"type:varchar(255)" json:"address"`
	TypeContract          string `gorm:"type:varchar(50)" json:"typeContract"`
	NumberAccountBank     string `gorm:"type:varchar(30)" json:"numberAccountBank"`
	Description           string `gorm:"type:text" json:"description"`
}
