package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateRoadCategoriesRequest request params
type CreateRoadCategoriesRequest struct {
	CodeCategoryRoad string `json:"codeCategoryRoad" binding:""`
	NameCategoryRoad string `json:"nameCategoryRoad" binding:""`
	Description      string `json:"description" binding:""`
}

// UpdateRoadCategoriesByIDRequest request params
type UpdateRoadCategoriesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	CodeCategoryRoad string `json:"codeCategoryRoad" binding:""`
	NameCategoryRoad string `json:"nameCategoryRoad" binding:""`
	Description      string `json:"description" binding:""`
}

// RoadCategoriesObjDetail detail
type RoadCategoriesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	CodeCategoryRoad string `json:"codeCategoryRoad"`
	NameCategoryRoad string `json:"nameCategoryRoad"`
	Description      string `json:"description"`
	CreatedAt        int    `json:"createdAt"`
}

// CreateRoadCategoriesReply only for api docs
type CreateRoadCategoriesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteRoadCategoriesByIDReply only for api docs
type DeleteRoadCategoriesByIDReply struct {
	Result
}

// UpdateRoadCategoriesByIDReply only for api docs
type UpdateRoadCategoriesByIDReply struct {
	Result
}

// GetRoadCategoriesByIDReply only for api docs
type GetRoadCategoriesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		RoadCategories RoadCategoriesObjDetail `json:"roadCategories"`
	} `json:"data"` // return data
}

// ListRoadCategoriessRequest request params
type ListRoadCategoriessRequest struct {
	query.Params
}

// ListRoadCategoriessReply only for api docs
type ListRoadCategoriessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		RoadCategoriess []RoadCategoriesObjDetail `json:"roadCategoriess"`
	} `json:"data"` // return data
}
