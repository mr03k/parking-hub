package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// rates business-level http error codes.
// the ratesNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	ratesNO       = 96
	ratesName     = "rates"
	ratesBaseCode = errcode.HCode(ratesNO)

	ErrCreateRates     = errcode.NewError(ratesBaseCode+1, "failed to create "+ratesName)
	ErrDeleteByIDRates = errcode.NewError(ratesBaseCode+2, "failed to delete "+ratesName)
	ErrUpdateByIDRates = errcode.NewError(ratesBaseCode+3, "failed to update "+ratesName)
	ErrGetByIDRates    = errcode.NewError(ratesBaseCode+4, "failed to get "+ratesName+" details")
	ErrListRates       = errcode.NewError(ratesBaseCode+5, "failed to list of "+ratesName)

	// error codes are globally unique, adding 1 to the previous error code
)
