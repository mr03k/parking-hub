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

var _ DistrictsHandler = (*districtsHandler)(nil)

// DistrictsHandler defining the handler interface
type DistrictsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type districtsHandler struct {
	iDao dao.DistrictsDao
}

// NewDistrictsHandler creating the handler interface
func NewDistrictsHandler() DistrictsHandler {
	return &districtsHandler{
		iDao: dao.NewDistrictsDao(
			model.GetDB(),
			cache.NewDistrictsCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create districts
// @Description submit information to create districts
// @Tags districts
// @accept json
// @Produce json
// @Param data body types.CreateDistrictsRequest true "districts information"
// @Success 200 {object} types.CreateDistrictsReply{}
// @Router /api/v1/districts [post]
// @Security BearerAuth
func (h *districtsHandler) Create(c *gin.Context) {
	form := &types.CreateDistrictsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	districts := &model.Districts{}
	err = copier.Copy(districts, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateDistricts)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, districts)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": districts.ID})
}

// DeleteByID delete a record by id
// @Summary delete districts
// @Description delete districts by id
// @Tags districts
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteDistrictsByIDReply{}
// @Router /api/v1/districts/{id} [delete]
// @Security BearerAuth
func (h *districtsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getDistrictsIDFromPath(c)
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
// @Summary update districts
// @Description update districts information by id
// @Tags districts
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateDistrictsByIDRequest true "districts information"
// @Success 200 {object} types.UpdateDistrictsByIDReply{}
// @Router /api/v1/districts/{id} [put]
// @Security BearerAuth
func (h *districtsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getDistrictsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateDistrictsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	districts := &model.Districts{}
	err = copier.Copy(districts, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDDistricts)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, districts)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get districts detail
// @Description get districts detail by id
// @Tags districts
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetDistrictsByIDReply{}
// @Router /api/v1/districts/{id} [get]
// @Security BearerAuth
func (h *districtsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getDistrictsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	districts, err := h.iDao.GetByID(ctx, id)
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

	data := &types.DistrictsObjDetail{}
	err = copier.Copy(data, districts)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDDistricts)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"districts": data})
}

// List of records by query parameters
// @Summary list of districtss by query parameters
// @Description list of districtss by paging and conditions
// @Tags districts
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListDistrictssReply{}
// @Router /api/v1/districts/list [post]
// @Security BearerAuth
func (h *districtsHandler) List(c *gin.Context) {
	form := &types.ListDistrictssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	districtss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertDistrictss(districtss)
	if err != nil {
		response.Error(c, ecode.ErrListDistricts)
		return
	}

	response.Success(c, gin.H{
		"districtss": data,
		"total":      total,
	})
}

func getDistrictsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertDistricts(districts *model.Districts) (*types.DistrictsObjDetail, error) {
	data := &types.DistrictsObjDetail{}
	err := copier.Copy(data, districts)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertDistrictss(fromValues []*model.Districts) ([]*types.DistrictsObjDetail, error) {
	toValues := []*types.DistrictsObjDetail{}
	for _, v := range fromValues {
		data, err := convertDistricts(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
