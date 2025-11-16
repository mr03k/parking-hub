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

type ContractorController struct {
	contractorService *service.ContractorService
	logger            *slog.Logger
	env               *godotenv.Env
}

func NewContractorController(logger *slog.Logger, us *service.ContractorService, env *godotenv.Env) *ContractorController {
	return &ContractorController{
		contractorService: us,
		logger:            logger.With("layer", "ContractorController"),
		env:               env,
	}
}

// CreateContractor godoc
// @Summary      Create a new contractor
// @Description  Create a new contractor by providing contractor details
// @Tags         contractors
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        contractor  body      dto.ContractorRequest  true  "Contractor Data"
// @Success      201   {object}  response.Response[dto.ContractorResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /contractors [post]
func (uc *ContractorController) CreateContractor(c *gin.Context) {
	lg := uc.logger.With("method", "CreateContractor")
	var contractorRequest dto.ContractorRequest

	if err := c.ShouldBindJSON(&contractorRequest); err != nil {
		lg.Error("failed to bind contractor data", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	createdContractor, err := uc.contractorService.CreateContractor(context.Background(), contractorRequest.ToEntity())
	if err != nil {
		lg.Error("failed to create contractor", "error", err)
		response.InternalError(c)
		return
	}
	contractorResponse := dto.ContractorResponse{}
	contractorResponse.FromEntity(createdContractor, uc.env)

	response.Created(c, contractorResponse)
}

// ListContractors godoc
// @Summary      List contractors
// @Description  Retrieve a list of contractors with optional filters and pagination
// @Tags         contractors
// @Accept       json
// @Produce      json
// @Param        sortField  query     string  false  "Field to sort by"  default(created_at)
// @Param        sortOrder  query     string  false  "Sort order (asc/desc)"  default(asc)
// @Param        page       query     int     false  "Page number"  default(1)
// @Param        pageSize   query     int     false  "Page size"  default(10)
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200        {object}  response.Response[dto.ContractorListResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /contractors [get]
func (uc *ContractorController) ListContractors(c *gin.Context) {
	lg := uc.logger.With("method", "ListContractors")
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

	contractors, total, err := uc.contractorService.ListContractors(context.Background(), filters, sortField, sortOrder, page, pageSize)
	if err != nil {
		lg.Error("failed to list contractors", "error", err)
		response.InternalError(c)
		return
	}
	contractorResponses := make([]dto.ContractorResponse, len(contractors))

	for i, contractor := range contractors {
		contractorResponse := dto.ContractorResponse{}
		contractorResponse.FromEntity(&contractor, uc.env)
		contractorResponses[i] = contractorResponse
	}

	response.Ok(c, dto.ContractorListResponse{
		Contractors: contractorResponses,
		Total:       total,
	}, "")
}

// UpdateContractor godoc
// @Summary      Update contractor details
// @Description  Update an existing contractor's details
// @Tags         contractors
// @Accept       json
// @Produce      json
// @Param        contractor  body      dto.ContractorRequest  true  "Updated Contractor Data"
// @Param        id  path      string  true  "id"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200   {object}  response.Response[dto.ContractorResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /contractors/{id} [put]
func (uc *ContractorController) UpdateContractor(c *gin.Context) {
	lg := uc.logger.With("method", "UpdateContractor")
	var contractorRequest dto.ContractorRequest

	if err := c.ShouldBindJSON(&contractorRequest); err != nil {
		lg.Error("failed to bind contractor data", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	contractor := contractorRequest.ToEntity()
	contractor.ID = c.Param("id")
	updatedContractor, err := uc.contractorService.UpdateContractor(context.Background(), contractor)
	if err != nil {
		lg.Error("failed to update contractor", "error", err)
		response.InternalError(c)
		return
	}

	contractorResponse := dto.ContractorResponse{}
	contractorResponse.FromEntity(updatedContractor, uc.env)

	response.Ok(c, contractorResponse, "")
}

// DeleteContractor godoc
// @Summary      Delete a contractor
// @Description  Delete a contractor by ID
// @Tags         contractors
// @Param        id   path      string  true  "Contractor ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200   {object}  response.Response[swagger.EmptyObject]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Failure      404   {object}  response.Response[swagger.EmptyObject]
// @Router       /contractors/{id} [delete]
func (uc *ContractorController) DeleteContractor(c *gin.Context) {
	lg := uc.logger.With("method", "DeleteContractor")
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		lg.Error("invalid contractor ID", "contractorID", id)
		response.BadRequest(c, "invalid request body")
		return
	}

	err := uc.contractorService.DeleteContractor(context.Background(), id)
	if err != nil {
		lg.Error("failed to delete contractor", "error", err)
		response.InternalError(c)
		return
	}

	response.Ok(c, nil, "contractor deleted successfully")
}

// GetContractorDetail godoc
// @Summary      Get contractor details
// @Description  Retrieve contractor details by ID
// @Tags         contractors
// @Param        id   path      string  true  "Contractor ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200   {object}  response.Response[dto.ContractorResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      404   {object}  response.Response[swagger.EmptyObject]
// @Router       /contractors/{id} [get]
func (uc *ContractorController) GetContractorDetail(c *gin.Context) {
	lg := uc.logger.With("method", "GetContractorDetail")
	id := c.Param("id")

	contractor, err := uc.contractorService.Detail(context.Background(), "id", id)
	if err != nil {
		lg.Error("failed to get contractor details", "error", err)
		response.NotFound(c)
		return
	}
	var contractorResponse dto.ContractorResponse
	contractorResponse.FromEntity(contractor, uc.env)

	response.Ok(c, contractorResponse, "")
}
