package entity

type Gender string

const (
	GenderMale   Gender = "Male"
	GenderFemale Gender = "Female"
)

type Status string

const (
	StatusActive   Status = "Active"
	StatusInactive Status = "Inactive"
)

type User struct {
	Base
	Username     string `gorm:"type:varchar(50);not null" json:"username"`
	Password     string `gorm:"type:varchar(255);not null" json:"-"`
	FirstName    string `gorm:"type:varchar(100);not null" json:"firstName"`
	LastName     string `gorm:"type:varchar(100);not null" json:"lastName"`
	Email        string `gorm:"type:varchar(100);not null;unique" json:"email"`
	PhoneNumber  string `gorm:"type:varchar(15)" json:"phoneNumber,omitempty"`
	NationalID   string `gorm:"type:char(10)" json:"nationalId,omitempty"`
	PostalCode   string `gorm:"type:char(10)" json:"postalCode,omitempty"`
	CompanyName  string `gorm:"type:varchar(100)" json:"companyName,omitempty"`
	ProfileImage string `gorm:"type:varchar(255)" json:"profileImage,omitempty"`
	Gender       Gender `gorm:"type:varchar(10)" json:"gender,omitempty"`
	Address      string `gorm:"type:varchar(255)" json:"address,omitempty"`
	Status       Status `gorm:"type:varchar(10)" json:"status,omitempty"`
	RoleID       string `gorm:"column:role_id;type:uuid" json:"roleId,omitempty"`
	Role         Role   `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}
