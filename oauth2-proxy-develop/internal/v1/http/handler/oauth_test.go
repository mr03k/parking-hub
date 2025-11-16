package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"application/config"
	"application/internal/v1/entity/rule"
	"application/internal/v1/http/response"
	mockRule "application/mock/datasource"
	mockHandler "application/mock/handler"
	mockOidc "application/mock/pkg/oidc"
	"application/pkg/utils"
	"github.com/coreos/go-oidc/v3/oidc"
	"go.uber.org/mock/gomock"
	"golang.org/x/oauth2"
)

func TestOAuthHandler_RegisterMuxRouter(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})
	rules := []*rule.Rule{
		{
			Name:   "auth_exact",
			Action: rule.ActionAuth,
			Path:   "/auth/exact",
		},
		{
			Name:   "allow_prefix",
			Action: rule.ActionAllow,
			Path:   "/allow",
		},
		{
			Name:   "deny_exact",
			Action: rule.ActionDeny,
			Path:   "/deny/exact",
		},
		{
			Name:   "auth_prefix",
			Action: rule.ActionAuth,
			Path:   "/auth",
		},
		{
			Name:   "allow_exact",
			Action: rule.ActionAllow,
			Path:   "/allow/exact",
		},
		{
			Name:   "deny_prefix",
			Action: rule.ActionDeny,
			Path:   "/deny",
		},
	}
	tests := []struct {
		name      string
		rulesRepo func() *mockRule.MockRule
		ctx       context.Context
		mux       func() *mockHandler.MockFuncHandler
	}{
		{
			name: "success",
			ctx:  context.Background(),
			rulesRepo: func() *mockRule.MockRule {
				r := mockRule.NewMockRule(ctrl)
				r.EXPECT().GetAll().Return(rules, nil)
				return r
			},
			mux: func() *mockHandler.MockFuncHandler {
				mux := mockHandler.NewMockFuncHandler(ctrl)
				for _, r := range rules {
					mux.EXPECT().HandleFunc(r.Path, gomock.Any()).Times(1)
				}
				mux.EXPECT().HandleFunc("//callback", gomock.Any()).Times(1)
				return mux
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(_ *testing.T) {
			rulesRepo := test.rulesRepo()
			handler := NewOauthHandler(slog.New(slog.NewTextHandler(os.Stdout, nil)), rulesRepo,
				mockOidc.NewMockConfig(ctrl), mockOidc.NewMockTokenVerifier(ctrl), &config.Security{
					Secret: "W+c%1*~(E)^2^#sxx^^7%b34",
				}, &config.HTTPServer{})
			mux := test.mux()
			handler.RegisterMuxRouter(mux)
			rulesRepo.EXPECT()
			mux.EXPECT()
		})
	}
}

