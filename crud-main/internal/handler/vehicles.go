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

var _ VehiclesHandler = (*vehiclesHandler)(nil)

// VehiclesHandler defining the handler interface
type VehiclesHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type vehiclesHandler struct {
	iDao dao.VehiclesDao
}

// NewVehiclesHandler creating the handler interface
func NewVehiclesHandler() VehiclesHandler {
	return &vehiclesHandler{
		iDao: dao.NewVehiclesDao(
			model.GetDB(),
			cache.NewVehiclesCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create vehicles
// @Description submit information to create vehicles
// @Tags vehicles
// @accept json
// @Produce json
// @Param data body types.CreateVehiclesRequest true "vehicles information"
// @Success 200 {object} types.CreateVehiclesReply{}
// @Router /api/v1/vehicles [post]
// @Security BearerAuth
func (h *vehiclesHandler) Create(c *gin.Context) {
	form := &types.CreateVehiclesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	vehicles := &model.Vehicles{}
	err = copier.Copy(vehicles, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateVehicles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, vehicles)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": vehicles.ID})
}

// DeleteByID delete a record by id
// @Summary delete vehicles
// @Description delete vehicles by id
// @Tags vehicles
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteVehiclesByIDReply{}
// @Router /api/v1/vehicles/{id} [delete]
// @Security BearerAuth
func (h *vehiclesHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getVehiclesIDFromPath(c)
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
// @Summary update vehicles
// @Description update vehicles information by id
// @Tags vehicles
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateVehiclesByIDRequest true "vehicles information"
// @Success 200 {object} types.UpdateVehiclesByIDReply{}
// @Router /api/v1/vehicles/{id} [put]
// @Security BearerAuth
func (h *vehiclesHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getVehiclesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateVehiclesByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	vehicles := &model.Vehicles{}
	err = copier.Copy(vehicles, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDVehicles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, vehicles)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get vehicles detail
// @Description get vehicles detail by id
// @Tags vehicles
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetVehiclesByIDReply{}
// @Router /api/v1/vehicles/{id} [get]
// @Security BearerAuth
func (h *vehiclesHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getVehiclesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	vehicles, err := h.iDao.GetByID(ctx, id)
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

	data := &types.VehiclesObjDetail{}
	err = copier.Copy(data, vehicles)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDVehicles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"vehicles": data})
}

// List of records by query parameters
// @Summary list of vehicless by query parameters
// @Description list of vehicless by paging and conditions
// @Tags vehicles
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListVehiclessReply{}
// @Router /api/v1/vehicles/list [post]
// @Security BearerAuth
func (h *vehiclesHandler) List(c *gin.Context) {
	form := &types.ListVehiclessRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	vehicless, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertVehicless(vehicless)
	if err != nil {
		response.Error(c, ecode.ErrListVehicles)
		return
	}

	response.Success(c, gin.H{
		"vehicless": data,
		"total":     total,
	})
}

func getVehiclesIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertVehicles(vehicles *model.Vehicles) (*types.VehiclesObjDetail, error) {
	data := &types.VehiclesObjDetail{}
	err := copier.Copy(data, vehicles)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertVehicless(fromValues []*model.Vehicles) ([]*types.VehiclesObjDetail, error) {
	toValues := []*types.VehiclesObjDetail{}
	for _, v := range fromValues {
		data, err := convertVehicles(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
