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

var _ PeakHourMultipliersHandler = (*peakHourMultipliersHandler)(nil)

// PeakHourMultipliersHandler defining the handler interface
type PeakHourMultipliersHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type peakHourMultipliersHandler struct {
	iDao dao.PeakHourMultipliersDao
}

// NewPeakHourMultipliersHandler creating the handler interface
func NewPeakHourMultipliersHandler() PeakHourMultipliersHandler {
	return &peakHourMultipliersHandler{
		iDao: dao.NewPeakHourMultipliersDao(
			model.GetDB(),
			cache.NewPeakHourMultipliersCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create peakHourMultipliers
// @Description submit information to create peakHourMultipliers
// @Tags peakHourMultipliers
// @accept json
// @Produce json
// @Param data body types.CreatePeakHourMultipliersRequest true "peakHourMultipliers information"
// @Success 200 {object} types.CreatePeakHourMultipliersReply{}
// @Router /api/v1/peakHourMultipliers [post]
// @Security BearerAuth
func (h *peakHourMultipliersHandler) Create(c *gin.Context) {
	form := &types.CreatePeakHourMultipliersRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	peakHourMultipliers := &model.PeakHourMultipliers{}
	err = copier.Copy(peakHourMultipliers, form)
	if err != nil {
		response.Error(c, ecode.ErrCreatePeakHourMultipliers)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, peakHourMultipliers)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": peakHourMultipliers.ID})
}

// DeleteByID delete a record by id
// @Summary delete peakHourMultipliers
// @Description delete peakHourMultipliers by id
// @Tags peakHourMultipliers
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeletePeakHourMultipliersByIDReply{}
// @Router /api/v1/peakHourMultipliers/{id} [delete]
// @Security BearerAuth
func (h *peakHourMultipliersHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getPeakHourMultipliersIDFromPath(c)
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
// @Summary update peakHourMultipliers
// @Description update peakHourMultipliers information by id
// @Tags peakHourMultipliers
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdatePeakHourMultipliersByIDRequest true "peakHourMultipliers information"
// @Success 200 {object} types.UpdatePeakHourMultipliersByIDReply{}
// @Router /api/v1/peakHourMultipliers/{id} [put]
// @Security BearerAuth
func (h *peakHourMultipliersHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getPeakHourMultipliersIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdatePeakHourMultipliersByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	peakHourMultipliers := &model.PeakHourMultipliers{}
	err = copier.Copy(peakHourMultipliers, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDPeakHourMultipliers)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, peakHourMultipliers)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get peakHourMultipliers detail
// @Description get peakHourMultipliers detail by id
// @Tags peakHourMultipliers
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetPeakHourMultipliersByIDReply{}
// @Router /api/v1/peakHourMultipliers/{id} [get]
// @Security BearerAuth
func (h *peakHourMultipliersHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getPeakHourMultipliersIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	peakHourMultipliers, err := h.iDao.GetByID(ctx, id)
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

	data := &types.PeakHourMultipliersObjDetail{}
	err = copier.Copy(data, peakHourMultipliers)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDPeakHourMultipliers)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"peakHourMultipliers": data})
}

// List of records by query parameters
// @Summary list of peakHourMultiplierss by query parameters
// @Description list of peakHourMultiplierss by paging and conditions
// @Tags peakHourMultipliers
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListPeakHourMultiplierssReply{}
// @Router /api/v1/peakHourMultipliers/list [post]
// @Security BearerAuth
func (h *peakHourMultipliersHandler) List(c *gin.Context) {
	form := &types.ListPeakHourMultiplierssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	peakHourMultiplierss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertPeakHourMultiplierss(peakHourMultiplierss)
	if err != nil {
		response.Error(c, ecode.ErrListPeakHourMultipliers)
		return
	}

	response.Success(c, gin.H{
		"peakHourMultiplierss": data,
		"total":                total,
	})
}

func getPeakHourMultipliersIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertPeakHourMultipliers(peakHourMultipliers *model.PeakHourMultipliers) (*types.PeakHourMultipliersObjDetail, error) {
	data := &types.PeakHourMultipliersObjDetail{}
	err := copier.Copy(data, peakHourMultipliers)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertPeakHourMultiplierss(fromValues []*model.PeakHourMultipliers) ([]*types.PeakHourMultipliersObjDetail, error) {
	toValues := []*types.PeakHourMultipliersObjDetail{}
	for _, v := range fromValues {
		data, err := convertPeakHourMultipliers(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
