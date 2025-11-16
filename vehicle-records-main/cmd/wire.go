//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	controller "git.abanppc.com/farin-project/vehicle-records/app/api/controllers"
	"git.abanppc.com/farin-project/vehicle-records/app/api/routes"
	rabbitApp "git.abanppc.com/farin-project/vehicle-records/app/rabbit"
	"git.abanppc.com/farin-project/vehicle-records/domain/opentelemetry"
	"git.abanppc.com/farin-project/vehicle-records/domain/repository"
	"git.abanppc.com/farin-project/vehicle-records/domain/service"
	"git.abanppc.com/farin-project/vehicle-records/infrastructure"
	gormdb "git.abanppc.com/farin-project/vehicle-records/infrastructure/gorm"
	"git.abanppc.com/farin-project/vehicle-records/infrastructure/rabbit"
	"github.com/google/wire"
	"github.com/mahdimehrabi/uploader/minio"
	"log/slog"
)

func wireApp(
	logger *slog.Logger,
	rabbit *rabbit.Rabbit,
	minio *minio.Minio,
	gorm *gormdb.GORMDB,
	ot *opentelemetry.OpenTelemetry,
) (*Boot, error) {
	panic(wire.Build(
		repository.ProviderSet,
		service.ProviderSet,
		infrastructure.ProviderSet,
		rabbitApp.ProviderSet,
		routes.ProviderSet,
		controller.ProviderSet,
		wire.NewSet(NewBoot),
	))
}
