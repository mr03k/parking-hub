package routes

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewHealthRouter, NewAuthRouter, CreateRouters, NewUserRouter)
