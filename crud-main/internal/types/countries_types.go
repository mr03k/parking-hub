package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateCountriesRequest request params
type CreateCountriesRequest struct {
	CountryName string `json:"countryName" binding:""`
	CountryCode string `json:"countryCode" binding:""`
	IsoCode     string `json:"isoCode" binding:""`
	Region      string `json:"region" binding:""`
	Capital     string `json:"capital" binding:""`
	PhoneCode   string `json:"phoneCode" binding:""`
	Currency    string `json:"currency" binding:""`
	Population  int64  `json:"population" binding:""`
	Area        string `json:"area" binding:""`
	GeoBoundary string `json:"geoBoundary" binding:""`
}

// UpdateCountriesByIDRequest request params
type UpdateCountriesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	CountryName string `json:"countryName" binding:""`
	CountryCode string `json:"countryCode" binding:""`
	IsoCode     string `json:"isoCode" binding:""`
	Region      string `json:"region" binding:""`
	Capital     string `json:"capital" binding:""`
	PhoneCode   string `json:"phoneCode" binding:""`
	Currency    string `json:"currency" binding:""`
	Population  int64  `json:"population" binding:""`
	Area        string `json:"area" binding:""`
	GeoBoundary string `json:"geoBoundary" binding:""`
}

// CountriesObjDetail detail
type CountriesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	CountryName string `json:"countryName"`
	CountryCode string `json:"countryCode"`
	IsoCode     string `json:"isoCode"`
	Region      string `json:"region"`
	Capital     string `json:"capital"`
	PhoneCode   string `json:"phoneCode"`
	Currency    string `json:"currency"`
	Population  int64  `json:"population"`
	Area        string `json:"area"`
	GeoBoundary string `json:"geoBoundary"`
	CreatedAt   int    `json:"createdAt"`
}

// CreateCountriesReply only for api docs
type CreateCountriesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteCountriesByIDReply only for api docs
type DeleteCountriesByIDReply struct {
	Result
}

// UpdateCountriesByIDReply only for api docs
type UpdateCountriesByIDReply struct {
	Result
}

// GetCountriesByIDReply only for api docs
type GetCountriesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Countries CountriesObjDetail `json:"countries"`
	} `json:"data"` // return data
}

// ListCountriessRequest request params
type ListCountriessRequest struct {
	query.Params
}

// ListCountriessReply only for api docs
type ListCountriessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Countriess []CountriesObjDetail `json:"countriess"`
	} `json:"data"` // return data
}
