package handler

import (
	"github.com/google/wire"
)

var HandlerProviderSet = wire.NewSet(
	NewMuxHealthzHandler,
	NewOauthHandler,
	NewHandlerList,
)

// New ServiceList
func NewHandlerList(healthzHandler *HealthzHandler, oAuthHandler *OAuthHandler) []Handler {
	return []Handler{
		healthzHandler, oAuthHandler,
	}
}

// Service Interface
type Handler interface {
	RegisterMuxRouter(handlerFunc FuncHandler)
}
