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

var _ VehicleCategoriesHandler = (*vehicleCategoriesHandler)(nil)

// VehicleCategoriesHandler defining the handler interface
type VehicleCategoriesHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type vehicleCategoriesHandler struct {
	iDao dao.VehicleCategoriesDao
}

// NewVehicleCategoriesHandler creating the handler interface
func NewVehicleCategoriesHandler() VehicleCategoriesHandler {
	return &vehicleCategoriesHandler{
		iDao: dao.NewVehicleCategoriesDao(
			model.GetDB(),
			cache.NewVehicleCategoriesCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create vehicleCategories
// @Description submit information to create vehicleCategories
// @Tags vehicleCategories
// @accept json
// @Produce json
// @Param data body types.CreateVehicleCategoriesRequest true "vehicleCategories information"
// @Success 200 {object} types.CreateVehicleCategoriesReply{}
// @Router /api/v1/vehicleCategories [post]
// @Security BearerAuth
func (h *vehicleCategoriesHandler) Create(c *gin.Context) {
	form := &types.CreateVehicleCategoriesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	vehicleCategories := &model.VehicleCategories{}
	err = copier.Copy(vehicleCategories, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateVehicleCategories)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, vehicleCategories)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": vehicleCategories.ID})
}

// DeleteByID delete a record by id
// @Summary delete vehicleCategories
// @Description delete vehicleCategories by id
// @Tags vehicleCategories
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteVehicleCategoriesByIDReply{}
// @Router /api/v1/vehicleCategories/{id} [delete]
// @Security BearerAuth
func (h *vehicleCategoriesHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getVehicleCategoriesIDFromPath(c)
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
// @Summary update vehicleCategories
// @Description update vehicleCategories information by id
// @Tags vehicleCategories
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateVehicleCategoriesByIDRequest true "vehicleCategories information"
// @Success 200 {object} types.UpdateVehicleCategoriesByIDReply{}
// @Router /api/v1/vehicleCategories/{id} [put]
// @Security BearerAuth
func (h *vehicleCategoriesHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getVehicleCategoriesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateVehicleCategoriesByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	vehicleCategories := &model.VehicleCategories{}
	err = copier.Copy(vehicleCategories, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDVehicleCategories)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, vehicleCategories)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get vehicleCategories detail
// @Description get vehicleCategories detail by id
// @Tags vehicleCategories
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetVehicleCategoriesByIDReply{}
// @Router /api/v1/vehicleCategories/{id} [get]
// @Security BearerAuth
func (h *vehicleCategoriesHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getVehicleCategoriesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	vehicleCategories, err := h.iDao.GetByID(ctx, id)
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

	data := &types.VehicleCategoriesObjDetail{}
	err = copier.Copy(data, vehicleCategories)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDVehicleCategories)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"vehicleCategories": data})
}

// List of records by query parameters
// @Summary list of vehicleCategoriess by query parameters
// @Description list of vehicleCategoriess by paging and conditions
// @Tags vehicleCategories
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListVehicleCategoriessReply{}
// @Router /api/v1/vehicleCategories/list [post]
// @Security BearerAuth
func (h *vehicleCategoriesHandler) List(c *gin.Context) {
	form := &types.ListVehicleCategoriessRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	vehicleCategoriess, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertVehicleCategoriess(vehicleCategoriess)
	if err != nil {
		response.Error(c, ecode.ErrListVehicleCategories)
		return
	}

	response.Success(c, gin.H{
		"vehicleCategoriess": data,
		"total":              total,
	})
}

func getVehicleCategoriesIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertVehicleCategories(vehicleCategories *model.VehicleCategories) (*types.VehicleCategoriesObjDetail, error) {
	data := &types.VehicleCategoriesObjDetail{}
	err := copier.Copy(data, vehicleCategories)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertVehicleCategoriess(fromValues []*model.VehicleCategories) ([]*types.VehicleCategoriesObjDetail, error) {
	toValues := []*types.VehicleCategoriesObjDetail{}
	for _, v := range fromValues {
		data, err := convertVehicleCategories(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
