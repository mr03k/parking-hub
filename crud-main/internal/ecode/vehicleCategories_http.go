package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// vehicleCategories business-level http error codes.
// the vehicleCategoriesNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	vehicleCategoriesNO       = 4
	vehicleCategoriesName     = "vehicleCategories"
	vehicleCategoriesBaseCode = errcode.HCode(vehicleCategoriesNO)

	ErrCreateVehicleCategories     = errcode.NewError(vehicleCategoriesBaseCode+1, "failed to create "+vehicleCategoriesName)
	ErrDeleteByIDVehicleCategories = errcode.NewError(vehicleCategoriesBaseCode+2, "failed to delete "+vehicleCategoriesName)
	ErrUpdateByIDVehicleCategories = errcode.NewError(vehicleCategoriesBaseCode+3, "failed to update "+vehicleCategoriesName)
	ErrGetByIDVehicleCategories    = errcode.NewError(vehicleCategoriesBaseCode+4, "failed to get "+vehicleCategoriesName+" details")
	ErrListVehicleCategories       = errcode.NewError(vehicleCategoriesBaseCode+5, "failed to list of "+vehicleCategoriesName)

	// error codes are globally unique, adding 1 to the previous error code
)
