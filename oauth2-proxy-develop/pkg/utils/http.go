package utils

import (
	"net/http"
	"time"
)

type HTTPError struct {
	Message string
	Code    uint
	Error   string
}

func NewHTTPError(code uint, message string, err error) HTTPError {
	if err == nil {
		return HTTPError{
			Message: message,
			Code:    code,
			Error:   "",
		}
	}
	return HTTPError{
		Message: message,
		Code:    code,
		Error:   err.Error(),
	}
}

func SetCookie1Hour(w http.ResponseWriter, r *http.Request, name, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, c)
}
