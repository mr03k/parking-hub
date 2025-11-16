package biz

import (
	"git.abanppc.com/farin-project/event/domain/biz/event"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(event.NewEventBiz)
