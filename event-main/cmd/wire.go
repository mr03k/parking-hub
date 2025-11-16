//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"git.abanppc.com/farin-project/event/config"
	"git.abanppc.com/farin-project/event/domain/biz"
	"git.abanppc.com/farin-project/event/domain/repository"
	"git.abanppc.com/farin-project/event/internal/rabbit/consumers"
	"git.abanppc.com/farin-project/event/pkg/mongo"
	"git.abanppc.com/farin-project/event/pkg/rabbit"
	"github.com/google/wire"
	"log/slog"
)

func wireApp(logger *slog.Logger, cfg *config.Config, rbt *rabbit.Rabbit, rbtConfig *config.RabbitMQ,
	mongo *mongo.Mongo,
) (*Bootstrap, error) {
	panic(wire.Build(
		consumers.ConsumerProvider,
		repository.ProviderSet,
		biz.ProviderSet,
		wire.NewSet(NewBootstrap)),
	)
}
