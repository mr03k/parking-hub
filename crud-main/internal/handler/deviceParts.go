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

var _ DevicePartsHandler = (*devicePartsHandler)(nil)

// DevicePartsHandler defining the handler interface
type DevicePartsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type devicePartsHandler struct {
	iDao dao.DevicePartsDao
}

// NewDevicePartsHandler creating the handler interface
func NewDevicePartsHandler() DevicePartsHandler {
	return &devicePartsHandler{
		iDao: dao.NewDevicePartsDao(
			model.GetDB(),
			cache.NewDevicePartsCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create deviceParts
// @Description submit information to create deviceParts
// @Tags deviceParts
// @accept json
// @Produce json
// @Param data body types.CreateDevicePartsRequest true "deviceParts information"
// @Success 200 {object} types.CreateDevicePartsReply{}
// @Router /api/v1/deviceParts [post]
// @Security BearerAuth
func (h *devicePartsHandler) Create(c *gin.Context) {
	form := &types.CreateDevicePartsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	deviceParts := &model.DeviceParts{}
	err = copier.Copy(deviceParts, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateDeviceParts)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, deviceParts)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": deviceParts.ID})
}

// DeleteByID delete a record by id
// @Summary delete deviceParts
// @Description delete deviceParts by id
// @Tags deviceParts
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteDevicePartsByIDReply{}
// @Router /api/v1/deviceParts/{id} [delete]
// @Security BearerAuth
func (h *devicePartsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getDevicePartsIDFromPath(c)
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
// @Summary update deviceParts
// @Description update deviceParts information by id
// @Tags deviceParts
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateDevicePartsByIDRequest true "deviceParts information"
// @Success 200 {object} types.UpdateDevicePartsByIDReply{}
// @Router /api/v1/deviceParts/{id} [put]
// @Security BearerAuth
func (h *devicePartsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getDevicePartsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateDevicePartsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	deviceParts := &model.DeviceParts{}
	err = copier.Copy(deviceParts, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDDeviceParts)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, deviceParts)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get deviceParts detail
// @Description get deviceParts detail by id
// @Tags deviceParts
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetDevicePartsByIDReply{}
// @Router /api/v1/deviceParts/{id} [get]
// @Security BearerAuth
func (h *devicePartsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getDevicePartsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	deviceParts, err := h.iDao.GetByID(ctx, id)
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

	data := &types.DevicePartsObjDetail{}
	err = copier.Copy(data, deviceParts)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDDeviceParts)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"deviceParts": data})
}

// List of records by query parameters
// @Summary list of devicePartss by query parameters
// @Description list of devicePartss by paging and conditions
// @Tags deviceParts
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListDevicePartssReply{}
// @Router /api/v1/deviceParts/list [post]
// @Security BearerAuth
func (h *devicePartsHandler) List(c *gin.Context) {
	form := &types.ListDevicePartssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	devicePartss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertDevicePartss(devicePartss)
	if err != nil {
		response.Error(c, ecode.ErrListDeviceParts)
		return
	}

	response.Success(c, gin.H{
		"devicePartss": data,
		"total":        total,
	})
}

func getDevicePartsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertDeviceParts(deviceParts *model.DeviceParts) (*types.DevicePartsObjDetail, error) {
	data := &types.DevicePartsObjDetail{}
	err := copier.Copy(data, deviceParts)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertDevicePartss(fromValues []*model.DeviceParts) ([]*types.DevicePartsObjDetail, error) {
	toValues := []*types.DevicePartsObjDetail{}
	for _, v := range fromValues {
		data, err := convertDeviceParts(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
