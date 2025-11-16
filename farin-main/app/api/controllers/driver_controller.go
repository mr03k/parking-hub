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

type DriverController struct {
	driverService *service.DriverService
	logger        *slog.Logger
	env           *godotenv.Env
}

func NewDriverController(logger *slog.Logger, us *service.DriverService, env *godotenv.Env) *DriverController {
	return &DriverController{
		driverService: us,
		logger:        logger.With("layer", "DriverController"),
		env:           env,
	}
}

// CreateDriver godoc
// @Summary      Create a new driver
// @Description  Create a new driver by providing driver details
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        driver  body      dto.DriverRequest  true  "Driver Data"
// @Success      201   {object}  response.Response[dto.DriverResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /drivers [post]
func (uc *DriverController) CreateDriver(c *gin.Context) {
	lg := uc.logger.With("method", "CreateDriver")
	var driverRequest dto.DriverRequest

	if err := c.ShouldBindJSON(&driverRequest); err != nil {
		lg.Error("failed to bind driver data", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	createdDriver, err := uc.driverService.CreateDriver(context.Background(), driverRequest.ToEntity(), driverRequest)
	if err != nil {
		lg.Error("failed to create driver", "error", err)
		response.InternalError(c)
		return
	}
	driverResponse := dto.DriverResponse{}
	driverResponse.FromEntity(createdDriver)

	response.Created(c, driverResponse)
}

// ListDrivers godoc
// @Summary      List drivers
// @Description  Retrieve a list of drivers with optional filters and pagination
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        sortField  query     string  false  "Field to sort by"  default(created_at)
// @Param        sortOrder  query     string  false  "Sort order (asc/desc)"  default(asc)
// @Param        page       query     int     false  "Page number"  default(1)
// @Param        pageSize   query     int     false  "Page size"  default(10)
// @Success      200        {object}  response.Response[dto.DriverListResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /drivers [get]
func (uc *DriverController) ListDrivers(c *gin.Context) {
	lg := uc.logger.With("method", "ListDrivers")
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

	drivers, total, err := uc.driverService.ListDrivers(context.Background(), filters, sortField, sortOrder, page, pageSize)
	if err != nil {
		lg.Error("failed to list drivers", "error", err)
		response.InternalError(c)
		return
	}
	driverResponses := make([]dto.DriverResponse, len(drivers))

	for i, driver := range drivers {
		driverResponse := dto.DriverResponse{}
		driverResponse.FromEntity(&driver)
		driverResponses[i] = driverResponse
	}

	response.Ok(c, dto.DriverListResponse{
		Drivers: driverResponses,
		Total:   total,
	}, "")
}

// UpdateDriver godoc
// @Summary      Update driver details
// @Description  Update an existing driver's details
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        driver  body      dto.DriverRequest  true  "Updated Driver Data"
// @Param        id  path      string  true  "id"
// @Success      200   {object}  response.Response[dto.DriverResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /drivers/{id} [put]
func (uc *DriverController) UpdateDriver(c *gin.Context) {
	lg := uc.logger.With("method", "UpdateDriver")
	var driverRequest dto.DriverRequest

	if err := c.ShouldBindJSON(&driverRequest); err != nil {
		lg.Error("failed to bind driver data", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	driver := driverRequest.ToEntity()
	driver.ID = c.Param("id")
	updatedDriver, err := uc.driverService.UpdateDriver(context.Background(), driver, driverRequest)
	if err != nil {
		lg.Error("failed to update driver", "error", err)
		response.InternalError(c)
		return
	}

	driverResponse := dto.DriverResponse{}
	driverResponse.FromEntity(updatedDriver)

	response.Ok(c, driverResponse, "")
}

// DeleteDriver godoc
// @Summary      Delete a driver
// @Description  Delete a driver by ID
// @Tags         drivers
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      string  true  "Driver ID"
// @Success      200   {object}  response.Response[swagger.EmptyObject]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Failure      404   {object}  response.Response[swagger.EmptyObject]
// @Router       /drivers/{id} [delete]
func (uc *DriverController) DeleteDriver(c *gin.Context) {
	lg := uc.logger.With("method", "DeleteDriver")
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		lg.Error("invalid driver ID", "driverID", id)
		response.BadRequest(c, "invalid request body")
		return
	}

	err := uc.driverService.DeleteDriver(context.Background(), id)
	if err != nil {
		lg.Error("failed to delete driver", "error", err)
		response.InternalError(c)
		return
	}

	response.Ok(c, nil, "driver deleted successfully")
}

// GetDriverDetail godoc
// @Summary      Get driver details
// @Description  Retrieve driver details by ID
// @Tags         drivers
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      string  true  "Driver ID"
// @Success      200   {object}  response.Response[dto.DriverResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      404   {object}  response.Response[swagger.EmptyObject]
// @Router       /drivers/{id} [get]
func (uc *DriverController) GetDriverDetail(c *gin.Context) {
	lg := uc.logger.With("method", "GetDriverDetail")
	id := c.Param("id")

	driver, err := uc.driverService.Detail(context.Background(), "id", id)
	if err != nil {
		lg.Error("failed to get driver details", "error", err)
		response.NotFound(c)
		return
	}
	var driverResponse dto.DriverResponse
	driverResponse.FromEntity(driver)

	response.Ok(c, driverResponse, "")
}
