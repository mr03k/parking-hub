package handler

import (
	"errors"
	"log/slog"
	"net/http"

	idpbiz "application/internal/biz/idp"
	"application/internal/http/dto"
	"application/internal/http/response"

	"github.com/go-playground/validator/v10"
	"github.com/swaggest/openapi-go"
)

type IDPHandler struct {
	logger    *slog.Logger
	uc        idpbiz.UserUseCaseInterface
	validator *validator.Validate
}

// NewIDPHandler initializes a new IDPHandler.
func NewIDPHandler(logger *slog.Logger, uc idpbiz.UserUseCaseInterface) *IDPHandler {
	return &IDPHandler{
		logger:    logger.With("layer", "IDPHandler"),
		uc:        uc,
		validator: validator.New(),
	}
}

// Ensure IDPHandler implements Handler interface.
var _ Handler = (*IDPHandler)(nil)

// CreateUser creates a new user.
func (h *IDPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "CreateUser")
	ctx := r.Context()

	// Parse request
	reg, err := dto.NewUserCreateFromRequest(r)
	if err != nil {
		logger.Error("error parse request", "error", err)
		response.BadRequest(w, "Invalid request payload")
		return
	}

	// Validate input
	err = reg.Validate(h.validator)
	if err != nil {
		logger.Error("error validate request", "error", err)
		response.GeneralError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create user
	_, user, err := h.uc.CreateUser(ctx, reg.ToEntity())
	if err != nil {
		switch err {
		case idpbiz.ErrorUserExist:
			response.GeneralError(w, http.StatusConflict, "User already exists")
			return
		default:
			logger.Error("error create user", "error", err)
			response.InternalError(w)
			return
		}
	}

	w.Header().Add("content-type", "application/json")

	// Respond
	resp := dto.NewUserCreateResponse(user)
	w.Header().Add("content-type", "application/json")
	response.Ok(w, resp, "")
}

// GetUser fetches a user by ID.
func (h *IDPHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "GetUser")
	ctx := r.Context()

	id := r.PathValue("user_id")

	user, err := h.uc.GetUserByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, idpbiz.ErrorUserNotFount):
			response.NotFound(w)
			return
		default:
			logger.Error("error get user", "error", err)
			response.InternalError(w)
			return
		}
	}
	w.Header().Add("content-type", "application/json")

	response.Pure(w, http.StatusOK, user)
}

// ListUser fetches all users.
func (h *IDPHandler) ListUser(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "ListUsers")
	ctx := r.Context()

	users, err := h.uc.ListUser(ctx)
	if err != nil {
		logger.Error("error list users", "error", err)
		response.InternalError(w)
		return
	}

	resp := dto.NewUserListResponse(users)

	w.Header().Add("content-type", "application/json")
	response.Pure(w, http.StatusOK, resp)
}

// DeleteUser deletes a user by ID.
func (h *IDPHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "DeleteUser")
	ctx := r.Context()

	id := r.PathValue("user_id")

	err := h.uc.DeleteUser(ctx, id)
	if err != nil {
		switch err {
		case idpbiz.ErrorUserNotFount:
			response.GeneralError(w, http.StatusNotFound, "User not found")
			return
		default:
			logger.Error("error delete user", "error", err)
			response.InternalError(w)
			return
		}
	}
	w.Header().Add("content-type", "application/json")

	w.WriteHeader(http.StatusOK)
}

// UpdateUser updates an existing user.
func (h *IDPHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "UpdateUser")
	ctx := r.Context()

	id := r.PathValue("user_id")

	// Parse request
	userUpdate, err := dto.NewUserUpdateFromRequest(r)
	if err != nil {
		logger.Error("error parse request", "error", err)
		response.BadRequest(w, "Invalid request payload")
		return
	}

	// Validate input
	err = userUpdate.Validate(h.validator)
	if err != nil {
		response.GeneralError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.uc.UpdateUser(ctx, id, userUpdate.ToEntity())
	if err != nil {
		switch err {
		case idpbiz.ErrorUserNotFount:
			response.GeneralError(w, http.StatusNotFound, "User not found")
			return
		default:
			logger.Error("error update user", "error", err)
			response.InternalError(w)
			return
		}
	}
	w.Header().Add("content-type", "application/json")

	w.WriteHeader(http.StatusOK)
}

