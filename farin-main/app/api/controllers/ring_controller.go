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

type RingController struct {
	ringService *service.RingService
	logger      *slog.Logger
	env         *godotenv.Env
}

func NewRingController(logger *slog.Logger, us *service.RingService, env *godotenv.Env) *RingController {
	return &RingController{
		ringService: us,
		logger:      logger.With("layer", "RingController"),
		env:         env,
	}
}

// CreateRing godoc
// @Summary      Create a new ring
// @Description  Create a new ring by providing ring details
// @Tags         rings
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        ring  body      dto.RingRequest  true  "Ring Data"
// @Success      201   {object}  response.Response[dto.RingResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /rings [post]
func (uc *RingController) CreateRing(c *gin.Context) {
	lg := uc.logger.With("method", "CreateRing")
	var ringRequest dto.RingRequest

	if err := c.ShouldBindJSON(&ringRequest); err != nil {
		lg.Error("failed to bind ring data", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	createdRing, err := uc.ringService.CreateRing(context.Background(), ringRequest.ToEntity())
	if err != nil {
		lg.Error("failed to create ring", "error", err)
		response.InternalError(c)
		return
	}
	ringResponse := dto.RingResponse{}
	ringResponse.FromEntity(createdRing)

	response.Created(c, ringResponse)
}

// ListRings godoc
// @Summary      List rings
// @Description  Retrieve a list of rings with optional filters and pagination
// @Tags         rings
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        sortField  query     string  false  "Field to sort by"  default(created_at)
// @Param        sortOrder  query     string  false  "Sort order (asc/desc)"  default(asc)
// @Param        page       query     int     false  "Page number"  default(1)
// @Param        pageSize   query     int     false  "Page size"  default(10)
// @Success      200        {object}  response.Response[dto.RingListResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /rings [get]
func (uc *RingController) ListRings(c *gin.Context) {
	lg := uc.logger.With("method", "ListRings")
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

	rings, total, err := uc.ringService.ListRings(context.Background(), filters, sortField, sortOrder, page, pageSize)
	if err != nil {
		lg.Error("failed to list rings", "error", err)
		response.InternalError(c)
		return
	}
	ringResponses := make([]dto.RingResponse, len(rings))

	for i, ring := range rings {
		ringResponse := dto.RingResponse{}
		ringResponse.FromEntity(&ring)
		ringResponses[i] = ringResponse
	}

	response.Ok(c, dto.RingListResponse{
		Rings: ringResponses,
		Total: total,
	}, "")
}

// UpdateRing godoc
// @Summary      Update ring details
// @Description  Update an existing ring's details
// @Tags         rings
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        ring  body      dto.RingRequest  true  "Updated Ring Data"
// @Param        id  path      string  true  "id"
// @Success      200   {object}  response.Response[dto.RingResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /rings/{id} [put]
func (uc *RingController) UpdateRing(c *gin.Context) {
	lg := uc.logger.With("method", "UpdateRing")
	var ringRequest dto.RingRequest

	if err := c.ShouldBindJSON(&ringRequest); err != nil {
		lg.Error("failed to bind ring data", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	ring := ringRequest.ToEntity()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		lg.Error("failed to convert id to int", "error", err)
	}
	ring.ID = int64(id)
	updatedRing, err := uc.ringService.UpdateRing(context.Background(), ring)
	if err != nil {
		lg.Error("failed to update ring", "error", err)
		response.InternalError(c)
		return
	}

	ringResponse := dto.RingResponse{}
	ringResponse.FromEntity(updatedRing)

	response.Ok(c, ringResponse, "")
}

// DeleteRing godoc
// @Summary      Delete a ring
// @Description  Delete a ring by ID
// @Tags         rings
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      string  true  "Ring ID"
// @Success      200   {object}  response.Response[swagger.EmptyObject]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Failure      404   {object}  response.Response[swagger.EmptyObject]
// @Router       /rings/{id} [delete]
func (uc *RingController) DeleteRing(c *gin.Context) {
	lg := uc.logger.With("method", "DeleteRing")
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		lg.Error("invalid ring ID", "ringID", id)
		response.BadRequest(c, "invalid request body")
		return
	}

	err := uc.ringService.DeleteRing(context.Background(), id)
	if err != nil {
		lg.Error("failed to delete ring", "error", err)
		response.InternalError(c)
		return
	}

	response.Ok(c, nil, "ring deleted successfully")
}

// GetRingDetail godoc
// @Summary      Get ring details
// @Description  Retrieve ring details by ID
// @Tags         rings
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      string  true  "Ring ID"
// @Success      200   {object}  response.Response[dto.RingResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      404   {object}  response.Response[swagger.EmptyObject]
// @Router       /rings/{id} [get]
func (uc *RingController) GetRingDetail(c *gin.Context) {
	lg := uc.logger.With("method", "GetRingDetail")
	id := c.Param("id")

	ring, err := uc.ringService.Detail(context.Background(), "id", id)
	if err != nil {
		lg.Error("failed to get ring details", "error", err)
		response.NotFound(c)
		return
	}
	var ringResponse dto.RingResponse
	ringResponse.FromEntity(ring)

	response.Ok(c, ringResponse, "")
}
