//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"
	"log/slog"
	"net/http"

	"application/config"
	"application/internal/v1/biz"
	"application/internal/v1/datasource"
	internalHttp "application/internal/v1/http"
	"application/internal/v1/http/handler"
	oidcpkg "application/pkg/oidc"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/wire"
	"github.com/knadh/koanf/v2"
	"golang.org/x/oauth2"
)

func wireApp(ctx context.Context, cfg config.Config, logger *slog.Logger, k *koanf.Koanf,
	oauthConfig *oauth2.Config, oauthVerifier *oidc.IDTokenVerifier, security *config.Security,
	httpConfig *config.HTTPServer,
) (http.Handler, error) {
	panic(wire.Build(
		wire.Bind(new(oidcpkg.Config), &oauthConfig),
		wire.Bind(new(oidcpkg.TokenVerifier), &oauthVerifier),
		datasource.DataProviderSet,
		biz.ProviderSet,
		internalHttp.ServerProviderSet,
		handler.HandlerProviderSet,
	))
}
