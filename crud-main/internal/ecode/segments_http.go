package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// segments business-level http error codes.
// the segmentsNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	segmentsNO       = 26
	segmentsName     = "segments"
	segmentsBaseCode = errcode.HCode(segmentsNO)

	ErrCreateSegments     = errcode.NewError(segmentsBaseCode+1, "failed to create "+segmentsName)
	ErrDeleteByIDSegments = errcode.NewError(segmentsBaseCode+2, "failed to delete "+segmentsName)
	ErrUpdateByIDSegments = errcode.NewError(segmentsBaseCode+3, "failed to update "+segmentsName)
	ErrGetByIDSegments    = errcode.NewError(segmentsBaseCode+4, "failed to get "+segmentsName+" details")
	ErrListSegments       = errcode.NewError(segmentsBaseCode+5, "failed to list of "+segmentsName)

	// error codes are globally unique, adding 1 to the previous error code
)
