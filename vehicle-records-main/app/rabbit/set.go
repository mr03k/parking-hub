package rabbit

import (
	"git.abanppc.com/farin-project/vehicle-records/app/rabbit/consumers"
	"git.abanppc.com/farin-project/vehicle-records/app/rabbit/handlers"
	"github.com/google/wire"
)

func Consumers(ec *consumers.EventConsumer, rc *consumers.RetryConsumer) []consumers.Consumer {
	return []consumers.Consumer{
		ec, rc,
	}
}

var ProviderSet = wire.NewSet(
	handlers.NewEventVehicleRecord,
	consumers.NewEventConsumer,
	consumers.NewRetryConsumer,
	handlers.NewRetry,
	Consumers,
)
