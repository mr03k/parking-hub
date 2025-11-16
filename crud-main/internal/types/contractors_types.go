package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateContractorsRequest request params
type CreateContractorsRequest struct {
	ContractorName        string `json:"contractorName" binding:""`
	CodeContractor        string `json:"codeContractor" binding:""`
	NumberRegistration    string `json:"numberRegistration" binding:""`
	PersonContact         string `json:"personContact" binding:""`
	CeoName               string `json:"ceoName" binding:""`
	SignatoriesAuthorized string `json:"signatoriesAuthorized" binding:""`
	PhoneNumber           string `json:"phoneNumber" binding:""`
	Email                 string `json:"email" binding:""`
	Address               string `json:"address" binding:""`
	TypeContract          string `json:"typeContract" binding:""`
	NumberAccountBank     string `json:"numberAccountBank" binding:""`
	Description           string `json:"description" binding:""`
}

// UpdateContractorsByIDRequest request params
type UpdateContractorsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	ContractorName        string `json:"contractorName" binding:""`
	CodeContractor        string `json:"codeContractor" binding:""`
	NumberRegistration    string `json:"numberRegistration" binding:""`
	PersonContact         string `json:"personContact" binding:""`
	CeoName               string `json:"ceoName" binding:""`
	SignatoriesAuthorized string `json:"signatoriesAuthorized" binding:""`
	PhoneNumber           string `json:"phoneNumber" binding:""`
	Email                 string `json:"email" binding:""`
	Address               string `json:"address" binding:""`
	TypeContract          string `json:"typeContract" binding:""`
	NumberAccountBank     string `json:"numberAccountBank" binding:""`
	Description           string `json:"description" binding:""`
}

// ContractorsObjDetail detail
type ContractorsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	ContractorName        string `json:"contractorName"`
	CodeContractor        string `json:"codeContractor"`
	NumberRegistration    string `json:"numberRegistration"`
	PersonContact         string `json:"personContact"`
	CeoName               string `json:"ceoName"`
	SignatoriesAuthorized string `json:"signatoriesAuthorized"`
	PhoneNumber           string `json:"phoneNumber"`
	Email                 string `json:"email"`
	Address               string `json:"address"`
	TypeContract          string `json:"typeContract"`
	NumberAccountBank     string `json:"numberAccountBank"`
	CreatedAt             int    `json:"createdAt"`
	Description           string `json:"description"`
}

// CreateContractorsReply only for api docs
type CreateContractorsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteContractorsByIDReply only for api docs
type DeleteContractorsByIDReply struct {
	Result
}

// UpdateContractorsByIDReply only for api docs
type UpdateContractorsByIDReply struct {
	Result
}

// GetContractorsByIDReply only for api docs
type GetContractorsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Contractors ContractorsObjDetail `json:"contractors"`
	} `json:"data"` // return data
}

// ListContractorssRequest request params
type ListContractorssRequest struct {
	query.Params
}

// ListContractorssReply only for api docs
type ListContractorssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Contractorss []ContractorsObjDetail `json:"contractorss"`
	} `json:"data"` // return data
}
