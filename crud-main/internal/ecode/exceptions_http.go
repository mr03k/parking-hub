package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// exceptions business-level http error codes.
// the exceptionsNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	exceptionsNO       = 19
	exceptionsName     = "exceptions"
	exceptionsBaseCode = errcode.HCode(exceptionsNO)

	ErrCreateExceptions     = errcode.NewError(exceptionsBaseCode+1, "failed to create "+exceptionsName)
	ErrDeleteByIDExceptions = errcode.NewError(exceptionsBaseCode+2, "failed to delete "+exceptionsName)
	ErrUpdateByIDExceptions = errcode.NewError(exceptionsBaseCode+3, "failed to update "+exceptionsName)
	ErrGetByIDExceptions    = errcode.NewError(exceptionsBaseCode+4, "failed to get "+exceptionsName+" details")
	ErrListExceptions       = errcode.NewError(exceptionsBaseCode+5, "failed to list of "+exceptionsName)

	// error codes are globally unique, adding 1 to the previous error code
)
