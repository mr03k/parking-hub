package idpbiz

import "errors"

var (
	ErrorUserExist        = errors.New("user exist")
	ErrorUserNotFount     = errors.New("user not found")
	ErrorValidationFailed = errors.New("validation failed")
	ErrorPasswordEmpty    = errors.New("password empty")
	ErrorEmailMsisdnEmpty = errors.New("email and msisdn empty")
	ErrorPasswordNotMatch = errors.New("password not match")
	ErrorNotFound         = errors.New("not found")
)
