//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	controller "farin/app/api/controllers"
	"farin/app/api/middleware"
	"farin/app/api/routes"
	"farin/app/api/validators"
	rabbitApp "farin/app/rabbit"
	"farin/domain/repository"
	"farin/domain/service"
	"farin/infrastructure"
	gormdb "farin/infrastructure/gorm"
	"farin/infrastructure/opentelemetry"
	"farin/infrastructure/rabbit"
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
		controller.ProviderSet,
		middleware.ProviderSet,
		routes.ProviderSet,
		repository.ProviderSet,
		service.ProviderSet,
		infrastructure.ProviderSet,
		validators.ProviderSet,
		validators.NewValidators,
		rabbitApp.ProviderSet,
		wire.NewSet(NewBoot),
	))
}
