package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// roadCategories business-level http error codes.
// the roadCategoriesNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	roadCategoriesNO       = 61
	roadCategoriesName     = "roadCategories"
	roadCategoriesBaseCode = errcode.HCode(roadCategoriesNO)

	ErrCreateRoadCategories     = errcode.NewError(roadCategoriesBaseCode+1, "failed to create "+roadCategoriesName)
	ErrDeleteByIDRoadCategories = errcode.NewError(roadCategoriesBaseCode+2, "failed to delete "+roadCategoriesName)
	ErrUpdateByIDRoadCategories = errcode.NewError(roadCategoriesBaseCode+3, "failed to update "+roadCategoriesName)
	ErrGetByIDRoadCategories    = errcode.NewError(roadCategoriesBaseCode+4, "failed to get "+roadCategoriesName+" details")
	ErrListRoadCategories       = errcode.NewError(roadCategoriesBaseCode+5, "failed to list of "+roadCategoriesName)

	// error codes are globally unique, adding 1 to the previous error code
)
