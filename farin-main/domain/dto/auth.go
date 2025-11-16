package dto

import (
	"farin/domain/entity"
	"strings"
)

type LoginRequest struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Remember    bool   `json:"remember"`
}

func (l LoginRequest) ToUser() *entity.User {
	return &entity.User{
		Email:    l.PhoneNumber,
		Password: l.Password,
	}
}

type LoginResponse struct {
	AccessToken     string       `json:"accessToken"`
	RefreshToken    string       `json:"refreshToken"`
	ExpRefreshToken string       `json:"expRefreshToken"`
	ExpAccessToken  string       `json:"expAccessToken"`
	User            UserResponse `json:"user"`
}

type DriverLoginRequest struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Password    string `json:"password" binding:"required"`
	VehicleID   string `json:"vehicleID" binding:"required"`
	Remember    bool   `json:"remember"`
}

func (l DriverLoginRequest) ToUser() *entity.User {
	return &entity.User{
		Email:    l.PhoneNumber,
		Password: l.Password,
	}
}

type DriverLoginResponse struct {
	LoginResponse
	Contract   ContractResponse   `json:"contract"`
	Calender   CalenderResponse   `json:"calender"`
	Vehicle    VehicleResponse    `json:"vehicle"`
	Driver     DriverResponse     `json:"driver"`
	Contractor ContractorResponse `json:"contractor"`
	Devices    []DeviceResponse   `json:"device"`
	Ring       RingResponse       `json:"ring"`
}

type AccessTokenReq struct {
	AccessToken string `json:"accessToken" binding:"required"`
}

type AccessTokenRes struct {
	AccessToken    string `json:"accessToken" binding:"required"`
	ExpAccessToken string `json:"expAccessToken" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type TokenRequest struct {
	Token string `json:"token" binding:"len=40,required"`
}

type TokenRequestNoLimit struct {
	Token string `json:"token" binding:"required"`
}

type OAuthData struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (o OAuthData) ToUser() *entity.User {
	names := strings.Split(o.Name, " ")
	return &entity.User{
		Email:     o.Email,
		FirstName: names[0],
		LastName:  names[1],
	}
}

type EmailFormRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func (e EmailFormRequest) ToUser() *entity.User {
	return &entity.User{
		Email: e.Email,
	}
}
