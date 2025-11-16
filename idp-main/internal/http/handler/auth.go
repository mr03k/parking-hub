package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	idpbiz "application/internal/biz/idp"

	authbiz "application/internal/biz/auth"
	"application/internal/http/dto"
	"application/internal/http/response"

	"github.com/swaggest/openapi-go"
)

type AuthHandler struct {
	logger *slog.Logger
	uc     authbiz.AutUsecaseInterface
}

// NewAuthHandler
func NewAuthHandler(logger *slog.Logger, uc authbiz.AutUsecaseInterface) *AuthHandler {
	return &AuthHandler{
		logger: logger.With("layer", "AuthHandler"),
		uc:     uc,
	}
}

var (
	_ Handler        = (*AuthHandler)(nil)
	_ OpenApiHandler = (*AuthHandler)(nil)
)

// Login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Login")
	ctx := r.Context()

	vehicle_id := r.PathValue("vehicle_id")
	// get data form request
	login, err := dto.NewDriverLoginFromRequest(r)
	if err != nil {
		logger.Error("error parse request", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, token, err := h.uc.DriverLogin(ctx, login.Msisdn, login.Password, vehicle_id)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, idpbiz.ErrorValidationFailed) || errors.Is(err, idpbiz.ErrorNotFound) {
			response.BadRequest(w, "invalid user credentials")
			return
		}
		logger.Error("error login", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response.Pure(w, http.StatusOK, dto.NewDriverLoginResponse(userID, token))
}

// register mux for auth
func (h *AuthHandler) RegisterMuxRouter(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/vehicles/{vehicle_id}/login", h.Login)
}

// login
func (s *AuthHandler) LoginSpec(op openapi.OperationContext) {
	op.SetSummary("Login")
	op.AddReqStructure(new(VehicalLogin), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(dto.DriverLoginResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusBadRequest))
}

func (s *AuthHandler) OpenApiSpec(api OAPI) {
	userTags := WithTags("Chapar")
	api.Register(http.MethodPost, "/api/v1/vehicles/{vehicle_id}/login", s.LoginSpec, userTags)
}

type AuthOpenApi struct {
	logger *slog.Logger
}

type VehicalLogin struct {
	VehicalID string `path:"vehicle_id"`
	dto.DriverLogin
}
