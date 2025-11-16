//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	controller "git.abanppc.com/farin-project/state/app/api/controllers"
	"git.abanppc.com/farin-project/state/app/api/routes"
	rbtApp "git.abanppc.com/farin-project/state/app/rabbit"
	"git.abanppc.com/farin-project/state/domain/repository"
	"git.abanppc.com/farin-project/state/domain/service"
	"git.abanppc.com/farin-project/state/infrastructure"
	gormdb "git.abanppc.com/farin-project/state/infrastructure/gorm"
	"git.abanppc.com/farin-project/state/infrastructure/rabbit"
	"github.com/google/wire"
	"log/slog"
)

func wireApp(
	logger *slog.Logger,
	rabbit *rabbit.Rabbit,
	gorm *gormdb.GORMDB,
) (*Boot, error) {
	panic(wire.Build(
		repository.ProviderSet,
		service.ProviderSet,
		infrastructure.ProviderSet,
		rbtApp.ProviderSet,
		routes.ProviderSet,
		controller.ProviderSet,
		wire.NewSet(NewBoot),
	))
}
