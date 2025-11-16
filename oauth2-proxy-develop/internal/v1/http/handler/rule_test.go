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
	"sort"
	"testing"
	"time"

	ruleEnt "application/internal/v1/entity/rule"
	"application/internal/v1/http/response"
	mockOidc "application/mock/pkg/oidc"
	"application/pkg/oidc"
	"application/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"golang.org/x/oauth2"
)

func TestRule_Handle_Auth_Rule(t *testing.T) {
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
		rule               *ruleEnt.Rule
		expectedResponse   func() http.ResponseWriter
		oauthConfig        func() oidc.Config
		oauthVerifier      func() oidc.TokenVerifier
	}{
		{
			name: "success_redirect_to_oidc",
			request: func() *http.Request {
				r := bytes.NewReader([]byte("test data"))
				req := httptest.NewRequest(http.MethodGet, "/admin", r)
				return req
			},
			ctx:                context.Background(),
			expectedStatusCode: http.StatusFound,
			expectedResponse: func() http.ResponseWriter {
				recorder := httptest.NewRecorder()
				r := bytes.NewReader([]byte("test data"))
				req := httptest.NewRequest(http.MethodGet, "/admin", r)
				http.SetCookie(recorder, &http.Cookie{
					Name:     "state",
					Value:    uuid.NewString(),
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
					Path:     "/",
				})
				http.SetCookie(recorder, &http.Cookie{
					Name:     "nonce",
					Value:    uuid.NewString(),
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
					Path:     "/",
				})
				http.SetCookie(recorder, &http.Cookie{
					Name:     "redirect_url",
					Value:    "/admin",
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
					Path:     "/",
				})

				http.Redirect(recorder, req, "https://example.com", http.StatusFound)
				return recorder
			},
			rule: &ruleEnt.Rule{
				Name:   "auth_exact",
				Action: ruleEnt.ActionAuth,
				Path:   "/auth/exact",
			},
			oauthConfig: func() oidc.Config {
				config := mockOidc.NewMockConfig(ctrl)
				config.EXPECT().AuthCodeURL(gomock.Any(), gomock.Any()).Return("https://example.com")
				return config
			},
			oauthVerifier: func() oidc.TokenVerifier {
				return mockOidc.NewMockTokenVerifier(ctrl)
			},
		},
		{
			name: "success_valid_cookie",
			request: func() *http.Request {
				r := bytes.NewReader([]byte("test data"))
				req := httptest.NewRequest(http.MethodGet, "/admin?sdgdsg=111x", r)
				token := &oauth2.Token{
					AccessToken:  "eyJhbGciOiJSUzI1NiIsImtpZCI6ImUzMmE2Y2UyNDgyMDhlOGFkMTVjNjc1N2ZlNjFiYjc5ZGZmNzZhYmMifQ.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjU1NTYvZGV4Iiwic3ViIjoiQ2lSaU9EZzBOemRqTmkxaU5HVm1MVFEzTnprdFlqVTBZaTB3T0RWalpEVTJNVEF3TVRVU0QyRjFkR2h2Y21OdmJtNWxZM1J2Y2ciLCJhdWQiOiJleGFtcGxlLWFwcCIsImV4cCI6MTcyNDc1MjgwNSwiaWF0IjoxNzI0NjY2NDA1LCJub25jZSI6Ild2a0JtRENNM3AzSlRkMjBsN0hYbmciLCJhdF9oYXNoIjoiS21YZG1FMW9jVk94elEwY2l3QThMQSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJncm91cHMiOlsidXNlciJdLCJuYW1lIjoiOTg5MTIwNDEwMTc2MSIsImZlZGVyYXRlZF9jbGFpbXMiOnsiY29ubmVjdG9yX2lkIjoiYXV0aG9yY29ubmVjdG9yIiwidXNlcl9pZCI6ImI4ODQ3N2M2LWI0ZWYtNDc3OS1iNTRiLTA4NWNkNTYxMDAxNSJ9fQ.QIm6ghZKm8mzrZQZTxhw25IUWZpAmhWHqS-ZhqkDdEB_xXslhfsREhice3A6wQereaFtsXMA82DX08-ZLvViU25fvB1QsL2piR5dg8bakXUhkHPGxOoGT4k4ocyo-isZdPMF0SrkJbpgLWiu9M4WCz9-9nXo-oOeSjs_6_CRDjj5bQOGVY1gpu8ah4yfDUT-xTe-bv7mu0h7pDfEGax7HPP4JPY-0UCNkNtSpVItPdOnTR_hcVWXW3Hl9MT0M9qeLTAHPkdAzaGlSKNi7EutW6Rj4hVcb5wiX08UzhRM_7nwpZjYrX7W3CszehNxnCSQFVjA9lS1Zq_XhhBhJBGe9g",
					RefreshToken: "ChlpMm1wbnMyazJ4bDZsYjJta3ppYmpkbWlsEhl0MnNlZ3d5cXh5ejVvZ3Jrc3llZGxncGJl",
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
				req.AddCookie(&http.Cookie{
					Name:     "auth",
					Value:    string(tokenJSON),
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
					Path:     "/",
				})

				return req
			},
			ctx:                context.Background(),
			expectedStatusCode: http.StatusOK,
			expectedResponse: func() http.ResponseWriter {
				recorder := httptest.NewRecorder()
				recorder.Header().Set("X-AUTH-NAME", "9891204101761")
				recorder.Header().Set("X-AUTH-VERIFIED", "true")
				recorder.Header().Set("X-AUTH-ACCESS-TOKEN", "eyJhbGciOiJSUzI1NiIsImtpZCI6ImUzMmE2Y2UyNDgyMDhlOGFkMTVjNjc1N2ZlNjFiYjc5ZGZmNzZhYmMifQ.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjU1NTYvZGV4Iiwic3ViIjoiQ2lSaU9EZzBOemRqTmkxaU5HVm1MVFEzTnprdFlqVTBZaTB3T0RWalpEVTJNVEF3TVRVU0QyRjFkR2h2Y21OdmJtNWxZM1J2Y2ciLCJhdWQiOiJleGFtcGxlLWFwcCIsImV4cCI6MTcyNDc1MjgwNSwiaWF0IjoxNzI0NjY2NDA1LCJub25jZSI6Ild2a0JtRENNM3AzSlRkMjBsN0hYbmciLCJhdF9oYXNoIjoiS21YZG1FMW9jVk94elEwY2l3QThMQSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJncm91cHMiOlsidXNlciJdLCJuYW1lIjoiOTg5MTIwNDEwMTc2MSIsImZlZGVyYXRlZF9jbGFpbXMiOnsiY29ubmVjdG9yX2lkIjoiYXV0aG9yY29ubmVjdG9yIiwidXNlcl9pZCI6ImI4ODQ3N2M2LWI0ZWYtNDc3OS1iNTRiLTA4NWNkNTYxMDAxNSJ9fQ.QIm6ghZKm8mzrZQZTxhw25IUWZpAmhWHqS-ZhqkDdEB_xXslhfsREhice3A6wQereaFtsXMA82DX08-ZLvViU25fvB1QsL2piR5dg8bakXUhkHPGxOoGT4k4ocyo-isZdPMF0SrkJbpgLWiu9M4WCz9-9nXo-oOeSjs_6_CRDjj5bQOGVY1gpu8ah4yfDUT-xTe-bv7mu0h7pDfEGax7HPP4JPY-0UCNkNtSpVItPdOnTR_hcVWXW3Hl9MT0M9qeLTAHPkdAzaGlSKNi7EutW6Rj4hVcb5wiX08UzhRM_7nwpZjYrX7W3CszehNxnCSQFVjA9lS1Zq_XhhBhJBGe9g")
				recorder.Header().Set("X-AUTH-GROUPS", "user")
				recorder.Header().Set("X-AUTH-SUB", "CiRiODg0NzdjNi1iNGVmLTQ3NzktYjU0Yi0wODVjZDU2MTAwMTUSD2F1dGhvcmNvbm5lY3Rvcg")
				response.Ok(recorder, nil, "access granted")
				return recorder
			},
			rule: &ruleEnt.Rule{
				Name:   "auth_exact",
				Action: ruleEnt.ActionAuth,
				Path:   "/admin?sdgdsg=111x",
			},
			oauthConfig: func() oidc.Config {
				config := mockOidc.NewMockConfig(ctrl)
				return config
			},
			oauthVerifier: func() oidc.TokenVerifier {
				m := mockOidc.NewMockTokenVerifier(ctrl)
				m.EXPECT().Verify(gomock.Any(), gomock.Any())
				return m
			},
		},
		{
			name: "valid_access_token_but_exceeded_expire_time",
			request: func() *http.Request {
				r := bytes.NewReader([]byte("test data"))
				req := httptest.NewRequest(http.MethodGet, "/admin?sdgdsg=111x", r)
				token := &oauth2.Token{
					AccessToken:  "eyJhbGciOiJSUzI1NiIsImtpZCI6IjViYjI2Y2U5Yjg5NGViYTBkNzhjZjQ3Y2YyZjc1YWY4ODk5ZjE3MTcifQ.eyJpc3MiOiJodHRwczovL2FwcC5ndy5rOHMuc2guYWJhbi5pby9kZXgiLCJzdWIiOiJDZzB3TFRNNE5TMHlPREE0T1Mwd0VnUnRiMk5yIiwiYXVkIjoiZXhhbXBsZS1hcHAiLCJleHAiOjE3MjQzMzUwOTQsImlhdCI6MTcyNDI0ODY5NCwibm9uY2UiOiIxYzFmOTU2My1iNWQwLTQwOWItOGRmNS1hYjhhYjdkYzE0ODgiLCJhdF9oYXNoIjoieU5qVFJTTWhZLXFKemJJRTlybG96ZyIsImVtYWlsIjoia2lsZ29yZUBraWxnb3JlLnRyb3V0IiwiZW1haWxfdmVyaWZpZWQiOnRydWUsIm5hbWUiOiJLaWxnb3JlIFRyb3V0In0.k-ajsvzlJfuwGkBy9lW8OgRtlzH0PXpU2HgrU4WaVssEdTc3pQWtw50Aqy5s-HWKUVxMwvLgU1qznTkakglNFAic-vRQzPf0cGQBwW7ojTaz-V9oktUhYIj5BXS001BvLqYGEWkKXS07IL2FYh0WEm7m6xDeOsN_5gwUf3NDpXyIiVfVWXEBRiKjsHK6nVA7YggzKF6K4PwK0BdzSy58AT0Q33zdoOPztb2fWvfjBQWZZFSog3Z0rvr0rurIWWlhSGveZcqkwZj9FeZZlIkGK1VDiQhmkha-lzeomgkh8X07Z8wxXS0tYVJ3S9E7bxMOkb7T-P4LzOo3Oj9GNxP7Zw",
					RefreshToken: "Chl4NjZxa2V4emJhMzJkNmdwZzZ3Y2tycHpoEhlydHp5dTc1d2ZxM29ncGxpcHZwYmtudTRy",
					TokenType:    "bearer",
					Expiry:       time.Now().Add(-1 * time.Hour),
				}
				oauth2Claims, err := json.Marshal(token)
				if err != nil {
					log.Fatal(err)
				}
				tokenJSON, err := encryptor.Encrypt(oauth2Claims)
				if err != nil {
					log.Fatal(err)
				}
				req.AddCookie(&http.Cookie{
					Name:     "auth",
					Value:    string(tokenJSON),
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
				})
				return req
			},
			ctx:                context.Background(),
			expectedStatusCode: http.StatusOK,
			expectedResponse: func() http.ResponseWriter {
				recorder := httptest.NewRecorder()
				recorder.Header().Set("X-AUTH-NAME", "9891204101761")
				recorder.Header().Set("X-AUTH-VERIFIED", "true")
				recorder.Header().Set("X-AUTH-ACCESS-TOKEN", "eyJhbGciOiJSUzI1NiIsImtpZCI6ImUzMmE2Y2UyNDgyMDhlOGFkMTVjNjc1N2ZlNjFiYjc5ZGZmNzZhYmMifQ.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjU1NTYvZGV4Iiwic3ViIjoiQ2lSaU9EZzBOemRqTmkxaU5HVm1MVFEzTnprdFlqVTBZaTB3T0RWalpEVTJNVEF3TVRVU0QyRjFkR2h2Y21OdmJtNWxZM1J2Y2ciLCJhdWQiOiJleGFtcGxlLWFwcCIsImV4cCI6MTcyNDc1MjgwNSwiaWF0IjoxNzI0NjY2NDA1LCJub25jZSI6Ild2a0JtRENNM3AzSlRkMjBsN0hYbmciLCJhdF9oYXNoIjoiS21YZG1FMW9jVk94elEwY2l3QThMQSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJncm91cHMiOlsidXNlciJdLCJuYW1lIjoiOTg5MTIwNDEwMTc2MSIsImZlZGVyYXRlZF9jbGFpbXMiOnsiY29ubmVjdG9yX2lkIjoiYXV0aG9yY29ubmVjdG9yIiwidXNlcl9pZCI6ImI4ODQ3N2M2LWI0ZWYtNDc3OS1iNTRiLTA4NWNkNTYxMDAxNSJ9fQ.QIm6ghZKm8mzrZQZTxhw25IUWZpAmhWHqS-ZhqkDdEB_xXslhfsREhice3A6wQereaFtsXMA82DX08-ZLvViU25fvB1QsL2piR5dg8bakXUhkHPGxOoGT4k4ocyo-isZdPMF0SrkJbpgLWiu9M4WCz9-9nXo-oOeSjs_6_CRDjj5bQOGVY1gpu8ah4yfDUT-xTe-bv7mu0h7pDfEGax7HPP4JPY-0UCNkNtSpVItPdOnTR_hcVWXW3Hl9MT0M9qeLTAHPkdAzaGlSKNi7EutW6Rj4hVcb5wiX08UzhRM_7nwpZjYrX7W3CszehNxnCSQFVjA9lS1Zq_XhhBhJBGe9g")
				recorder.Header().Set("X-AUTH-GROUPS", "user")
				recorder.Header().Set("X-AUTH-SUB", "CiRiODg0NzdjNi1iNGVmLTQ3NzktYjU0Yi0wODVjZDU2MTAwMTUSD2F1dGhvcmNvbm5lY3Rvcg")
				token := &oauth2.Token{
					AccessToken:  "eyJhbGciOiJSUzI1NiIsImtpZCI6ImUzMmE2Y2UyNDgyMDhlOGFkMTVjNjc1N2ZlNjFiYjc5ZGZmNzZhYmMifQ.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjU1NTYvZGV4Iiwic3ViIjoiQ2lSaU9EZzBOemRqTmkxaU5HVm1MVFEzTnprdFlqVTBZaTB3T0RWalpEVTJNVEF3TVRVU0QyRjFkR2h2Y21OdmJtNWxZM1J2Y2ciLCJhdWQiOiJleGFtcGxlLWFwcCIsImV4cCI6MTcyNDc1MjgwNSwiaWF0IjoxNzI0NjY2NDA1LCJub25jZSI6Ild2a0JtRENNM3AzSlRkMjBsN0hYbmciLCJhdF9oYXNoIjoiS21YZG1FMW9jVk94elEwY2l3QThMQSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJncm91cHMiOlsidXNlciJdLCJuYW1lIjoiOTg5MTIwNDEwMTc2MSIsImZlZGVyYXRlZF9jbGFpbXMiOnsiY29ubmVjdG9yX2lkIjoiYXV0aG9yY29ubmVjdG9yIiwidXNlcl9pZCI6ImI4ODQ3N2M2LWI0ZWYtNDc3OS1iNTRiLTA4NWNkNTYxMDAxNSJ9fQ.QIm6ghZKm8mzrZQZTxhw25IUWZpAmhWHqS-ZhqkDdEB_xXslhfsREhice3A6wQereaFtsXMA82DX08-ZLvViU25fvB1QsL2piR5dg8bakXUhkHPGxOoGT4k4ocyo-isZdPMF0SrkJbpgLWiu9M4WCz9-9nXo-oOeSjs_6_CRDjj5bQOGVY1gpu8ah4yfDUT-xTe-bv7mu0h7pDfEGax7HPP4JPY-0UCNkNtSpVItPdOnTR_hcVWXW3Hl9MT0M9qeLTAHPkdAzaGlSKNi7EutW6Rj4hVcb5wiX08UzhRM_7nwpZjYrX7W3CszehNxnCSQFVjA9lS1Zq_XhhBhJBGe9g",
					RefreshToken: "Chl4NjZxa2V4emJhMzJkNmdwZzZ3Y2tycHpoEhlydHp5dTc1d2ZxM29ncGxpcHZwYmtudTRy",
					TokenType:    "bearer",
					Expiry:       time.Now().Add(1 * time.Hour),
				}
				oauth2Claims, err := json.Marshal(token)
				if err != nil {
					t.Error(err.Error())
				}

				tokenJSON, err := encryptor.Encrypt(oauth2Claims)
				if err != nil {
					t.Error(err.Error())
				}

				http.SetCookie(recorder, &http.Cookie{
					Name:     "auth",
					Value:    string(tokenJSON),
					Secure:   false,
					HttpOnly: true,
					Path:     "/",
				})

				response.Ok(recorder, nil, "access granted")
				return recorder
			},
			rule: &ruleEnt.Rule{
				Name:   "auth_exact",
				Action: ruleEnt.ActionAuth,
				Path:   "/admin?sdgdsg=111x",
			},
			oauthConfig: func() oidc.Config {
				token := &oauth2.Token{
					AccessToken:  "eyJhbGciOiJSUzI1NiIsImtpZCI6ImUzMmE2Y2UyNDgyMDhlOGFkMTVjNjc1N2ZlNjFiYjc5ZGZmNzZhYmMifQ.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjU1NTYvZGV4Iiwic3ViIjoiQ2lSaU9EZzBOemRqTmkxaU5HVm1MVFEzTnprdFlqVTBZaTB3T0RWalpEVTJNVEF3TVRVU0QyRjFkR2h2Y21OdmJtNWxZM1J2Y2ciLCJhdWQiOiJleGFtcGxlLWFwcCIsImV4cCI6MTcyNDc1MjgwNSwiaWF0IjoxNzI0NjY2NDA1LCJub25jZSI6Ild2a0JtRENNM3AzSlRkMjBsN0hYbmciLCJhdF9oYXNoIjoiS21YZG1FMW9jVk94elEwY2l3QThMQSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJncm91cHMiOlsidXNlciJdLCJuYW1lIjoiOTg5MTIwNDEwMTc2MSIsImZlZGVyYXRlZF9jbGFpbXMiOnsiY29ubmVjdG9yX2lkIjoiYXV0aG9yY29ubmVjdG9yIiwidXNlcl9pZCI6ImI4ODQ3N2M2LWI0ZWYtNDc3OS1iNTRiLTA4NWNkNTYxMDAxNSJ9fQ.QIm6ghZKm8mzrZQZTxhw25IUWZpAmhWHqS-ZhqkDdEB_xXslhfsREhice3A6wQereaFtsXMA82DX08-ZLvViU25fvB1QsL2piR5dg8bakXUhkHPGxOoGT4k4ocyo-isZdPMF0SrkJbpgLWiu9M4WCz9-9nXo-oOeSjs_6_CRDjj5bQOGVY1gpu8ah4yfDUT-xTe-bv7mu0h7pDfEGax7HPP4JPY-0UCNkNtSpVItPdOnTR_hcVWXW3Hl9MT0M9qeLTAHPkdAzaGlSKNi7EutW6Rj4hVcb5wiX08UzhRM_7nwpZjYrX7W3CszehNxnCSQFVjA9lS1Zq_XhhBhJBGe9g",
					RefreshToken: "Chl4NjZxa2V4emJhMzJkNmdwZzZ3Y2tycHpoEhlydHp5dTc1d2ZxM29ncGxpcHZwYmtudTRy",
					TokenType:    "bearer",
					Expiry:       time.Now().Add(1 * time.Hour),
				}
				tokenSource := mockOidc.NewMockTokenSource(ctrl)
				tokenSource.EXPECT().Token().Return(token, nil)
				config := mockOidc.NewMockConfig(ctrl)
				config.EXPECT().TokenSource(gomock.Any(), gomock.Any()).Return(tokenSource)
				return config
			},
			oauthVerifier: func() oidc.TokenVerifier {
				m := mockOidc.NewMockTokenVerifier(ctrl)
				m.EXPECT().Verify(gomock.Any(), gomock.Any())
				return m
			},
		},
		{
			name: "invalid_access_token_invalid_refresh_token",
			request: func() *http.Request {
				r := bytes.NewReader([]byte("test data"))
				req := httptest.NewRequest(http.MethodGet, "/admin?sdgdsg=111x", r)
				token := &oauth2.Token{
					AccessToken:  "eyJhbGciOiJSUzI1NiIsImtpZCI6IjViYjI2Y2U5Yjg5NGViYTBkNzhjZjQ3Y2YyZjc1YWY4ODk5ZjE3MTcifQ.eyJpc3MiOiJodHRwczovL2FwcC5ndy5rOHMuc2guYWJhbi5pby9kZXgiLCJzdWIiOiJDZzB3TFRNNE5TMHlPREE0T1Mwd0VnUnRiMk5yIiwiYXVkIjoiZXhhbXBsZS1hcHAiLCJleHAiOjE3MjQzMzUwOTQsImlhdCI6MTcyNDI0ODY5NCwibm9uY2UiOiIxYzFmOTU2My1iNWQwLTQwOWItOGRmNS1hYjhhYjdkYzE0ODgiLCJhdF9oYXNoIjoieU5qVFJTTWhZLXFKemJJRTlybG96ZyIsImVtYWlsIjoia2lsZ29yZUBraWxnb3JlLnRyb3V0IiwiZW1haWxfdmVyaWZpZWQiOnRydWUsIm5hbWUiOiJLaWxnb3JlIFRyb3V0In0.k-ajsvzlJfuwGkBy9lW8OgRtlzH0PXpU2HgrU4WaVssEdTc3pQWtw50Aqy5s-HWKUVxMwvLgU1qznTkakglNFAic-vRQzPf0cGQBwW7ojTaz-V9oktUhYIj5BXS001BvLqYGEWkKXS07IL2FYh0WEm7m6xDeOsN_5gwUf3NDpXyIiVfVWXEBRiKjsHK6nVA7YggzKF6K4PwK0BdzSy58AT0Q33zdoOPztb2fWvfjBQWZZFSog3Z0rvr0rurIWWlhSGveZcqkwZj9FeZZlIkGK1VDiQhmkha-lzeomgkh8X07Z8wxXS0tYVJ3S9E7bxMOkb7T-P4LzOo3Oj9GNxP7Zw",
					RefreshToken: "Chl4NjZxa2V4emJhMzJkNmdwZzZ3Y2tycHpoEhlydHp5dTc1d2ZxM29ncGxpcHZwYmtudTRy",
					TokenType:    "bearer",
					Expiry:       time.Now().Add(-1 * time.Hour),
				}
				oauth2Claims, err := json.Marshal(token)
				if err != nil {
					log.Fatal(err)
				}
				tokenJSON, err := encryptor.Encrypt(oauth2Claims)
				if err != nil {
					log.Fatal(err)
				}
				req.AddCookie(&http.Cookie{
					Name:     "auth",
					Value:    string(tokenJSON),
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
					Path:     "/",
				})
				return req
			},
			ctx:                context.Background(),
			expectedStatusCode: http.StatusForbidden,
			expectedResponse: func() http.ResponseWriter {
				recorder := httptest.NewRecorder()
				http.SetCookie(recorder, &http.Cookie{
					Name:     "auth",
					Value:    "",
					Path:     "/",
					Expires:  time.Unix(0, 0),
					HttpOnly: true,
				})
				response.Custom(recorder, http.StatusForbidden, nil, "access denied")
				return recorder
			},
			rule: &ruleEnt.Rule{
				Name:   "auth_exact",
				Action: ruleEnt.ActionAuth,
				Path:   "/admin?sdgdsg=111x",
			},
			oauthConfig: func() oidc.Config {
				tokenSource := mockOidc.NewMockTokenSource(ctrl)
				tokenSource.EXPECT().Token().Return(nil, errors.New("invalid refresh token"))
				config := mockOidc.NewMockConfig(ctrl)
				config.EXPECT().TokenSource(gomock.Any(), gomock.Any()).Return(tokenSource)
				return config
			},
			oauthVerifier: func() oidc.TokenVerifier {
				m := mockOidc.NewMockTokenVerifier(ctrl)
				return m
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			oauthConfig := test.oauthConfig()
			oauthVerifier := test.oauthVerifier()
			recorder := httptest.NewRecorder()
			rh := NewRuleHandler(test.rule, oauthConfig, oauthVerifier,
				slog.New(slog.NewTextHandler(os.Stdin, nil)), encryptor)
			rh.Handle(recorder, test.request())

			// Get result and the expected result
			result := recorder.Result()
			expectedResult := test.expectedResponse().(*httptest.ResponseRecorder).Result()

			// Compare status codes
			if result.StatusCode != test.expectedStatusCode {
				t.Errorf("response status code:%d is not match expected value:%d",
					result.StatusCode, test.expectedStatusCode)
			}

			// Compare cookies
			resultCookies := result.Cookies()
			expectedCookies := expectedResult.Cookies()

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

			// Compare headers
			for k, expectedValues := range expectedResult.Header {
				// skip cookies we check them separately
				if k == "Set-Cookie" {
					continue
				}
				if resultValues, ok := result.Header[k]; ok {
					if !isStringSliceEqual(resultValues, expectedValues) {
						t.Error("Header:", k, "Mismatch in result header value(s). Got: ", resultValues, "expected: ", expectedValues)
						return
					}
				} else {
					t.Errorf("Expected header %s is missing from the result", k)
					return
				}
			}
		})
	}
}

