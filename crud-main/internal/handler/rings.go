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

var _ RingsHandler = (*ringsHandler)(nil)

// RingsHandler defining the handler interface
type RingsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type ringsHandler struct {
	iDao dao.RingsDao
}

// NewRingsHandler creating the handler interface
func NewRingsHandler() RingsHandler {
	return &ringsHandler{
		iDao: dao.NewRingsDao(
			model.GetDB(),
			cache.NewRingsCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create rings
// @Description submit information to create rings
// @Tags rings
// @accept json
// @Produce json
// @Param data body types.CreateRingsRequest true "rings information"
// @Success 200 {object} types.CreateRingsReply{}
// @Router /api/v1/rings [post]
// @Security BearerAuth
func (h *ringsHandler) Create(c *gin.Context) {
	form := &types.CreateRingsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	rings := &model.Rings{}
	err = copier.Copy(rings, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateRings)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, rings)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": rings.ID})
}

// DeleteByID delete a record by id
// @Summary delete rings
// @Description delete rings by id
// @Tags rings
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteRingsByIDReply{}
// @Router /api/v1/rings/{id} [delete]
// @Security BearerAuth
func (h *ringsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getRingsIDFromPath(c)
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
// @Summary update rings
// @Description update rings information by id
// @Tags rings
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateRingsByIDRequest true "rings information"
// @Success 200 {object} types.UpdateRingsByIDReply{}
// @Router /api/v1/rings/{id} [put]
// @Security BearerAuth
func (h *ringsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getRingsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateRingsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	rings := &model.Rings{}
	err = copier.Copy(rings, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDRings)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, rings)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get rings detail
// @Description get rings detail by id
// @Tags rings
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetRingsByIDReply{}
// @Router /api/v1/rings/{id} [get]
// @Security BearerAuth
func (h *ringsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getRingsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	rings, err := h.iDao.GetByID(ctx, id)
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

	data := &types.RingsObjDetail{}
	err = copier.Copy(data, rings)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDRings)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"rings": data})
}

// List of records by query parameters
// @Summary list of ringss by query parameters
// @Description list of ringss by paging and conditions
// @Tags rings
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListRingssReply{}
// @Router /api/v1/rings/list [post]
// @Security BearerAuth
func (h *ringsHandler) List(c *gin.Context) {
	form := &types.ListRingssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	ringss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertRingss(ringss)
	if err != nil {
		response.Error(c, ecode.ErrListRings)
		return
	}

	response.Success(c, gin.H{
		"ringss": data,
		"total":  total,
	})
}

func getRingsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertRings(rings *model.Rings) (*types.RingsObjDetail, error) {
	data := &types.RingsObjDetail{}
	err := copier.Copy(data, rings)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertRingss(fromValues []*model.Rings) ([]*types.RingsObjDetail, error) {
	toValues := []*types.RingsObjDetail{}
	for _, v := range fromValues {
		data, err := convertRings(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
