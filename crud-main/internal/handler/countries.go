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

var _ CountriesHandler = (*countriesHandler)(nil)

// CountriesHandler defining the handler interface
type CountriesHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type countriesHandler struct {
	iDao dao.CountriesDao
}

// NewCountriesHandler creating the handler interface
func NewCountriesHandler() CountriesHandler {
	return &countriesHandler{
		iDao: dao.NewCountriesDao(
			model.GetDB(),
			cache.NewCountriesCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create countries
// @Description submit information to create countries
// @Tags countries
// @accept json
// @Produce json
// @Param data body types.CreateCountriesRequest true "countries information"
// @Success 200 {object} types.CreateCountriesReply{}
// @Router /api/v1/countries [post]
// @Security BearerAuth
func (h *countriesHandler) Create(c *gin.Context) {
	form := &types.CreateCountriesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	countries := &model.Countries{}
	err = copier.Copy(countries, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateCountries)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, countries)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": countries.ID})
}

// DeleteByID delete a record by id
// @Summary delete countries
// @Description delete countries by id
// @Tags countries
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteCountriesByIDReply{}
// @Router /api/v1/countries/{id} [delete]
// @Security BearerAuth
func (h *countriesHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getCountriesIDFromPath(c)
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
// @Summary update countries
// @Description update countries information by id
// @Tags countries
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateCountriesByIDRequest true "countries information"
// @Success 200 {object} types.UpdateCountriesByIDReply{}
// @Router /api/v1/countries/{id} [put]
// @Security BearerAuth
func (h *countriesHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getCountriesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateCountriesByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	countries := &model.Countries{}
	err = copier.Copy(countries, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDCountries)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, countries)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get countries detail
// @Description get countries detail by id
// @Tags countries
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetCountriesByIDReply{}
// @Router /api/v1/countries/{id} [get]
// @Security BearerAuth
func (h *countriesHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getCountriesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	countries, err := h.iDao.GetByID(ctx, id)
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

	data := &types.CountriesObjDetail{}
	err = copier.Copy(data, countries)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDCountries)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"countries": data})
}

// List of records by query parameters
// @Summary list of countriess by query parameters
// @Description list of countriess by paging and conditions
// @Tags countries
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListCountriessReply{}
// @Router /api/v1/countries/list [post]
// @Security BearerAuth
func (h *countriesHandler) List(c *gin.Context) {
	form := &types.ListCountriessRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	countriess, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertCountriess(countriess)
	if err != nil {
		response.Error(c, ecode.ErrListCountries)
		return
	}

	response.Success(c, gin.H{
		"countriess": data,
		"total":      total,
	})
}

func getCountriesIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertCountries(countries *model.Countries) (*types.CountriesObjDetail, error) {
	data := &types.CountriesObjDetail{}
	err := copier.Copy(data, countries)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertCountriess(fromValues []*model.Countries) ([]*types.CountriesObjDetail, error) {
	toValues := []*types.CountriesObjDetail{}
	for _, v := range fromValues {
		data, err := convertCountries(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
