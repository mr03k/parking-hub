package types

import (
	"encoding/json"
	"net/http"
)

type VerifyRequest struct {
	Msisdn   string `json:"msisdn,omitempty" validate:"required_without=Email,omitempty,msisdn" `
	Password string `json:"password,omitempty" validate:"omitempty"`
}

func NewCustomerVerifyFromRequest(r *http.Request) (*VerifyRequest, error) {
	reg := new(VerifyRequest)
	err := json.NewDecoder(r.Body).Decode(reg)
	if err != nil {
		return nil, err
	}
	return reg, nil
}
