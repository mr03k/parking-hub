package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// drivers business-level http error codes.
// the driversNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	driversNO       = 36
	driversName     = "drivers"
	driversBaseCode = errcode.HCode(driversNO)

	ErrCreateDrivers     = errcode.NewError(driversBaseCode+1, "failed to create "+driversName)
	ErrDeleteByIDDrivers = errcode.NewError(driversBaseCode+2, "failed to delete "+driversName)
	ErrUpdateByIDDrivers = errcode.NewError(driversBaseCode+3, "failed to update "+driversName)
	ErrGetByIDDrivers    = errcode.NewError(driversBaseCode+4, "failed to get "+driversName+" details")
	ErrListDrivers       = errcode.NewError(driversBaseCode+5, "failed to list of "+driversName)

	// error codes are globally unique, adding 1 to the previous error code
)
