package idpentities

import (
	"regexp"

	validator "github.com/go-playground/validator/v10"
)

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleDriver      = "driver"
)

type User struct {
	UserID                string  `json:"userId"`
	Msisdn                *string `json:"msisdn,omitempty" validate:"required_without=Email,omitempty,msisdn"`
	MsisdnVerified        bool    `json:"msisdnVerified"`
	ClearPassword         []byte  `json:"-" gorm:"-" validate:"omitempty,min=6"`
	EncryptedPassword     []byte  `json:"password"`
	ContractorName        string  `json:"contractorName"`
	ContractorCode        string  `json:"contractorCode"`
	RegistrationNumber    string  `json:"registrationNumber"`
	ContactPerson         string  `json:"contactPerson"`
	CeoName               string  `json:"ceoName"`
	AuthorizedSignatories string  `json:"authorizedSignatories"`
	PhoneNumber           string  `json:"phoneNumber"`
	Email                 string  `json:"email"`
	Address               string  `json:"address"`
	ContractType          string  `json:"contractType"`
	BankAccountNumber     string  `json:"bankAccountNumber"`
	Description           string  `json:"description"`
	Role                  Role    `json:"role"`
	CreatedAt             int64   `json:"createdAt"`
	UpdatedAt             int64   `json:"updatedAt"`
}

// validate
func (u *User) Validate(v *validator.Validate) error {
	err := v.RegisterValidation("msisdn", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^98\d{10}$`).MatchString(fl.Field().String())
	})
	if err != nil {
		return err
	}

	err = v.Struct(u)
	if err != nil {
		return err
	}
	return nil
}
