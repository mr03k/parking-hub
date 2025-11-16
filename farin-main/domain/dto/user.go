package dto

import (
	"farin/domain/entity"
	"farin/infrastructure/godotenv"
	"farin/util"
)

type UserRequest struct {
	ID           string `json:"-"` //need for fields uniqueness check in update
	Username     string `json:"username" binding:"required,min=3,max=50"`
	Firstname    string `json:"firstName" binding:"required"`
	Lastname     string `json:"lastName" binding:"required"`
	Email        string `json:"email" binding:"required,email,uniqueGorm=users&email"`
	Password     string `json:"password" binding:"required,min=8,max=32"`
	PhoneNumber  string `json:"phoneNumber,omitempty" binding:"required,uniqueGorm=users&phone_number"`
	NationalID   string `json:"nationalId,omitempty" binding:"required,uniqueGorm=users&national_id"`
	PostalCode   string `json:"postalCode,omitempty"`
	ProfileImage []byte `json:"profileImage,omitempty" binding:"fileData=image/jpeg&image/png;2048000"`
	CompanyName  string `json:"companyName,omitempty"`
	Gender       string `json:"gender,omitempty"`
	Address      string `json:"address,omitempty"`
	Status       string `json:"status,omitempty"`
	RoleID       string `json:"roleID"`
}

func (req *UserRequest) ToEntity() *entity.User {
	var gender entity.Gender
	var status entity.Status

	if req.Gender == string(entity.GenderMale) {
		gender = entity.GenderMale
	} else if req.Gender == string(entity.GenderFemale) {
		gender = entity.GenderFemale
	}

	if req.Status == string(entity.StatusActive) {
		status = entity.StatusActive
	} else if req.Status == string(entity.StatusInactive) {
		status = entity.StatusInactive
	}

	return &entity.User{
		Username:    req.Username,
		FirstName:   req.Firstname,
		LastName:    req.Lastname,
		Email:       req.Email,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
		NationalID:  req.NationalID,
		PostalCode:  req.PostalCode,
		CompanyName: req.CompanyName,
		Gender:      gender,
		Address:     req.Address,
		Status:      status,
		RoleID:      req.RoleID,
	}
}

type UserResponse struct {
	ID           string      `json:"id"`
	Username     string      `json:"username"`
	Firstname    string      `json:"firstName"`
	Lastname     string      `json:"lastName"`
	Email        string      `json:"email"`
	PhoneNumber  string      `json:"phoneNumber,omitempty"`
	NationalID   string      `json:"nationalId,omitempty"`
	PostalCode   string      `json:"postalCode,omitempty"`
	CompanyName  string      `json:"companyName,omitempty"`
	ProfileImage string      `json:"profileImage,omitempty"`
	Gender       string      `json:"gender,omitempty"`
	Address      string      `json:"address,omitempty"`
	Status       string      `json:"status,omitempty"`
	Role         entity.Role `json:"role" binding:"required,oneof=admin driver"`
	RoleID       string      `json:"roleID"`
	CreatedAt    int64       `json:"createdAt"`
	UpdatedAt    int64       `json:"updatedAt"`
}

func (resp *UserResponse) FromEntity(user *entity.User, env *godotenv.Env) {
	resp.ID = user.ID
	resp.Username = user.Username
	resp.Firstname = user.FirstName
	resp.Lastname = user.LastName
	resp.Email = user.Email
	resp.PhoneNumber = user.PhoneNumber
	resp.NationalID = user.NationalID
	resp.PostalCode = user.PostalCode
	resp.CompanyName = user.CompanyName
	resp.ProfileImage = user.ProfileImage
	resp.Gender = string(user.Gender)
	resp.Address = user.Address
	resp.CreatedAt = user.CreatedAt
	resp.UpdatedAt = user.UpdatedAt
	resp.Status = string(user.Status)
	resp.Role = user.Role
	resp.ProfileImage = util.GeneratePublicURL(user.ProfileImage, env)
}

type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total"`
}
type NamesRequest struct {
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastName"`
}

func (o NamesRequest) ToEntity() *entity.User {
	return &entity.User{
		FirstName: o.Firstname,
		LastName:  o.Lastname,
	}
}

type PasswordRequest struct {
	Password       string `json:"password" binding:"required,min=8"`
	RepeatPassword string `json:"repeatPassword" binding:"required,eqfield=Password"`
}

func (o PasswordRequest) ToEntity() *entity.User {
	return &entity.User{
		Password: o.Password,
	}
}
