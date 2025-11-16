package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateCitiesRequest request params
type CreateCitiesRequest struct {
	CityName  string `json:"cityName" binding:""`
	CodeCity  string `json:"codeCity" binding:""`
	IDCountry string `json:"iDCountry" binding:""`
	Boundary  string `json:"boundary" binding:""`
}

// UpdateCitiesByIDRequest request params
type UpdateCitiesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	CityName  string `json:"cityName" binding:""`
	CodeCity  string `json:"codeCity" binding:""`
	IDCountry string `json:"iDCountry" binding:""`
	Boundary  string `json:"boundary" binding:""`
}

// CitiesObjDetail detail
type CitiesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	CityName  string `json:"cityName"`
	CodeCity  string `json:"codeCity"`
	IDCountry string `json:"iDCountry"`
	Boundary  string `json:"boundary"`
	CreatedAt int    `json:"createdAt"`
}

// CreateCitiesReply only for api docs
type CreateCitiesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteCitiesByIDReply only for api docs
type DeleteCitiesByIDReply struct {
	Result
}

// UpdateCitiesByIDReply only for api docs
type UpdateCitiesByIDReply struct {
	Result
}

// GetCitiesByIDReply only for api docs
type GetCitiesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Cities CitiesObjDetail `json:"cities"`
	} `json:"data"` // return data
}

// ListCitiessRequest request params
type ListCitiessRequest struct {
	query.Params
}

// ListCitiessReply only for api docs
type ListCitiessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Citiess []CitiesObjDetail `json:"citiess"`
	} `json:"data"` // return data
}
