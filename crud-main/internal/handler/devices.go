package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"

	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/utils"

	"git.abanppc.com/farin-project/crud/internal/cache"
	"git.abanppc.com/farin-project/crud/internal/dao"
	"git.abanppc.com/farin-project/crud/internal/ecode"
	"git.abanppc.com/farin-project/crud/internal/model"
	"git.abanppc.com/farin-project/crud/internal/types"
)

var _ DevicesHandler = (*devicesHandler)(nil)

// DevicesHandler defining the handler interface
type DevicesHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type devicesHandler struct {
	iDao dao.DevicesDao
}

// NewDevicesHandler creating the handler interface
func NewDevicesHandler() DevicesHandler {
	return &devicesHandler{
		iDao: dao.NewDevicesDao(
			model.GetDB(),
			cache.NewDevicesCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create devices
// @Description submit information to create devices
// @Tags devices
// @accept json
// @Produce json
// @Param data body types.CreateDevicesRequest true "devices information"
// @Success 200 {object} types.CreateDevicesReply{}
// @Router /api/v1/devices [post]
// @Security BearerAuth
func (h *devicesHandler) Create(c *gin.Context) {
	form := &types.CreateDevicesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	devices := &model.Devices{}
	err = copier.Copy(devices, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateDevices)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, devices)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": devices.ID})
}

// DeleteByID delete a record by id
// @Summary delete devices
// @Description delete devices by id
// @Tags devices
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteDevicesByIDReply{}
// @Router /api/v1/devices/{id} [delete]
// @Security BearerAuth
func (h *devicesHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getDevicesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err := h.iDao.DeleteByID(ctx, id)
	if err != nil {
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// UpdateByID update information by id
// @Summary update devices
// @Description update devices information by id
// @Tags devices
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateDevicesByIDRequest true "devices information"
// @Success 200 {object} types.UpdateDevicesByIDReply{}
// @Router /api/v1/devices/{id} [put]
// @Security BearerAuth
func (h *devicesHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getDevicesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateDevicesByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	devices := &model.Devices{}
	err = copier.Copy(devices, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDDevices)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, devices)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get devices detail
// @Description get devices detail by id
// @Tags devices
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetDevicesByIDReply{}
// @Router /api/v1/devices/{id} [get]
// @Security BearerAuth
func (h *devicesHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getDevicesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	devices, err := h.iDao.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			logger.Warn("GetByID not found", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.DevicesObjDetail{}
	err = copier.Copy(data, devices)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDDevices)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"devices": data})
}

// List of records by query parameters
// @Summary list of devicess by query parameters
// @Description list of devicess by paging and conditions
// @Tags devices
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListDevicessReply{}
// @Router /api/v1/devices/list [post]
// @Security BearerAuth
func (h *devicesHandler) List(c *gin.Context) {
	form := &types.ListDevicessRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	devicess, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertDevicess(devicess)
	if err != nil {
		response.Error(c, ecode.ErrListDevices)
		return
	}

	response.Success(c, gin.H{
		"devicess": data,
		"total":    total,
	})
}

func getDevicesIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertDevices(devices *model.Devices) (*types.DevicesObjDetail, error) {
	data := &types.DevicesObjDetail{}
	err := copier.Copy(data, devices)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertDevicess(fromValues []*model.Devices) ([]*types.DevicesObjDetail, error) {
	toValues := []*types.DevicesObjDetail{}
	for _, v := range fromValues {
		data, err := convertDevices(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
