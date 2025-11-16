package datasource

import (
	"application/internal/v1/datasource/healthz/memory"
	"application/internal/v1/datasource/rule/koanf"
	"github.com/google/wire"
)

var DataProviderSet = wire.NewSet(memory.NewHealthzDS, koanf.NewRuleDS)
