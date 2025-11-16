package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// countries business-level http error codes.
// the countriesNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	countriesNO       = 51
	countriesName     = "countries"
	countriesBaseCode = errcode.HCode(countriesNO)

	ErrCreateCountries     = errcode.NewError(countriesBaseCode+1, "failed to create "+countriesName)
	ErrDeleteByIDCountries = errcode.NewError(countriesBaseCode+2, "failed to delete "+countriesName)
	ErrUpdateByIDCountries = errcode.NewError(countriesBaseCode+3, "failed to update "+countriesName)
	ErrGetByIDCountries    = errcode.NewError(countriesBaseCode+4, "failed to get "+countriesName+" details")
	ErrListCountries       = errcode.NewError(countriesBaseCode+5, "failed to list of "+countriesName)

	// error codes are globally unique, adding 1 to the previous error code
)
