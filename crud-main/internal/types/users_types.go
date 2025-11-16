package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateUsersRequest request params
type CreateUsersRequest struct {
	Username     string `json:"username" binding:""`
	Password     string `json:"password" binding:""`
	FirstName    string `json:"firstName" binding:""`
	LastName     string `json:"lastName" binding:""`
	Email        string `json:"email" binding:""`
	NumberPhone  string `json:"numberPhone" binding:""`
	NumberMobile string `json:"numberMobile" binding:""`
	IDNational   string `json:"iDNational" binding:""`
	CodePostal   string `json:"codePostal" binding:""`
	NameCompany  string `json:"nameCompany" binding:""`
	ImageProfile string `json:"imageProfile" binding:""`
	Gender       string `json:"gender" binding:""`
	Address      string `json:"address" binding:""`
	Status       string `json:"status" binding:""`
}

// UpdateUsersByIDRequest request params
type UpdateUsersByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	Username     string `json:"username" binding:""`
	Password     string `json:"password" binding:""`
	FirstName    string `json:"firstName" binding:""`
	LastName     string `json:"lastName" binding:""`
	Email        string `json:"email" binding:""`
	NumberPhone  string `json:"numberPhone" binding:""`
	NumberMobile string `json:"numberMobile" binding:""`
	IDNational   string `json:"iDNational" binding:""`
	CodePostal   string `json:"codePostal" binding:""`
	NameCompany  string `json:"nameCompany" binding:""`
	ImageProfile string `json:"imageProfile" binding:""`
	Gender       string `json:"gender" binding:""`
	Address      string `json:"address" binding:""`
	Status       string `json:"status" binding:""`
}

// UsersObjDetail detail
type UsersObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	Username     string `json:"username"`
	Password     string `json:"password"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	NumberPhone  string `json:"numberPhone"`
	NumberMobile string `json:"numberMobile"`
	IDNational   string `json:"iDNational"`
	CodePostal   string `json:"codePostal"`
	NameCompany  string `json:"nameCompany"`
	ImageProfile string `json:"imageProfile"`
	Gender       string `json:"gender"`
	Address      string `json:"address"`
	Status       string `json:"status"`
	CreatedAt    int    `json:"createdAt"`
}

// CreateUsersReply only for api docs
type CreateUsersReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteUsersByIDReply only for api docs
type DeleteUsersByIDReply struct {
	Result
}

// UpdateUsersByIDReply only for api docs
type UpdateUsersByIDReply struct {
	Result
}

// GetUsersByIDReply only for api docs
type GetUsersByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Users UsersObjDetail `json:"users"`
	} `json:"data"` // return data
}

// ListUserssRequest request params
type ListUserssRequest struct {
	query.Params
}

// ListUserssReply only for api docs
type ListUserssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Userss []UsersObjDetail `json:"userss"`
	} `json:"data"` // return data
}
