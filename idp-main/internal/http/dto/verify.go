package dto

import (
	"encoding/json"
	"net/http"

	idpentities "application/internal/entity/idp"
)

type VerifyRequest struct {
	Msisdn   string `json:"msisdn,omitempty" validate:"required_without=Email,omitempty,msisdn" `
	Password string `json:"password,omitempty" validate:"omitempty"`
}

type VerifyResponse struct {
	User UserResponse `json:"user"`
}

// Get verify request from http
func NewUserVerifyFromRequest(r *http.Request) (*VerifyRequest, error) {
	reg := new(VerifyRequest)
	err := json.NewDecoder(r.Body).Decode(reg)
	if err != nil {
		return nil, err
	}
	return reg, nil
}

// VerifyResponse
func NewUserVerifyResponse(user *idpentities.User) *VerifyResponse {
	resp := &VerifyResponse{}
	if user == nil {
		return nil
	}
	resp.User = *NewUserFromEntity(user)
	return resp
}
