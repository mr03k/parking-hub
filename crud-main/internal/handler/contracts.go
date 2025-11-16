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

var _ ContractsHandler = (*contractsHandler)(nil)

// ContractsHandler defining the handler interface
type ContractsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type contractsHandler struct {
	iDao dao.ContractsDao
}

// NewContractsHandler creating the handler interface
func NewContractsHandler() ContractsHandler {
	return &contractsHandler{
		iDao: dao.NewContractsDao(
			model.GetDB(),
			cache.NewContractsCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create contracts
// @Description submit information to create contracts
// @Tags contracts
// @accept json
// @Produce json
// @Param data body types.CreateContractsRequest true "contracts information"
// @Success 200 {object} types.CreateContractsReply{}
// @Router /api/v1/contracts [post]
// @Security BearerAuth
func (h *contractsHandler) Create(c *gin.Context) {
	form := &types.CreateContractsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	contracts := &model.Contracts{}
	err = copier.Copy(contracts, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateContracts)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, contracts)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": contracts.ID})
}

// DeleteByID delete a record by id
// @Summary delete contracts
// @Description delete contracts by id
// @Tags contracts
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteContractsByIDReply{}
// @Router /api/v1/contracts/{id} [delete]
// @Security BearerAuth
func (h *contractsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getContractsIDFromPath(c)
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
// @Summary update contracts
// @Description update contracts information by id
// @Tags contracts
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateContractsByIDRequest true "contracts information"
// @Success 200 {object} types.UpdateContractsByIDReply{}
// @Router /api/v1/contracts/{id} [put]
// @Security BearerAuth
func (h *contractsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getContractsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateContractsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	contracts := &model.Contracts{}
	err = copier.Copy(contracts, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDContracts)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, contracts)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get contracts detail
// @Description get contracts detail by id
// @Tags contracts
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetContractsByIDReply{}
// @Router /api/v1/contracts/{id} [get]
// @Security BearerAuth
func (h *contractsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getContractsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	contracts, err := h.iDao.GetByID(ctx, id)
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

	data := &types.ContractsObjDetail{}
	err = copier.Copy(data, contracts)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDContracts)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"contracts": data})
}

// List of records by query parameters
// @Summary list of contractss by query parameters
// @Description list of contractss by paging and conditions
// @Tags contracts
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListContractssReply{}
// @Router /api/v1/contracts/list [post]
// @Security BearerAuth
func (h *contractsHandler) List(c *gin.Context) {
	form := &types.ListContractssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	contractss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertContractss(contractss)
	if err != nil {
		response.Error(c, ecode.ErrListContracts)
		return
	}

	response.Success(c, gin.H{
		"contractss": data,
		"total":      total,
	})
}

func getContractsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertContracts(contracts *model.Contracts) (*types.ContractsObjDetail, error) {
	data := &types.ContractsObjDetail{}
	err := copier.Copy(data, contracts)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertContractss(fromValues []*model.Contracts) ([]*types.ContractsObjDetail, error) {
	toValues := []*types.ContractsObjDetail{}
	for _, v := range fromValues {
		data, err := convertContracts(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
