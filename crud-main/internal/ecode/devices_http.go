package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// devices business-level http error codes.
// the devicesNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	devicesNO       = 66
	devicesName     = "devices"
	devicesBaseCode = errcode.HCode(devicesNO)

	ErrCreateDevices     = errcode.NewError(devicesBaseCode+1, "failed to create "+devicesName)
	ErrDeleteByIDDevices = errcode.NewError(devicesBaseCode+2, "failed to delete "+devicesName)
	ErrUpdateByIDDevices = errcode.NewError(devicesBaseCode+3, "failed to update "+devicesName)
	ErrGetByIDDevices    = errcode.NewError(devicesBaseCode+4, "failed to get "+devicesName+" details")
	ErrListDevices       = errcode.NewError(devicesBaseCode+5, "failed to list of "+devicesName)

	// error codes are globally unique, adding 1 to the previous error code
)
