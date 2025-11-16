package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// forms business-level http error codes.
// the formsNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	formsNO       = 93
	formsName     = "forms"
	formsBaseCode = errcode.HCode(formsNO)

	ErrCreateForms     = errcode.NewError(formsBaseCode+1, "failed to create "+formsName)
	ErrDeleteByIDForms = errcode.NewError(formsBaseCode+2, "failed to delete "+formsName)
	ErrUpdateByIDForms = errcode.NewError(formsBaseCode+3, "failed to update "+formsName)
	ErrGetByIDForms    = errcode.NewError(formsBaseCode+4, "failed to get "+formsName+" details")
	ErrListForms       = errcode.NewError(formsBaseCode+5, "failed to list of "+formsName)

	// error codes are globally unique, adding 1 to the previous error code
)
