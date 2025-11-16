package repository

import (
	"git.abanppc.com/farin-project/event/domain/repository/v1/event"
	"git.abanppc.com/farin-project/event/domain/repository/v1/event/mongo"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	wire.Bind(new(event.EventRepoInterface), new(*mongo.EventRepository)),
	mongo.NewEventRepository)
