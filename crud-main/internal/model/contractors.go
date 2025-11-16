package model

type Contractors struct {
	ID                    uint64 `gorm:"column:id;type:uuid;primary_key" json:"id"`
	ContractorName        string `gorm:"column:contractor_name;type:varchar(100);NOT NULL" json:"contractorName"`
	CodeContractor        string `gorm:"column:code_contractor;type:varchar(10);NOT NULL" json:"codeContractor"`
	NumberRegistration    string `gorm:"column:number_registration;type:varchar(50)" json:"numberRegistration"`
	PersonContact         string `gorm:"column:person_contact;type:varchar(100)" json:"personContact"`
	CeoName               string `gorm:"column:ceo_name;type:varchar(100)" json:"ceoName"`
	SignatoriesAuthorized string `gorm:"column:signatories_authorized;type:text" json:"signatoriesAuthorized"`
	PhoneNumber           string `gorm:"column:phone_number;type:varchar(15)" json:"phoneNumber"`
	Email                 string `gorm:"column:email;type:varchar(100)" json:"email"`
	Address               string `gorm:"column:address;type:varchar(255)" json:"address"`
	TypeContract          string `gorm:"column:type_contract;type:varchar(50)" json:"typeContract"`
	NumberAccountBank     string `gorm:"column:number_account_bank;type:varchar(30)" json:"numberAccountBank"`
	CreatedAt             int    `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
	Description           string `gorm:"column:description;type:text" json:"description"`
}
