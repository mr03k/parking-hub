package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// calendars business-level http error codes.
// the calendarsNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	calendarsNO       = 37
	calendarsName     = "calendars"
	calendarsBaseCode = errcode.HCode(calendarsNO)

	ErrCreateCalendars     = errcode.NewError(calendarsBaseCode+1, "failed to create "+calendarsName)
	ErrDeleteByIDCalendars = errcode.NewError(calendarsBaseCode+2, "failed to delete "+calendarsName)
	ErrUpdateByIDCalendars = errcode.NewError(calendarsBaseCode+3, "failed to update "+calendarsName)
	ErrGetByIDCalendars    = errcode.NewError(calendarsBaseCode+4, "failed to get "+calendarsName+" details")
	ErrListCalendars       = errcode.NewError(calendarsBaseCode+5, "failed to list of "+calendarsName)

	// error codes are globally unique, adding 1 to the previous error code
)
