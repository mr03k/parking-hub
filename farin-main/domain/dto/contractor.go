package dto

import (
	"farin/domain/entity"
	"farin/infrastructure/godotenv"
)

type ContractorRequest struct {
	ID                    string `json:"-"` //for update unique check
	ContractorName        string `json:"contractorName" binding:"required,min=3,max=100"`
	CodeContractor        string `json:"codeContractor" binding:"required,min=3,max=10,uniqueGorm=contractors&code_contractor"`
	NumberRegistration    string `json:"numberRegistration,omitempty"`
	PersonContact         string `json:"personContact,omitempty"`
	CEOName               string `json:"ceoName,omitempty"`
	SignatoriesAuthorized string `json:"signatoriesAuthorized,omitempty"`
	PhoneNumber           string `json:"phoneNumber,omitempty" binding:"omitempty,min=10,max=15"`
	Email                 string `json:"email,omitempty" binding:"omitempty,email"`
	Address               string `json:"address,omitempty"`
	TypeContract          string `json:"typeContract,omitempty"`
	NumberAccountBank     string `json:"numberAccountBank,omitempty"`
	Description           string `json:"description,omitempty"`
}

func (req *ContractorRequest) ToEntity() *entity.Contractor {
	return &entity.Contractor{
		ContractorName:        req.ContractorName,
		CodeContractor:        req.CodeContractor,
		NumberRegistration:    req.NumberRegistration,
		PersonContact:         req.PersonContact,
		CEOName:               req.CEOName,
		SignatoriesAuthorized: req.SignatoriesAuthorized,
		PhoneNumber:           req.PhoneNumber,
		Email:                 req.Email,
		Address:               req.Address,
		TypeContract:          req.TypeContract,
		NumberAccountBank:     req.NumberAccountBank,
		Description:           req.Description,
	}
}

type ContractorResponse struct {
	ID                    string `json:"id"`
	ContractorName        string `json:"contractorName"`
	CodeContractor        string `json:"codeContractor"`
	NumberRegistration    string `json:"numberRegistration,omitempty"`
	PersonContact         string `json:"personContact,omitempty"`
	CEOName               string `json:"ceoName,omitempty"`
	SignatoriesAuthorized string `json:"signatoriesAuthorized,omitempty"`
	PhoneNumber           string `json:"phoneNumber,omitempty"`
	Email                 string `json:"email,omitempty"`
	Address               string `json:"address,omitempty"`
	TypeContract          string `json:"typeContract,omitempty"`
	NumberAccountBank     string `json:"numberAccountBank,omitempty"`
	Description           string `json:"description,omitempty"`
	CreatedAt             int64  `json:"createdAt"`
	UpdatedAt             int64  `json:"updatedAt"`
}

func (resp *ContractorResponse) FromEntity(contractor *entity.Contractor, env *godotenv.Env) {
	resp.ID = contractor.ID
	resp.ContractorName = contractor.ContractorName
	resp.CodeContractor = contractor.CodeContractor
	resp.NumberRegistration = contractor.NumberRegistration
	resp.PersonContact = contractor.PersonContact
	resp.CEOName = contractor.CEOName
	resp.SignatoriesAuthorized = contractor.SignatoriesAuthorized
	resp.PhoneNumber = contractor.PhoneNumber
	resp.Email = contractor.Email
	resp.Address = contractor.Address
	resp.TypeContract = contractor.TypeContract
	resp.NumberAccountBank = contractor.NumberAccountBank
	resp.Description = contractor.Description
	resp.CreatedAt = contractor.CreatedAt
	resp.UpdatedAt = contractor.UpdatedAt
}

type ContractorListResponse struct {
	Contractors []ContractorResponse `json:"contractors"`
	Total       int64                `json:"total"`
}
