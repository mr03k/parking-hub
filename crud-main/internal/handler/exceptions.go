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

var _ ExceptionsHandler = (*exceptionsHandler)(nil)

// ExceptionsHandler defining the handler interface
type ExceptionsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type exceptionsHandler struct {
	iDao dao.ExceptionsDao
}

// NewExceptionsHandler creating the handler interface
func NewExceptionsHandler() ExceptionsHandler {
	return &exceptionsHandler{
		iDao: dao.NewExceptionsDao(
			model.GetDB(),
			cache.NewExceptionsCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create exceptions
// @Description submit information to create exceptions
// @Tags exceptions
// @accept json
// @Produce json
// @Param data body types.CreateExceptionsRequest true "exceptions information"
// @Success 200 {object} types.CreateExceptionsReply{}
// @Router /api/v1/exceptions [post]
// @Security BearerAuth
func (h *exceptionsHandler) Create(c *gin.Context) {
	form := &types.CreateExceptionsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	exceptions := &model.Exceptions{}
	err = copier.Copy(exceptions, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateExceptions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, exceptions)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": exceptions.ID})
}

// DeleteByID delete a record by id
// @Summary delete exceptions
// @Description delete exceptions by id
// @Tags exceptions
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteExceptionsByIDReply{}
// @Router /api/v1/exceptions/{id} [delete]
// @Security BearerAuth
func (h *exceptionsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getExceptionsIDFromPath(c)
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
// @Summary update exceptions
// @Description update exceptions information by id
// @Tags exceptions
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateExceptionsByIDRequest true "exceptions information"
// @Success 200 {object} types.UpdateExceptionsByIDReply{}
// @Router /api/v1/exceptions/{id} [put]
// @Security BearerAuth
func (h *exceptionsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getExceptionsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateExceptionsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	exceptions := &model.Exceptions{}
	err = copier.Copy(exceptions, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDExceptions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, exceptions)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get exceptions detail
// @Description get exceptions detail by id
// @Tags exceptions
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetExceptionsByIDReply{}
// @Router /api/v1/exceptions/{id} [get]
// @Security BearerAuth
func (h *exceptionsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getExceptionsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	exceptions, err := h.iDao.GetByID(ctx, id)
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

	data := &types.ExceptionsObjDetail{}
	err = copier.Copy(data, exceptions)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDExceptions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"exceptions": data})
}

// List of records by query parameters
// @Summary list of exceptionss by query parameters
// @Description list of exceptionss by paging and conditions
// @Tags exceptions
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListExceptionssReply{}
// @Router /api/v1/exceptions/list [post]
// @Security BearerAuth
func (h *exceptionsHandler) List(c *gin.Context) {
	form := &types.ListExceptionssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	exceptionss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertExceptionss(exceptionss)
	if err != nil {
		response.Error(c, ecode.ErrListExceptions)
		return
	}

	response.Success(c, gin.H{
		"exceptionss": data,
		"total":       total,
	})
}

func getExceptionsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertExceptions(exceptions *model.Exceptions) (*types.ExceptionsObjDetail, error) {
	data := &types.ExceptionsObjDetail{}
	err := copier.Copy(data, exceptions)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertExceptionss(fromValues []*model.Exceptions) ([]*types.ExceptionsObjDetail, error) {
	toValues := []*types.ExceptionsObjDetail{}
	for _, v := range fromValues {
		data, err := convertExceptions(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
