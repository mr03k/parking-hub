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

var _ PartsHandler = (*partsHandler)(nil)

// PartsHandler defining the handler interface
type PartsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type partsHandler struct {
	iDao dao.PartsDao
}

// NewPartsHandler creating the handler interface
func NewPartsHandler() PartsHandler {
	return &partsHandler{
		iDao: dao.NewPartsDao(
			model.GetDB(),
			cache.NewPartsCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create parts
// @Description submit information to create parts
// @Tags parts
// @accept json
// @Produce json
// @Param data body types.CreatePartsRequest true "parts information"
// @Success 200 {object} types.CreatePartsReply{}
// @Router /api/v1/parts [post]
// @Security BearerAuth
func (h *partsHandler) Create(c *gin.Context) {
	form := &types.CreatePartsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	parts := &model.Parts{}
	err = copier.Copy(parts, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateParts)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, parts)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": parts.ID})
}

// DeleteByID delete a record by id
// @Summary delete parts
// @Description delete parts by id
// @Tags parts
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeletePartsByIDReply{}
// @Router /api/v1/parts/{id} [delete]
// @Security BearerAuth
func (h *partsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getPartsIDFromPath(c)
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
// @Summary update parts
// @Description update parts information by id
// @Tags parts
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdatePartsByIDRequest true "parts information"
// @Success 200 {object} types.UpdatePartsByIDReply{}
// @Router /api/v1/parts/{id} [put]
// @Security BearerAuth
func (h *partsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getPartsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdatePartsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	parts := &model.Parts{}
	err = copier.Copy(parts, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDParts)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, parts)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get parts detail
// @Description get parts detail by id
// @Tags parts
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetPartsByIDReply{}
// @Router /api/v1/parts/{id} [get]
// @Security BearerAuth
func (h *partsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getPartsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	parts, err := h.iDao.GetByID(ctx, id)
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

	data := &types.PartsObjDetail{}
	err = copier.Copy(data, parts)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDParts)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"parts": data})
}

// List of records by query parameters
// @Summary list of partss by query parameters
// @Description list of partss by paging and conditions
// @Tags parts
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListPartssReply{}
// @Router /api/v1/parts/list [post]
// @Security BearerAuth
func (h *partsHandler) List(c *gin.Context) {
	form := &types.ListPartssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	partss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertPartss(partss)
	if err != nil {
		response.Error(c, ecode.ErrListParts)
		return
	}

	response.Success(c, gin.H{
		"partss": data,
		"total":  total,
	})
}

func getPartsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertParts(parts *model.Parts) (*types.PartsObjDetail, error) {
	data := &types.PartsObjDetail{}
	err := copier.Copy(data, parts)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertPartss(fromValues []*model.Parts) ([]*types.PartsObjDetail, error) {
	toValues := []*types.PartsObjDetail{}
	for _, v := range fromValues {
		data, err := convertParts(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