func TestOAuthHandler_Callback(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})
	encryptor := utils.NewAESEncryptor([]byte("W+c%1*~(E)^2^#sxx^^7%b34"), []byte("E)^2^#sxx^^7%b34"))

	tests := []struct {
		name               string
		request            func() *http.Request
		expectedStatusCode int
		ctx                context.Context
		expectedResponse   func() http.ResponseWriter
		oauthConfig        func() *mockOidc.MockConfig
		oauthVerifier      func() *mockOidc.MockTokenVerifier
		security           *config.Security
		noDeadLine         bool
	}{
		{
			name: "state_not_found",
			request: func() *http.Request {
				r := bytes.NewReader([]byte("test data"))
				req := httptest.NewRequest(http.MethodGet, `/admin?state=state_fake`, r)
				req.Header.Add("test", "test value")
				req.AddCookie(&http.Cookie{
					Name:     "nonce",
					Value:    "nonce_fake",
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
				})
				req.AddCookie(&http.Cookie{
					Name:     "redirect_url",
					Value:    "/admin",
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
				})
				return req
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: func() http.ResponseWriter {
				recorder := httptest.NewRecorder()
				response.BadRequest(recorder, "state not found")
				return recorder
			},
			oauthConfig: func() *mockOidc.MockConfig {
				config := mockOidc.NewMockConfig(ctrl)
				return config
			},
			oauthVerifier: func() *mockOidc.MockTokenVerifier {
				mv := mockOidc.NewMockTokenVerifier(ctrl)
				return mv
			},
			security: &config.Security{
				Secret: "W+c%1*~(E)^2^#sxx^^7%b34",
			},
		},
		{
			name: "exchange error",
			request: func() *http.Request {
				r := bytes.NewReader([]byte("test data"))
				req := httptest.NewRequest(http.MethodGet, `/admin?state=state_fake`, r)
				req.Header.Add("test", "test value")
				req.AddCookie(&http.Cookie{
					Name:     "nonce",
					Value:    "nonce_fake",
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
				})
				req.AddCookie(&http.Cookie{
					Name:     "state",
					Value:    "state_fake",
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
				})
				req.AddCookie(&http.Cookie{
					Name:     "redirect_url",
					Value:    "/admin",
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
				})
				return req
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: func() http.ResponseWriter {
				recorder := httptest.NewRecorder()
				response.InternalError(recorder)
				return recorder
			},
			oauthConfig: func() *mockOidc.MockConfig {
				config := mockOidc.NewMockConfig(ctrl)
				config.EXPECT().Exchange(gomock.Any(), gomock.Any()).DoAndReturn(func(any, any, ...any) (*oauth2.Token, error) {
					return nil, errors.New("error")
				})
				return config
			},
			oauthVerifier: func() *mockOidc.MockTokenVerifier {
				mv := mockOidc.NewMockTokenVerifier(ctrl)
				return mv
			},
			security: &config.Security{
				Secret: "W+c%1*~(E)^2^#sxx^^7%b34",
			},
		},
		{
			name: "nonce_not_found",
			request: func() *http.Request {
				r := bytes.NewReader([]byte("test data"))
				req := httptest.NewRequest(http.MethodGet, `/admin?state=state_fake`, r)
				req.Header.Add("test", "test value")
				req.AddCookie(&http.Cookie{
					Name:     "state",
					Value:    "state_fake",
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
				})
				req.AddCookie(&http.Cookie{
					Name:     "redirect_url",
					Value:    "/admin",
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
				})
				return req
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: func() http.ResponseWriter {
				recorder := httptest.NewRecorder()
				response.BadRequest(recorder, "nonce did not match")
				return recorder
			},
			oauthConfig: func() *mockOidc.MockConfig {
				config := mockOidc.NewMockConfig(ctrl)
				extra := map[string]interface{}{
					"access_token":  "eyJhbGciOiJSUzI1NiIsImtpZCI6IjViYjI2Y2U5Yjg5NGViYTBkNzhjZjQ3Y2YyZjc1YWY4ODk5ZjE3MTcifQ.eyJpc3MiOiJodHRwczovL2FwcC5ndy5rOHMuc2guYWJhbi5pby9kZXgiLCJzdWIiOiJDZzB3TFRNNE5TMHlPREE0T1Mwd0VnUnRiMk5yIiwiYXVkIjoiZXhhbXBsZS1hcHAiLCJleHAiOjE3MjQzMzU2ODgsImlhdCI6MTcyNDI0OTI4OCwibm9uY2UiOiJmYTA5MjZhNy1hYWJjLTRiOTYtYWNlNy01OGY4NTNkYWMxYWQiLCJhdF9oYXNoIjoiNEpqWWxlREhRMmNLcXFmbWsyRUczZyIsImVtYWlsIjoia2lsZ29yZUBraWxnb3JlLnRyb3V0IiwiZW1haWxfdmVyaWZpZWQiOnRydWUsIm5hbWUiOiJLaWxnb3JlIFRyb3V0In0.B7n3v40XmdfMjMW2UB1NmFLZQyyRORwOPSknF3FcOD1juoJmql-1T_NvuHpTCmFqGFDtBKvZa9_RKoisC2LxVvJoUrlDUzDxUafdu3gL9uIwWAtvfAhAjcOujgOupdcwDL-omfoghV83fHdmUr20vOsitLbiKWzWSjzsMRXT4c6OWiQetTrJrQxWQbqFBm-7iqJjXazzDbmhxBpXg6s_BurnJZ0MjXDgSDWDcUxp4-dgvyd_gxmJS_MgnWEgZAaV_y7qG0CHwLyLRl9FHMbSmHUaIRy8KHbUnp9k0-RncTlEM90E09l7kcEbX-74IsqhiyowqPHFHcS9mAZeNlgeaQ", //nolint:lll
					"token_type":    "bearer",
					"expires_in":    "86399",
					"refresh_token": "ChlyZnlsa2k0eXVscHZzaWs2dGZjYWF5YTRkEhl1cnUzN2VpczZsMzVrZjd6dTJvaHQ3YzZq",
					"id_token":      "eyJhbGciOiJSUzI1NiIsImtpZCI6IjViYjI2Y2U5Yjg5NGViYTBkNzhjZjQ3Y2YyZjc1YWY4ODk5ZjE3MTcifQ.eyJpc3MiOiJodHRwczovL2FwcC5ndy5rOHMuc2guYWJhbi5pby9kZXgiLCJzdWIiOiJDZzB3TFRNNE5TMHlPREE0T1Mwd0VnUnRiMk5yIiwiYXVkIjoiZXhhbXBsZS1hcHAiLCJleHAiOjE3MjQzMzU2ODgsImlhdCI6MTcyNDI0OTI4OCwibm9uY2UiOiJmYTA5MjZhNy1hYWJjLTRiOTYtYWNlNy01OGY4NTNkYWMxYWQiLCJhdF9oYXNoIjoiSC12VDNKTGd1R3BPVHczTm9qM3FCdyIsImNfaGFzaCI6Ii1OMHRJTllsWWt0ZlNaRnp0MFVuUWciLCJlbWFpbCI6ImtpbGdvcmVAa2lsZ29yZS50cm91dCIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiS2lsZ29yZSBUcm91dCJ9.QH9mbZrWX-P7DBfR4V93DhTSuEo8PWWlT7J-YTbZqcIberfkyucYyNoCJynKIGWLNb9yhxtbl4Ar5YhMnyOjzYOJXtPnNy2nTDp1SoOTEtfAb7FBcY1Gq9c-hl7e1nX5PaJlLM57KRhotWt4et8B1VY2e4MgYcj3touMGUV69kHZ0nv18KV7hUkAtsJ3OEQcVcw94wccMPFh7oOJnDIGoVxLAwXNkk-3h3PG0U_G7FoHwqewOfXiMKABSCZ3jH4P7gjTjFQ919iUNcMErp5q3xgcSJ61YDx9blv4AvE2SKq9oDZVQy7EYRPekT2NTrqJFmQU5teKA1SajzEvItBGyw", //nolint:lll
				}
				token := &oauth2.Token{
					AccessToken:  "eyJhbGciOiJSUzI1NiIsImtpZCI6IjViYjI2Y2U5Yjg5NGViYTBkNzhjZjQ3Y2YyZjc1YWY4ODk5ZjE3MTcifQ.eyJpc3MiOiJodHRwczovL2FwcC5ndy5rOHMuc2guYWJhbi5pby9kZXgiLCJzdWIiOiJDZzB3TFRNNE5TMHlPREE0T1Mwd0VnUnRiMk5yIiwiYXVkIjoiZXhhbXBsZS1hcHAiLCJleHAiOjE3MjQzMzUwOTQsImlhdCI6MTcyNDI0ODY5NCwibm9uY2UiOiIxYzFmOTU2My1iNWQwLTQwOWItOGRmNS1hYjhhYjdkYzE0ODgiLCJhdF9oYXNoIjoieU5qVFJTTWhZLXFKemJJRTlybG96ZyIsImVtYWlsIjoia2lsZ29yZUBraWxnb3JlLnRyb3V0IiwiZW1haWxfdmVyaWZpZWQiOnRydWUsIm5hbWUiOiJLaWxnb3JlIFRyb3V0In0.k-ajsvzlJfuwGkBy9lW8OgRtlzH0PXpU2HgrU4WaVssEdTc3pQWtw50Aqy5s-HWKUVxMwvLgU1qznTkakglNFAic-vRQzPf0cGQBwW7ojTaz-V9oktUhYIj5BXS001BvLqYGEWkKXS07IL2FYh0WEm7m6xDeOsN_5gwUf3NDpXyIiVfVWXEBRiKjsHK6nVA7YggzKF6K4PwK0BdzSy58AT0Q33zdoOPztb2fWvfjBQWZZFSog3Z0rvr0rurIWWlhSGveZcqkwZj9FeZZlIkGK1VDiQhmkha-lzeomgkh8X07Z8wxXS0tYVJ3S9E7bxMOkb7T-P4LzOo3Oj9GNxP7Zw", //nolint:lll
					RefreshToken: "Chl4NjZxa2V4emJhMzJkNmdwZzZ3Y2tycHpoEhlydHp5dTc1d2ZxM29ncGxpcHZwYmtudTRy",
					TokenType:    "bearer",
					Expiry:       time.Now().Add(1 * time.Hour),
				}
				config.EXPECT().Exchange(gomock.Any(), gomock.Any()).DoAndReturn(func(any, any, ...any) (*oauth2.Token, error) {
					return token.WithExtra(extra), nil
				})
				return config
			},
			oauthVerifier: func() *mockOidc.MockTokenVerifier {
				mv := mockOidc.NewMockTokenVerifier(ctrl)
				mv.EXPECT().Verify(gomock.Any(), gomock.Any()).Return(&oidc.IDToken{
					Nonce: "nonce_fake",
				}, nil)
				return mv
			},
			security: &config.Security{
				Secret: "W+c%1*~(E)^2^#sxx^^7%b34",
			},
		},
		{
			name: "wrong_idtoken",
			request: func() *http.Request {
				r := bytes.NewReader([]byte("test data"))
				req := httptest.NewRequest(http.MethodGet, `/admin?state=state_fake`, r)
				req.Header.Add("test", "test value")
				req.AddCookie(&http.Cookie{
					Name:     "nonce",
					Value:    "nonce_fake",
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
				})
				req.AddCookie(&http.Cookie{
					Name:     "state",
					Value:    "state_fake",
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
				})
				req.AddCookie(&http.Cookie{
					Name:     "redirect_url",
					Value:    "/admin",
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
				})
				return req
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: func() http.ResponseWriter {
				recorder := httptest.NewRecorder()
				response.BadRequest(recorder, "nonce did not match")
				return recorder
			},
			oauthConfig: func() *mockOidc.MockConfig {
				config := mockOidc.NewMockConfig(ctrl)
				extra := map[string]interface{}{
					"access_token":  "eyJhbGciOiJSUzI1NiIsImtpZCI6IjViYjI2Y2U5Yjg5NGViYTBkNzhjZjQ3Y2YyZjc1YWY4ODk5ZjE3MTcifQ.eyJpc3MiOiJodHRwczovL2FwcC5ndy5rOHMuc2guYWJhbi5pby9kZXgiLCJzdWIiOiJDZzB3TFRNNE5TMHlPREE0T1Mwd0VnUnRiMk5yIiwiYXVkIjoiZXhhbXBsZS1hcHAiLCJleHAiOjE3MjQzMzU2ODgsImlhdCI6MTcyNDI0OTI4OCwibm9uY2UiOiJmYTA5MjZhNy1hYWJjLTRiOTYtYWNlNy01OGY4NTNkYWMxYWQiLCJhdF9oYXNoIjoiNEpqWWxlREhRMmNLcXFmbWsyRUczZyIsImVtYWlsIjoia2lsZ29yZUBraWxnb3JlLnRyb3V0IiwiZW1haWxfdmVyaWZpZWQiOnRydWUsIm5hbWUiOiJLaWxnb3JlIFRyb3V0In0.B7n3v40XmdfMjMW2UB1NmFLZQyyRORwOPSknF3FcOD1juoJmql-1T_NvuHpTCmFqGFDtBKvZa9_RKoisC2LxVvJoUrlDUzDxUafdu3gL9uIwWAtvfAhAjcOujgOupdcwDL-omfoghV83fHdmUr20vOsitLbiKWzWSjzsMRXT4c6OWiQetTrJrQxWQbqFBm-7iqJjXazzDbmhxBpXg6s_BurnJZ0MjXDgSDWDcUxp4-dgvyd_gxmJS_MgnWEgZAaV_y7qG0CHwLyLRl9FHMbSmHUaIRy8KHbUnp9k0-RncTlEM90E09l7kcEbX-74IsqhiyowqPHFHcS9mAZeNlgeaQ", //nolint:lll
					"token_type":    "bearer",
					"expires_in":    "86399",
					"refresh_token": "ChlyZnlsa2k0eXVscHZzaWs2dGZjYWF5YTRkEhl1cnUzN2VpczZsMzVrZjd6dTJvaHQ3YzZq",
					"id_token":      "eyJhbGciOiJSUzI1NiIsImtpZCI6IjViYjI2Y2U5Yjg5NGViYTBkNzhjZjQ3Y2YyZjc1YWY4ODk5ZjE3MTcifQ.eyJpc3MiOiJodHRwczovL2FwcC5ndy5rOHMuc2guYWJhbi5pby9kZXgiLCJzdWIiOiJDZzB3TFRNNE5TMHlPREE0T1Mwd0VnUnRiMk5yIiwiYXVkIjoiZXhhbXBsZS1hcHAiLCJleHAiOjE3MjQzMzU2ODgsImlhdCI6MTcyNDI0OTI4OCwibm9uY2UiOiJmYTA5MjZhNy1hYWJjLTRiOTYtYWNlNy01OGY4NTNkYWMxYWQiLCJhdF9oYXNoIjoiSC12VDNKTGd1R3BPVHczTm9qM3FCdyIsImNfaGFzaCI6Ii1OMHRJTllsWWt0ZlNaRnp0MFVuUWciLCJlbWFpbCI6ImtpbGdvcmVAa2lsZ29yZS50cm91dCIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiS2lsZ29yZSBUcm91dCJ9.QH9mbZrWX-P7DBfR4V93DhTSuEo8PWWlT7J-YTbZqcIberfkyucYyNoCJynKIGWLNb9yhxtbl4Ar5YhMnyOjzYOJXtPnNy2nTDp1SoOTEtfAb7FBcY1Gq9c-hl7e1nX5PaJlLM57KRhotWt4et8B1VY2e4MgYcj3touMGUV69kHZ0nv18KV7hUkAtsJ3OEQcVcw94wccMPFh7oOJnDIGoVxLAwXNkk-3h3PG0U_G7FoHwqewOfXiMKABSCZ3jH4P7gjTjFQ919iUNcMErp5q3xgcSJ61YDx9blv4AvE2SKq9oDZVQy7EYRPekT2NTrqJFmQU5teKA1SajzEvItBGyw", //nolint:lll
				}
				token := &oauth2.Token{
					AccessToken:  "eyJhbGciOiJSUzI1NiIsImpZCI6IjViYjI2Y2U5Yjg5NGViYTBkNzhjZjQ3Y2YyZjc1YWY4ODk5ZjE3MTcifQ.eyJpc3MiOiJodHRwczovL2FwcC5ndy5rOHMuc2guYWJhbi5pby9kZXgiLCJzdWIiOiJDZzB3TFRNNE5TMHlPREE0T1Mwd0VnUnRiMk5yIiwiYXVkIjoiZXhhbXBsZS1hcHAiLCJleHAiOjE3MjQzMzUwOTQsImlhdCI6MTcyNDI0ODY5NCwibm9uY2UiOiIxYzFmOTU2My1iNWQwLTQwOWItOGRmNS1hYjhhYjdkYzE0ODgiLCJhdF9oYXNoIjoieU5qVFJTTWhZLXFKemJJRTlybG96ZyIsImVtYWlsIjoia2lsZ29yZUBraWxnb3JlLnRyb3V0IiwiZW1haWxfdmVyaWZpZWQiOnRydWUsIm5hbWUiOiJLaWxnb3JlIFRyb3V0In0.k-ajsvzlJfuwGkBy9lW8OgRtlzH0PXpU2HgrU4WaVssEdTc3pQWtw50Aqy5s-HWKUVxMwvLgU1qznTkakglNFAic-vRQzPf0cGQBwW7ojTaz-V9oktUhYIj5BXS001BvLqYGEWkKXS07IL2FYh0WEm7m6xDeOsN_5gwUf3NDpXyIiVfVWXEBRiKjsHK6nVA7YggzKF6K4PwK0BdzSy58AT0Q33zdoOPztb2fWvfjBQWZZFSog3Z0rvr0rurIWWlhSGveZcqkwZj9FeZZlIkGK1VDiQhmkha-lzeomgkh8X07Z8wxXS0tYVJ3S9E7bxMOkb7T-P4LzOo3Oj9GNxP7Zw",
					RefreshToken: "Chl4NjZxa2V4emJhMzJkNmdwZzZ3Y2tycHpoEhlydHp5dTc1d2ZxM29ncGxpcHZwYmtudTRy",
					TokenType:    "bearer",
					Expiry:       time.Now().Add(1 * time.Hour),
				}
				config.EXPECT().Exchange(gomock.Any(), gomock.Any()).DoAndReturn(func(any, any, ...any) (*oauth2.Token, error) {
					return token.WithExtra(extra), nil
				})
				return config
			},
			oauthVerifier: func() *mockOidc.MockTokenVerifier {
				mv := mockOidc.NewMockTokenVerifier(ctrl)
				mv.EXPECT().Verify(gomock.Any(), gomock.Any()).Return(&oidc.IDToken{}, nil)
				return mv
			},
			security: &config.Security{
				Secret: "W+c%1*~(E)^2^#sxx^^7%b34",
			},
		},
		{
			name: "ctx_deadline_exceeded_and_verify_error",
			request: func() *http.Request {
				r := bytes.NewReader([]byte("test data"))
				req := httptest.NewRequest(http.MethodGet, `/admin?state=state_fake`, r)
				req.Header.Add("test", "test value")
				req.AddCookie(&http.Cookie{
					Name:     "nonce",
					Value:    "nonce_fake",
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
				})
				req.AddCookie(&http.Cookie{
					Name:     "state",
					Value:    "state_fake",
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
				})
				req.AddCookie(&http.Cookie{
					Name:     "redirect_url",
					Value:    "/admin",
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
				})
				return req
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: func() http.ResponseWriter {
				recorder := httptest.NewRecorder()
				response.InternalError(recorder)
				return recorder
			},
			oauthConfig: func() *mockOidc.MockConfig {
				config := mockOidc.NewMockConfig(ctrl)
				extra := map[string]interface{}{
					"access_token":  "eyJhbGciOiJSUzI1NiIsImtpZCI6IjViYjI2Y2U5Yjg5NGViYTBkNzhjZjQ3Y2YyZjc1YWY4ODk5ZjE3MTcifQ.eyJpc3MiOiJodHRwczovL2FwcC5ndy5rOHMuc2guYWJhbi5pby9kZXgiLCJzdWIiOiJDZzB3TFRNNE5TMHlPREE0T1Mwd0VnUnRiMk5yIiwiYXVkIjoiZXhhbXBsZS1hcHAiLCJleHAiOjE3MjQzMzU2ODgsImlhdCI6MTcyNDI0OTI4OCwibm9uY2UiOiJmYTA5MjZhNy1hYWJjLTRiOTYtYWNlNy01OGY4NTNkYWMxYWQiLCJhdF9oYXNoIjoiNEpqWWxlREhRMmNLcXFmbWsyRUczZyIsImVtYWlsIjoia2lsZ29yZUBraWxnb3JlLnRyb3V0IiwiZW1haWxfdmVyaWZpZWQiOnRydWUsIm5hbWUiOiJLaWxnb3JlIFRyb3V0In0.B7n3v40XmdfMjMW2UB1NmFLZQyyRORwOPSknF3FcOD1juoJmql-1T_NvuHpTCmFqGFDtBKvZa9_RKoisC2LxVvJoUrlDUzDxUafdu3gL9uIwWAtvfAhAjcOujgOupdcwDL-omfoghV83fHdmUr20vOsitLbiKWzWSjzsMRXT4c6OWiQetTrJrQxWQbqFBm-7iqJjXazzDbmhxBpXg6s_BurnJZ0MjXDgSDWDcUxp4-dgvyd_gxmJS_MgnWEgZAaV_y7qG0CHwLyLRl9FHMbSmHUaIRy8KHbUnp9k0-RncTlEM90E09l7kcEbX-74IsqhiyowqPHFHcS9mAZeNlgeaQ", //nolint:lll
					"token_type":    "bearer",
					"expires_in":    "86399",
					"refresh_token": "ChlyZnlsa2k0eXVscHZzaWs2dGZjYWF5YTRkEhl1cnUzN2VpczZsMzVrZjd6dTJvaHQ3YzZq",
					"id_token":      "eyJhbGciOiJSUzI1NiIsImtpZCI6IjViYjI2Y2U5Yjg5NGViYTBkNzhjZjQ3Y2YyZjc1YWY4ODk5ZjE3MTcifQ.eyJpc3MiOiJodHRwczovL2FwcC5ndy5rOHMuc2guYWJhbi5pby9kZXgiLCJzdWIiOiJDZzB3TFRNNE5TMHlPREE0T1Mwd0VnUnRiMk5yIiwiYXVkIjoiZXhhbXBsZS1hcHAiLCJleHAiOjE3MjQzMzU2ODgsImlhdCI6MTcyNDI0OTI4OCwibm9uY2UiOiJmYTA5MjZhNy1hYWJjLTRiOTYtYWNlNy01OGY4NTNkYWMxYWQiLCJhdF9oYXNoIjoiSC12VDNKTGd1R3BPVHczTm9qM3FCdyIsImNfaGFzaCI6Ii1OMHRJTllsWWt0ZlNaRnp0MFVuUWciLCJlbWFpbCI6ImtpbGdvcmVAa2lsZ29yZS50cm91dCIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiS2lsZ29yZSBUcm91dCJ9.QH9mbZrWX-P7DBfR4V93DhTSuEo8PWWlT7J-YTbZqcIberfkyucYyNoCJynKIGWLNb9yhxtbl4Ar5YhMnyOjzYOJXtPnNy2nTDp1SoOTEtfAb7FBcY1Gq9c-hl7e1nX5PaJlLM57KRhotWt4et8B1VY2e4MgYcj3touMGUV69kHZ0nv18KV7hUkAtsJ3OEQcVcw94wccMPFh7oOJnDIGoVxLAwXNkk-3h3PG0U_G7FoHwqewOfXiMKABSCZ3jH4P7gjTjFQ919iUNcMErp5q3xgcSJ61YDx9blv4AvE2SKq9oDZVQy7EYRPekT2NTrqJFmQU5teKA1SajzEvItBGyw", //nolint:lll
				}
				token := &oauth2.Token{
					AccessToken:  "eyJhbGciOiJSUzI1NiIsImtpZCI6IjViYjI2Y2U5Yjg5NGViYTBkNzhjZjQ3Y2YyZjc1YWY4ODk5ZjE3MTcifQ.eyJpc3MiOiJodHRwczovL2FwcC5ndy5rOHMuc2guYWJhbi5pby9kZXgiLCJzdWIiOiJDZzB3TFRNNE5TMHlPREE0T1Mwd0VnUnRiMk5yIiwiYXVkIjoiZXhhbXBsZS1hcHAiLCJleHAiOjE3MjQzMzUwOTQsImlhdCI6MTcyNDI0ODY5NCwibm9uY2UiOiIxYzFmOTU2My1iNWQwLTQwOWItOGRmNS1hYjhhYjdkYzE0ODgiLCJhdF9oYXNoIjoieU5qVFJTTWhZLXFKemJJRTlybG96ZyIsImVtYWlsIjoia2lsZ29yZUBraWxnb3JlLnRyb3V0IiwiZW1haWxfdmVyaWZpZWQiOnRydWUsIm5hbWUiOiJLaWxnb3JlIFRyb3V0In0.k-ajsvzlJfuwGkBy9lW8OgRtlzH0PXpU2HgrU4WaVssEdTc3pQWtw50Aqy5s-HWKUVxMwvLgU1qznTkakglNFAic-vRQzPf0cGQBwW7ojTaz-V9oktUhYIj5BXS001BvLqYGEWkKXS07IL2FYh0WEm7m6xDeOsN_5gwUf3NDpXyIiVfVWXEBRiKjsHK6nVA7YggzKF6K4PwK0BdzSy58AT0Q33zdoOPztb2fWvfjBQWZZFSog3Z0rvr0rurIWWlhSGveZcqkwZj9FeZZlIkGK1VDiQhmkha-lzeomgkh8X07Z8wxXS0tYVJ3S9E7bxMOkb7T-P4LzOo3Oj9GNxP7Zw", //nolint:lll
					RefreshToken: "Chl4NjZxa2V4emJhMzJkNmdwZzZ3Y2tycHpoEhlydHp5dTc1d2ZxM29ncGxpcHZwYmtudTRy",
					TokenType:    "bearer",
					Expiry:       time.Now().Add(1 * time.Hour),
				}
				config.EXPECT().Exchange(gomock.Any(), gomock.Any()).DoAndReturn(func(any, any, ...any) (*oauth2.Token, error) {
					time.Sleep(15 * time.Second)
					return token.WithExtra(extra), nil
				})
				return config
			},
			oauthVerifier: func() *mockOidc.MockTokenVerifier {
				mv := mockOidc.NewMockTokenVerifier(ctrl)
				mv.EXPECT().Verify(gomock.Any(), gomock.Any()).Return(&oidc.IDToken{
					Nonce: "nonce_fake",
				}, context.DeadlineExceeded)
				return mv
			},
			security: &config.Security{
				Secret: "W+c%1*~(E)^2^#sxx^^7%b34",
			},
			noDeadLine: true,
		},
		{
			name: "success",
			request: func() *http.Request {
				r := bytes.NewReader([]byte("test data"))
				req := httptest.NewRequest(http.MethodGet, `/admin?state=state_fake`, r)
				req.Header.Add("test", "test value")
				req.AddCookie(&http.Cookie{
					Name:     "nonce",
					Value:    "nonce_fake",
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
				})
				req.AddCookie(&http.Cookie{
					Name:     "state",
					Value:    "state_fake",
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
				})
				req.AddCookie(&http.Cookie{
					Name:     "redirect_url",
					Value:    "/admin",
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
				})
				return req
			},
			expectedStatusCode: http.StatusFound,
			expectedResponse: func() http.ResponseWriter {
				recorder := httptest.NewRecorder()
				r := bytes.NewReader([]byte("test data"))
				req := httptest.NewRequest(http.MethodGet, "/admin", r)
				token := &oauth2.Token{
					AccessToken:  "eyJhbGciOiJSUzI1NiIsImtpZCI6IjViYjI2Y2U5Yjg5NGViYTBkNzhjZjQ3Y2YyZjc1YWY4ODk5ZjE3MTcifQ.eyJpc3MiOiJodHRwczovL2FwcC5ndy5rOHMuc2guYWJhbi5pby9kZXgiLCJzdWIiOiJDZzB3TFRNNE5TMHlPREE0T1Mwd0VnUnRiMk5yIiwiYXVkIjoiZXhhbXBsZS1hcHAiLCJleHAiOjE3MjQzMzUwOTQsImlhdCI6MTcyNDI0ODY5NCwibm9uY2UiOiIxYzFmOTU2My1iNWQwLTQwOWItOGRmNS1hYjhhYjdkYzE0ODgiLCJhdF9oYXNoIjoieU5qVFJTTWhZLXFKemJJRTlybG96ZyIsImVtYWlsIjoia2lsZ29yZUBraWxnb3JlLnRyb3V0IiwiZW1haWxfdmVyaWZpZWQiOnRydWUsIm5hbWUiOiJLaWxnb3JlIFRyb3V0In0.k-ajsvzlJfuwGkBy9lW8OgRtlzH0PXpU2HgrU4WaVssEdTc3pQWtw50Aqy5s-HWKUVxMwvLgU1qznTkakglNFAic-vRQzPf0cGQBwW7ojTaz-V9oktUhYIj5BXS001BvLqYGEWkKXS07IL2FYh0WEm7m6xDeOsN_5gwUf3NDpXyIiVfVWXEBRiKjsHK6nVA7YggzKF6K4PwK0BdzSy58AT0Q33zdoOPztb2fWvfjBQWZZFSog3Z0rvr0rurIWWlhSGveZcqkwZj9FeZZlIkGK1VDiQhmkha-lzeomgkh8X07Z8wxXS0tYVJ3S9E7bxMOkb7T-P4LzOo3Oj9GNxP7Zw", //nolint:lll
					RefreshToken: "Chl4NjZxa2V4emJhMzJkNmdwZzZ3Y2tycHpoEhlydHp5dTc1d2ZxM29ncGxpcHZwYmtudTRy",
					TokenType:    "bearer",
					Expiry:       time.Now().Add(1 * time.Hour),
				}

				oauth2Claims, err := json.Marshal(token)
				if err != nil {
					log.Fatal(err)
				}
				tokenJSON, err := encryptor.Encrypt(oauth2Claims)
				if err != nil {
					log.Fatal(err)
				}

				http.SetCookie(recorder, &http.Cookie{
					Name:     "auth",
					Value:    string(tokenJSON),
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
				})
				http.SetCookie(recorder, &http.Cookie{
					Name:     "nonce",
					Value:    "",
					MaxAge:   0,
					Secure:   false,
					HttpOnly: true,
				})
				http.SetCookie(recorder, &http.Cookie{
					Name:     "state",
					Value:    "",
					MaxAge:   0,
					Secure:   false,
					HttpOnly: true,
				})

				http.Redirect(recorder, req, "/admin", http.StatusFound)
				return recorder
			},
			oauthConfig: func() *mockOidc.MockConfig {
				config := mockOidc.NewMockConfig(ctrl)
				extra := map[string]interface{}{
					"access_token":  "eyJhbGciOiJSUzI1NiIsImtpZCI6IjViYjI2Y2U5Yjg5NGViYTBkNzhjZjQ3Y2YyZjc1YWY4ODk5ZjE3MTcifQ.eyJpc3MiOiJodHRwczovL2FwcC5ndy5rOHMuc2guYWJhbi5pby9kZXgiLCJzdWIiOiJDZzB3TFRNNE5TMHlPREE0T1Mwd0VnUnRiMk5yIiwiYXVkIjoiZXhhbXBsZS1hcHAiLCJleHAiOjE3MjQzMzU2ODgsImlhdCI6MTcyNDI0OTI4OCwibm9uY2UiOiJmYTA5MjZhNy1hYWJjLTRiOTYtYWNlNy01OGY4NTNkYWMxYWQiLCJhdF9oYXNoIjoiNEpqWWxlREhRMmNLcXFmbWsyRUczZyIsImVtYWlsIjoia2lsZ29yZUBraWxnb3JlLnRyb3V0IiwiZW1haWxfdmVyaWZpZWQiOnRydWUsIm5hbWUiOiJLaWxnb3JlIFRyb3V0In0.B7n3v40XmdfMjMW2UB1NmFLZQyyRORwOPSknF3FcOD1juoJmql-1T_NvuHpTCmFqGFDtBKvZa9_RKoisC2LxVvJoUrlDUzDxUafdu3gL9uIwWAtvfAhAjcOujgOupdcwDL-omfoghV83fHdmUr20vOsitLbiKWzWSjzsMRXT4c6OWiQetTrJrQxWQbqFBm-7iqJjXazzDbmhxBpXg6s_BurnJZ0MjXDgSDWDcUxp4-dgvyd_gxmJS_MgnWEgZAaV_y7qG0CHwLyLRl9FHMbSmHUaIRy8KHbUnp9k0-RncTlEM90E09l7kcEbX-74IsqhiyowqPHFHcS9mAZeNlgeaQ", //nolint:lll
					"token_type":    "bearer",
					"expires_in":    "86399",
					"refresh_token": "ChlyZnlsa2k0eXVscHZzaWs2dGZjYWF5YTRkEhl1cnUzN2VpczZsMzVrZjd6dTJvaHQ3YzZq",
					"id_token":      "eyJhbGciOiJSUzI1NiIsImtpZCI6IjViYjI2Y2U5Yjg5NGViYTBkNzhjZjQ3Y2YyZjc1YWY4ODk5ZjE3MTcifQ.eyJpc3MiOiJodHRwczovL2FwcC5ndy5rOHMuc2guYWJhbi5pby9kZXgiLCJzdWIiOiJDZzB3TFRNNE5TMHlPREE0T1Mwd0VnUnRiMk5yIiwiYXVkIjoiZXhhbXBsZS1hcHAiLCJleHAiOjE3MjQzMzU2ODgsImlhdCI6MTcyNDI0OTI4OCwibm9uY2UiOiJmYTA5MjZhNy1hYWJjLTRiOTYtYWNlNy01OGY4NTNkYWMxYWQiLCJhdF9oYXNoIjoiSC12VDNKTGd1R3BPVHczTm9qM3FCdyIsImNfaGFzaCI6Ii1OMHRJTllsWWt0ZlNaRnp0MFVuUWciLCJlbWFpbCI6ImtpbGdvcmVAa2lsZ29yZS50cm91dCIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiS2lsZ29yZSBUcm91dCJ9.QH9mbZrWX-P7DBfR4V93DhTSuEo8PWWlT7J-YTbZqcIberfkyucYyNoCJynKIGWLNb9yhxtbl4Ar5YhMnyOjzYOJXtPnNy2nTDp1SoOTEtfAb7FBcY1Gq9c-hl7e1nX5PaJlLM57KRhotWt4et8B1VY2e4MgYcj3touMGUV69kHZ0nv18KV7hUkAtsJ3OEQcVcw94wccMPFh7oOJnDIGoVxLAwXNkk-3h3PG0U_G7FoHwqewOfXiMKABSCZ3jH4P7gjTjFQ919iUNcMErp5q3xgcSJ61YDx9blv4AvE2SKq9oDZVQy7EYRPekT2NTrqJFmQU5teKA1SajzEvItBGyw", //nolint:lll
				}
				token := &oauth2.Token{
					AccessToken:  "eyJhbGciOiJSUzI1NiIsImtpZCI6IjViYjI2Y2U5Yjg5NGViYTBkNzhjZjQ3Y2YyZjc1YWY4ODk5ZjE3MTcifQ.eyJpc3MiOiJodHRwczovL2FwcC5ndy5rOHMuc2guYWJhbi5pby9kZXgiLCJzdWIiOiJDZzB3TFRNNE5TMHlPREE0T1Mwd0VnUnRiMk5yIiwiYXVkIjoiZXhhbXBsZS1hcHAiLCJleHAiOjE3MjQzMzUwOTQsImlhdCI6MTcyNDI0ODY5NCwibm9uY2UiOiIxYzFmOTU2My1iNWQwLTQwOWItOGRmNS1hYjhhYjdkYzE0ODgiLCJhdF9oYXNoIjoieU5qVFJTTWhZLXFKemJJRTlybG96ZyIsImVtYWlsIjoia2lsZ29yZUBraWxnb3JlLnRyb3V0IiwiZW1haWxfdmVyaWZpZWQiOnRydWUsIm5hbWUiOiJLaWxnb3JlIFRyb3V0In0.k-ajsvzlJfuwGkBy9lW8OgRtlzH0PXpU2HgrU4WaVssEdTc3pQWtw50Aqy5s-HWKUVxMwvLgU1qznTkakglNFAic-vRQzPf0cGQBwW7ojTaz-V9oktUhYIj5BXS001BvLqYGEWkKXS07IL2FYh0WEm7m6xDeOsN_5gwUf3NDpXyIiVfVWXEBRiKjsHK6nVA7YggzKF6K4PwK0BdzSy58AT0Q33zdoOPztb2fWvfjBQWZZFSog3Z0rvr0rurIWWlhSGveZcqkwZj9FeZZlIkGK1VDiQhmkha-lzeomgkh8X07Z8wxXS0tYVJ3S9E7bxMOkb7T-P4LzOo3Oj9GNxP7Zw", //nolint:lll
					RefreshToken: "Chl4NjZxa2V4emJhMzJkNmdwZzZ3Y2tycHpoEhlydHp5dTc1d2ZxM29ncGxpcHZwYmtudTRy",
					TokenType:    "bearer",
					Expiry:       time.Now().Add(1 * time.Hour),
				}
				config.EXPECT().Exchange(gomock.Any(), gomock.Any()).Return(token.WithExtra(extra), nil)
				return config
			},
			oauthVerifier: func() *mockOidc.MockTokenVerifier {
				mv := mockOidc.NewMockTokenVerifier(ctrl)
				mv.EXPECT().Verify(gomock.Any(), gomock.Any()).Return(&oidc.IDToken{
					Nonce: "nonce_fake",
				}, nil)
				return mv
			},
			security: &config.Security{
				Secret: "W+c%1*~(E)^2^#sxx^^7%b34",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			deadline, _ := t.Deadline()
			if deadline.Compare(time.Now().Add(2*time.Second)) == -1 && test.noDeadLine {
				t.Skip()
			}
			oauthConfig := test.oauthConfig()
			oauthVerifier := test.oauthVerifier()

			recorder := httptest.NewRecorder()
			oh := NewOauthHandler(slog.New(slog.NewTextHandler(os.Stdin, nil)),
				mockRule.NewMockRule(ctrl), oauthConfig, oauthVerifier, test.security,
				&config.HTTPServer{
					BasePath: "",
				},
			)
			oh.Callback(recorder, test.request())
			result := recorder.Result()

			// Compare status codes
			if result.StatusCode != test.expectedStatusCode {
				t.Errorf("response status code:%d is not match expected value:%d",
					result.StatusCode, test.expectedStatusCode)
			}
			loc, _ := recorder.Result().Location()
			expectRes := test.expectedResponse().(*httptest.ResponseRecorder).Result()
			expectLoc, _ := expectRes.Location()

			if !gomock.Eq(loc).Matches(expectLoc) {
				t.Errorf("location %v is not match expected value: %v", loc, expectLoc)
			}

			// Compare cookies
			resultCookies := result.Cookies()
			expectedCookies := expectRes.Cookies()

			if len(resultCookies) != len(expectedCookies) {
				t.Errorf("Mismatch in number of cookies. Got: %v, expected: %v",
					len(resultCookies), len(expectedCookies))
				return
			}

			for i, c := range resultCookies {
				if c.Name != expectedCookies[i].Name {
					t.Errorf("Unexpected cookie name. Got: %v, expected: %v", c.Name, expectedCookies[i].Name)
				}
			}

			oauthConfig.EXPECT()
			oauthVerifier.EXPECT()
		})
	}
}
