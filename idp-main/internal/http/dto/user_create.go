package dto

import (
	"encoding/json"
	"net/http"
	"regexp"

	idpentities "application/internal/entity/idp"

	validator "github.com/go-playground/validator/v10"
)

const (
	msisdnRegex = `^98\d.*$`
)

type User struct {
	Msisdn                string `json:"msisdn,omitempty" validate:"required_without=email,omitempty,msisdn"`
	MsisdnVerified        bool   `json:"msisdnVerified" default:"true"`
	ContractorName        string `json:"contractorName,omitempty"`
	ContractorCode        string `json:"contractorCode,omitempty"`
	RegistrationNumber    string `json:"registrationNumber,omitempty"`
	ContactPerson         string `json:"contactPerson,omitempty"`
	CeoName               string `json:"ceoName,omitempty"`
	AuthorizedSignatories string `json:"authorizedSignatories,omitempty"`
	PhoneNumber           string `json:"phoneNumber,omitempty"`
	Email                 string `json:"email,omitempty"`
	Address               string `json:"address,omitempty"`
	ContractType          string `json:"contractType,omitempty"`
	BankAccountNumber     string `json:"bankAccountNumber,omitempty"`
	Description           string `json:"description,omitempty"`
	Role                  string `json:"role,omitempty"`
}

type UserRequest struct {
	Msisdn                string `json:"msisdn,omitempty" validate:"required_without=email,omitempty,msisdn"`
	MsisdnVerified        bool   `json:"msisdnVerified" default:"true"`
	Password              string `json:"password,omitempty" validate:"omitempty,min=6"`
	ContractorName        string `json:"contractorName,omitempty"`
	ContractorCode        string `json:"contractorCode,omitempty"`
	RegistrationNumber    string `json:"registrationNumber,omitempty"`
	ContactPerson         string `json:"contactPerson,omitempty"`
	CeoName               string `json:"ceoName,omitempty"`
	AuthorizedSignatories string `json:"authorizedSignatories,omitempty"`
	PhoneNumber           string `json:"phoneNumber,omitempty"`
	Email                 string `json:"email,omitempty"`
	Address               string `json:"address,omitempty"`
	ContractType          string `json:"contractType,omitempty"`
	BankAccountNumber     string `json:"bankAccountNumber,omitempty"`
	Description           string `json:"description,omitempty"`
	Role                  string `json:"role,omitempty"`
}

type UserResponse struct {
	User
	UserID string `json:"userId"`
}

func (u *UserRequest) ToEntity() *idpentities.User {
	if u == nil {
		return nil
	}

	user := new(idpentities.User)
	if u.Msisdn != "" {
		user.Msisdn = &u.Msisdn
	} else {
		user.Msisdn = nil
	}
	user.MsisdnVerified = u.MsisdnVerified

	if u.Password != "" {
		user.ClearPassword = []byte(u.Password)
	} else {
		user.ClearPassword = nil
	}

	user.ContractorName = u.ContractorName
	user.ContractorCode = u.ContractorCode
	user.RegistrationNumber = u.RegistrationNumber
	user.ContactPerson = u.ContactPerson
	user.CeoName = u.CeoName
	user.AuthorizedSignatories = u.AuthorizedSignatories
	user.PhoneNumber = u.PhoneNumber
	user.Email = u.Email
	user.Address = u.Address
	user.ContractType = u.ContractType
	user.BankAccountNumber = u.BankAccountNumber
	user.Description = u.Description
	user.Role = idpentities.Role(u.Role)

	return user
}

func NewUserCreateFromRequest(r *http.Request) (*UserRequest, error) {
	decoder := json.NewDecoder(r.Body)
	req := new(UserRequest)
	if err := decoder.Decode(req); err != nil {
		return nil, err
	}
	return req, nil
}

func (u *UserRequest) Validate(v *validator.Validate) error {
	err := v.RegisterValidation("msisdn", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(msisdnRegex).MatchString(fl.Field().String())
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

type UpdateUserRequest struct {
	Msisdn                string `json:"msisdn,omitempty" validate:"required_without=email,omitempty,msisdn"`
	MsisdnVerified        bool   `json:"msisdnVerified" default:"true"`
	Password              string `json:"password,omitempty" validate:"omitempty,min=6"`
	ContractorName        string `json:"contractorName,omitempty"`
	ContractorCode        string `json:"contractorCode,omitempty"`
	RegistrationNumber    string `json:"registrationNumber,omitempty"`
	ContactPerson         string `json:"contactPerson,omitempty"`
	CeoName               string `json:"ceoName,omitempty"`
	AuthorizedSignatories string `json:"authorizedSignatories,omitempty"`
	PhoneNumber           string `json:"phoneNumber,omitempty"`
	Email                 string `json:"email,omitempty"`
	Address               string `json:"address,omitempty"`
	ContractType          string `json:"contractType,omitempty"`
	BankAccountNumber     string `json:"bankAccountNumber,omitempty"`
	Description           string `json:"description,omitempty"`
	Role                  string `json:"role,omitempty"`
}

func (u *UpdateUserRequest) ToEntity() *idpentities.User {
	if u == nil {
		return nil
	}

	user := new(idpentities.User)
	if u.Msisdn != "" {
		user.Msisdn = &u.Msisdn
	} else {
		user.Msisdn = nil
	}
	user.MsisdnVerified = u.MsisdnVerified

	if u.Password != "" {
		user.ClearPassword = []byte(u.Password)
	} else {
		user.ClearPassword = nil
	}

	user.ContractorName = u.ContractorName
	user.ContractorCode = u.ContractorCode
	user.RegistrationNumber = u.RegistrationNumber
	user.ContactPerson = u.ContactPerson
	user.CeoName = u.CeoName
	user.AuthorizedSignatories = u.AuthorizedSignatories
	user.PhoneNumber = u.PhoneNumber
	user.Email = u.Email
	user.Address = u.Address
	user.ContractType = u.ContractType
	user.BankAccountNumber = u.BankAccountNumber
	user.Description = u.Description
	user.Role = idpentities.Role(u.Role)

	return user
}

func NewUserUpdateFromRequest(r *http.Request) (*UserRequest, error) {
	decoder := json.NewDecoder(r.Body)
	req := new(UserRequest)
	if err := decoder.Decode(req); err != nil {
		return nil, err
	}
	return req, nil
}

func NewUserCreateResponse(user *idpentities.User) *UserResponse {
	if user == nil {
		return nil
	}

	resp := NewUserFromEntity(user)
	return resp
}

func NewUserFromEntity(user *idpentities.User) *UserResponse {
	if user == nil {
		return nil
	}

	resp := new(UserResponse)

	if user.Msisdn != nil {
		resp.Msisdn = *user.Msisdn
	}

	resp.MsisdnVerified = user.MsisdnVerified
	resp.ContractorName = user.ContractorName
	resp.ContractorCode = user.ContractorCode
	resp.RegistrationNumber = user.RegistrationNumber
	resp.ContactPerson = user.ContactPerson
	resp.CeoName = user.CeoName
	resp.AuthorizedSignatories = user.AuthorizedSignatories
	resp.PhoneNumber = user.PhoneNumber
	resp.Email = user.Email
	resp.Address = user.Address
	resp.ContractType = user.ContractType
	resp.BankAccountNumber = user.BankAccountNumber
	resp.Description = user.Description
	resp.UserID = user.UserID
	resp.Role = string(user.Role)

	return resp
}
