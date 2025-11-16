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

var _ BaseRatesHandler = (*baseRatesHandler)(nil)

// BaseRatesHandler defining the handler interface
type BaseRatesHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type baseRatesHandler struct {
	iDao dao.BaseRatesDao
}

// NewBaseRatesHandler creating the handler interface
func NewBaseRatesHandler() BaseRatesHandler {
	return &baseRatesHandler{
		iDao: dao.NewBaseRatesDao(
			model.GetDB(),
			cache.NewBaseRatesCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create baseRates
// @Description submit information to create baseRates
// @Tags baseRates
// @accept json
// @Produce json
// @Param data body types.CreateBaseRatesRequest true "baseRates information"
// @Success 200 {object} types.CreateBaseRatesReply{}
// @Router /api/v1/baseRates [post]
// @Security BearerAuth
func (h *baseRatesHandler) Create(c *gin.Context) {
	form := &types.CreateBaseRatesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	baseRates := &model.BaseRates{}
	err = copier.Copy(baseRates, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateBaseRates)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, baseRates)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": baseRates.ID})
}

// DeleteByID delete a record by id
// @Summary delete baseRates
// @Description delete baseRates by id
// @Tags baseRates
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteBaseRatesByIDReply{}
// @Router /api/v1/baseRates/{id} [delete]
// @Security BearerAuth
func (h *baseRatesHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getBaseRatesIDFromPath(c)
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
// @Summary update baseRates
// @Description update baseRates information by id
// @Tags baseRates
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateBaseRatesByIDRequest true "baseRates information"
// @Success 200 {object} types.UpdateBaseRatesByIDReply{}
// @Router /api/v1/baseRates/{id} [put]
// @Security BearerAuth
func (h *baseRatesHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getBaseRatesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateBaseRatesByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	baseRates := &model.BaseRates{}
	err = copier.Copy(baseRates, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDBaseRates)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, baseRates)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get baseRates detail
// @Description get baseRates detail by id
// @Tags baseRates
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetBaseRatesByIDReply{}
// @Router /api/v1/baseRates/{id} [get]
// @Security BearerAuth
func (h *baseRatesHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getBaseRatesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	baseRates, err := h.iDao.GetByID(ctx, id)
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

	data := &types.BaseRatesObjDetail{}
	err = copier.Copy(data, baseRates)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDBaseRates)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"baseRates": data})
}

// List of records by query parameters
// @Summary list of baseRatess by query parameters
// @Description list of baseRatess by paging and conditions
// @Tags baseRates
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListBaseRatessReply{}
// @Router /api/v1/baseRates/list [post]
// @Security BearerAuth
func (h *baseRatesHandler) List(c *gin.Context) {
	form := &types.ListBaseRatessRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	baseRatess, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertBaseRatess(baseRatess)
	if err != nil {
		response.Error(c, ecode.ErrListBaseRates)
		return
	}

	response.Success(c, gin.H{
		"baseRatess": data,
		"total":      total,
	})
}

func getBaseRatesIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertBaseRates(baseRates *model.BaseRates) (*types.BaseRatesObjDetail, error) {
	data := &types.BaseRatesObjDetail{}
	err := copier.Copy(data, baseRates)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertBaseRatess(fromValues []*model.BaseRates) ([]*types.BaseRatesObjDetail, error) {
	toValues := []*types.BaseRatesObjDetail{}
	for _, v := range fromValues {
		data, err := convertBaseRates(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
