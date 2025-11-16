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

var _ RoadsHandler = (*roadsHandler)(nil)

// RoadsHandler defining the handler interface
type RoadsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type roadsHandler struct {
	iDao dao.RoadsDao
}

// NewRoadsHandler creating the handler interface
func NewRoadsHandler() RoadsHandler {
	return &roadsHandler{
		iDao: dao.NewRoadsDao(
			model.GetDB(),
			cache.NewRoadsCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create roads
// @Description submit information to create roads
// @Tags roads
// @accept json
// @Produce json
// @Param data body types.CreateRoadsRequest true "roads information"
// @Success 200 {object} types.CreateRoadsReply{}
// @Router /api/v1/roads [post]
// @Security BearerAuth
func (h *roadsHandler) Create(c *gin.Context) {
	form := &types.CreateRoadsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	roads := &model.Roads{}
	err = copier.Copy(roads, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateRoads)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, roads)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": roads.ID})
}

// DeleteByID delete a record by id
// @Summary delete roads
// @Description delete roads by id
// @Tags roads
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteRoadsByIDReply{}
// @Router /api/v1/roads/{id} [delete]
// @Security BearerAuth
func (h *roadsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getRoadsIDFromPath(c)
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
// @Summary update roads
// @Description update roads information by id
// @Tags roads
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateRoadsByIDRequest true "roads information"
// @Success 200 {object} types.UpdateRoadsByIDReply{}
// @Router /api/v1/roads/{id} [put]
// @Security BearerAuth
func (h *roadsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getRoadsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateRoadsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	roads := &model.Roads{}
	err = copier.Copy(roads, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDRoads)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, roads)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get roads detail
// @Description get roads detail by id
// @Tags roads
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetRoadsByIDReply{}
// @Router /api/v1/roads/{id} [get]
// @Security BearerAuth
func (h *roadsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getRoadsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	roads, err := h.iDao.GetByID(ctx, id)
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

	data := &types.RoadsObjDetail{}
	err = copier.Copy(data, roads)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDRoads)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"roads": data})
}

// List of records by query parameters
// @Summary list of roadss by query parameters
// @Description list of roadss by paging and conditions
// @Tags roads
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListRoadssReply{}
// @Router /api/v1/roads/list [post]
// @Security BearerAuth
func (h *roadsHandler) List(c *gin.Context) {
	form := &types.ListRoadssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	roadss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertRoadss(roadss)
	if err != nil {
		response.Error(c, ecode.ErrListRoads)
		return
	}

	response.Success(c, gin.H{
		"roadss": data,
		"total":  total,
	})
}

func getRoadsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertRoads(roads *model.Roads) (*types.RoadsObjDetail, error) {
	data := &types.RoadsObjDetail{}
	err := copier.Copy(data, roads)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertRoadss(fromValues []*model.Roads) ([]*types.RoadsObjDetail, error) {
	toValues := []*types.RoadsObjDetail{}
	for _, v := range fromValues {
		data, err := convertRoads(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
