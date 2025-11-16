package oidc

import (
	"context"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// Config defines the methods that the Config struct implements.
type Config interface {
	// AuthCodeURL generates an OAuth 2.0 authorization URL.
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string

	// PasswordCredentialsToken exchanges a username and password for an access token.
	PasswordCredentialsToken(ctx context.Context, username, password string) (*oauth2.Token, error)

	// Exchange converts an authorization code into an access token.
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)

	// Client returns an HTTP client configured with the provided token.
	Client(ctx context.Context, t *oauth2.Token) *http.Client

	// TokenSource returns a TokenSource that returns the provided token until it expires.
	TokenSource(ctx context.Context, t *oauth2.Token) oauth2.TokenSource
}

type TokenVerifier interface {
	Verify(ctx context.Context, rawIDToken string) (*oidc.IDToken, error)
}

// A TokenSource is anything that can return a token.
type TokenSource interface {
	// Token returns a token or an error.
	// Token must be safe for concurrent use by multiple goroutines.
	// The returned Token must not be modified.
	Token() (*oauth2.Token, error)
}
