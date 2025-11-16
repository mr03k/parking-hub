package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateRolesRequest request params
type CreateRolesRequest struct {
	RoleName    string `json:"roleName" binding:""`
	Description string `json:"description" binding:""`
}

// UpdateRolesByIDRequest request params
type UpdateRolesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	RoleName    string `json:"roleName" binding:""`
	Description string `json:"description" binding:""`
}

// RolesObjDetail detail
type RolesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	RoleName    string `json:"roleName"`
	Description string `json:"description"`
	CreatedAt   int    `json:"createdAt"`
}

// CreateRolesReply only for api docs
type CreateRolesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteRolesByIDReply only for api docs
type DeleteRolesByIDReply struct {
	Result
}

// UpdateRolesByIDReply only for api docs
type UpdateRolesByIDReply struct {
	Result
}

// GetRolesByIDReply only for api docs
type GetRolesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Roles RolesObjDetail `json:"roles"`
	} `json:"data"` // return data
}

// ListRolessRequest request params
type ListRolessRequest struct {
	query.Params
}

// ListRolessReply only for api docs
type ListRolessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Roless []RolesObjDetail `json:"roless"`
	} `json:"data"` // return data
}
