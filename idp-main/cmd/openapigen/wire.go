//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"application/internal/http/oapi"
	"context"
	"log/slog"

	"github.com/google/wire"
	"github.com/swaggest/openapi-go/openapi3"
)

func wireApp(ctx context.Context, logger *slog.Logger, r *openapi3.Reflector) (*openapi3.Spec, error) {
	panic(wire.Build(
		// datasource.DataProviderSet,
		// biz.BizProviderSet,
		// rest_api.ServerProviderSet,
		// handler.HandlerProviderSet,
		oapi.OapiProviderSet,
	))
}
