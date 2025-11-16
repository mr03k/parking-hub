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

type VehicleController struct {
	vehicleService *service.VehicleService
	logger         *slog.Logger
	env            *godotenv.Env
}

func NewVehicleController(logger *slog.Logger, us *service.VehicleService, env *godotenv.Env) *VehicleController {
	return &VehicleController{
		vehicleService: us,
		logger:         logger.With("layer", "VehicleController"),
		env:            env,
	}
}

// CreateVehicle godoc
// @Summary      Create a new vehicle
// @Description  Create a new vehicle by providing vehicle details
// @Tags         vehicles
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        vehicle  body      dto.VehicleRequest  true  "Vehicle Data"
// @Success      201   {object}  response.Response[dto.VehicleResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /vehicles [post]
func (uc *VehicleController) CreateVehicle(c *gin.Context) {
	lg := uc.logger.With("method", "CreateVehicle")
	var vehicleRequest dto.VehicleRequest

	if err := c.ShouldBindJSON(&vehicleRequest); err != nil {
		lg.Error("failed to bind vehicle data", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	createdVehicle, err := uc.vehicleService.CreateVehicle(context.Background(), vehicleRequest.ToEntity(), &vehicleRequest)
	if err != nil {
		lg.Error("failed to create vehicle", "error", err)
		response.InternalError(c)
		return
	}
	vehicleResponse := dto.VehicleResponse{}
	vehicleResponse.FromEntity(createdVehicle)

	response.Created(c, vehicleResponse)
}

// ListVehicles godoc
// @Summary      List vehicles
// @Description  Retrieve a list of vehicles with optional filters and pagination
// @Tags         vehicles
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        sortField  query     string  false  "Field to sort by"  default(created_at)
// @Param        sortOrder  query     string  false  "Sort order (asc/desc)"  default(asc)
// @Param        page       query     int     false  "Page number"  default(1)
// @Param        pageSize   query     int     false  "Page size"  default(10)
// @Success      200        {object}  response.Response[dto.VehicleListResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /vehicles [get]
func (uc *VehicleController) ListVehicles(c *gin.Context) {
	lg := uc.logger.With("method", "ListVehicles")
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

	vehicles, total, err := uc.vehicleService.ListVehicles(context.Background(), filters, sortField, sortOrder, page, pageSize)
	if err != nil {
		lg.Error("failed to list vehicles", "error", err)
		response.InternalError(c)
		return
	}
	vehicleResponses := make([]dto.VehicleResponse, len(vehicles))

	for i, vehicle := range vehicles {
		vehicleResponse := dto.VehicleResponse{}
		vehicleResponse.FromEntity(&vehicle)
		vehicleResponses[i] = vehicleResponse
	}

	response.Ok(c, dto.VehicleListResponse{
		Vehicles: vehicleResponses,
		Total:    total,
	}, "")
}

// UpdateVehicle godoc
// @Summary      Update vehicle details
// @Description  Update an existing vehicle's details
// @Tags         vehicles
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        vehicle  body      dto.VehicleRequest  true  "Updated Vehicle Data"
// @Param        id  path      string  true  "id"
// @Success      200   {object}  response.Response[dto.VehicleResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /vehicles/{id} [put]
func (uc *VehicleController) UpdateVehicle(c *gin.Context) {
	lg := uc.logger.With("method", "UpdateVehicle")
	var vehicleRequest dto.VehicleRequest

	if err := c.ShouldBindJSON(&vehicleRequest); err != nil {
		lg.Error("failed to bind vehicle data", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	vehicle := vehicleRequest.ToEntity()
	vehicle.ID = c.Param("id")
	updatedVehicle, err := uc.vehicleService.UpdateVehicle(context.Background(), vehicle, vehicleRequest)
	if err != nil {
		lg.Error("failed to update vehicle", "error", err)
		response.InternalError(c)
		return
	}

	vehicleResponse := dto.VehicleResponse{}
	vehicleResponse.FromEntity(updatedVehicle)

	response.Ok(c, vehicleResponse, "")
}

// DeleteVehicle godoc
// @Summary      Delete a vehicle
// @Description  Delete a vehicle by ID
// @Tags         vehicles
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      string  true  "Vehicle ID"
// @Success      200   {object}  response.Response[swagger.EmptyObject]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Failure      404   {object}  response.Response[swagger.EmptyObject]
// @Router       /vehicles/{id} [delete]
func (uc *VehicleController) DeleteVehicle(c *gin.Context) {
	lg := uc.logger.With("method", "DeleteVehicle")
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		lg.Error("invalid vehicle ID", "vehicleID", id)
		response.BadRequest(c, "invalid request body")
		return
	}

	err := uc.vehicleService.DeleteVehicle(context.Background(), id)
	if err != nil {
		lg.Error("failed to delete vehicle", "error", err)
		response.InternalError(c)
		return
	}

	response.Ok(c, nil, "vehicle deleted successfully")
}

// GetVehicleDetail godoc
// @Summary      Get vehicle details
// @Description  Retrieve vehicle details by ID
// @Tags         vehicles
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      string  true  "Vehicle ID"
// @Success      200   {object}  response.Response[dto.VehicleResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      404   {object}  response.Response[swagger.EmptyObject]
// @Router       /vehicles/{id} [get]
func (uc *VehicleController) GetVehicleDetail(c *gin.Context) {
	lg := uc.logger.With("method", "GetVehicleDetail")
	id := c.Param("id")

	vehicle, err := uc.vehicleService.Detail(context.Background(), "id", id)
	if err != nil {
		lg.Error("failed to get vehicle details", "error", err)
		response.NotFound(c)
		return
	}
	var vehicleResponse dto.VehicleResponse
	vehicleResponse.FromEntity(vehicle)

	response.Ok(c, vehicleResponse, "")
}
