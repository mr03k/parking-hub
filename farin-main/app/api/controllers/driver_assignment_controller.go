package controller

import (
	"context"
	"farin/app/api/response"
	"farin/domain/dto"
	"farin/domain/service"
	"farin/infrastructure/godotenv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"strconv"
)

type DriverAssignmentController struct {
	driverAssignmentService *service.DriverAssignmentService
	logger                  *slog.Logger
	env                     *godotenv.Env
}

func NewDriverAssignmentController(logger *slog.Logger, us *service.DriverAssignmentService, env *godotenv.Env) *DriverAssignmentController {
	return &DriverAssignmentController{
		driverAssignmentService: us,
		logger:                  logger.With("layer", "DriverAssignmentController"),
		env:                     env,
	}
}

// CreateDriverAssignment godoc
// @Summary      Create a new driverAssignment
// @Description  Create a new driverAssignment by providing driverAssignment details
// @Tags         driverAssignments
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        driverAssignment  body      dto.DriverAssignmentRequest  true  "DriverAssignment Data"
// @Success      201   {object}  response.Response[dto.DriverAssignmentResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /driver-assignments [post]
func (uc *DriverAssignmentController) CreateDriverAssignment(c *gin.Context) {
	lg := uc.logger.With("method", "CreateDriverAssignment")
	var driverAssignmentRequest dto.DriverAssignmentRequest

	if err := c.ShouldBindJSON(&driverAssignmentRequest); err != nil {
		lg.Error("failed to bind driverAssignment data", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	createdDriverAssignment, err := uc.driverAssignmentService.CreateDriverAssignment(context.Background(), driverAssignmentRequest.ToEntity())
	if err != nil {
		lg.Error("failed to create driverAssignment", "error", err)
		response.InternalError(c)
		return
	}
	driverAssignmentResponse := dto.DriverAssignmentResponse{}
	driverAssignmentResponse.FromEntity(createdDriverAssignment)

	response.Created(c, driverAssignmentResponse)
}

// ListDriverAssignments godoc
// @Summary      List driverAssignments
// @Description  Retrieve a list of driverAssignments with optional filters and pagination
// @Tags         driverAssignments
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        sortField  query     string  false  "Field to sort by"  default(created_at)
// @Param        sortOrder  query     string  false  "Sort order (asc/desc)"  default(asc)
// @Param        page       query     int     false  "Page number"  default(1)
// @Param        pageSize   query     int     false  "Page size"  default(10)
// @Success      200        {object}  response.Response[dto.DriverAssignmentListResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /driver-assignments [get]
func (uc *DriverAssignmentController) ListDriverAssignments(c *gin.Context) {
	lg := uc.logger.With("method", "ListDriverAssignments")
	filters := make(map[string]interface{})
	sortField := c.DefaultQuery("sortField", "created_at")
	sortOrder := c.DefaultQuery("sortOrder", "asc")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		response.BadRequest(c, "invalid page param")
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		response.BadRequest(c, "invalid page size param")
	}

	driverAssignments, total, err := uc.driverAssignmentService.ListDriverAssignments(context.Background(), filters, sortField, sortOrder, page, pageSize)
	if err != nil {
		lg.Error("failed to list driverAssignments", "error", err)
		response.InternalError(c)
		return
	}
	driverAssignmentResponses := make([]dto.DriverAssignmentResponse, len(driverAssignments))

	for i, driverAssignment := range driverAssignments {
		driverAssignmentResponse := dto.DriverAssignmentResponse{}
		driverAssignmentResponse.FromEntity(&driverAssignment)
		driverAssignmentResponses[i] = driverAssignmentResponse
	}

	response.Ok(c, dto.DriverAssignmentListResponse{
		DriverAssignments: driverAssignmentResponses,
		Total:             total,
	}, "")
}

// UpdateDriverAssignment godoc
// @Summary      Update driverAssignment details
// @Description  Update an existing driverAssignment's details
// @Tags         driverAssignments
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        driverAssignment  body      dto.DriverAssignmentRequest  true  "Updated DriverAssignment Data"
// @Param        id  path      string  true  "id"
// @Success      200   {object}  response.Response[dto.DriverAssignmentResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /driver-assignments/{id} [put]
func (uc *DriverAssignmentController) UpdateDriverAssignment(c *gin.Context) {
	lg := uc.logger.With("method", "UpdateDriverAssignment")
	var driverAssignmentRequest dto.DriverAssignmentRequest

	if err := c.ShouldBindJSON(&driverAssignmentRequest); err != nil {
		lg.Error("failed to bind driverAssignment data", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	driverAssignment := driverAssignmentRequest.ToEntity()
	driverAssignment.ID = c.Param("id")
	updatedDriverAssignment, err := uc.driverAssignmentService.UpdateDriverAssignment(context.Background(), driverAssignment)
	if err != nil {
		lg.Error("failed to update driverAssignment", "error", err)
		response.InternalError(c)
		return
	}

	driverAssignmentResponse := dto.DriverAssignmentResponse{}
	driverAssignmentResponse.FromEntity(updatedDriverAssignment)

	response.Ok(c, driverAssignmentResponse, "")
}

// DeleteDriverAssignment godoc
// @Summary      Delete a driverAssignment
// @Description  Delete a driverAssignment by ID
// @Tags         driverAssignments
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      string  true  "DriverAssignment ID"
// @Success      200   {object}  response.Response[swagger.EmptyObject]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Failure      404   {object}  response.Response[swagger.EmptyObject]
// @Router       /driver-assignments/{id} [delete]
func (uc *DriverAssignmentController) DeleteDriverAssignment(c *gin.Context) {
	lg := uc.logger.With("method", "DeleteDriverAssignment")
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		lg.Error("invalid driverAssignment ID", "driverAssignmentID", id)
		response.BadRequest(c, "invalid request body")
		return
	}

	err := uc.driverAssignmentService.DeleteDriverAssignment(context.Background(), id)
	if err != nil {
		lg.Error("failed to delete driverAssignment", "error", err)
		response.InternalError(c)
		return
	}

	response.Ok(c, nil, "driverAssignment deleted successfully")
}

// GetDriverAssignmentDetail godoc
// @Summary      Get driverAssignment details
// @Description  Retrieve driverAssignment details by ID
// @Tags         driverAssignments
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      string  true  "DriverAssignment ID"
// @Success      200   {object}  response.Response[dto.DriverAssignmentResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      404   {object}  response.Response[swagger.EmptyObject]
// @Router       /driver-assignments/{id} [get]
func (uc *DriverAssignmentController) GetDriverAssignmentDetail(c *gin.Context) {
	lg := uc.logger.With("method", "GetDriverAssignmentDetail")
	id := c.Param("id")

	driverAssignment, err := uc.driverAssignmentService.Detail(context.Background(), "id", id)
	if err != nil {
		lg.Error("failed to get driverAssignment details", "error", err)
		response.NotFound(c)
		return
	}
	var driverAssignmentResponse dto.DriverAssignmentResponse
	driverAssignmentResponse.FromEntity(driverAssignment)

	response.Ok(c, driverAssignmentResponse, "")
}
