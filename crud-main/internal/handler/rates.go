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

var _ RatesHandler = (*ratesHandler)(nil)

// RatesHandler defining the handler interface
type RatesHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type ratesHandler struct {
	iDao dao.RatesDao
}

// NewRatesHandler creating the handler interface
func NewRatesHandler() RatesHandler {
	return &ratesHandler{
		iDao: dao.NewRatesDao(
			model.GetDB(),
			cache.NewRatesCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create rates
// @Description submit information to create rates
// @Tags rates
// @accept json
// @Produce json
// @Param data body types.CreateRatesRequest true "rates information"
// @Success 200 {object} types.CreateRatesReply{}
// @Router /api/v1/rates [post]
// @Security BearerAuth
func (h *ratesHandler) Create(c *gin.Context) {
	form := &types.CreateRatesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	rates := &model.Rates{}
	err = copier.Copy(rates, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateRates)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, rates)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": rates.ID})
}

// DeleteByID delete a record by id
// @Summary delete rates
// @Description delete rates by id
// @Tags rates
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteRatesByIDReply{}
// @Router /api/v1/rates/{id} [delete]
// @Security BearerAuth
func (h *ratesHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getRatesIDFromPath(c)
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
// @Summary update rates
// @Description update rates information by id
// @Tags rates
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateRatesByIDRequest true "rates information"
// @Success 200 {object} types.UpdateRatesByIDReply{}
// @Router /api/v1/rates/{id} [put]
// @Security BearerAuth
func (h *ratesHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getRatesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateRatesByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	rates := &model.Rates{}
	err = copier.Copy(rates, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDRates)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, rates)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get rates detail
// @Description get rates detail by id
// @Tags rates
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetRatesByIDReply{}
// @Router /api/v1/rates/{id} [get]
// @Security BearerAuth
func (h *ratesHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getRatesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	rates, err := h.iDao.GetByID(ctx, id)
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

	data := &types.RatesObjDetail{}
	err = copier.Copy(data, rates)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDRates)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"rates": data})
}

// List of records by query parameters
// @Summary list of ratess by query parameters
// @Description list of ratess by paging and conditions
// @Tags rates
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListRatessReply{}
// @Router /api/v1/rates/list [post]
// @Security BearerAuth
func (h *ratesHandler) List(c *gin.Context) {
	form := &types.ListRatessRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	ratess, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertRatess(ratess)
	if err != nil {
		response.Error(c, ecode.ErrListRates)
		return
	}

	response.Success(c, gin.H{
		"ratess": data,
		"total":  total,
	})
}

func getRatesIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertRates(rates *model.Rates) (*types.RatesObjDetail, error) {
	data := &types.RatesObjDetail{}
	err := copier.Copy(data, rates)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertRatess(fromValues []*model.Rates) ([]*types.RatesObjDetail, error) {
	toValues := []*types.RatesObjDetail{}
	for _, v := range fromValues {
		data, err := convertRates(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
