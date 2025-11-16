package rabbit

import (
	"git.abanppc.com/farin-project/state/app/rabbit/consumers"
	"git.abanppc.com/farin-project/state/app/rabbit/handlers"
	"github.com/google/wire"
)

func Consumers(ec *consumers.StateConsumer) []consumers.Consumer {
	return []consumers.Consumer{
		ec,
	}
}

var ProviderSet = wire.NewSet(
	handlers.NewState,
	consumers.NewEventConsumer,
	Consumers,
)
