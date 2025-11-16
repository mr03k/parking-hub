package rabbit

import (
	"farin/app/rabbit/consumers"
	"farin/app/rabbit/handlers"
	"github.com/google/wire"
)

func Consumers(ec *consumers.EventConsumer) []consumers.Consumer {
	return []consumers.Consumer{
		ec,
	}
}

var ProviderSet = wire.NewSet(
	handlers.NewVehicleRecord,
	consumers.NewEventConsumer,
	Consumers,
)