func TestRule_Handle_Allow_Rule(t *testing.T) {
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
		rule               *ruleEnt.Rule
		expectedResponse   func() http.ResponseWriter
		oauthConfig        func() oidc.Config
		oauthVerifier      func() oidc.TokenVerifier
	}{
		{
			name: "success_right_cookie_with_headers",
			rule: &ruleEnt.Rule{
				Name:   "allow_rule",
				Action: ruleEnt.ActionAllow,
				Path:   "/",
			},
			request: func() *http.Request {
				r := bytes.NewReader([]byte("test data"))
				req := httptest.NewRequest(http.MethodGet, "/admin?sdgdsg=111x", r)
				token := &oauth2.Token{
					AccessToken:  "eyJhbGciOiJSUzI1NiIsImtpZCI6ImUzMmE2Y2UyNDgyMDhlOGFkMTVjNjc1N2ZlNjFiYjc5ZGZmNzZhYmMifQ.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjU1NTYvZGV4Iiwic3ViIjoiQ2lSaU9EZzBOemRqTmkxaU5HVm1MVFEzTnprdFlqVTBZaTB3T0RWalpEVTJNVEF3TVRVU0QyRjFkR2h2Y21OdmJtNWxZM1J2Y2ciLCJhdWQiOiJleGFtcGxlLWFwcCIsImV4cCI6MTcyNDc1MjgwNSwiaWF0IjoxNzI0NjY2NDA1LCJub25jZSI6Ild2a0JtRENNM3AzSlRkMjBsN0hYbmciLCJhdF9oYXNoIjoiS21YZG1FMW9jVk94elEwY2l3QThMQSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJncm91cHMiOlsidXNlciJdLCJuYW1lIjoiOTg5MTIwNDEwMTc2MSIsImZlZGVyYXRlZF9jbGFpbXMiOnsiY29ubmVjdG9yX2lkIjoiYXV0aG9yY29ubmVjdG9yIiwidXNlcl9pZCI6ImI4ODQ3N2M2LWI0ZWYtNDc3OS1iNTRiLTA4NWNkNTYxMDAxNSJ9fQ.QIm6ghZKm8mzrZQZTxhw25IUWZpAmhWHqS-ZhqkDdEB_xXslhfsREhice3A6wQereaFtsXMA82DX08-ZLvViU25fvB1QsL2piR5dg8bakXUhkHPGxOoGT4k4ocyo-isZdPMF0SrkJbpgLWiu9M4WCz9-9nXo-oOeSjs_6_CRDjj5bQOGVY1gpu8ah4yfDUT-xTe-bv7mu0h7pDfEGax7HPP4JPY-0UCNkNtSpVItPdOnTR_hcVWXW3Hl9MT0M9qeLTAHPkdAzaGlSKNi7EutW6Rj4hVcb5wiX08UzhRM_7nwpZjYrX7W3CszehNxnCSQFVjA9lS1Zq_XhhBhJBGe9g",
					RefreshToken: "ChlpMm1wbnMyazJ4bDZsYjJta3ppYmpkbWlsEhl0MnNlZ3d5cXh5ejVvZ3Jrc3llZGxncGJl",
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
				req.AddCookie(&http.Cookie{
					Name:     "auth",
					Value:    string(tokenJSON),
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
					Path:     "/",
				})

				return req
			},
			ctx:                context.Background(),
			expectedStatusCode: http.StatusOK,
			expectedResponse: func() http.ResponseWriter {
				recorder := httptest.NewRecorder()
				recorder.Header().Set("X-AUTH-NAME", "9891204101761")
				recorder.Header().Set("X-AUTH-VERIFIED", "true")
				recorder.Header().Set("X-AUTH-ACCESS-TOKEN", "eyJhbGciOiJSUzI1NiIsImtpZCI6ImUzMmE2Y2UyNDgyMDhlOGFkMTVjNjc1N2ZlNjFiYjc5ZGZmNzZhYmMifQ.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjU1NTYvZGV4Iiwic3ViIjoiQ2lSaU9EZzBOemRqTmkxaU5HVm1MVFEzTnprdFlqVTBZaTB3T0RWalpEVTJNVEF3TVRVU0QyRjFkR2h2Y21OdmJtNWxZM1J2Y2ciLCJhdWQiOiJleGFtcGxlLWFwcCIsImV4cCI6MTcyNDc1MjgwNSwiaWF0IjoxNzI0NjY2NDA1LCJub25jZSI6Ild2a0JtRENNM3AzSlRkMjBsN0hYbmciLCJhdF9oYXNoIjoiS21YZG1FMW9jVk94elEwY2l3QThMQSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJncm91cHMiOlsidXNlciJdLCJuYW1lIjoiOTg5MTIwNDEwMTc2MSIsImZlZGVyYXRlZF9jbGFpbXMiOnsiY29ubmVjdG9yX2lkIjoiYXV0aG9yY29ubmVjdG9yIiwidXNlcl9pZCI6ImI4ODQ3N2M2LWI0ZWYtNDc3OS1iNTRiLTA4NWNkNTYxMDAxNSJ9fQ.QIm6ghZKm8mzrZQZTxhw25IUWZpAmhWHqS-ZhqkDdEB_xXslhfsREhice3A6wQereaFtsXMA82DX08-ZLvViU25fvB1QsL2piR5dg8bakXUhkHPGxOoGT4k4ocyo-isZdPMF0SrkJbpgLWiu9M4WCz9-9nXo-oOeSjs_6_CRDjj5bQOGVY1gpu8ah4yfDUT-xTe-bv7mu0h7pDfEGax7HPP4JPY-0UCNkNtSpVItPdOnTR_hcVWXW3Hl9MT0M9qeLTAHPkdAzaGlSKNi7EutW6Rj4hVcb5wiX08UzhRM_7nwpZjYrX7W3CszehNxnCSQFVjA9lS1Zq_XhhBhJBGe9g")
				recorder.Header().Set("X-AUTH-GROUPS", "user")
				recorder.Header().Set("X-AUTH-SUB", "CiRiODg0NzdjNi1iNGVmLTQ3NzktYjU0Yi0wODVjZDU2MTAwMTUSD2F1dGhvcmNvbm5lY3Rvcg")
				response.Ok(recorder, nil, "access granted")
				return recorder
			},
			oauthConfig: func() oidc.Config {
				return mockOidc.NewMockConfig(ctrl)
			},
			oauthVerifier: func() oidc.TokenVerifier {
				m := mockOidc.NewMockTokenVerifier(ctrl)
				m.EXPECT().Verify(gomock.Any(), gomock.Any())
				return m
			},
		},
		{
			name: "success_expired_access_token_with_headers",
			rule: &ruleEnt.Rule{
				Name:   "allow_rule",
				Action: ruleEnt.ActionAllow,
				Path:   "/",
			},
			request: func() *http.Request {
				r := bytes.NewReader([]byte("test data"))
				req := httptest.NewRequest(http.MethodGet, "/admin?sdgdsg=111x", r)
				token := &oauth2.Token{
					AccessToken:  "eyJhbGciOiJSUzI1NiIsImtpZCI6IjViYjI2Y2U5Yjg5NGViYTBkNzhjZjQ3Y2YyZjc1YWY4ODk5ZjE3MTcifQ.eyJpc3MiOiJodHRwczovL2FwcC5ndy5rOHMuc2guYWJhbi5pby9kZXgiLCJzdWIiOiJDZzB3TFRNNE5TMHlPREE0T1Mwd0VnUnRiMk5yIiwiYXVkIjoiZXhhbXBsZS1hcHAiLCJleHAiOjE3MjQzMzUwOTQsImlhdCI6MTcyNDI0ODY5NCwibm9uY2UiOiIxYzFmOTU2My1iNWQwLTQwOWItOGRmNS1hYjhhYjdkYzE0ODgiLCJhdF9oYXNoIjoieU5qVFJTTWhZLXFKemJJRTlybG96ZyIsImVtYWlsIjoia2lsZ29yZUBraWxnb3JlLnRyb3V0IiwiZW1haWxfdmVyaWZpZWQiOnRydWUsIm5hbWUiOiJLaWxnb3JlIFRyb3V0In0.k-ajsvzlJfuwGkBy9lW8OgRtlzH0PXpU2HgrU4WaVssEdTc3pQWtw50Aqy5s-HWKUVxMwvLgU1qznTkakglNFAic-vRQzPf0cGQBwW7ojTaz-V9oktUhYIj5BXS001BvLqYGEWkKXS07IL2FYh0WEm7m6xDeOsN_5gwUf3NDpXyIiVfVWXEBRiKjsHK6nVA7YggzKF6K4PwK0BdzSy58AT0Q33zdoOPztb2fWvfjBQWZZFSog3Z0rvr0rurIWWlhSGveZcqkwZj9FeZZlIkGK1VDiQhmkha-lzeomgkh8X07Z8wxXS0tYVJ3S9E7bxMOkb7T-P4LzOo3Oj9GNxP7Zw",
					RefreshToken: "Chl4NjZxa2V4emJhMzJkNmdwZzZ3Y2tycHpoEhlydHp5dTc1d2ZxM29ncGxpcHZwYmtudTRy",
					TokenType:    "bearer",
					Expiry:       time.Now().Add(-1 * time.Hour),
				}
				oauth2Claims, err := json.Marshal(token)
				if err != nil {
					log.Fatal(err)
				}
				tokenJSON, err := encryptor.Encrypt(oauth2Claims)
				if err != nil {
					log.Fatal(err)
				}
				req.AddCookie(&http.Cookie{
					Name:     "auth",
					Value:    string(tokenJSON),
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
				})
				return req
			},
			ctx:                context.Background(),
			expectedStatusCode: http.StatusOK,
			expectedResponse: func() http.ResponseWriter {
				recorder := httptest.NewRecorder()
				recorder.Header().Set("X-AUTH-NAME", "9891204101761")
				recorder.Header().Set("X-AUTH-VERIFIED", "true")
				recorder.Header().Set("X-AUTH-ACCESS-TOKEN", "eyJhbGciOiJSUzI1NiIsImtpZCI6ImUzMmE2Y2UyNDgyMDhlOGFkMTVjNjc1N2ZlNjFiYjc5ZGZmNzZhYmMifQ.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjU1NTYvZGV4Iiwic3ViIjoiQ2lSaU9EZzBOemRqTmkxaU5HVm1MVFEzTnprdFlqVTBZaTB3T0RWalpEVTJNVEF3TVRVU0QyRjFkR2h2Y21OdmJtNWxZM1J2Y2ciLCJhdWQiOiJleGFtcGxlLWFwcCIsImV4cCI6MTcyNDc1MjgwNSwiaWF0IjoxNzI0NjY2NDA1LCJub25jZSI6Ild2a0JtRENNM3AzSlRkMjBsN0hYbmciLCJhdF9oYXNoIjoiS21YZG1FMW9jVk94elEwY2l3QThMQSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJncm91cHMiOlsidXNlciJdLCJuYW1lIjoiOTg5MTIwNDEwMTc2MSIsImZlZGVyYXRlZF9jbGFpbXMiOnsiY29ubmVjdG9yX2lkIjoiYXV0aG9yY29ubmVjdG9yIiwidXNlcl9pZCI6ImI4ODQ3N2M2LWI0ZWYtNDc3OS1iNTRiLTA4NWNkNTYxMDAxNSJ9fQ.QIm6ghZKm8mzrZQZTxhw25IUWZpAmhWHqS-ZhqkDdEB_xXslhfsREhice3A6wQereaFtsXMA82DX08-ZLvViU25fvB1QsL2piR5dg8bakXUhkHPGxOoGT4k4ocyo-isZdPMF0SrkJbpgLWiu9M4WCz9-9nXo-oOeSjs_6_CRDjj5bQOGVY1gpu8ah4yfDUT-xTe-bv7mu0h7pDfEGax7HPP4JPY-0UCNkNtSpVItPdOnTR_hcVWXW3Hl9MT0M9qeLTAHPkdAzaGlSKNi7EutW6Rj4hVcb5wiX08UzhRM_7nwpZjYrX7W3CszehNxnCSQFVjA9lS1Zq_XhhBhJBGe9g")
				recorder.Header().Set("X-AUTH-GROUPS", "user")
				recorder.Header().Set("X-AUTH-SUB", "CiRiODg0NzdjNi1iNGVmLTQ3NzktYjU0Yi0wODVjZDU2MTAwMTUSD2F1dGhvcmNvbm5lY3Rvcg")
				token := &oauth2.Token{
					AccessToken:  "eyJhbGciOiJSUzI1NiIsImtpZCI6ImUzMmE2Y2UyNDgyMDhlOGFkMTVjNjc1N2ZlNjFiYjc5ZGZmNzZhYmMifQ.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjU1NTYvZGV4Iiwic3ViIjoiQ2lSaU9EZzBOemRqTmkxaU5HVm1MVFEzTnprdFlqVTBZaTB3T0RWalpEVTJNVEF3TVRVU0QyRjFkR2h2Y21OdmJtNWxZM1J2Y2ciLCJhdWQiOiJleGFtcGxlLWFwcCIsImV4cCI6MTcyNDc1MjgwNSwiaWF0IjoxNzI0NjY2NDA1LCJub25jZSI6Ild2a0JtRENNM3AzSlRkMjBsN0hYbmciLCJhdF9oYXNoIjoiS21YZG1FMW9jVk94elEwY2l3QThMQSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJncm91cHMiOlsidXNlciJdLCJuYW1lIjoiOTg5MTIwNDEwMTc2MSIsImZlZGVyYXRlZF9jbGFpbXMiOnsiY29ubmVjdG9yX2lkIjoiYXV0aG9yY29ubmVjdG9yIiwidXNlcl9pZCI6ImI4ODQ3N2M2LWI0ZWYtNDc3OS1iNTRiLTA4NWNkNTYxMDAxNSJ9fQ.QIm6ghZKm8mzrZQZTxhw25IUWZpAmhWHqS-ZhqkDdEB_xXslhfsREhice3A6wQereaFtsXMA82DX08-ZLvViU25fvB1QsL2piR5dg8bakXUhkHPGxOoGT4k4ocyo-isZdPMF0SrkJbpgLWiu9M4WCz9-9nXo-oOeSjs_6_CRDjj5bQOGVY1gpu8ah4yfDUT-xTe-bv7mu0h7pDfEGax7HPP4JPY-0UCNkNtSpVItPdOnTR_hcVWXW3Hl9MT0M9qeLTAHPkdAzaGlSKNi7EutW6Rj4hVcb5wiX08UzhRM_7nwpZjYrX7W3CszehNxnCSQFVjA9lS1Zq_XhhBhJBGe9g",
					RefreshToken: "Chl4NjZxa2V4emJhMzJkNmdwZzZ3Y2tycHpoEhlydHp5dTc1d2ZxM29ncGxpcHZwYmtudTRy",
					TokenType:    "bearer",
					Expiry:       time.Now().Add(1 * time.Hour),
				}
				oauth2Claims, err := json.Marshal(token)
				if err != nil {
					t.Error(err.Error())
				}

				tokenJSON, err := encryptor.Encrypt(oauth2Claims)
				if err != nil {
					t.Error(err.Error())
				}

				http.SetCookie(recorder, &http.Cookie{
					Name:     "auth",
					Value:    string(tokenJSON),
					Secure:   false,
					HttpOnly: true,
					Path:     "/",
				})

				response.Ok(recorder, nil, "access granted")
				return recorder
			},
			oauthConfig: func() oidc.Config {
				token := &oauth2.Token{
					AccessToken:  "eyJhbGciOiJSUzI1NiIsImtpZCI6ImUzMmE2Y2UyNDgyMDhlOGFkMTVjNjc1N2ZlNjFiYjc5ZGZmNzZhYmMifQ.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjU1NTYvZGV4Iiwic3ViIjoiQ2lSaU9EZzBOemRqTmkxaU5HVm1MVFEzTnprdFlqVTBZaTB3T0RWalpEVTJNVEF3TVRVU0QyRjFkR2h2Y21OdmJtNWxZM1J2Y2ciLCJhdWQiOiJleGFtcGxlLWFwcCIsImV4cCI6MTcyNDc1MjgwNSwiaWF0IjoxNzI0NjY2NDA1LCJub25jZSI6Ild2a0JtRENNM3AzSlRkMjBsN0hYbmciLCJhdF9oYXNoIjoiS21YZG1FMW9jVk94elEwY2l3QThMQSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJncm91cHMiOlsidXNlciJdLCJuYW1lIjoiOTg5MTIwNDEwMTc2MSIsImZlZGVyYXRlZF9jbGFpbXMiOnsiY29ubmVjdG9yX2lkIjoiYXV0aG9yY29ubmVjdG9yIiwidXNlcl9pZCI6ImI4ODQ3N2M2LWI0ZWYtNDc3OS1iNTRiLTA4NWNkNTYxMDAxNSJ9fQ.QIm6ghZKm8mzrZQZTxhw25IUWZpAmhWHqS-ZhqkDdEB_xXslhfsREhice3A6wQereaFtsXMA82DX08-ZLvViU25fvB1QsL2piR5dg8bakXUhkHPGxOoGT4k4ocyo-isZdPMF0SrkJbpgLWiu9M4WCz9-9nXo-oOeSjs_6_CRDjj5bQOGVY1gpu8ah4yfDUT-xTe-bv7mu0h7pDfEGax7HPP4JPY-0UCNkNtSpVItPdOnTR_hcVWXW3Hl9MT0M9qeLTAHPkdAzaGlSKNi7EutW6Rj4hVcb5wiX08UzhRM_7nwpZjYrX7W3CszehNxnCSQFVjA9lS1Zq_XhhBhJBGe9g",
					RefreshToken: "Chl4NjZxa2V4emJhMzJkNmdwZzZ3Y2tycHpoEhlydHp5dTc1d2ZxM29ncGxpcHZwYmtudTRy",
					TokenType:    "bearer",
					Expiry:       time.Now().Add(1 * time.Hour),
				}
				tokenSource := mockOidc.NewMockTokenSource(ctrl)
				tokenSource.EXPECT().Token().Return(token, nil)
				config := mockOidc.NewMockConfig(ctrl)
				config.EXPECT().TokenSource(gomock.Any(), gomock.Any()).Return(tokenSource)
				return config
			},
			oauthVerifier: func() oidc.TokenVerifier {
				m := mockOidc.NewMockTokenVerifier(ctrl)
				m.EXPECT().Verify(gomock.Any(), gomock.Any())
				return m
			},
		},
		{
			name: "success_without_cookie",
			rule: &ruleEnt.Rule{
				Name:   "allow_rule",
				Action: ruleEnt.ActionAllow,
				Path:   "/",
			},
			request: func() *http.Request {
				r := bytes.NewReader([]byte("test data"))
				req := httptest.NewRequest(http.MethodGet, "/admin?sdgdsg=111x", r)
				return req
			},
			ctx:                context.Background(),
			expectedStatusCode: http.StatusOK,
			expectedResponse: func() http.ResponseWriter {
				recorder := httptest.NewRecorder()
				response.Ok(recorder, nil, "access granted")
				return recorder
			},
			oauthConfig: func() oidc.Config {
				return mockOidc.NewMockConfig(ctrl)
			},
			oauthVerifier: func() oidc.TokenVerifier {
				m := mockOidc.NewMockTokenVerifier(ctrl)
				return m
			},
		},
		{
			name: "success_with_invalid_access_refresh_tokens",
			rule: &ruleEnt.Rule{
				Name:   "allow_rule",
				Action: ruleEnt.ActionAllow,
				Path:   "/",
			},
			request: func() *http.Request {
				r := bytes.NewReader([]byte("test data"))
				req := httptest.NewRequest(http.MethodGet, "/admin?sdgdsg=111x", r)
				token := &oauth2.Token{
					AccessToken:  "eyJhbGciOiJSUzI1NiIsImtpZCI6IjViYjI2Y2U5Yjg5NGViYTBkNzhjZjQ3Y2YyZjc1YWY4ODk5ZjE3MTcifQ.eyJpc3MiOiJodHRwczovL2FwcC5ndy5rOHMuc2guYWJhbi5pby9kZXgiLCJzdWIiOiJDZzB3TFRNNE5TMHlPREE0T1Mwd0VnUnRiMk5yIiwiYXVkIjoiZXhhbXBsZS1hcHAiLCJleHAiOjE3MjQzMzUwOTQsImlhdCI6MTcyNDI0ODY5NCwibm9uY2UiOiIxYzFmOTU2My1iNWQwLTQwOWItOGRmNS1hYjhhYjdkYzE0ODgiLCJhdF9oYXNoIjoieU5qVFJTTWhZLXFKemJJRTlybG96ZyIsImVtYWlsIjoia2lsZ29yZUBraWxnb3JlLnRyb3V0IiwiZW1haWxfdmVyaWZpZWQiOnRydWUsIm5hbWUiOiJLaWxnb3JlIFRyb3V0In0.k-ajsvzlJfuwGkBy9lW8OgRtlzH0PXpU2HgrU4WaVssEdTc3pQWtw50Aqy5s-HWKUVxMwvLgU1qznTkakglNFAic-vRQzPf0cGQBwW7ojTaz-V9oktUhYIj5BXS001BvLqYGEWkKXS07IL2FYh0WEm7m6xDeOsN_5gwUf3NDpXyIiVfVWXEBRiKjsHK6nVA7YggzKF6K4PwK0BdzSy58AT0Q33zdoOPztb2fWvfjBQWZZFSog3Z0rvr0rurIWWlhSGveZcqkwZj9FeZZlIkGK1VDiQhmkha-lzeomgkh8X07Z8wxXS0tYVJ3S9E7bxMOkb7T-P4LzOo3Oj9GNxP7Zw",
					RefreshToken: "Chl4NjZxa2V4emJhMzJkNmdwZzZ3Y2tycHpoEhlydHp5dTc1d2ZxM29ncGxpcHZwYmtudTRy",
					TokenType:    "bearer",
					Expiry:       time.Now().Add(-1 * time.Hour),
				}
				oauth2Claims, err := json.Marshal(token)
				if err != nil {
					log.Fatal(err)
				}
				tokenJSON, err := encryptor.Encrypt(oauth2Claims)
				if err != nil {
					log.Fatal(err)
				}
				req.AddCookie(&http.Cookie{
					Name:     "auth",
					Value:    string(tokenJSON),
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
					Path:     "/",
				})
				return req
			},
			ctx:                context.Background(),
			expectedStatusCode: http.StatusOK,
			expectedResponse: func() http.ResponseWriter {
				recorder := httptest.NewRecorder()
				http.SetCookie(recorder, &http.Cookie{
					Name:     "auth",
					Value:    "",
					Path:     "/",
					Expires:  time.Unix(0, 0),
					HttpOnly: true,
				})
				response.Ok(recorder, nil, "access granted")
				return recorder
			},
			oauthConfig: func() oidc.Config {
				tokenSource := mockOidc.NewMockTokenSource(ctrl)
				tokenSource.EXPECT().Token().Return(nil, errors.New("invalid refresh token"))
				config := mockOidc.NewMockConfig(ctrl)
				config.EXPECT().TokenSource(gomock.Any(), gomock.Any()).Return(tokenSource)
				return config
			},
			oauthVerifier: func() oidc.TokenVerifier {
				m := mockOidc.NewMockTokenVerifier(ctrl)
				return m
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			oauthConfig := test.oauthConfig()
			oauthVerifier := test.oauthVerifier()
			recorder := httptest.NewRecorder()
			rh := NewRuleHandler(test.rule, oauthConfig, oauthVerifier,
				slog.New(slog.NewTextHandler(os.Stdin, nil)), encryptor)
			rh.Handle(recorder, test.request())

			// Get result and the expected result
			result := recorder.Result()
			expectedResult := test.expectedResponse().(*httptest.ResponseRecorder).Result()

			// Compare status codes
			if result.StatusCode != test.expectedStatusCode {
				t.Errorf("response status code:%d is not match expected value:%d",
					result.StatusCode, test.expectedStatusCode)
			}

			// Compare cookies
			resultCookies := result.Cookies()
			expectedCookies := expectedResult.Cookies()

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

			// Compare headers
			for k, expectedValues := range expectedResult.Header {
				// skip cookies we check them separately
				if k == "Set-Cookie" {
					continue
				}
				if resultValues, ok := result.Header[k]; ok {
					if !isStringSliceEqual(resultValues, expectedValues) {
						t.Error("Header:", k, "Mismatch in result header value(s). Got: ", resultValues, "expected: ", expectedValues)
						return
					}
				} else {
					t.Errorf("Expected header %s is missing from the result", k)
					return
				}
			}
		})
	}
}

