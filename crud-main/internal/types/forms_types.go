package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateFormsRequest request params
type CreateFormsRequest struct {
	FormName    string `json:"formName" binding:""`
	Description string `json:"description" binding:""`
}

// UpdateFormsByIDRequest request params
type UpdateFormsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	FormName    string `json:"formName" binding:""`
	Description string `json:"description" binding:""`
}

// FormsObjDetail detail
type FormsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	FormName    string `json:"formName"`
	Description string `json:"description"`
	CreatedAt   int    `json:"createdAt"`
}

// CreateFormsReply only for api docs
type CreateFormsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteFormsByIDReply only for api docs
type DeleteFormsByIDReply struct {
	Result
}

// UpdateFormsByIDReply only for api docs
type UpdateFormsByIDReply struct {
	Result
}

// GetFormsByIDReply only for api docs
type GetFormsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Forms FormsObjDetail `json:"forms"`
	} `json:"data"` // return data
}

// ListFormssRequest request params
type ListFormssRequest struct {
	query.Params
}

// ListFormssReply only for api docs
type ListFormssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Formss []FormsObjDetail `json:"formss"`
	} `json:"data"` // return data
}
