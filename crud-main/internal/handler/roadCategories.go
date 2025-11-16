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

var _ RoadCategoriesHandler = (*roadCategoriesHandler)(nil)

// RoadCategoriesHandler defining the handler interface
type RoadCategoriesHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type roadCategoriesHandler struct {
	iDao dao.RoadCategoriesDao
}

// NewRoadCategoriesHandler creating the handler interface
func NewRoadCategoriesHandler() RoadCategoriesHandler {
	return &roadCategoriesHandler{
		iDao: dao.NewRoadCategoriesDao(
			model.GetDB(),
			cache.NewRoadCategoriesCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create roadCategories
// @Description submit information to create roadCategories
// @Tags roadCategories
// @accept json
// @Produce json
// @Param data body types.CreateRoadCategoriesRequest true "roadCategories information"
// @Success 200 {object} types.CreateRoadCategoriesReply{}
// @Router /api/v1/roadCategories [post]
// @Security BearerAuth
func (h *roadCategoriesHandler) Create(c *gin.Context) {
	form := &types.CreateRoadCategoriesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	roadCategories := &model.RoadCategories{}
	err = copier.Copy(roadCategories, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateRoadCategories)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, roadCategories)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": roadCategories.ID})
}

// DeleteByID delete a record by id
// @Summary delete roadCategories
// @Description delete roadCategories by id
// @Tags roadCategories
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteRoadCategoriesByIDReply{}
// @Router /api/v1/roadCategories/{id} [delete]
// @Security BearerAuth
func (h *roadCategoriesHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getRoadCategoriesIDFromPath(c)
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
// @Summary update roadCategories
// @Description update roadCategories information by id
// @Tags roadCategories
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateRoadCategoriesByIDRequest true "roadCategories information"
// @Success 200 {object} types.UpdateRoadCategoriesByIDReply{}
// @Router /api/v1/roadCategories/{id} [put]
// @Security BearerAuth
func (h *roadCategoriesHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getRoadCategoriesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateRoadCategoriesByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	roadCategories := &model.RoadCategories{}
	err = copier.Copy(roadCategories, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDRoadCategories)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, roadCategories)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get roadCategories detail
// @Description get roadCategories detail by id
// @Tags roadCategories
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetRoadCategoriesByIDReply{}
// @Router /api/v1/roadCategories/{id} [get]
// @Security BearerAuth
func (h *roadCategoriesHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getRoadCategoriesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	roadCategories, err := h.iDao.GetByID(ctx, id)
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

	data := &types.RoadCategoriesObjDetail{}
	err = copier.Copy(data, roadCategories)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDRoadCategories)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"roadCategories": data})
}

// List of records by query parameters
// @Summary list of roadCategoriess by query parameters
// @Description list of roadCategoriess by paging and conditions
// @Tags roadCategories
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListRoadCategoriessReply{}
// @Router /api/v1/roadCategories/list [post]
// @Security BearerAuth
func (h *roadCategoriesHandler) List(c *gin.Context) {
	form := &types.ListRoadCategoriessRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	roadCategoriess, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertRoadCategoriess(roadCategoriess)
	if err != nil {
		response.Error(c, ecode.ErrListRoadCategories)
		return
	}

	response.Success(c, gin.H{
		"roadCategoriess": data,
		"total":           total,
	})
}

func getRoadCategoriesIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertRoadCategories(roadCategories *model.RoadCategories) (*types.RoadCategoriesObjDetail, error) {
	data := &types.RoadCategoriesObjDetail{}
	err := copier.Copy(data, roadCategories)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertRoadCategoriess(fromValues []*model.RoadCategories) ([]*types.RoadCategoriesObjDetail, error) {
	toValues := []*types.RoadCategoriesObjDetail{}
	for _, v := range fromValues {
		data, err := convertRoadCategories(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
