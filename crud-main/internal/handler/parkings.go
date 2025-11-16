package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"

	"git.abanppc.com/farin-project/crud/internal/cache"
	"git.abanppc.com/farin-project/crud/internal/dao"
	"git.abanppc.com/farin-project/crud/internal/ecode"
	"git.abanppc.com/farin-project/crud/internal/model"
	"git.abanppc.com/farin-project/crud/internal/types"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/logger"
)

var _ ParkingsHandler = (*parkingsHandler)(nil)

// ParkingsHandler defining the handler interface
type ParkingsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type parkingsHandler struct {
	iDao dao.ParkingsDao
}

// NewParkingsHandler creating the handler interface
func NewParkingsHandler() ParkingsHandler {
	return &parkingsHandler{
		iDao: dao.NewParkingsDao(
			model.GetDB(),
			cache.NewParkingsCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create parkings
// @Description submit information to create parkings
// @Tags parkings
// @accept json
// @Produce json
// @Param data body types.CreateParkingsRequest true "parkings information"
// @Success 200 {object} types.CreateParkingsReply{}
// @Router /api/v1/parkings [post]
// @Security BearerAuth
func (h *parkingsHandler) Create(c *gin.Context) {
	form := &types.CreateParkingsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	parkings := &model.Parkings{}
	err = copier.Copy(parkings, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateParkings)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, parkings)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": parkings.ID})
}

// DeleteByID delete a record by id
// @Summary delete parkings
// @Description delete parkings by id
// @Tags parkings
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteParkingsByIDReply{}
// @Router /api/v1/parkings/{id} [delete]
// @Security BearerAuth
func (h *parkingsHandler) DeleteByID(c *gin.Context) {
	id := c.Param("id")

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
// @Summary update parkings
// @Description update parkings information by id
// @Tags parkings
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateParkingsByIDRequest true "parkings information"
// @Success 200 {object} types.UpdateParkingsByIDReply{}
// @Router /api/v1/parkings/{id} [put]
// @Security BearerAuth
func (h *parkingsHandler) UpdateByID(c *gin.Context) {
	id := c.Param("id")

	form := &types.UpdateParkingsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	parkings := &model.Parkings{}
	err = copier.Copy(parkings, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDParkings)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, parkings)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get parkings detail
// @Description get parkings detail by id
// @Tags parkings
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetParkingsByIDReply{}
// @Router /api/v1/parkings/{id} [get]
// @Security BearerAuth
func (h *parkingsHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	ctx := middleware.WrapCtx(c)
	parkings, err := h.iDao.GetByID(ctx, id)
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

	data := &types.ParkingsObjDetail{}
	err = copier.Copy(data, parkings)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDParkings)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"parkings": data})
}

// List of records by query parameters
// @Summary list of parkingss by query parameters
// @Description list of parkingss by paging and conditions
// @Tags parkings
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListParkingssReply{}
// @Router /api/v1/parkings/list [post]
// @Security BearerAuth
func (h *parkingsHandler) List(c *gin.Context) {
	form := &types.ListParkingssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	parkingss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertParkingss(parkingss)
	if err != nil {
		response.Error(c, ecode.ErrListParkings)
		return
	}

	response.Success(c, gin.H{
		"parkingss": data,
		"total":     total,
	})
}

func convertParkings(parkings *model.Parkings) (*types.ParkingsObjDetail, error) {
	data := &types.ParkingsObjDetail{}
	err := copier.Copy(data, parkings)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertParkingss(fromValues []*model.Parkings) ([]*types.ParkingsObjDetail, error) {
	toValues := []*types.ParkingsObjDetail{}
	for _, v := range fromValues {
		data, err := convertParkings(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
