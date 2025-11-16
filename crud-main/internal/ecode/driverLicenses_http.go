package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// driverLicenses business-level http error codes.
// the driverLicensesNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	driverLicensesNO       = 17
	driverLicensesName     = "driverLicenses"
	driverLicensesBaseCode = errcode.HCode(driverLicensesNO)

	ErrCreateDriverLicenses     = errcode.NewError(driverLicensesBaseCode+1, "failed to create "+driverLicensesName)
	ErrDeleteByIDDriverLicenses = errcode.NewError(driverLicensesBaseCode+2, "failed to delete "+driverLicensesName)
	ErrUpdateByIDDriverLicenses = errcode.NewError(driverLicensesBaseCode+3, "failed to update "+driverLicensesName)
	ErrGetByIDDriverLicenses    = errcode.NewError(driverLicensesBaseCode+4, "failed to get "+driverLicensesName+" details")
	ErrListDriverLicenses       = errcode.NewError(driverLicensesBaseCode+5, "failed to list of "+driverLicensesName)

	// error codes are globally unique, adding 1 to the previous error code
)
