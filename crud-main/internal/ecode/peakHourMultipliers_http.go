package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// peakHourMultipliers business-level http error codes.
// the peakHourMultipliersNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	peakHourMultipliersNO       = 85
	peakHourMultipliersName     = "peakHourMultipliers"
	peakHourMultipliersBaseCode = errcode.HCode(peakHourMultipliersNO)

	ErrCreatePeakHourMultipliers     = errcode.NewError(peakHourMultipliersBaseCode+1, "failed to create "+peakHourMultipliersName)
	ErrDeleteByIDPeakHourMultipliers = errcode.NewError(peakHourMultipliersBaseCode+2, "failed to delete "+peakHourMultipliersName)
	ErrUpdateByIDPeakHourMultipliers = errcode.NewError(peakHourMultipliersBaseCode+3, "failed to update "+peakHourMultipliersName)
	ErrGetByIDPeakHourMultipliers    = errcode.NewError(peakHourMultipliersBaseCode+4, "failed to get "+peakHourMultipliersName+" details")
	ErrListPeakHourMultipliers       = errcode.NewError(peakHourMultipliersBaseCode+5, "failed to list of "+peakHourMultipliersName)

	// error codes are globally unique, adding 1 to the previous error code
)
