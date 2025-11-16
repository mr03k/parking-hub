package consumers

import "github.com/google/wire"

var ConsumerProvider = wire.NewSet(Consumers)
