package biz

import (
	biz "application/internal/v1/biz/healthz"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(biz.NewHealthzBiz)
