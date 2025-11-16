package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// assignments business-level http error codes.
// the assignmentsNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	assignmentsNO       = 45
	assignmentsName     = "assignments"
	assignmentsBaseCode = errcode.HCode(assignmentsNO)

	ErrCreateAssignments     = errcode.NewError(assignmentsBaseCode+1, "failed to create "+assignmentsName)
	ErrDeleteByIDAssignments = errcode.NewError(assignmentsBaseCode+2, "failed to delete "+assignmentsName)
	ErrUpdateByIDAssignments = errcode.NewError(assignmentsBaseCode+3, "failed to update "+assignmentsName)
	ErrGetByIDAssignments    = errcode.NewError(assignmentsBaseCode+4, "failed to get "+assignmentsName+" details")
	ErrListAssignments       = errcode.NewError(assignmentsBaseCode+5, "failed to list of "+assignmentsName)

	// error codes are globally unique, adding 1 to the previous error code
)
