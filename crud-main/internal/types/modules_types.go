package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateModulesRequest request params
type CreateModulesRequest struct {
	ModuleName  string `json:"moduleName" binding:""`
	Description string `json:"description" binding:""`
}

// UpdateModulesByIDRequest request params
type UpdateModulesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	ModuleName  string `json:"moduleName" binding:""`
	Description string `json:"description" binding:""`
}

// ModulesObjDetail detail
type ModulesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	ModuleName  string `json:"moduleName"`
	Description string `json:"description"`
	CreatedAt   int    `json:"createdAt"`
}

// CreateModulesReply only for api docs
type CreateModulesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteModulesByIDReply only for api docs
type DeleteModulesByIDReply struct {
	Result
}

// UpdateModulesByIDReply only for api docs
type UpdateModulesByIDReply struct {
	Result
}

// GetModulesByIDReply only for api docs
type GetModulesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Modules ModulesObjDetail `json:"modules"`
	} `json:"data"` // return data
}

// ListModulessRequest request params
type ListModulessRequest struct {
	query.Params
}

// ListModulessReply only for api docs
type ListModulessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Moduless []ModulesObjDetail `json:"moduless"`
	} `json:"data"` // return data
}
