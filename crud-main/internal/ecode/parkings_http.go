package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// parkings business-level http error codes.
// the parkingsNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	parkingsNO       = 54
	parkingsName     = "parkings"
	parkingsBaseCode = errcode.HCode(parkingsNO)

	ErrCreateParkings     = errcode.NewError(parkingsBaseCode+1, "failed to create "+parkingsName)
	ErrDeleteByIDParkings = errcode.NewError(parkingsBaseCode+2, "failed to delete "+parkingsName)
	ErrUpdateByIDParkings = errcode.NewError(parkingsBaseCode+3, "failed to update "+parkingsName)
	ErrGetByIDParkings    = errcode.NewError(parkingsBaseCode+4, "failed to get "+parkingsName+" details")
	ErrListParkings       = errcode.NewError(parkingsBaseCode+5, "failed to list of "+parkingsName)

	// error codes are globally unique, adding 1 to the previous error code
)
