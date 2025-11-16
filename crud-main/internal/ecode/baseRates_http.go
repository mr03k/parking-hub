package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// baseRates business-level http error codes.
// the baseRatesNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	baseRatesNO       = 46
	baseRatesName     = "baseRates"
	baseRatesBaseCode = errcode.HCode(baseRatesNO)

	ErrCreateBaseRates     = errcode.NewError(baseRatesBaseCode+1, "failed to create "+baseRatesName)
	ErrDeleteByIDBaseRates = errcode.NewError(baseRatesBaseCode+2, "failed to delete "+baseRatesName)
	ErrUpdateByIDBaseRates = errcode.NewError(baseRatesBaseCode+3, "failed to update "+baseRatesName)
	ErrGetByIDBaseRates    = errcode.NewError(baseRatesBaseCode+4, "failed to get "+baseRatesName+" details")
	ErrListBaseRates       = errcode.NewError(baseRatesBaseCode+5, "failed to list of "+baseRatesName)

	// error codes are globally unique, adding 1 to the previous error code
)
