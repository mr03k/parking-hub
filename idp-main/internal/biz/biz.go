package biz

import (
	authbiz "application/internal/biz/auth"
	"application/internal/biz/device"
	healthzusecase "application/internal/biz/healthz"
	idpbiz "application/internal/biz/idp"

	"github.com/google/wire"
)

var BizProviderSet = wire.NewSet(healthzusecase.NewHealthzBiz, idpbiz.NewUserUseCase,
	authbiz.NewAuthUsecase, device.NewVehicleService)