// RegisterMuxRouter registers all routes for the handler.
func (h *IDPHandler) RegisterMuxRouter(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/idp/v1/users", h.ListUser)
	mux.HandleFunc("GET /api/idp/v1/users/{user_id}", h.GetUser)
	mux.HandleFunc("POST /api/idp/v1/users", h.CreateUser)
	mux.HandleFunc("PUT /api/idp/v1/users/{user_id}", h.UpdateUser)
	mux.HandleFunc("DELETE /api/idp/v1/users/{user_id}", h.DeleteUser)
}

type GetUserByIDRequest struct {
	UserID string `path:"user_id" example:"a0c3abdf-500c-4662-ac42-e78990f67431"`
}

type UserRoleRequest struct {
	UserID string `path:"user_id" example:"a0c3abdf-500c-4662-ac42-e78990f67431"`
	RoleID string `path:"role_id" example:"a0c3abdf-500c-4662-ac42-e78990f67431"`
}

func (s *IDPHandler) CreateUserOAPI(op openapi.OperationContext) {
	op.SetSummary("Create User")
	op.AddReqStructure(new(dto.UserRequest), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(response.Response[dto.UserResponse]), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusBadRequest))
}

func (s *IDPHandler) ListUserOAPI(op openapi.OperationContext) {
	op.SetSummary("List Users")
	op.AddRespStructure(new(dto.UserListResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

func (s *IDPHandler) GetUserByIDOAPI(op openapi.OperationContext) {
	op.SetSummary("Get User By ID")
	op.AddReqStructure(new(GetUserByIDRequest), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(dto.UserResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusNotFound))
}

func (s *IDPHandler) DeleteUserOAPI(op openapi.OperationContext) {
	op.SetSummary("Delete User By ID")
	op.AddReqStructure(new(GetUserByIDRequest), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusNotFound))
}

// func (s *UserOpenApi) GetUserRoles(op openapi.OperationContext) {
// 	op.SetSummary("Get User Roles")
// 	op.AddReqStructure(new(GetUserByIDRequest), openapi.WithHTTPStatus(http.StatusOK))
// 	op.AddRespStructure(new(dto.RoleListResponse), openapi.WithHTTPStatus(http.StatusOK))
// 	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
// }

func (s *IDPHandler) AddUserToRoleOAPI(op openapi.OperationContext) {
	op.SetSummary("Add User to Role")
	op.AddReqStructure(new(UserRoleRequest), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusCreated))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusBadRequest))
}

func (s *IDPHandler) DeleteUserFromRoleOAPI(op openapi.OperationContext) {
	op.SetSummary("Delete User from Role")
	op.AddReqStructure(new(UserRoleRequest), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusBadRequest))
}

func (s *IDPHandler) VerifyUserOAPI(op openapi.OperationContext) {
	op.SetSummary("Verify User")
	op.AddReqStructure(new(dto.VerifyRequest), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(dto.VerifyResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusBadRequest))
}

func (s *IDPHandler) OpenApiSpec(api OAPI) {
	userTags := WithTags("User")
	api.Register(http.MethodPost, "/api/idp/v1/users", s.CreateUserOAPI, userTags)
	api.Register(http.MethodGet, "/api/idp/v1/users", s.ListUserOAPI, userTags)
	api.Register(http.MethodGet, "/api/idp/v1/users/{user_id}", s.GetUserByIDOAPI, userTags)
	api.Register(http.MethodDelete, "/api/idp/v1/users/{user_id}", s.DeleteUserOAPI, userTags)
	// api.Register(http.MethodGet, "/api/idp/v1/users/{user_id}/roles", s.GetUserRoles, userTags)
	api.Register(http.MethodPost, "/api/idp/v1/users/{user_id}/roles/{role_id}", s.AddUserToRoleOAPI, userTags)
	api.Register(http.MethodDelete, "/api/idp/v1/users/{user_id}/roles/{role_id}", s.DeleteUserFromRoleOAPI, userTags)
}
