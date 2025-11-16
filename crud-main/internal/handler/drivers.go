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

var _ DriversHandler = (*driversHandler)(nil)

// DriversHandler defining the handler interface
type DriversHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type driversHandler struct {
	iDao dao.DriversDao
}

// NewDriversHandler creating the handler interface
func NewDriversHandler() DriversHandler {
	return &driversHandler{
		iDao: dao.NewDriversDao(
			model.GetDB(),
			cache.NewDriversCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create drivers
// @Description submit information to create drivers
// @Tags drivers
// @accept json
// @Produce json
// @Param data body types.CreateDriversRequest true "drivers information"
// @Success 200 {object} types.CreateDriversReply{}
// @Router /api/v1/drivers [post]
// @Security BearerAuth
func (h *driversHandler) Create(c *gin.Context) {
	form := &types.CreateDriversRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	drivers := &model.Drivers{}
	err = copier.Copy(drivers, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateDrivers)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, drivers)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": drivers.ID})
}

// DeleteByID delete a record by id
// @Summary delete drivers
// @Description delete drivers by id
// @Tags drivers
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteDriversByIDReply{}
// @Router /api/v1/drivers/{id} [delete]
// @Security BearerAuth
func (h *driversHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getDriversIDFromPath(c)
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
// @Summary update drivers
// @Description update drivers information by id
// @Tags drivers
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateDriversByIDRequest true "drivers information"
// @Success 200 {object} types.UpdateDriversByIDReply{}
// @Router /api/v1/drivers/{id} [put]
// @Security BearerAuth
func (h *driversHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getDriversIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateDriversByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	drivers := &model.Drivers{}
	err = copier.Copy(drivers, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDDrivers)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, drivers)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get drivers detail
// @Description get drivers detail by id
// @Tags drivers
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetDriversByIDReply{}
// @Router /api/v1/drivers/{id} [get]
// @Security BearerAuth
func (h *driversHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getDriversIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	drivers, err := h.iDao.GetByID(ctx, id)
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

	data := &types.DriversObjDetail{}
	err = copier.Copy(data, drivers)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDDrivers)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"drivers": data})
}

// List of records by query parameters
// @Summary list of driverss by query parameters
// @Description list of driverss by paging and conditions
// @Tags drivers
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListDriverssReply{}
// @Router /api/v1/drivers/list [post]
// @Security BearerAuth
func (h *driversHandler) List(c *gin.Context) {
	form := &types.ListDriverssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	driverss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertDriverss(driverss)
	if err != nil {
		response.Error(c, ecode.ErrListDrivers)
		return
	}

	response.Success(c, gin.H{
		"driverss": data,
		"total":    total,
	})
}

func getDriversIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertDrivers(drivers *model.Drivers) (*types.DriversObjDetail, error) {
	data := &types.DriversObjDetail{}
	err := copier.Copy(data, drivers)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertDriverss(fromValues []*model.Drivers) ([]*types.DriversObjDetail, error) {
	toValues := []*types.DriversObjDetail{}
	for _, v := range fromValues {
		data, err := convertDrivers(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
