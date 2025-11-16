package handler

import "net/http"

type FuncHandler interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}
