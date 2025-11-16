package dto

import (
	"application/internal/entity/device"
	"encoding/json"
	"net/http"
)

type ContractorRequest struct {
	Name                    string `json:"contractor_name"`
	Code                    string `json:"contractor_code"`
	RegisterNumber          string `json:"register_number"`
	ContactPerson           string `json:"contact_person"`
	CeoName                 string `json:"ceo_name"`
	AutorizationSignatories string `json:"autorization_signatories"`
	PhoneNumbers            string `json:"phone_numbers"`
	Email                   string `json:"email"`
	Address                 string `json:"address"`
	ContractType            string `json:"contract_type"`
	BankAccountNumber       string `json:"bank_account_number"`
	Description             string `json:"description"`
}

type ContractorResponse struct {
	ID string `json:"id"`
	ContractorRequest
}

type ConttactorCreateRequest struct {
	ContractorRequest
}

type ConttactorCreateResponse struct {
	ContractorResponse
}

// list response
type ConttactorListResponse struct {
	Count       int                 `json:"count"`
	Contractors []device.Contractor `json:"contractors"`
}

func NewConttactorCreateRequestFromRequest(r *http.Request) (*ConttactorCreateRequest, error) {
	req := &ConttactorCreateRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func NewConttactorListResponse(count int, contractors []device.Contractor) *ConttactorListResponse {
	return &ConttactorListResponse{
		Count:       count,
		Contractors: contractors,
	}
}
