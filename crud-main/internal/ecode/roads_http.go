package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// roads business-level http error codes.
// the roadsNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	roadsNO       = 88
	roadsName     = "roads"
	roadsBaseCode = errcode.HCode(roadsNO)

	ErrCreateRoads     = errcode.NewError(roadsBaseCode+1, "failed to create "+roadsName)
	ErrDeleteByIDRoads = errcode.NewError(roadsBaseCode+2, "failed to delete "+roadsName)
	ErrUpdateByIDRoads = errcode.NewError(roadsBaseCode+3, "failed to update "+roadsName)
	ErrGetByIDRoads    = errcode.NewError(roadsBaseCode+4, "failed to get "+roadsName+" details")
	ErrListRoads       = errcode.NewError(roadsBaseCode+5, "failed to list of "+roadsName)

	// error codes are globally unique, adding 1 to the previous error code
)
