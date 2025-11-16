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

var _ SegmentsHandler = (*segmentsHandler)(nil)

// SegmentsHandler defining the handler interface
type SegmentsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type segmentsHandler struct {
	iDao dao.SegmentsDao
}

// NewSegmentsHandler creating the handler interface
func NewSegmentsHandler() SegmentsHandler {
	return &segmentsHandler{
		iDao: dao.NewSegmentsDao(
			model.GetDB(),
			cache.NewSegmentsCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create segments
// @Description submit information to create segments
// @Tags segments
// @accept json
// @Produce json
// @Param data body types.CreateSegmentsRequest true "segments information"
// @Success 200 {object} types.CreateSegmentsReply{}
// @Router /api/v1/segments [post]
// @Security BearerAuth
func (h *segmentsHandler) Create(c *gin.Context) {
	form := &types.CreateSegmentsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	segments := &model.Segments{}
	err = copier.Copy(segments, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateSegments)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, segments)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": segments.ID})
}

// DeleteByID delete a record by id
// @Summary delete segments
// @Description delete segments by id
// @Tags segments
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteSegmentsByIDReply{}
// @Router /api/v1/segments/{id} [delete]
// @Security BearerAuth
func (h *segmentsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getSegmentsIDFromPath(c)
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
// @Summary update segments
// @Description update segments information by id
// @Tags segments
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateSegmentsByIDRequest true "segments information"
// @Success 200 {object} types.UpdateSegmentsByIDReply{}
// @Router /api/v1/segments/{id} [put]
// @Security BearerAuth
func (h *segmentsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getSegmentsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateSegmentsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	segments := &model.Segments{}
	err = copier.Copy(segments, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDSegments)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, segments)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get segments detail
// @Description get segments detail by id
// @Tags segments
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetSegmentsByIDReply{}
// @Router /api/v1/segments/{id} [get]
// @Security BearerAuth
func (h *segmentsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getSegmentsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	segments, err := h.iDao.GetByID(ctx, id)
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

	data := &types.SegmentsObjDetail{}
	err = copier.Copy(data, segments)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDSegments)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"segments": data})
}

// List of records by query parameters
// @Summary list of segmentss by query parameters
// @Description list of segmentss by paging and conditions
// @Tags segments
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListSegmentssReply{}
// @Router /api/v1/segments/list [post]
// @Security BearerAuth
func (h *segmentsHandler) List(c *gin.Context) {
	form := &types.ListSegmentssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	segmentss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertSegmentss(segmentss)
	if err != nil {
		response.Error(c, ecode.ErrListSegments)
		return
	}

	response.Success(c, gin.H{
		"segmentss": data,
		"total":     total,
	})
}

func getSegmentsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertSegments(segments *model.Segments) (*types.SegmentsObjDetail, error) {
	data := &types.SegmentsObjDetail{}
	err := copier.Copy(data, segments)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertSegmentss(fromValues []*model.Segments) ([]*types.SegmentsObjDetail, error) {
	toValues := []*types.SegmentsObjDetail{}
	for _, v := range fromValues {
		data, err := convertSegments(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
