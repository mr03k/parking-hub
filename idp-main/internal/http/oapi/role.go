package oapi

import (
	"log/slog"
	"net/http"

	"application/internal/http/response"

	"github.com/swaggest/openapi-go"
)

type RoleOpenApi struct {
	logger *slog.Logger
}

type RoleID struct {
	RoleID string `json:"role_id" path:"role_id"`
}

func NewRoleOpenApi(logger *slog.Logger) *RoleOpenApi {
	return &RoleOpenApi{
		logger: logger.With("layer", "RoleOpenApi"),
	}
}

// Create Role
// func (s *RoleOpenApi) CreateRole(op openapi.OperationContext) {
// 	op.SetSummary("Create Role")

// 	op.AddReqStructure(new(dto.RoleRequest), openapi.WithContentType("application/json"))
// 	op.AddRespStructure(new(dto.RoleResponse), openapi.WithHTTPStatus(http.StatusOK))
// 	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
// 	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusBadRequest))
// }

// List Roles
// func (s *RoleOpenApi) ListRoles(op openapi.OperationContext) {
// 	op.SetSummary("List Roles")

// 	op.AddRespStructure(new(dto.RoleListResponse), openapi.WithHTTPStatus(http.StatusOK))
// 	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
// }

// Delete Role
func (s *RoleOpenApi) DeleteRole(op openapi.OperationContext) {
	op.SetSummary("Delete Role")
	op.AddReqStructure(new(RoleID), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

// Get Role By ID
// func (s *RoleOpenApi) GetRole(op openapi.OperationContext) {
// 	op.SetSummary("Get Role By ID")
// 	op.AddReqStructure(new(RoleID), openapi.WithContentType("application/json"))
// 	op.AddRespStructure(new(dto.RoleResponse), openapi.WithHTTPStatus(http.StatusOK))
// 	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
// }

func (s *RoleOpenApi) OpenApiSpec(api OAPI) {
	roleTags := WithTags("Role")
	// api.Register("POST", "/api/idp/v1/roles", s.CreateRole, roleTags)
	// api.Register("GET", "/api/idp/v1/roles", s.ListRoles, roleTags)
	api.Register("DELETE", "/api/idp/v1/roles/{role_id}", s.DeleteRole, roleTags)
	// api.Register("GET", "/api/idp/v1/roles/{role_id}", s.GetRole, roleTags)
}
