package dto

import (
	"encoding/json"
	"net/http"
	"time"
)

type DriverLogin struct {
	VerifyRequest
}

type DriverLoginResponse struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

// NewDriverLoginFromRequest
func NewDriverLoginFromRequest(r *http.Request) (*DriverLogin, error) {
	reg := new(DriverLogin)
	err := json.NewDecoder(r.Body).Decode(reg)
	if err != nil {
		return nil, err
	}
	return reg, nil
}

// NewDriverLoginResponse
func NewDriverLoginResponse(userID, token string) *DriverLoginResponse {
	return &DriverLoginResponse{
		Token:  token,
		UserID: userID,
	}
}

// DriverResponse represents the summary response for a driver.
type DriverResponse struct {
	ID               string `json:"id"`
	Address          string `json:"address"`
	DriverType       string `json:"driverType"`
	ShiftType        string `json:"shiftType"`
	EmploymentStatus string `json:"employmentStatus"`
}

// DriverDetailResponse represents the detailed response for a driver.
type DriverDetailResponse struct {
	ID                          string     `json:"id"`
	Address                     string     `json:"address"`
	DriverType                  string     `json:"driverType"`
	ShiftType                   string     `json:"shiftType"`
	EmploymentStatus            string     `json:"employmentStatus"`
	EmploymentStartDate         *time.Time `json:"employmentStartDate,omitempty"`
	EmploymentEndDate           *time.Time `json:"employmentEndDate,omitempty"`
	DriverPhotoURL              string     `json:"driverPhotoUrl,omitempty"`
	IDCardImageURL              string     `json:"idCardImageUrl,omitempty"`
	BirthCertificateImageURL    string     `json:"birthCertificateImageUrl,omitempty"`
	MilitaryServiceCardImageURL string     `json:"militaryServiceCardImageUrl,omitempty"`
	HealthCertificateImageURL   string     `json:"healthCertificateImageUrl,omitempty"`
	CriminalRecordImageURL      string     `json:"criminalRecordImageUrl,omitempty"`
	Description                 string     `json:"description,omitempty"`
}

// DriverListResponse represents the response for a list of drivers.
type DriverListResponse struct {
	Drivers []DriverResponse `json:"drivers"`
}
