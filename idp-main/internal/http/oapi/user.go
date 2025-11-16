package oapi

import (
	"log/slog"
	"net/http"

	"application/internal/http/dto"
	"application/internal/http/response"

	"github.com/swaggest/openapi-go"
)

type UserOpenApi struct {
	logger *slog.Logger
}

func NewUserOpenApi(logger *slog.Logger) *UserOpenApi {
	return &UserOpenApi{
		logger: logger.With("layer", "UserOpenApi"),
	}
}

type GetUserByIDRequest struct {
	UserID string `path:"user_id" example:"a0c3abdf-500c-4662-ac42-e78990f67431"`
}

type UserRoleRequest struct {
	UserID string `path:"user_id" example:"a0c3abdf-500c-4662-ac42-e78990f67431"`
	RoleID string `path:"role_id" example:"a0c3abdf-500c-4662-ac42-e78990f67431"`
}

func (s *UserOpenApi) CreateUser(op openapi.OperationContext) {
	op.SetSummary("Create User")
	op.AddReqStructure(new(dto.UserRequest), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(response.Response[dto.UserResponse]), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusBadRequest))
}

func (s *UserOpenApi) ListUser(op openapi.OperationContext) {
	op.SetSummary("List Users")
	op.AddRespStructure(new(dto.UserListResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

func (s *UserOpenApi) GetUserByID(op openapi.OperationContext) {
	op.SetSummary("Get User By ID")
	op.AddReqStructure(new(GetUserByIDRequest), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(dto.UserResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusNotFound))
}

func (s *UserOpenApi) DeleteUser(op openapi.OperationContext) {
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

func (s *UserOpenApi) AddUserToRole(op openapi.OperationContext) {
	op.SetSummary("Add User to Role")
	op.AddReqStructure(new(UserRoleRequest), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusCreated))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusBadRequest))
}

func (s *UserOpenApi) DeleteUserFromRole(op openapi.OperationContext) {
	op.SetSummary("Delete User from Role")
	op.AddReqStructure(new(UserRoleRequest), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusBadRequest))
}

func (s *UserOpenApi) VerifyUser(op openapi.OperationContext) {
	op.SetSummary("Verify User")
	op.AddReqStructure(new(dto.VerifyRequest), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(dto.VerifyResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusBadRequest))
}

func (s *UserOpenApi) OpenApiSpec(api OAPI) {
	userTags := WithTags("User")
	api.Register(http.MethodPost, "/api/idp/v1/users", s.CreateUser, userTags)
	api.Register(http.MethodGet, "/api/idp/v1/users", s.ListUser, userTags)
	api.Register(http.MethodGet, "/api/idp/v1/users/{user_id}", s.GetUserByID, userTags)
	api.Register(http.MethodDelete, "/api/idp/v1/users/{user_id}", s.DeleteUser, userTags)
	// api.Register(http.MethodGet, "/api/idp/v1/users/{user_id}/roles", s.GetUserRoles, userTags)
	api.Register(http.MethodPost, "/api/idp/v1/users/{user_id}/roles/{role_id}", s.AddUserToRole, userTags)
	api.Register(http.MethodDelete, "/api/idp/v1/users/{user_id}/roles/{role_id}", s.DeleteUserFromRole, userTags)
}
