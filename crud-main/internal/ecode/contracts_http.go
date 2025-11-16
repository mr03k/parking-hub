package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// contracts business-level http error codes.
// the contractsNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	contractsNO       = 70
	contractsName     = "contracts"
	contractsBaseCode = errcode.HCode(contractsNO)

	ErrCreateContracts     = errcode.NewError(contractsBaseCode+1, "failed to create "+contractsName)
	ErrDeleteByIDContracts = errcode.NewError(contractsBaseCode+2, "failed to delete "+contractsName)
	ErrUpdateByIDContracts = errcode.NewError(contractsBaseCode+3, "failed to update "+contractsName)
	ErrGetByIDContracts    = errcode.NewError(contractsBaseCode+4, "failed to get "+contractsName+" details")
	ErrListContracts       = errcode.NewError(contractsBaseCode+5, "failed to list of "+contractsName)

	// error codes are globally unique, adding 1 to the previous error code
)
