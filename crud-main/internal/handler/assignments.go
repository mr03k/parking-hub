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

var _ AssignmentsHandler = (*assignmentsHandler)(nil)

// AssignmentsHandler defining the handler interface
type AssignmentsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type assignmentsHandler struct {
	iDao dao.AssignmentsDao
}

// NewAssignmentsHandler creating the handler interface
func NewAssignmentsHandler() AssignmentsHandler {
	return &assignmentsHandler{
		iDao: dao.NewAssignmentsDao(
			model.GetDB(),
			cache.NewAssignmentsCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create assignments
// @Description submit information to create assignments
// @Tags assignments
// @accept json
// @Produce json
// @Param data body types.CreateAssignmentsRequest true "assignments information"
// @Success 200 {object} types.CreateAssignmentsReply{}
// @Router /api/v1/assignments [post]
// @Security BearerAuth
func (h *assignmentsHandler) Create(c *gin.Context) {
	form := &types.CreateAssignmentsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	assignments := &model.Assignments{}
	err = copier.Copy(assignments, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateAssignments)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, assignments)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": assignments.ID})
}

// DeleteByID delete a record by id
// @Summary delete assignments
// @Description delete assignments by id
// @Tags assignments
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteAssignmentsByIDReply{}
// @Router /api/v1/assignments/{id} [delete]
// @Security BearerAuth
func (h *assignmentsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getAssignmentsIDFromPath(c)
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
// @Summary update assignments
// @Description update assignments information by id
// @Tags assignments
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateAssignmentsByIDRequest true "assignments information"
// @Success 200 {object} types.UpdateAssignmentsByIDReply{}
// @Router /api/v1/assignments/{id} [put]
// @Security BearerAuth
func (h *assignmentsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getAssignmentsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateAssignmentsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	assignments := &model.Assignments{}
	err = copier.Copy(assignments, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDAssignments)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, assignments)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get assignments detail
// @Description get assignments detail by id
// @Tags assignments
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetAssignmentsByIDReply{}
// @Router /api/v1/assignments/{id} [get]
// @Security BearerAuth
func (h *assignmentsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getAssignmentsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	assignments, err := h.iDao.GetByID(ctx, id)
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

	data := &types.AssignmentsObjDetail{}
	err = copier.Copy(data, assignments)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDAssignments)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"assignments": data})
}

// List of records by query parameters
// @Summary list of assignmentss by query parameters
// @Description list of assignmentss by paging and conditions
// @Tags assignments
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListAssignmentssReply{}
// @Router /api/v1/assignments/list [post]
// @Security BearerAuth
func (h *assignmentsHandler) List(c *gin.Context) {
	form := &types.ListAssignmentssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	assignmentss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertAssignmentss(assignmentss)
	if err != nil {
		response.Error(c, ecode.ErrListAssignments)
		return
	}

	response.Success(c, gin.H{
		"assignmentss": data,
		"total":        total,
	})
}

func getAssignmentsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertAssignments(assignments *model.Assignments) (*types.AssignmentsObjDetail, error) {
	data := &types.AssignmentsObjDetail{}
	err := copier.Copy(data, assignments)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertAssignmentss(fromValues []*model.Assignments) ([]*types.AssignmentsObjDetail, error) {
	toValues := []*types.AssignmentsObjDetail{}
	for _, v := range fromValues {
		data, err := convertAssignments(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
