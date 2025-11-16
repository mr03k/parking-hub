package datasource

import (
	"github.com/google/wire"
)

var DataProviderSet = wire.NewSet(NewDataSource, NewRingRepository, NewDriverRepository)
