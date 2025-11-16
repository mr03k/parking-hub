package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateContractsRequest request params
type CreateContractsRequest struct {
	NumberContract  string    `json:"numberContract" binding:""`
	DateContract    time.Time `json:"dateContract" binding:""`
	DateStart       time.Time `json:"dateStart" binding:""`
	DateEnd         time.Time `json:"dateEnd" binding:""`
	AmountContract  int64     `json:"amountContract" binding:""`
	TypeContract    string    `json:"typeContract" binding:""`
	IDContractor    string    `json:"iDContractor" binding:""`
	PeriodOperation int       `json:"periodOperation" binding:""`
	PeriodEquipment int       `json:"periodEquipment" binding:""`
	Description     string    `json:"description" binding:""`
}

// UpdateContractsByIDRequest request params
type UpdateContractsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	NumberContract  string    `json:"numberContract" binding:""`
	DateContract    time.Time `json:"dateContract" binding:""`
	DateStart       time.Time `json:"dateStart" binding:""`
	DateEnd         time.Time `json:"dateEnd" binding:""`
	AmountContract  int64     `json:"amountContract" binding:""`
	TypeContract    string    `json:"typeContract" binding:""`
	IDContractor    string    `json:"iDContractor" binding:""`
	PeriodOperation int       `json:"periodOperation" binding:""`
	PeriodEquipment int       `json:"periodEquipment" binding:""`
	Description     string    `json:"description" binding:""`
}

// ContractsObjDetail detail
type ContractsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	NumberContract  string    `json:"numberContract"`
	DateContract    time.Time `json:"dateContract"`
	DateStart       time.Time `json:"dateStart"`
	DateEnd         time.Time `json:"dateEnd"`
	AmountContract  int64     `json:"amountContract"`
	TypeContract    string    `json:"typeContract"`
	IDContractor    string    `json:"iDContractor"`
	PeriodOperation int       `json:"periodOperation"`
	PeriodEquipment int       `json:"periodEquipment"`
	Description     string    `json:"description"`
	CreatedAt       int       `json:"createdAt"`
}

// CreateContractsReply only for api docs
type CreateContractsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteContractsByIDReply only for api docs
type DeleteContractsByIDReply struct {
	Result
}

// UpdateContractsByIDReply only for api docs
type UpdateContractsByIDReply struct {
	Result
}

// GetContractsByIDReply only for api docs
type GetContractsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Contracts ContractsObjDetail `json:"contracts"`
	} `json:"data"` // return data
}

// ListContractssRequest request params
type ListContractssRequest struct {
	query.Params
}

// ListContractssReply only for api docs
type ListContractssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Contractss []ContractsObjDetail `json:"contractss"`
	} `json:"data"` // return data
}
