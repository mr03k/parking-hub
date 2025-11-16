package dto

import (
	"farin/domain/entity"
)

type RoleRequest struct {
	Title string `json:"title" binding:"required,min=1,max=100"`
}

func (req *RoleRequest) ToEntity() *entity.Role {
	return &entity.Role{
		Title: req.Title,
	}
}

type RoleResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

func (resp *RoleResponse) FromEntity(role *entity.Role) {
	resp.ID = role.ID
	resp.Title = role.Title
	resp.CreatedAt = role.CreatedAt
	resp.UpdatedAt = role.UpdatedAt
}

type RoleListResponse struct {
	Roles []RoleResponse `json:"roles"`
	Total int64          `json:"total"`
}
