package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// deviceParts business-level http error codes.
// the devicePartsNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	devicePartsNO       = 81
	devicePartsName     = "deviceParts"
	devicePartsBaseCode = errcode.HCode(devicePartsNO)

	ErrCreateDeviceParts     = errcode.NewError(devicePartsBaseCode+1, "failed to create "+devicePartsName)
	ErrDeleteByIDDeviceParts = errcode.NewError(devicePartsBaseCode+2, "failed to delete "+devicePartsName)
	ErrUpdateByIDDeviceParts = errcode.NewError(devicePartsBaseCode+3, "failed to update "+devicePartsName)
	ErrGetByIDDeviceParts    = errcode.NewError(devicePartsBaseCode+4, "failed to get "+devicePartsName+" details")
	ErrListDeviceParts       = errcode.NewError(devicePartsBaseCode+5, "failed to list of "+devicePartsName)

	// error codes are globally unique, adding 1 to the previous error code
)
