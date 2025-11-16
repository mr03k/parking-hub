package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// cities business-level http error codes.
// the citiesNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	citiesNO       = 29
	citiesName     = "cities"
	citiesBaseCode = errcode.HCode(citiesNO)

	ErrCreateCities     = errcode.NewError(citiesBaseCode+1, "failed to create "+citiesName)
	ErrDeleteByIDCities = errcode.NewError(citiesBaseCode+2, "failed to delete "+citiesName)
	ErrUpdateByIDCities = errcode.NewError(citiesBaseCode+3, "failed to update "+citiesName)
	ErrGetByIDCities    = errcode.NewError(citiesBaseCode+4, "failed to get "+citiesName+" details")
	ErrListCities       = errcode.NewError(citiesBaseCode+5, "failed to list of "+citiesName)

	// error codes are globally unique, adding 1 to the previous error code
)
