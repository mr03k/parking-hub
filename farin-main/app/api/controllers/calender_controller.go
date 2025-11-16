package controller

import (
	"context"
	"farin/app/api/response"
	"farin/domain/dto"
	"farin/domain/service"
	"farin/infrastructure/godotenv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"strconv"
)

type CalenderController struct {
	calenderService *service.CalenderService
	logger          *slog.Logger
	env             *godotenv.Env
}

func NewCalenderController(logger *slog.Logger, us *service.CalenderService, env *godotenv.Env) *CalenderController {
	return &CalenderController{
		calenderService: us,
		logger:          logger.With("layer", "CalenderController"),
		env:             env,
	}
}

// CreateCalender godoc
// @Summary      Create a new calender
// @Description  Create a new calender by providing calender details
// @Tags         calenders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        calender  body      dto.CalenderRequest  true  "Calender Data"
// @Success      201   {object}  response.Response[dto.CalenderResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /calenders [post]
func (uc *CalenderController) CreateCalender(c *gin.Context) {
	lg := uc.logger.With("method", "CreateCalender")
	var calenderRequest dto.CalenderRequest

	if err := c.ShouldBindJSON(&calenderRequest); err != nil {
		lg.Error("failed to bind calender data", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	createdCalender, err := uc.calenderService.CreateCalender(context.Background(), calenderRequest.ToEntity())
	if err != nil {
		lg.Error("failed to create calender", "error", err)
		response.InternalError(c)
		return
	}
	calenderResponse := dto.CalenderResponse{}
	calenderResponse.FromEntity(createdCalender)

	response.Created(c, calenderResponse)
}

// ListCalenders godoc
// @Summary      List calenders
// @Description  Retrieve a list of calenders with optional filters and pagination
// @Tags         calenders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        sortField  query     string  false  "Field to sort by"  default(created_at)
// @Param        sortOrder  query     string  false  "Sort order (asc/desc)"  default(asc)
// @Param        page       query     int     false  "Page number"  default(1)
// @Param        pageSize   query     int     false  "Page size"  default(10)
// @Success      200        {object}  response.Response[dto.CalenderListResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /calenders [get]
func (uc *CalenderController) ListCalenders(c *gin.Context) {
	lg := uc.logger.With("method", "ListCalenders")
	filters := make(map[string]interface{})
	sortField := c.DefaultQuery("sortField", "created_at")
	sortOrder := c.DefaultQuery("sortOrder", "asc")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		response.BadRequest(c, "invalid page param")
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		response.BadRequest(c, "invalid page size param")
	}

	calenders, total, err := uc.calenderService.ListCalenders(context.Background(), filters, sortField, sortOrder, page, pageSize)
	if err != nil {
		lg.Error("failed to list calenders", "error", err)
		response.InternalError(c)
		return
	}
	calenderResponses := make([]dto.CalenderResponse, len(calenders))

	for i, calender := range calenders {
		calenderResponse := dto.CalenderResponse{}
		calenderResponse.FromEntity(&calender)
		calenderResponses[i] = calenderResponse
	}

	response.Ok(c, dto.CalenderListResponse{
		Calenders: calenderResponses,
		Total:     total,
	}, "")
}

// UpdateCalender godoc
// @Summary      Update calender details
// @Description  Update an existing calender's details
// @Tags         calenders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        calender  body      dto.CalenderRequest  true  "Updated Calender Data"
// @Param        id  path      string  true  "id"
// @Success      200   {object}  response.Response[dto.CalenderResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /calenders/{id} [put]
func (uc *CalenderController) UpdateCalender(c *gin.Context) {
	lg := uc.logger.With("method", "UpdateCalender")
	var calenderRequest dto.CalenderRequest

	if err := c.ShouldBindJSON(&calenderRequest); err != nil {
		lg.Error("failed to bind calender data", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	calender := calenderRequest.ToEntity()
	calender.ID = c.Param("id")
	updatedCalender, err := uc.calenderService.UpdateCalender(context.Background(), calender)
	if err != nil {
		lg.Error("failed to update calender", "error", err)
		response.InternalError(c)
		return
	}

	calenderResponse := dto.CalenderResponse{}
	calenderResponse.FromEntity(updatedCalender)

	response.Ok(c, calenderResponse, "")
}

// DeleteCalender godoc
// @Summary      Delete a calender
// @Description  Delete a calender by ID
// @Tags         calenders
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      string  true  "Calender ID"
// @Success      200   {object}  response.Response[swagger.EmptyObject]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Failure      404   {object}  response.Response[swagger.EmptyObject]
// @Router       /calenders/{id} [delete]
func (uc *CalenderController) DeleteCalender(c *gin.Context) {
	lg := uc.logger.With("method", "DeleteCalender")
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		lg.Error("invalid calender ID", "calenderID", id)
		response.BadRequest(c, "invalid request body")
		return
	}

	err := uc.calenderService.DeleteCalender(context.Background(), id)
	if err != nil {
		lg.Error("failed to delete calender", "error", err)
		response.InternalError(c)
		return
	}

	response.Ok(c, nil, "calender deleted successfully")
}

// GetCalenderDetail godoc
// @Summary      Get calender details
// @Description  Retrieve calender details by ID
// @Tags         calenders
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      string  true  "Calender ID"
// @Success      200   {object}  response.Response[dto.CalenderResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      404   {object}  response.Response[swagger.EmptyObject]
// @Router       /calenders/{id} [get]
func (uc *CalenderController) GetCalenderDetail(c *gin.Context) {
	lg := uc.logger.With("method", "GetCalenderDetail")
	id := c.Param("id")

	calender, err := uc.calenderService.Detail(context.Background(), "id", id)
	if err != nil {
		lg.Error("failed to get calender details", "error", err)
		response.NotFound(c)
		return
	}
	var calenderResponse dto.CalenderResponse
	calenderResponse.FromEntity(calender)

	response.Ok(c, calenderResponse, "")
}
