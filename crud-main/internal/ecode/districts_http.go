package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// districts business-level http error codes.
// the districtsNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	districtsNO       = 99
	districtsName     = "districts"
	districtsBaseCode = errcode.HCode(districtsNO)

	ErrCreateDistricts     = errcode.NewError(districtsBaseCode+1, "failed to create "+districtsName)
	ErrDeleteByIDDistricts = errcode.NewError(districtsBaseCode+2, "failed to delete "+districtsName)
	ErrUpdateByIDDistricts = errcode.NewError(districtsBaseCode+3, "failed to update "+districtsName)
	ErrGetByIDDistricts    = errcode.NewError(districtsBaseCode+4, "failed to get "+districtsName+" details")
	ErrListDistricts       = errcode.NewError(districtsBaseCode+5, "failed to list of "+districtsName)

	// error codes are globally unique, adding 1 to the previous error code
)
