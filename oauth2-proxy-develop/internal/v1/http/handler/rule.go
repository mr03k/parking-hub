package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	entity "application/internal/v1/entity/rule"
	"application/internal/v1/http/response"
	oidcpkg "application/pkg/oidc"
	"application/pkg/utils"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

// create handler function
// handle response base on rule
type RuleHandler struct {
	*entity.Rule
	actionFunc    http.HandlerFunc
	oauthConfig   oidcpkg.Config
	oauthVerifier oidcpkg.TokenVerifier
	logger        *slog.Logger
	encryptor     *utils.AESEncryptor
}

func NewRuleHandler(rule *entity.Rule, oauthConfig oidcpkg.Config,
	verifier oidcpkg.TokenVerifier, logger *slog.Logger, encryptor *utils.AESEncryptor,
) *RuleHandler {
	r := &RuleHandler{
		Rule:          rule,
		oauthConfig:   oauthConfig,
		oauthVerifier: verifier,
		logger:        logger.With("layer", "RuleHandler"),
		encryptor:     encryptor,
	}
	switch rule.Action {
	case entity.ActionAuth:
		r.actionFunc = r.actionAuth
	case entity.ActionAllow:
		r.actionFunc = r.actionAllow
	default:
		r.actionFunc = r.actionDeny
	}
	return r
}

func (h *RuleHandler) Handle(rw http.ResponseWriter, req *http.Request) {
	h.actionFunc(rw, req)
}

func (h *RuleHandler) actionDeny(rw http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	logger := h.logger.With("method", "Callback", "ctx", utils.GetLoggerContext(ctx))
	defer cancel()

	cookie, err := req.Cookie("auth")
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		h.logger.Error("get cookie err", "err", err)
		return
	}
	if cookie != nil {
		h.handleAuthCookie(ctx, rw, req, cookie, logger)
		return
	}

	response.Custom(rw, http.StatusForbidden, nil, "access denied")
}

func (h *RuleHandler) actionAllow(rw http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	logger := h.logger.With("method", "Callback", "ctx", utils.GetLoggerContext(ctx))
	defer cancel()

	cookie, err := req.Cookie("auth")
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		h.logger.Error("get cookie err", "err", err)
		return
	}
	if cookie != nil {
		h.handleAuthCookie(ctx, rw, req, cookie, logger)
		return
	}

	response.Ok(rw, nil, "access granted")
}

func (h *RuleHandler) actionAuth(rw http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	logger := h.logger.With("method", "Callback", "ctx", utils.GetLoggerContext(ctx))
	defer cancel()

	cookie, err := req.Cookie("auth")
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		h.logger.Error("get cookie err", "err", err)
		return
	}
	if cookie != nil {
		h.handleAuthCookie(ctx, rw, req, cookie, logger)
		return
	}
	state := uuid.NewString()
	utils.SetCookie1Hour(rw, req, "state", state)
	nonce := uuid.NewString()
	utils.SetCookie1Hour(rw, req, "nonce", nonce)
	utils.SetCookie1Hour(rw, req, "redirect_url", req.RequestURI)

	http.Redirect(rw, req, h.oauthConfig.AuthCodeURL(state, oidc.Nonce(nonce)), http.StatusFound)
}

func (h *RuleHandler) handleAuthCookie(ctx context.Context, rw http.ResponseWriter,
	r *http.Request, cookie *http.Cookie, logger *slog.Logger,
) {
	tokenBytes, err := h.encryptor.Decrypt([]byte(cookie.Value))
	if err != nil {
		h.logger.Error("decrypt cookie err", "err", err)
		h.handleAuthCookieError(rw)
		return
	}

	oauthToken := &oauth2.Token{}
	if err := json.Unmarshal(tokenBytes, &oauthToken); err != nil {
		h.logger.Error("unmarshal cookie err", "err", err)
		h.handleAuthCookieError(rw)
		return
	}

	// Refresh token if access token is expired
	oauthToken, err = h.renewAccessToken(ctx, rw, r, oauthToken)
	if err != nil {
		logger.Warn("error in renewing token", "err", err)
		h.handleAuthCookieError(rw)
		return
	}

	_, err = h.oauthVerifier.Verify(ctx, oauthToken.AccessToken)
	if err != nil {
		logger.Error("verify access err", "err", err)
		h.handleAuthCookieError(rw)
		return
	}

	accessArr := strings.Split(oauthToken.AccessToken, ".")
	if len(accessArr) != 3 {
		logger.Warn("access token is not in JWT format")
		h.handleAuthCookieError(rw)
		return
	}

	claimsJSON, err := utils.DecodeBase64([]byte(accessArr[1]))
	if err != nil {
		logger.Warn("decode access token err", "err", err)
		h.handleAuthCookieError(rw)
		return
	}

	// Handle errors from claimsToHeaders
	if err := h.claimsToHeaders(rw, claimsJSON, oauthToken); err != nil {
		logger.Warn("error setting headers from claims", "err", err)
		h.handleAuthCookieError(rw)
		return
	}

	response.Ok(rw, nil, "")
}

