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

var _ FormsHandler = (*formsHandler)(nil)

// FormsHandler defining the handler interface
type FormsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type formsHandler struct {
	iDao dao.FormsDao
}

// NewFormsHandler creating the handler interface
func NewFormsHandler() FormsHandler {
	return &formsHandler{
		iDao: dao.NewFormsDao(
			model.GetDB(),
			cache.NewFormsCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create forms
// @Description submit information to create forms
// @Tags forms
// @accept json
// @Produce json
// @Param data body types.CreateFormsRequest true "forms information"
// @Success 200 {object} types.CreateFormsReply{}
// @Router /api/v1/forms [post]
// @Security BearerAuth
func (h *formsHandler) Create(c *gin.Context) {
	form := &types.CreateFormsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	forms := &model.Forms{}
	err = copier.Copy(forms, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateForms)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, forms)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": forms.ID})
}

// DeleteByID delete a record by id
// @Summary delete forms
// @Description delete forms by id
// @Tags forms
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteFormsByIDReply{}
// @Router /api/v1/forms/{id} [delete]
// @Security BearerAuth
func (h *formsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getFormsIDFromPath(c)
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
// @Summary update forms
// @Description update forms information by id
// @Tags forms
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateFormsByIDRequest true "forms information"
// @Success 200 {object} types.UpdateFormsByIDReply{}
// @Router /api/v1/forms/{id} [put]
// @Security BearerAuth
func (h *formsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getFormsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateFormsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	forms := &model.Forms{}
	err = copier.Copy(forms, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDForms)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, forms)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get forms detail
// @Description get forms detail by id
// @Tags forms
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetFormsByIDReply{}
// @Router /api/v1/forms/{id} [get]
// @Security BearerAuth
func (h *formsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getFormsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	forms, err := h.iDao.GetByID(ctx, id)
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

	data := &types.FormsObjDetail{}
	err = copier.Copy(data, forms)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDForms)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"forms": data})
}

// List of records by query parameters
// @Summary list of formss by query parameters
// @Description list of formss by paging and conditions
// @Tags forms
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListFormssReply{}
// @Router /api/v1/forms/list [post]
// @Security BearerAuth
func (h *formsHandler) List(c *gin.Context) {
	form := &types.ListFormssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	formss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertFormss(formss)
	if err != nil {
		response.Error(c, ecode.ErrListForms)
		return
	}

	response.Success(c, gin.H{
		"formss": data,
		"total":  total,
	})
}

func getFormsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertForms(forms *model.Forms) (*types.FormsObjDetail, error) {
	data := &types.FormsObjDetail{}
	err := copier.Copy(data, forms)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertFormss(fromValues []*model.Forms) ([]*types.FormsObjDetail, error) {
	toValues := []*types.FormsObjDetail{}
	for _, v := range fromValues {
		data, err := convertForms(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
