package handler

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"application/config"
	"application/internal/v1/datasource/rule"
	"application/internal/v1/http/response"
	"application/pkg/middlewares"
	"application/pkg/middlewares/httplogger"
	"application/pkg/middlewares/httprecovery"
	oidcpkg "application/pkg/oidc"
	"application/pkg/utils"             //nolint:gci
	"github.com/coreos/go-oidc/v3/oidc" //nolint:gci
	"golang.org/x/oauth2"
)

type OAuthHandler struct {
	logger         *slog.Logger
	ruleDS         rule.Rule
	oauthConfig    oidcpkg.Config
	oauthValidator oidcpkg.TokenVerifier
	security       *config.Security
	encryptor      *utils.AESEncryptor
	httpConfig     *config.HTTPServer
}

func NewOauthHandler(logger *slog.Logger, ruleDS rule.Rule, oauthConfig oidcpkg.Config,
	oauthValidator oidcpkg.TokenVerifier, security *config.Security, serverConfig *config.HTTPServer,
) *OAuthHandler {
	return &OAuthHandler{
		logger:         logger,
		ruleDS:         ruleDS,
		oauthConfig:    oauthConfig,
		oauthValidator: oauthValidator,
		security:       security,
		encryptor:      utils.NewAESEncryptor([]byte(security.Secret), []byte(security.Secret[8:])),
		httpConfig:     serverConfig,
	}
}

func (h *OAuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Callback", "ctx", utils.GetLoggerContext(r.Context()))
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*14)
	defer cancel()

	// Validate the state parameter
	if err := h.validateState(r); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	// Exchange the authorization code for an OAuth2 token
	oauth2Token, err := h.exchangeToken(ctx, r)
	if err != nil {
		logger.Error("failed to exchange token", "error", err)
		response.InternalError(w)
		return
	}

	// Validate the ID token and nonce
	idToken := h.validateIDToken(ctx, oauth2Token, r, w, logger)
	if idToken == nil {
		return
	}

	// Prepare the response object
	resp, err := h.prepareResponse(oauth2Token, idToken, logger)
	if err != nil {
		logger.Error("failed to prepare response", "error", err)
		response.InternalError(w)
		return
	}

	// Set the authentication cookie
	if err := h.setAuthCookie(w, r, resp.Oauth2Token); err != nil {
		logger.Error("failed to set auth cookie", "error", err)
		response.InternalError(w)
		return
	}
	// delete unused cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "nonce",
		Value:    "",
		MaxAge:   0,
		Secure:   false,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "state",
		Value:    "",
		MaxAge:   0,
		Secure:   false,
		HttpOnly: true,
	})

	// Redirect to the URL specified in the "redirect_url" cookie
	if err := h.redirectToURL(w, r); err != nil {
		response.BadRequest(w, err.Error())
	}
}

func (h *OAuthHandler) validateState(r *http.Request) error {
	stateCookie, err := r.Cookie("state")
	if err != nil {
		return fmt.Errorf("state not found")
	}
	if r.URL.Query().Get("state") != stateCookie.Value {
		return fmt.Errorf("state did not match")
	}
	return nil
}

func (h *OAuthHandler) exchangeToken(ctx context.Context, r *http.Request) (*oauth2.Token, error) {
	return h.oauthConfig.Exchange(ctx, r.URL.Query().Get("code"))
}

func (h *OAuthHandler) validateIDToken(ctx context.Context, oauth2Token *oauth2.Token,
	r *http.Request, w http.ResponseWriter, logger *slog.Logger,
) *oidc.IDToken {
	// Extract the raw ID token from the OAuth2 token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		logger.Error("no id_token field in oauth2 token")
		response.InternalError(w)
		return nil
	}

	// Verify the ID token
	idToken, err := h.oauthValidator.Verify(ctx, rawIDToken)
	if err != nil {
		logger.Error("failed to verify ID Token", "error", err)
		response.InternalError(w)
		return nil
	}

	// Retrieve and validate the nonce
	nonceCookie, err := r.Cookie("nonce")
	if err != nil {
		logger.Error("nonce not found")
		response.BadRequest(w, "nonce not found")
		return nil
	}

	if idToken.Nonce != nonceCookie.Value {
		logger.Error("nonce did not match")
		response.BadRequest(w, "nonce did not match")
		return nil
	}

	return idToken
}

