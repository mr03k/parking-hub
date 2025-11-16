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

var _ CalendarsHandler = (*calendarsHandler)(nil)

// CalendarsHandler defining the handler interface
type CalendarsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type calendarsHandler struct {
	iDao dao.CalendarsDao
}

// NewCalendarsHandler creating the handler interface
func NewCalendarsHandler() CalendarsHandler {
	return &calendarsHandler{
		iDao: dao.NewCalendarsDao(
			model.GetDB(),
			cache.NewCalendarsCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create calendars
// @Description submit information to create calendars
// @Tags calendars
// @accept json
// @Produce json
// @Param data body types.CreateCalendarsRequest true "calendars information"
// @Success 200 {object} types.CreateCalendarsReply{}
// @Router /api/v1/calendars [post]
// @Security BearerAuth
func (h *calendarsHandler) Create(c *gin.Context) {
	form := &types.CreateCalendarsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	calendars := &model.Calendars{}
	err = copier.Copy(calendars, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateCalendars)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, calendars)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": calendars.ID})
}

// DeleteByID delete a record by id
// @Summary delete calendars
// @Description delete calendars by id
// @Tags calendars
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteCalendarsByIDReply{}
// @Router /api/v1/calendars/{id} [delete]
// @Security BearerAuth
func (h *calendarsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getCalendarsIDFromPath(c)
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
// @Summary update calendars
// @Description update calendars information by id
// @Tags calendars
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateCalendarsByIDRequest true "calendars information"
// @Success 200 {object} types.UpdateCalendarsByIDReply{}
// @Router /api/v1/calendars/{id} [put]
// @Security BearerAuth
func (h *calendarsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getCalendarsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateCalendarsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	calendars := &model.Calendars{}
	err = copier.Copy(calendars, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDCalendars)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, calendars)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get calendars detail
// @Description get calendars detail by id
// @Tags calendars
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetCalendarsByIDReply{}
// @Router /api/v1/calendars/{id} [get]
// @Security BearerAuth
func (h *calendarsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getCalendarsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	calendars, err := h.iDao.GetByID(ctx, id)
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

	data := &types.CalendarsObjDetail{}
	err = copier.Copy(data, calendars)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDCalendars)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"calendars": data})
}

// List of records by query parameters
// @Summary list of calendarss by query parameters
// @Description list of calendarss by paging and conditions
// @Tags calendars
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListCalendarssReply{}
// @Router /api/v1/calendars/list [post]
// @Security BearerAuth
func (h *calendarsHandler) List(c *gin.Context) {
	form := &types.ListCalendarssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	calendarss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertCalendarss(calendarss)
	if err != nil {
		response.Error(c, ecode.ErrListCalendars)
		return
	}

	response.Success(c, gin.H{
		"calendarss": data,
		"total":      total,
	})
}

func getCalendarsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertCalendars(calendars *model.Calendars) (*types.CalendarsObjDetail, error) {
	data := &types.CalendarsObjDetail{}
	err := copier.Copy(data, calendars)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertCalendarss(fromValues []*model.Calendars) ([]*types.CalendarsObjDetail, error) {
	toValues := []*types.CalendarsObjDetail{}
	for _, v := range fromValues {
		data, err := convertCalendars(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
