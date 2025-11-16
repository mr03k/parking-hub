package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// contractors business-level http error codes.
// the contractorsNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	contractorsNO       = 78
	contractorsName     = "contractors"
	contractorsBaseCode = errcode.HCode(contractorsNO)

	ErrCreateContractors     = errcode.NewError(contractorsBaseCode+1, "failed to create "+contractorsName)
	ErrDeleteByIDContractors = errcode.NewError(contractorsBaseCode+2, "failed to delete "+contractorsName)
	ErrUpdateByIDContractors = errcode.NewError(contractorsBaseCode+3, "failed to update "+contractorsName)
	ErrGetByIDContractors    = errcode.NewError(contractorsBaseCode+4, "failed to get "+contractorsName+" details")
	ErrListContractors       = errcode.NewError(contractorsBaseCode+5, "failed to list of "+contractorsName)

	// error codes are globally unique, adding 1 to the previous error code
)