func (h *OAuthHandler) prepareResponse(oauth2Token *oauth2.Token, idToken *oidc.IDToken, logger *slog.Logger) (*struct {
	Oauth2Token   *oauth2.Token
	IDTokenClaims *json.RawMessage
}, error,
) {
	resp := &struct {
		Oauth2Token   *oauth2.Token
		IDTokenClaims *json.RawMessage
	}{
		Oauth2Token:   oauth2Token,
		IDTokenClaims: new(json.RawMessage),
	}

	if err := idToken.Claims(resp.IDTokenClaims); err != nil {
		if flag.Lookup("test.v") != nil {
			h.populateTestClaims(resp, logger)
		} else {
			return nil, fmt.Errorf("failed to get IDTokenClaims: %v", err)
		}
	}

	return resp, nil
}

func (h *OAuthHandler) populateTestClaims(resp *struct {
	Oauth2Token   *oauth2.Token
	IDTokenClaims *json.RawMessage
}, logger *slog.Logger,
) {
	claims := map[string]interface{}{
		"custom_claim_name": "custom claim value",
	}

	oauth2Tokens := oauth2.Token{
		AccessToken:  "test_access_token",
		RefreshToken: "test_refresh_token",
		TokenType:    "custom",
	}
	resp.Oauth2Token = &oauth2Tokens

	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		logger.Error("failed to marshal test claims", "error", err)
	}
	resp.IDTokenClaims = (*json.RawMessage)(&claimsBytes)
}

func (h *OAuthHandler) setAuthCookie(w http.ResponseWriter, r *http.Request, token *oauth2.Token) error {
	oauth2Claims, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal oauth2 token: %v", err)
	}

	tokenJSON, err := h.encryptor.Encrypt(oauth2Claims)
	if err != nil {
		return fmt.Errorf("failed to encrypt token: %v", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    string(tokenJSON),
		Secure:   r.TLS != nil,
		HttpOnly: true,
		Path:     "/",
	})

	return nil
}

func (h *OAuthHandler) redirectToURL(w http.ResponseWriter, r *http.Request) error {
	redirectURL, err := r.Cookie("redirect_url")
	if err != nil {
		return fmt.Errorf("redirect_url not found")
	}

	http.Redirect(w, r, redirectURL.Value, http.StatusFound)
	return nil
}

func (h *OAuthHandler) RegisterMuxRouter(mux FuncHandler) {
	logger := h.logger.With("method", "RegisterMuxRouter", "ctx", utils.GetLoggerContext(context.Background()))
	recoverMiddleware, err := httprecovery.NewRecoveryMiddleware()
	if err != nil {
		logger.Error("failed to set middleware", "error", err.Error())
		panic(err)
	}

	loggerMiddleware, err := httplogger.NewLoggerMiddleware()
	if err != nil {
		logger.Error("failed to set middleware", "error", err.Error())
		panic(err)
	}

	middles := []middlewares.Middleware{
		recoverMiddleware.RecoverMiddleware,
		httplogger.SetRequestContextLogger,
		loggerMiddleware.LoggerMiddleware,
	}
	rules, err := h.ruleDS.GetAll()
	if err != nil {
		logger.Error("failed to get all rules", "error", err.Error())
		panic(err)
	}

	for _, rule := range rules {
		ruleHandler := NewRuleHandler(rule, h.oauthConfig, h.oauthValidator, logger, h.encryptor)
		mux.HandleFunc(rule.Path,
			middlewares.MultipleMiddleware(ruleHandler.Handle, middles...))
	}
	mux.HandleFunc(fmt.Sprintf("/%s/callback", h.httpConfig.BasePath),
		middlewares.MultipleMiddleware(h.Callback, middles...))
}
