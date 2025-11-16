//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"application/config"
	"application/internal/biz"
	"application/internal/datasource"
	"application/internal/http/handler"
	"application/internal/repo"
	"context"
	"log/slog"
	"net/http"

	rest_api "application/internal/http"

	"github.com/google/wire"
	"github.com/swaggest/openapi-go/openapi3"
)

func wireApp(ctx context.Context, cfg config.Config, logger *slog.Logger, oapi3r *openapi3.Reflector) (http.Handler, error) {
	panic(wire.Build(
		datasource.DataProviderSet,
		biz.BizProviderSet,
		rest_api.ServerProviderSet,
		handler.HandlerProviderSet,
		repo.RepoProviderSet,
	))
}
