package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// rings business-level http error codes.
// the ringsNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	ringsNO       = 87
	ringsName     = "rings"
	ringsBaseCode = errcode.HCode(ringsNO)

	ErrCreateRings     = errcode.NewError(ringsBaseCode+1, "failed to create "+ringsName)
	ErrDeleteByIDRings = errcode.NewError(ringsBaseCode+2, "failed to delete "+ringsName)
	ErrUpdateByIDRings = errcode.NewError(ringsBaseCode+3, "failed to update "+ringsName)
	ErrGetByIDRings    = errcode.NewError(ringsBaseCode+4, "failed to get "+ringsName+" details")
	ErrListRings       = errcode.NewError(ringsBaseCode+5, "failed to list of "+ringsName)

	// error codes are globally unique, adding 1 to the previous error code
)