func TestRule_Handle_Deny_Rule(t *testing.T) {
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
		rule               *ruleEnt.Rule
		expectedResponse   func() http.ResponseWriter
		oauthConfig        func() oidc.Config
		oauthVerifier      func() oidc.TokenVerifier
	}{
		{
			name: "success_right_cookie_with_headers",
			rule: &ruleEnt.Rule{
				Name:   "deny_rule",
				Action: ruleEnt.ActionDeny,
				Path:   "/",
			},
			request: func() *http.Request {
				r := bytes.NewReader([]byte("test data"))
				req := httptest.NewRequest(http.MethodGet, "/admin?sdgdsg=111x", r)
				token := &oauth2.Token{
					AccessToken:  "eyJhbGciOiJSUzI1NiIsImtpZCI6ImUzMmE2Y2UyNDgyMDhlOGFkMTVjNjc1N2ZlNjFiYjc5ZGZmNzZhYmMifQ.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjU1NTYvZGV4Iiwic3ViIjoiQ2lSaU9EZzBOemRqTmkxaU5HVm1MVFEzTnprdFlqVTBZaTB3T0RWalpEVTJNVEF3TVRVU0QyRjFkR2h2Y21OdmJtNWxZM1J2Y2ciLCJhdWQiOiJleGFtcGxlLWFwcCIsImV4cCI6MTcyNDc1MjgwNSwiaWF0IjoxNzI0NjY2NDA1LCJub25jZSI6Ild2a0JtRENNM3AzSlRkMjBsN0hYbmciLCJhdF9oYXNoIjoiS21YZG1FMW9jVk94elEwY2l3QThMQSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJncm91cHMiOlsidXNlciJdLCJuYW1lIjoiOTg5MTIwNDEwMTc2MSIsImZlZGVyYXRlZF9jbGFpbXMiOnsiY29ubmVjdG9yX2lkIjoiYXV0aG9yY29ubmVjdG9yIiwidXNlcl9pZCI6ImI4ODQ3N2M2LWI0ZWYtNDc3OS1iNTRiLTA4NWNkNTYxMDAxNSJ9fQ.QIm6ghZKm8mzrZQZTxhw25IUWZpAmhWHqS-ZhqkDdEB_xXslhfsREhice3A6wQereaFtsXMA82DX08-ZLvViU25fvB1QsL2piR5dg8bakXUhkHPGxOoGT4k4ocyo-isZdPMF0SrkJbpgLWiu9M4WCz9-9nXo-oOeSjs_6_CRDjj5bQOGVY1gpu8ah4yfDUT-xTe-bv7mu0h7pDfEGax7HPP4JPY-0UCNkNtSpVItPdOnTR_hcVWXW3Hl9MT0M9qeLTAHPkdAzaGlSKNi7EutW6Rj4hVcb5wiX08UzhRM_7nwpZjYrX7W3CszehNxnCSQFVjA9lS1Zq_XhhBhJBGe9g",
					RefreshToken: "ChlpMm1wbnMyazJ4bDZsYjJta3ppYmpkbWlsEhl0MnNlZ3d5cXh5ejVvZ3Jrc3llZGxncGJl",
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
				req.AddCookie(&http.Cookie{
					Name:     "auth",
					Value:    string(tokenJSON),
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
					Path:     "/",
				})

				return req
			},
			ctx:                context.Background(),
			expectedStatusCode: http.StatusOK,
			expectedResponse: func() http.ResponseWriter {
				recorder := httptest.NewRecorder()
				recorder.Header().Set("X-AUTH-NAME", "9891204101761")
				recorder.Header().Set("X-AUTH-VERIFIED", "true")
				recorder.Header().Set("X-AUTH-ACCESS-TOKEN", "eyJhbGciOiJSUzI1NiIsImtpZCI6ImUzMmE2Y2UyNDgyMDhlOGFkMTVjNjc1N2ZlNjFiYjc5ZGZmNzZhYmMifQ.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjU1NTYvZGV4Iiwic3ViIjoiQ2lSaU9EZzBOemRqTmkxaU5HVm1MVFEzTnprdFlqVTBZaTB3T0RWalpEVTJNVEF3TVRVU0QyRjFkR2h2Y21OdmJtNWxZM1J2Y2ciLCJhdWQiOiJleGFtcGxlLWFwcCIsImV4cCI6MTcyNDc1MjgwNSwiaWF0IjoxNzI0NjY2NDA1LCJub25jZSI6Ild2a0JtRENNM3AzSlRkMjBsN0hYbmciLCJhdF9oYXNoIjoiS21YZG1FMW9jVk94elEwY2l3QThMQSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJncm91cHMiOlsidXNlciJdLCJuYW1lIjoiOTg5MTIwNDEwMTc2MSIsImZlZGVyYXRlZF9jbGFpbXMiOnsiY29ubmVjdG9yX2lkIjoiYXV0aG9yY29ubmVjdG9yIiwidXNlcl9pZCI6ImI4ODQ3N2M2LWI0ZWYtNDc3OS1iNTRiLTA4NWNkNTYxMDAxNSJ9fQ.QIm6ghZKm8mzrZQZTxhw25IUWZpAmhWHqS-ZhqkDdEB_xXslhfsREhice3A6wQereaFtsXMA82DX08-ZLvViU25fvB1QsL2piR5dg8bakXUhkHPGxOoGT4k4ocyo-isZdPMF0SrkJbpgLWiu9M4WCz9-9nXo-oOeSjs_6_CRDjj5bQOGVY1gpu8ah4yfDUT-xTe-bv7mu0h7pDfEGax7HPP4JPY-0UCNkNtSpVItPdOnTR_hcVWXW3Hl9MT0M9qeLTAHPkdAzaGlSKNi7EutW6Rj4hVcb5wiX08UzhRM_7nwpZjYrX7W3CszehNxnCSQFVjA9lS1Zq_XhhBhJBGe9g")
				recorder.Header().Set("X-AUTH-GROUPS", "user")
				recorder.Header().Set("X-AUTH-SUB", "CiRiODg0NzdjNi1iNGVmLTQ3NzktYjU0Yi0wODVjZDU2MTAwMTUSD2F1dGhvcmNvbm5lY3Rvcg")
				response.Ok(recorder, nil, "access granted")
				return recorder
			},
			oauthConfig: func() oidc.Config {
				return mockOidc.NewMockConfig(ctrl)
			},
			oauthVerifier: func() oidc.TokenVerifier {
				m := mockOidc.NewMockTokenVerifier(ctrl)
				m.EXPECT().Verify(gomock.Any(), gomock.Any())
				return m
			},
		},
		{
			name: "success_expired_access_token_with_headers",
			rule: &ruleEnt.Rule{
				Name:   "deny_rule",
				Action: ruleEnt.ActionAllow,
				Path:   "/",
			},
			request: func() *http.Request {
				r := bytes.NewReader([]byte("test data"))
				req := httptest.NewRequest(http.MethodGet, "/admin?sdgdsg=111x", r)
				token := &oauth2.Token{
					AccessToken:  "eyJhbGciOiJSUzI1NiIsImtpZCI6IjViYjI2Y2U5Yjg5NGViYTBkNzhjZjQ3Y2YyZjc1YWY4ODk5ZjE3MTcifQ.eyJpc3MiOiJodHRwczovL2FwcC5ndy5rOHMuc2guYWJhbi5pby9kZXgiLCJzdWIiOiJDZzB3TFRNNE5TMHlPREE0T1Mwd0VnUnRiMk5yIiwiYXVkIjoiZXhhbXBsZS1hcHAiLCJleHAiOjE3MjQzMzUwOTQsImlhdCI6MTcyNDI0ODY5NCwibm9uY2UiOiIxYzFmOTU2My1iNWQwLTQwOWItOGRmNS1hYjhhYjdkYzE0ODgiLCJhdF9oYXNoIjoieU5qVFJTTWhZLXFKemJJRTlybG96ZyIsImVtYWlsIjoia2lsZ29yZUBraWxnb3JlLnRyb3V0IiwiZW1haWxfdmVyaWZpZWQiOnRydWUsIm5hbWUiOiJLaWxnb3JlIFRyb3V0In0.k-ajsvzlJfuwGkBy9lW8OgRtlzH0PXpU2HgrU4WaVssEdTc3pQWtw50Aqy5s-HWKUVxMwvLgU1qznTkakglNFAic-vRQzPf0cGQBwW7ojTaz-V9oktUhYIj5BXS001BvLqYGEWkKXS07IL2FYh0WEm7m6xDeOsN_5gwUf3NDpXyIiVfVWXEBRiKjsHK6nVA7YggzKF6K4PwK0BdzSy58AT0Q33zdoOPztb2fWvfjBQWZZFSog3Z0rvr0rurIWWlhSGveZcqkwZj9FeZZlIkGK1VDiQhmkha-lzeomgkh8X07Z8wxXS0tYVJ3S9E7bxMOkb7T-P4LzOo3Oj9GNxP7Zw",
					RefreshToken: "Chl4NjZxa2V4emJhMzJkNmdwZzZ3Y2tycHpoEhlydHp5dTc1d2ZxM29ncGxpcHZwYmtudTRy",
					TokenType:    "bearer",
					Expiry:       time.Now().Add(-1 * time.Hour),
				}
				oauth2Claims, err := json.Marshal(token)
				if err != nil {
					log.Fatal(err)
				}
				tokenJSON, err := encryptor.Encrypt(oauth2Claims)
				if err != nil {
					log.Fatal(err)
				}
				req.AddCookie(&http.Cookie{
					Name:     "auth",
					Value:    string(tokenJSON),
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
				})
				return req
			},
			ctx:                context.Background(),
			expectedStatusCode: http.StatusOK,
			expectedResponse: func() http.ResponseWriter {
				recorder := httptest.NewRecorder()
				recorder.Header().Set("X-AUTH-NAME", "9891204101761")
				recorder.Header().Set("X-AUTH-VERIFIED", "true")
				recorder.Header().Set("X-AUTH-ACCESS-TOKEN", "eyJhbGciOiJSUzI1NiIsImtpZCI6ImUzMmE2Y2UyNDgyMDhlOGFkMTVjNjc1N2ZlNjFiYjc5ZGZmNzZhYmMifQ.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjU1NTYvZGV4Iiwic3ViIjoiQ2lSaU9EZzBOemRqTmkxaU5HVm1MVFEzTnprdFlqVTBZaTB3T0RWalpEVTJNVEF3TVRVU0QyRjFkR2h2Y21OdmJtNWxZM1J2Y2ciLCJhdWQiOiJleGFtcGxlLWFwcCIsImV4cCI6MTcyNDc1MjgwNSwiaWF0IjoxNzI0NjY2NDA1LCJub25jZSI6Ild2a0JtRENNM3AzSlRkMjBsN0hYbmciLCJhdF9oYXNoIjoiS21YZG1FMW9jVk94elEwY2l3QThMQSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJncm91cHMiOlsidXNlciJdLCJuYW1lIjoiOTg5MTIwNDEwMTc2MSIsImZlZGVyYXRlZF9jbGFpbXMiOnsiY29ubmVjdG9yX2lkIjoiYXV0aG9yY29ubmVjdG9yIiwidXNlcl9pZCI6ImI4ODQ3N2M2LWI0ZWYtNDc3OS1iNTRiLTA4NWNkNTYxMDAxNSJ9fQ.QIm6ghZKm8mzrZQZTxhw25IUWZpAmhWHqS-ZhqkDdEB_xXslhfsREhice3A6wQereaFtsXMA82DX08-ZLvViU25fvB1QsL2piR5dg8bakXUhkHPGxOoGT4k4ocyo-isZdPMF0SrkJbpgLWiu9M4WCz9-9nXo-oOeSjs_6_CRDjj5bQOGVY1gpu8ah4yfDUT-xTe-bv7mu0h7pDfEGax7HPP4JPY-0UCNkNtSpVItPdOnTR_hcVWXW3Hl9MT0M9qeLTAHPkdAzaGlSKNi7EutW6Rj4hVcb5wiX08UzhRM_7nwpZjYrX7W3CszehNxnCSQFVjA9lS1Zq_XhhBhJBGe9g")
				recorder.Header().Set("X-AUTH-GROUPS", "user")
				recorder.Header().Set("X-AUTH-SUB", "CiRiODg0NzdjNi1iNGVmLTQ3NzktYjU0Yi0wODVjZDU2MTAwMTUSD2F1dGhvcmNvbm5lY3Rvcg")
				token := &oauth2.Token{
					AccessToken:  "eyJhbGciOiJSUzI1NiIsImtpZCI6ImUzMmE2Y2UyNDgyMDhlOGFkMTVjNjc1N2ZlNjFiYjc5ZGZmNzZhYmMifQ.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjU1NTYvZGV4Iiwic3ViIjoiQ2lSaU9EZzBOemRqTmkxaU5HVm1MVFEzTnprdFlqVTBZaTB3T0RWalpEVTJNVEF3TVRVU0QyRjFkR2h2Y21OdmJtNWxZM1J2Y2ciLCJhdWQiOiJleGFtcGxlLWFwcCIsImV4cCI6MTcyNDc1MjgwNSwiaWF0IjoxNzI0NjY2NDA1LCJub25jZSI6Ild2a0JtRENNM3AzSlRkMjBsN0hYbmciLCJhdF9oYXNoIjoiS21YZG1FMW9jVk94elEwY2l3QThMQSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJncm91cHMiOlsidXNlciJdLCJuYW1lIjoiOTg5MTIwNDEwMTc2MSIsImZlZGVyYXRlZF9jbGFpbXMiOnsiY29ubmVjdG9yX2lkIjoiYXV0aG9yY29ubmVjdG9yIiwidXNlcl9pZCI6ImI4ODQ3N2M2LWI0ZWYtNDc3OS1iNTRiLTA4NWNkNTYxMDAxNSJ9fQ.QIm6ghZKm8mzrZQZTxhw25IUWZpAmhWHqS-ZhqkDdEB_xXslhfsREhice3A6wQereaFtsXMA82DX08-ZLvViU25fvB1QsL2piR5dg8bakXUhkHPGxOoGT4k4ocyo-isZdPMF0SrkJbpgLWiu9M4WCz9-9nXo-oOeSjs_6_CRDjj5bQOGVY1gpu8ah4yfDUT-xTe-bv7mu0h7pDfEGax7HPP4JPY-0UCNkNtSpVItPdOnTR_hcVWXW3Hl9MT0M9qeLTAHPkdAzaGlSKNi7EutW6Rj4hVcb5wiX08UzhRM_7nwpZjYrX7W3CszehNxnCSQFVjA9lS1Zq_XhhBhJBGe9g",
					RefreshToken: "Chl4NjZxa2V4emJhMzJkNmdwZzZ3Y2tycHpoEhlydHp5dTc1d2ZxM29ncGxpcHZwYmtudTRy",
					TokenType:    "bearer",
					Expiry:       time.Now().Add(1 * time.Hour),
				}
				oauth2Claims, err := json.Marshal(token)
				if err != nil {
					t.Error(err.Error())
				}

				tokenJSON, err := encryptor.Encrypt(oauth2Claims)
				if err != nil {
					t.Error(err.Error())
				}

				http.SetCookie(recorder, &http.Cookie{
					Name:     "auth",
					Value:    string(tokenJSON),
					Secure:   false,
					HttpOnly: true,
					Path:     "/",
				})

				response.Ok(recorder, nil, "access granted")
				return recorder
			},
			oauthConfig: func() oidc.Config {
				token := &oauth2.Token{
					AccessToken:  "eyJhbGciOiJSUzI1NiIsImtpZCI6ImUzMmE2Y2UyNDgyMDhlOGFkMTVjNjc1N2ZlNjFiYjc5ZGZmNzZhYmMifQ.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjU1NTYvZGV4Iiwic3ViIjoiQ2lSaU9EZzBOemRqTmkxaU5HVm1MVFEzTnprdFlqVTBZaTB3T0RWalpEVTJNVEF3TVRVU0QyRjFkR2h2Y21OdmJtNWxZM1J2Y2ciLCJhdWQiOiJleGFtcGxlLWFwcCIsImV4cCI6MTcyNDc1MjgwNSwiaWF0IjoxNzI0NjY2NDA1LCJub25jZSI6Ild2a0JtRENNM3AzSlRkMjBsN0hYbmciLCJhdF9oYXNoIjoiS21YZG1FMW9jVk94elEwY2l3QThMQSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJncm91cHMiOlsidXNlciJdLCJuYW1lIjoiOTg5MTIwNDEwMTc2MSIsImZlZGVyYXRlZF9jbGFpbXMiOnsiY29ubmVjdG9yX2lkIjoiYXV0aG9yY29ubmVjdG9yIiwidXNlcl9pZCI6ImI4ODQ3N2M2LWI0ZWYtNDc3OS1iNTRiLTA4NWNkNTYxMDAxNSJ9fQ.QIm6ghZKm8mzrZQZTxhw25IUWZpAmhWHqS-ZhqkDdEB_xXslhfsREhice3A6wQereaFtsXMA82DX08-ZLvViU25fvB1QsL2piR5dg8bakXUhkHPGxOoGT4k4ocyo-isZdPMF0SrkJbpgLWiu9M4WCz9-9nXo-oOeSjs_6_CRDjj5bQOGVY1gpu8ah4yfDUT-xTe-bv7mu0h7pDfEGax7HPP4JPY-0UCNkNtSpVItPdOnTR_hcVWXW3Hl9MT0M9qeLTAHPkdAzaGlSKNi7EutW6Rj4hVcb5wiX08UzhRM_7nwpZjYrX7W3CszehNxnCSQFVjA9lS1Zq_XhhBhJBGe9g",
					RefreshToken: "Chl4NjZxa2V4emJhMzJkNmdwZzZ3Y2tycHpoEhlydHp5dTc1d2ZxM29ncGxpcHZwYmtudTRy",
					TokenType:    "bearer",
					Expiry:       time.Now().Add(1 * time.Hour),
				}
				tokenSource := mockOidc.NewMockTokenSource(ctrl)
				tokenSource.EXPECT().Token().Return(token, nil)
				config := mockOidc.NewMockConfig(ctrl)
				config.EXPECT().TokenSource(gomock.Any(), gomock.Any()).Return(tokenSource)
				return config
			},
			oauthVerifier: func() oidc.TokenVerifier {
				m := mockOidc.NewMockTokenVerifier(ctrl)
				m.EXPECT().Verify(gomock.Any(), gomock.Any())
				return m
			},
		},
		{
			name: "forbidden_without_cookie",
			rule: &ruleEnt.Rule{
				Name:   "deny_rule",
				Action: ruleEnt.ActionDeny,
				Path:   "/",
			},
			request: func() *http.Request {
				r := bytes.NewReader([]byte("test data"))
				req := httptest.NewRequest(http.MethodGet, "/admin?sdgdsg=111x", r)
				return req
			},
			ctx:                context.Background(),
			expectedStatusCode: http.StatusForbidden,
			expectedResponse: func() http.ResponseWriter {
				recorder := httptest.NewRecorder()
				response.Custom(recorder, http.StatusForbidden, nil, "access denied")
				return recorder
			},
			oauthConfig: func() oidc.Config {
				return mockOidc.NewMockConfig(ctrl)
			},
			oauthVerifier: func() oidc.TokenVerifier {
				m := mockOidc.NewMockTokenVerifier(ctrl)
				return m
			},
		},
		{
			name: "forbidden_with_invalid_access_refresh_tokens",
			rule: &ruleEnt.Rule{
				Name:   "allow_rule",
				Action: ruleEnt.ActionDeny,
				Path:   "/",
			},
			request: func() *http.Request {
				r := bytes.NewReader([]byte("test data"))
				req := httptest.NewRequest(http.MethodGet, "/admin?sdgdsg=111x", r)
				token := &oauth2.Token{
					AccessToken:  "eyJhbGciOiJSUzI1NiIsImtpZCI6IjViYjI2Y2U5Yjg5NGViYTBkNzhjZjQ3Y2YyZjc1YWY4ODk5ZjE3MTcifQ.eyJpc3MiOiJodHRwczovL2FwcC5ndy5rOHMuc2guYWJhbi5pby9kZXgiLCJzdWIiOiJDZzB3TFRNNE5TMHlPREE0T1Mwd0VnUnRiMk5yIiwiYXVkIjoiZXhhbXBsZS1hcHAiLCJleHAiOjE3MjQzMzUwOTQsImlhdCI6MTcyNDI0ODY5NCwibm9uY2UiOiIxYzFmOTU2My1iNWQwLTQwOWItOGRmNS1hYjhhYjdkYzE0ODgiLCJhdF9oYXNoIjoieU5qVFJTTWhZLXFKemJJRTlybG96ZyIsImVtYWlsIjoia2lsZ29yZUBraWxnb3JlLnRyb3V0IiwiZW1haWxfdmVyaWZpZWQiOnRydWUsIm5hbWUiOiJLaWxnb3JlIFRyb3V0In0.k-ajsvzlJfuwGkBy9lW8OgRtlzH0PXpU2HgrU4WaVssEdTc3pQWtw50Aqy5s-HWKUVxMwvLgU1qznTkakglNFAic-vRQzPf0cGQBwW7ojTaz-V9oktUhYIj5BXS001BvLqYGEWkKXS07IL2FYh0WEm7m6xDeOsN_5gwUf3NDpXyIiVfVWXEBRiKjsHK6nVA7YggzKF6K4PwK0BdzSy58AT0Q33zdoOPztb2fWvfjBQWZZFSog3Z0rvr0rurIWWlhSGveZcqkwZj9FeZZlIkGK1VDiQhmkha-lzeomgkh8X07Z8wxXS0tYVJ3S9E7bxMOkb7T-P4LzOo3Oj9GNxP7Zw",
					RefreshToken: "Chl4NjZxa2V4emJhMzJkNmdwZzZ3Y2tycHpoEhlydHp5dTc1d2ZxM29ncGxpcHZwYmtudTRy",
					TokenType:    "bearer",
					Expiry:       time.Now().Add(-1 * time.Hour),
				}
				oauth2Claims, err := json.Marshal(token)
				if err != nil {
					log.Fatal(err)
				}
				tokenJSON, err := encryptor.Encrypt(oauth2Claims)
				if err != nil {
					log.Fatal(err)
				}
				req.AddCookie(&http.Cookie{
					Name:     "auth",
					Value:    string(tokenJSON),
					MaxAge:   int(time.Hour.Seconds()),
					Secure:   false,
					HttpOnly: true,
					Path:     "/",
				})
				return req
			},
			ctx:                context.Background(),
			expectedStatusCode: http.StatusForbidden,
			expectedResponse: func() http.ResponseWriter {
				recorder := httptest.NewRecorder()
				http.SetCookie(recorder, &http.Cookie{
					Name:     "auth",
					Value:    "",
					Path:     "/",
					Expires:  time.Unix(0, 0),
					HttpOnly: true,
				})
				response.Custom(recorder, http.StatusForbidden, nil, "access denied")
				return recorder
			},
			oauthConfig: func() oidc.Config {
				tokenSource := mockOidc.NewMockTokenSource(ctrl)
				tokenSource.EXPECT().Token().Return(nil, errors.New("invalid refresh token"))
				config := mockOidc.NewMockConfig(ctrl)
				config.EXPECT().TokenSource(gomock.Any(), gomock.Any()).Return(tokenSource)
				return config
			},
			oauthVerifier: func() oidc.TokenVerifier {
				m := mockOidc.NewMockTokenVerifier(ctrl)
				return m
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			oauthConfig := test.oauthConfig()
			oauthVerifier := test.oauthVerifier()
			recorder := httptest.NewRecorder()
			rh := NewRuleHandler(test.rule, oauthConfig, oauthVerifier,
				slog.New(slog.NewTextHandler(os.Stdin, nil)), encryptor)
			rh.Handle(recorder, test.request())

			// Get result and the expected result
			result := recorder.Result()
			expectedResult := test.expectedResponse().(*httptest.ResponseRecorder).Result()

			// Compare status codes
			if result.StatusCode != test.expectedStatusCode {
				t.Errorf("response status code:%d is not match expected value:%d",
					result.StatusCode, test.expectedStatusCode)
			}

			// Compare cookies
			resultCookies := result.Cookies()
			expectedCookies := expectedResult.Cookies()

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

			// Compare headers
			for k, expectedValues := range expectedResult.Header {
				// skip cookies we check them separately
				if k == "Set-Cookie" {
					continue
				}
				if resultValues, ok := result.Header[k]; ok {
					if !isStringSliceEqual(resultValues, expectedValues) {
						t.Error("Header:", k, "Mismatch in result header value(s). Got: ", resultValues, "expected: ", expectedValues)
						return
					}
				} else {
					t.Errorf("Expected header %s is missing from the result", k)
					return
				}
			}
		})
	}
}

// Add this helper function to compare String slices
func isStringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Strings(a)
	sort.Strings(b)
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
