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

var _ CitiesHandler = (*citiesHandler)(nil)

// CitiesHandler defining the handler interface
type CitiesHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type citiesHandler struct {
	iDao dao.CitiesDao
}

// NewCitiesHandler creating the handler interface
func NewCitiesHandler() CitiesHandler {
	return &citiesHandler{
		iDao: dao.NewCitiesDao(
			model.GetDB(),
			cache.NewCitiesCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create cities
// @Description submit information to create cities
// @Tags cities
// @accept json
// @Produce json
// @Param data body types.CreateCitiesRequest true "cities information"
// @Success 200 {object} types.CreateCitiesReply{}
// @Router /api/v1/cities [post]
// @Security BearerAuth
func (h *citiesHandler) Create(c *gin.Context) {
	form := &types.CreateCitiesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	cities := &model.Cities{}
	err = copier.Copy(cities, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateCities)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, cities)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": cities.ID})
}

// DeleteByID delete a record by id
// @Summary delete cities
// @Description delete cities by id
// @Tags cities
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteCitiesByIDReply{}
// @Router /api/v1/cities/{id} [delete]
// @Security BearerAuth
func (h *citiesHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getCitiesIDFromPath(c)
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
// @Summary update cities
// @Description update cities information by id
// @Tags cities
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateCitiesByIDRequest true "cities information"
// @Success 200 {object} types.UpdateCitiesByIDReply{}
// @Router /api/v1/cities/{id} [put]
// @Security BearerAuth
func (h *citiesHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getCitiesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateCitiesByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	cities := &model.Cities{}
	err = copier.Copy(cities, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDCities)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, cities)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get cities detail
// @Description get cities detail by id
// @Tags cities
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetCitiesByIDReply{}
// @Router /api/v1/cities/{id} [get]
// @Security BearerAuth
func (h *citiesHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getCitiesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	cities, err := h.iDao.GetByID(ctx, id)
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

	data := &types.CitiesObjDetail{}
	err = copier.Copy(data, cities)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDCities)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"cities": data})
}

// List of records by query parameters
// @Summary list of citiess by query parameters
// @Description list of citiess by paging and conditions
// @Tags cities
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListCitiessReply{}
// @Router /api/v1/cities/list [post]
// @Security BearerAuth
func (h *citiesHandler) List(c *gin.Context) {
	form := &types.ListCitiessRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	citiess, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertCitiess(citiess)
	if err != nil {
		response.Error(c, ecode.ErrListCities)
		return
	}

	response.Success(c, gin.H{
		"citiess": data,
		"total":   total,
	})
}

func getCitiesIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertCities(cities *model.Cities) (*types.CitiesObjDetail, error) {
	data := &types.CitiesObjDetail{}
	err := copier.Copy(data, cities)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertCitiess(fromValues []*model.Cities) ([]*types.CitiesObjDetail, error) {
	toValues := []*types.CitiesObjDetail{}
	for _, v := range fromValues {
		data, err := convertCities(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
