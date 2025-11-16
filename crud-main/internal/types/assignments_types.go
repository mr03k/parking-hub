package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateAssignmentsRequest request params
type CreateAssignmentsRequest struct {
	IDUser        string `json:"iDUser" binding:""`
	IDRole        string `json:"iDRole" binding:""`
	IDModule      string `json:"iDModule" binding:""`
	IDForm        string `json:"iDForm" binding:""`
	AccessEndDate int    `json:"accessEndDate" binding:""`
}

// UpdateAssignmentsByIDRequest request params
type UpdateAssignmentsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	IDUser        string `json:"iDUser" binding:""`
	IDRole        string `json:"iDRole" binding:""`
	IDModule      string `json:"iDModule" binding:""`
	IDForm        string `json:"iDForm" binding:""`
	AccessEndDate int    `json:"accessEndDate" binding:""`
}

// AssignmentsObjDetail detail
type AssignmentsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	IDUser        string `json:"iDUser"`
	IDRole        string `json:"iDRole"`
	IDModule      string `json:"iDModule"`
	IDForm        string `json:"iDForm"`
	CreatedAt     int    `json:"createdAt"`
	AccessEndDate int    `json:"accessEndDate"`
}

// CreateAssignmentsReply only for api docs
type CreateAssignmentsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteAssignmentsByIDReply only for api docs
type DeleteAssignmentsByIDReply struct {
	Result
}

// UpdateAssignmentsByIDReply only for api docs
type UpdateAssignmentsByIDReply struct {
	Result
}

// GetAssignmentsByIDReply only for api docs
type GetAssignmentsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Assignments AssignmentsObjDetail `json:"assignments"`
	} `json:"data"` // return data
}

// ListAssignmentssRequest request params
type ListAssignmentssRequest struct {
	query.Params
}

// ListAssignmentssReply only for api docs
type ListAssignmentssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Assignmentss []AssignmentsObjDetail `json:"assignmentss"`
	} `json:"data"` // return data
}
