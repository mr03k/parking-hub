package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// parts business-level http error codes.
// the partsNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	partsNO       = 34
	partsName     = "parts"
	partsBaseCode = errcode.HCode(partsNO)

	ErrCreateParts     = errcode.NewError(partsBaseCode+1, "failed to create "+partsName)
	ErrDeleteByIDParts = errcode.NewError(partsBaseCode+2, "failed to delete "+partsName)
	ErrUpdateByIDParts = errcode.NewError(partsBaseCode+3, "failed to update "+partsName)
	ErrGetByIDParts    = errcode.NewError(partsBaseCode+4, "failed to get "+partsName+" details")
	ErrListParts       = errcode.NewError(partsBaseCode+5, "failed to list of "+partsName)

	// error codes are globally unique, adding 1 to the previous error code
)
