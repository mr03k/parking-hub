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

var _ ModulesHandler = (*modulesHandler)(nil)

// ModulesHandler defining the handler interface
type ModulesHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type modulesHandler struct {
	iDao dao.ModulesDao
}

// NewModulesHandler creating the handler interface
func NewModulesHandler() ModulesHandler {
	return &modulesHandler{
		iDao: dao.NewModulesDao(
			model.GetDB(),
			cache.NewModulesCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create modules
// @Description submit information to create modules
// @Tags modules
// @accept json
// @Produce json
// @Param data body types.CreateModulesRequest true "modules information"
// @Success 200 {object} types.CreateModulesReply{}
// @Router /api/v1/modules [post]
// @Security BearerAuth
func (h *modulesHandler) Create(c *gin.Context) {
	form := &types.CreateModulesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	modules := &model.Modules{}
	err = copier.Copy(modules, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateModules)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, modules)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": modules.ID})
}

// DeleteByID delete a record by id
// @Summary delete modules
// @Description delete modules by id
// @Tags modules
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteModulesByIDReply{}
// @Router /api/v1/modules/{id} [delete]
// @Security BearerAuth
func (h *modulesHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getModulesIDFromPath(c)
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
// @Summary update modules
// @Description update modules information by id
// @Tags modules
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateModulesByIDRequest true "modules information"
// @Success 200 {object} types.UpdateModulesByIDReply{}
// @Router /api/v1/modules/{id} [put]
// @Security BearerAuth
func (h *modulesHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getModulesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateModulesByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	modules := &model.Modules{}
	err = copier.Copy(modules, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDModules)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, modules)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get modules detail
// @Description get modules detail by id
// @Tags modules
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetModulesByIDReply{}
// @Router /api/v1/modules/{id} [get]
// @Security BearerAuth
func (h *modulesHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getModulesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	modules, err := h.iDao.GetByID(ctx, id)
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

	data := &types.ModulesObjDetail{}
	err = copier.Copy(data, modules)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDModules)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"modules": data})
}

// List of records by query parameters
// @Summary list of moduless by query parameters
// @Description list of moduless by paging and conditions
// @Tags modules
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListModulessReply{}
// @Router /api/v1/modules/list [post]
// @Security BearerAuth
func (h *modulesHandler) List(c *gin.Context) {
	form := &types.ListModulessRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	moduless, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertModuless(moduless)
	if err != nil {
		response.Error(c, ecode.ErrListModules)
		return
	}

	response.Success(c, gin.H{
		"moduless": data,
		"total":    total,
	})
}

func getModulesIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertModules(modules *model.Modules) (*types.ModulesObjDetail, error) {
	data := &types.ModulesObjDetail{}
	err := copier.Copy(data, modules)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertModuless(fromValues []*model.Modules) ([]*types.ModulesObjDetail, error) {
	toValues := []*types.ModulesObjDetail{}
	for _, v := range fromValues {
		data, err := convertModules(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
