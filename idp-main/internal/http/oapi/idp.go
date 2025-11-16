package oapi

import (
	"application/internal/http/dto"
	"application/internal/http/response"
	"log/slog"
	"net/http"

	"github.com/swaggest/openapi-go"
)

// IDPOpenApi represents the OpenAPI specification for IDP endpoints.
type IDPOpenApi struct {
	logger *slog.Logger
}

// UserID represents the path parameter for user operations.
type UserID struct {
	ID string `json:"id" path:"id"`
}

// NewIDPOpenApi creates a new instance of IDPOpenApi.
func NewIDPOpenApi(logger *slog.Logger) *IDPOpenApi {
	return &IDPOpenApi{
		logger: logger.With("layer", "IDPOpenApi"),
	}
}

// ListUsers describes the OpenAPI specification for listing users.
func (s *IDPOpenApi) ListUsers(op openapi.OperationContext) {
	op.SetSummary("List Users")
	op.SetDescription("Fetch a list of all users with basic information.")

	op.AddRespStructure(new(dto.UserListResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
}

// GetUser describes the OpenAPI specification for fetching a single user's details.
func (s *IDPOpenApi) GetUser(op openapi.OperationContext) {
	op.SetSummary("Get User")
	op.SetDescription("Fetch detailed information for a specific user by their unique ID.")

	op.AddReqStructure(new(UserID), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(dto.UserResponse), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusNotFound))
}

// CreateUser describes the OpenAPI specification for creating a new user.
func (s *IDPOpenApi) CreateUser(op openapi.OperationContext) {
	op.SetSummary("Create User")
	op.SetDescription("Create a new user with the provided information.")

	op.AddReqStructure(new(dto.UserRequest), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(dto.UserResponse), openapi.WithHTTPStatus(http.StatusCreated))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusConflict))
}

// UpdateUser describes the OpenAPI specification for updating an existing user.
func (s *IDPOpenApi) UpdateUser(op openapi.OperationContext) {
	op.SetSummary("Update User")
	op.SetDescription("Update the details of an existing user by their unique ID.")

	op.AddReqStructure(new(UserID), openapi.WithContentType("application/json"))
	op.AddReqStructure(new(dto.UpdateUserRequest), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusNotFound))
}

// DeleteUser describes the OpenAPI specification for deleting a user.
func (s *IDPOpenApi) DeleteUser(op openapi.OperationContext) {
	op.SetSummary("Delete User")
	op.SetDescription("Delete an existing user by their unique ID.")

	op.AddReqStructure(new(UserID), openapi.WithContentType("application/json"))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusOK))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusInternalServerError))
	op.AddRespStructure(new(response.Response[string]), openapi.WithHTTPStatus(http.StatusNotFound))
}

// OpenApiSpec registers the IDP OpenAPI endpoints.
func (s *IDPOpenApi) OpenApiSpec(api OAPI) {
	idpTags := WithTags("IDP")

	api.Register("GET", "/api/idp/v1/users", s.ListUsers, idpTags)          // List Users
	api.Register("GET", "/api/idp/v1/users/{id}", s.GetUser, idpTags)       // Get User
	api.Register("POST", "/api/idp/v1/users", s.CreateUser, idpTags)        // Create User
	api.Register("PUT", "/api/idp/v1/users/{id}", s.UpdateUser, idpTags)    // Update User
	api.Register("DELETE", "/api/idp/v1/users/{id}", s.DeleteUser, idpTags) // Delete User
}
