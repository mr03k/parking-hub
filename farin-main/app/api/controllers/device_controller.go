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

type DeviceController struct {
	deviceService *service.DeviceService
	logger        *slog.Logger
	env           *godotenv.Env
}

func NewDeviceController(logger *slog.Logger, us *service.DeviceService, env *godotenv.Env) *DeviceController {
	return &DeviceController{
		deviceService: us,
		logger:        logger.With("layer", "DeviceController"),
		env:           env,
	}
}

// CreateDevice godoc
// @Summary      Create a new device
// @Description  Create a new device by providing device details
// @Tags         devices
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        device  body      dto.DeviceRequest  true  "Device Data"
// @Success      201   {object}  response.Response[dto.DeviceResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /devices [post]
func (uc *DeviceController) CreateDevice(c *gin.Context) {
	lg := uc.logger.With("method", "CreateDevice")
	var deviceRequest dto.DeviceRequest

	if err := c.ShouldBindJSON(&deviceRequest); err != nil {
		lg.Error("failed to bind device data", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	createdDevice, err := uc.deviceService.CreateDevice(context.Background(), deviceRequest.ToEntity(), &deviceRequest)
	if err != nil {
		lg.Error("failed to create device", "error", err)
		response.InternalError(c)
		return
	}
	deviceResponse := dto.DeviceResponse{}
	deviceResponse.FromEntity(uc.env, createdDevice)

	response.Created(c, deviceResponse)
}

// ListDevices godoc
// @Summary      List devices
// @Description  Retrieve a list of devices with optional filters and pagination
// @Tags         devices
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        sortField  query     string  false  "Field to sort by"  default(created_at)
// @Param        sortOrder  query     string  false  "Sort order (asc/desc)"  default(asc)
// @Param        page       query     int     false  "Page number"  default(1)
// @Param        pageSize   query     int     false  "Page size"  default(10)
// @Success      200        {object}  response.Response[dto.DeviceListResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /devices [get]
func (uc *DeviceController) ListDevices(c *gin.Context) {
	lg := uc.logger.With("method", "ListDevices")
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

	devices, total, err := uc.deviceService.ListDevices(context.Background(), filters, sortField, sortOrder, page, pageSize)
	if err != nil {
		lg.Error("failed to list devices", "error", err)
		response.InternalError(c)
		return
	}
	deviceResponses := make([]dto.DeviceResponse, len(devices))

	for i, device := range devices {
		deviceResponse := dto.DeviceResponse{}
		deviceResponse.FromEntity(uc.env, &device)
		deviceResponses[i] = deviceResponse
	}

	response.Ok(c, dto.DeviceListResponse{
		Devices: deviceResponses,
		Total:   total,
	}, "")
}

// UpdateDevice godoc
// @Summary      Update device details
// @Description  Update an existing device's details
// @Tags         devices
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        device  body      dto.DeviceRequest  true  "Updated Device Data"
// @Param        id  path      string  true  "id"
// @Success      200   {object}  response.Response[dto.DeviceResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /devices/{id} [put]
func (uc *DeviceController) UpdateDevice(c *gin.Context) {
	lg := uc.logger.With("method", "UpdateDevice")
	var deviceRequest dto.DeviceRequest

	if err := c.ShouldBindJSON(&deviceRequest); err != nil {
		lg.Error("failed to bind device data", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	device := deviceRequest.ToEntity()
	device.ID = c.Param("id")
	updatedDevice, err := uc.deviceService.UpdateDevice(context.Background(), device, deviceRequest)
	if err != nil {
		lg.Error("failed to update device", "error", err)
		response.InternalError(c)
		return
	}

	deviceResponse := dto.DeviceResponse{}
	deviceResponse.FromEntity(uc.env, updatedDevice)

	response.Ok(c, deviceResponse, "")
}

// DeleteDevice godoc
// @Summary      Delete a device
// @Description  Delete a device by ID
// @Tags         devices
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      string  true  "Device ID"
// @Success      200   {object}  response.Response[swagger.EmptyObject]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Failure      404   {object}  response.Response[swagger.EmptyObject]
// @Router       /devices/{id} [delete]
func (uc *DeviceController) DeleteDevice(c *gin.Context) {
	lg := uc.logger.With("method", "DeleteDevice")
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		lg.Error("invalid device ID", "deviceID", id)
		response.BadRequest(c, "invalid request body")
		return
	}

	err := uc.deviceService.DeleteDevice(context.Background(), id)
	if err != nil {
		lg.Error("failed to delete device", "error", err)
		response.InternalError(c)
		return
	}

	response.Ok(c, nil, "device deleted successfully")
}

// GetDeviceDetail godoc
// @Summary      Get device details
// @Description  Retrieve device details by ID
// @Tags         devices
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      string  true  "Device ID"
// @Success      200   {object}  response.Response[dto.DeviceResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      404   {object}  response.Response[swagger.EmptyObject]
// @Router       /devices/{id} [get]
func (uc *DeviceController) GetDeviceDetail(c *gin.Context) {
	lg := uc.logger.With("method", "GetDeviceDetail")
	id := c.Param("id")

	device, err := uc.deviceService.Detail(context.Background(), "id", id)
	if err != nil {
		lg.Error("failed to get device details", "error", err)
		response.NotFound(c)
		return
	}
	var deviceResponse dto.DeviceResponse
	deviceResponse.FromEntity(uc.env, device)

	response.Ok(c, deviceResponse, "")
}