func (h *RuleHandler) claimsToHeaders(rw http.ResponseWriter, claimsJSON []byte, oauthToken *oauth2.Token) error {
	claims := make(map[string]interface{})
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return fmt.Errorf("unmarshal access token err: %w", err)
	}

	rw.Header().Add("X-AUTH-ACCESS-TOKEN", oauthToken.AccessToken)

	if err := h.strClaimToHeader(rw, claims, "sub", "X-AUTH-SUB"); err != nil {
		return fmt.Errorf("access token does not contain specific claim 'sub': %w", err)
	}

	if err := h.strClaimToHeader(rw, claims, "name", "X-AUTH-NAME"); err != nil {
		return fmt.Errorf("access token does not contain specific claim 'name': %w", err)
	}

	// Email is optional, so no need to return an error if it's not present
	h.strClaimToHeader(rw, claims, "email", "X-AUTH-EMAIL") //nolint:all

	if err := h.boolClaimToHeader(rw, claims, "email_verified", "X-AUTH-VERIFIED"); err != nil {
		return fmt.Errorf("access token does not contain specific claim 'email_verified': %w", err)
	}

	if err := h.setUserIDHeader(claims, rw); err != nil {
		return fmt.Errorf("error extracting user id and setting header: %w", err)
	}

	if err := h.arrClaimToHeader(rw, claims, "groups", "X-AUTH-GROUPS"); err != nil {
		return fmt.Errorf("access token does not contain specific claim 'groups': %w", err)
	}

	return nil
}

func (h *RuleHandler) renewAccessToken(ctx context.Context, rw http.ResponseWriter, r *http.Request,
	oauthToken *oauth2.Token,
) (*oauth2.Token, error) {
	now := time.Now().Add(-3 * time.Second).Unix()
	exp := oauthToken.Expiry.Unix()
	if exp < now {
		ts := h.oauthConfig.TokenSource(ctx, &oauth2.Token{
			RefreshToken: oauthToken.RefreshToken,
		})
		token, err := ts.Token()
		if err != nil {
			return nil, err
		}
		oauthToken = token
		if err := h.setAuthCookie(rw, r, oauthToken); err != nil {
			return nil, err
		}
	}
	return oauthToken, nil
}

func (h *RuleHandler) setAuthCookie(w http.ResponseWriter, r *http.Request, token *oauth2.Token) error {
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

func (h *RuleHandler) handleAuthCookieError(rw http.ResponseWriter) {
	// delete auth cookie
	http.SetCookie(rw, &http.Cookie{
		Name:     "auth",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
	if h.Action == entity.ActionAllow {
		response.Ok(rw, nil, "access granted")
	}
	response.Custom(rw, http.StatusForbidden, nil, "access denied")
}

func (h *RuleHandler) strClaimToHeader(rw http.ResponseWriter, claims map[string]interface{},
	claimName, headerName string,
) error {
	claim, ok := claims[claimName]
	if !ok {
		return fmt.Errorf("claim %s not found", claimName)
	}
	claimString, ok := claim.(string)
	if !ok {
		return fmt.Errorf("claim %s is not a string", claimName)
	}
	rw.Header().Add(headerName, claimString)
	return nil
}

func (h *RuleHandler) boolClaimToHeader(rw http.ResponseWriter, claims map[string]interface{},
	claimName, headerName string,
) error {
	claim, ok := claims[claimName]
	if !ok {
		return fmt.Errorf("claim %s not found", claimName)
	}
	// Ensure claim is a boolean
	boolClaim, ok := claim.(bool)
	if !ok {
		return fmt.Errorf("claim %s is not a boolean", claimName)
	}
	rw.Header().Add(headerName, strconv.FormatBool(boolClaim))
	return nil
}

func (h *RuleHandler) arrClaimToHeader(rw http.ResponseWriter, claims map[string]interface{},
	claimName, headerName string,
) error {
	claim, ok := claims[claimName]
	if !ok {
		return fmt.Errorf("claim %s not found", claimName)
	}
	claimString, ok := claim.([]interface{})
	if !ok {
		return fmt.Errorf("claim %s is not a []string", claimName)
	}
	claimStrings := make([]string, len(claimString))
	for i, v := range claimString {
		claimStrings[i] = fmt.Sprintf("%v", v)
	}
	rw.Header().Add(headerName, strings.Join(claimStrings, ","))
	return nil
}

func (h *RuleHandler) setUserIDHeader(claims map[string]interface{}, rw http.ResponseWriter) error {
	federatedClaims, ok := claims["federated_claims"]
	if !ok {
		return fmt.Errorf("federated claims not found")
	}

	federatedClaimsMap, ok := federatedClaims.(map[string]interface{})
	if !ok {
		return fmt.Errorf("federated claims are not a map")
	}

	userID, ok := federatedClaimsMap["user_id"]
	if !ok {
		return fmt.Errorf("user id not found in federated claims")
	}

	userIDStr := userID.(string)
	rw.Header().Add("X-AUTH-USER-ID", userIDStr)
	return nil
}
