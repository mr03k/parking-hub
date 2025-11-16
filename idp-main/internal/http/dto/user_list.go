package dto

import (
	idpentities "application/internal/entity/idp"
)

type UserListResponse struct {
	Total int             `json:"total"`
	Users []*UserResponse `json:"users"`
}

func NewUserListResponse(users []idpentities.User) *UserListResponse {
	var res UserListResponse
	res.Total = len(users)
	res.Users = make([]*UserResponse, res.Total)

	for i, user := range users {
		res.Users[i] = NewUserFromEntity(&user)
	}

	return &res
}
