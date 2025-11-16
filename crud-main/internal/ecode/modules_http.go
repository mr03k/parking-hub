package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// modules business-level http error codes.
// the modulesNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	modulesNO       = 69
	modulesName     = "modules"
	modulesBaseCode = errcode.HCode(modulesNO)

	ErrCreateModules     = errcode.NewError(modulesBaseCode+1, "failed to create "+modulesName)
	ErrDeleteByIDModules = errcode.NewError(modulesBaseCode+2, "failed to delete "+modulesName)
	ErrUpdateByIDModules = errcode.NewError(modulesBaseCode+3, "failed to update "+modulesName)
	ErrGetByIDModules    = errcode.NewError(modulesBaseCode+4, "failed to get "+modulesName+" details")
	ErrListModules       = errcode.NewError(modulesBaseCode+5, "failed to list of "+modulesName)

	// error codes are globally unique, adding 1 to the previous error code
)
