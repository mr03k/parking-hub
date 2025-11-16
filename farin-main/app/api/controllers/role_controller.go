package controller

import (
	"context"
	"errors"
	"farin/app/api/response"
	"farin/domain/dto"
	"farin/domain/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"strconv"
)

type RoleController struct {
	roleService *service.RoleService
	logger      *slog.Logger
}

func NewRoleController(logger *slog.Logger, rs *service.RoleService) *RoleController {
	return &RoleController{
		roleService: rs,
		logger:      logger.With("layer", "RoleController"),
	}
}

// CreateRole godoc
// @Summary      Create a new role
// @Description  Create a new role by providing role details
// @Tags         roles
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        role  body      dto.RoleRequest  true  "Role Data"
// @Success      201   {object}  response.Response[dto.RoleResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /roles [post]
func (rc *RoleController) CreateRole(c *gin.Context) {
	lg := rc.logger.With("method", "CreateRole")
	var roleRequest dto.RoleRequest

	if err := c.ShouldBindJSON(&roleRequest); err != nil {
		lg.Error("failed to bind role data", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	createdRole, err := rc.roleService.CreateRole(context.Background(), roleRequest.ToEntity(), &roleRequest)
	if err != nil {
		lg.Error("failed to create role", "error", err)
		response.InternalError(c)
		return
	}
	roleResponse := dto.RoleResponse{}
	roleResponse.FromEntity(createdRole)

	response.Created(c, roleResponse)
}

// ListRoles godoc
// @Summary      List roles
// @Description  Retrieve a list of roles with optional filters and pagination
// @Tags         roles
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        sortField  query     string  false  "Field to sort by"  default(created_at)
// @Param        sortOrder  query     string  false  "Sort order (asc/desc)"  default(asc)
// @Param        page       query     int     false  "Page number"  default(1)
// @Param        pageSize   query     int     false  "Page size"  default(10)
// @Success      200        {object}  response.Response[dto.RoleListResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /roles [get]
func (rc *RoleController) ListRoles(c *gin.Context) {
	lg := rc.logger.With("method", "ListRoles")
	filters := make(map[string]interface{})
	sortField := c.DefaultQuery("sortField", "created_at")
	sortOrder := c.DefaultQuery("sortOrder", "asc")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		response.BadRequest(c, "invalid page param")
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		response.BadRequest(c, "invalid page size param")
		return
	}

	roles, total, err := rc.roleService.ListRoles(context.Background(), filters, sortField, sortOrder, page, pageSize)
	if err != nil {
		lg.Error("failed to list roles", "error", err)
		response.InternalError(c)
		return
	}
	roleResponses := make([]dto.RoleResponse, len(roles))

	for i, role := range roles {
		roleResponse := dto.RoleResponse{}
		roleResponse.FromEntity(&role)
		roleResponses[i] = roleResponse
	}

	response.Ok(c, dto.RoleListResponse{
		Roles: roleResponses,
		Total: total,
	}, "")
}

// UpdateRole godoc
// @Summary      Update role details
// @Description  Update an existing role's details
// @Tags         roles
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        role  body      dto.RoleRequest  true  "Updated Role Data"
// @Param        id  path      string  true  "id"
// @Success      200   {object}  response.Response[dto.RoleResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /roles/{id} [put]
func (rc *RoleController) UpdateRole(c *gin.Context) {
	lg := rc.logger.With("method", "UpdateRole")
	var roleRequest dto.RoleRequest

	if err := c.ShouldBindJSON(&roleRequest); err != nil {
		lg.Error("failed to bind role data", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	role := roleRequest.ToEntity()
	role.ID = c.Param("id")
	updatedRole, err := rc.roleService.UpdateRole(context.Background(), role, roleRequest)
	if err != nil {
		if errors.Is(err, service.ErrModifingSystemRole) {
			response.BadRequest(c, "you cannot modify system role")
			return
		}
		lg.Error("failed to update role", "error", err)
		response.InternalError(c)
		return
	}

	roleResponse := dto.RoleResponse{}
	roleResponse.FromEntity(updatedRole)

	response.Ok(c, roleResponse, "")
}

// DeleteRole godoc
// @Summary      Delete a role
// @Description  Delete a role by ID
// @Tags         roles
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      string  true  "Role ID"
// @Success      200   {object}  response.Response[swagger.EmptyObject]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Failure      404   {object}  response.Response[swagger.EmptyObject]
// @Router       /roles/{id} [delete]
func (rc *RoleController) DeleteRole(c *gin.Context) {
	lg := rc.logger.With("method", "DeleteRole")
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		lg.Error("invalid role ID", "roleID", id)
		response.BadRequest(c, "invalid request body")
		return
	}

	err := rc.roleService.DeleteRole(context.Background(), id)
	if err != nil {
		if errors.Is(err, service.ErrModifingSystemRole) {
			response.BadRequest(c, "you cannot delete system role")
			return
		}
		lg.Error("failed to delete role", "error", err)
		response.InternalError(c)
		return
	}

	response.Ok(c, nil, "role deleted successfully")
}

// GetRoleDetail godoc
// @Summary      Get role details
// @Description  Retrieve role details by ID
// @Tags         roles
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      string  true  "Role ID"
// @Success      200   {object}  response.Response[dto.RoleResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      404   {object}  response.Response[swagger.EmptyObject]
// @Router       /roles/{id} [get]
func (rc *RoleController) GetRoleDetail(c *gin.Context) {
	lg := rc.logger.With("method", "GetRoleDetail")
	id := c.Param("id")

	role, err := rc.roleService.Detail(context.Background(), "id", id)
	if err != nil {
		lg.Error("failed to get role details", "error", err)
		response.NotFound(c)
		return
	}
	var roleResponse dto.RoleResponse
	roleResponse.FromEntity(role)

	response.Ok(c, roleResponse, "")
}
