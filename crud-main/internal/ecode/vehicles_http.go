package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// vehicles business-level http error codes.
// the vehiclesNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	vehiclesNO       = 6
	vehiclesName     = "vehicles"
	vehiclesBaseCode = errcode.HCode(vehiclesNO)

	ErrCreateVehicles     = errcode.NewError(vehiclesBaseCode+1, "failed to create "+vehiclesName)
	ErrDeleteByIDVehicles = errcode.NewError(vehiclesBaseCode+2, "failed to delete "+vehiclesName)
	ErrUpdateByIDVehicles = errcode.NewError(vehiclesBaseCode+3, "failed to update "+vehiclesName)
	ErrGetByIDVehicles    = errcode.NewError(vehiclesBaseCode+4, "failed to get "+vehiclesName+" details")
	ErrListVehicles       = errcode.NewError(vehiclesBaseCode+5, "failed to list of "+vehiclesName)

	// error codes are globally unique, adding 1 to the previous error code
)
