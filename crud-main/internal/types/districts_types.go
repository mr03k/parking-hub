package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateDistrictsRequest request params
type CreateDistrictsRequest struct {
	DistrictName string `json:"districtName" binding:""`
	CodeDistrict string `json:"codeDistrict" binding:""`
	IDCity       string `json:"iDCity" binding:""`
	BoundaryGeo  string `json:"boundaryGeo" binding:""`
	Population   int64  `json:"population" binding:""`
	Area         string `json:"area" binding:""`
}

// UpdateDistrictsByIDRequest request params
type UpdateDistrictsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	DistrictName string `json:"districtName" binding:""`
	CodeDistrict string `json:"codeDistrict" binding:""`
	IDCity       string `json:"iDCity" binding:""`
	BoundaryGeo  string `json:"boundaryGeo" binding:""`
	Population   int64  `json:"population" binding:""`
	Area         string `json:"area" binding:""`
}

// DistrictsObjDetail detail
type DistrictsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	DistrictName string `json:"districtName"`
	CodeDistrict string `json:"codeDistrict"`
	IDCity       string `json:"iDCity"`
	BoundaryGeo  string `json:"boundaryGeo"`
	Population   int64  `json:"population"`
	Area         string `json:"area"`
	CreatedAt    int    `json:"createdAt"`
}

// CreateDistrictsReply only for api docs
type CreateDistrictsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteDistrictsByIDReply only for api docs
type DeleteDistrictsByIDReply struct {
	Result
}

// UpdateDistrictsByIDReply only for api docs
type UpdateDistrictsByIDReply struct {
	Result
}

// GetDistrictsByIDReply only for api docs
type GetDistrictsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Districts DistrictsObjDetail `json:"districts"`
	} `json:"data"` // return data
}

// ListDistrictssRequest request params
type ListDistrictssRequest struct {
	query.Params
}

// ListDistrictssReply only for api docs
type ListDistrictssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Districtss []DistrictsObjDetail `json:"districtss"`
	} `json:"data"` // return data
}
