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

type ContractController struct {
	contractService *service.ContractService
	logger          *slog.Logger
	env             *godotenv.Env
}

func NewContractController(logger *slog.Logger, us *service.ContractService, env *godotenv.Env) *ContractController {
	return &ContractController{
		contractService: us,
		logger:          logger.With("layer", "ContractController"),
		env:             env,
	}
}

// CreateContract godoc
// @Summary      Create a new contract
// @Description  Create a new contract by providing contract details
// @Tags         contracts
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        contract  body      dto.ContractRequest  true  "Contract Data"
// @Success      201   {object}  response.Response[dto.ContractResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /contracts [post]
func (uc *ContractController) CreateContract(c *gin.Context) {
	lg := uc.logger.With("method", "CreateContract")
	var contractRequest dto.ContractRequest

	if err := c.ShouldBindJSON(&contractRequest); err != nil {
		lg.Error("failed to bind contract data", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	createdContract, err := uc.contractService.CreateContract(context.Background(), contractRequest.ToEntity())
	if err != nil {
		lg.Error("failed to create contract", "error", err)
		response.InternalError(c)
		return
	}
	contractResponse := dto.ContractResponse{}
	contractResponse.FromEntity(createdContract)

	response.Created(c, contractResponse)
}

// ListContracts godoc
// @Summary      List contracts
// @Description  Retrieve a list of contracts with optional filters and pagination
// @Tags         contracts
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        sortField  query     string  false  "Field to sort by"  default(created_at)
// @Param        sortOrder  query     string  false  "Sort order (asc/desc)"  default(asc)
// @Param        page       query     int     false  "Page number"  default(1)
// @Param        pageSize   query     int     false  "Page size"  default(10)
// @Success      200        {object}  response.Response[dto.ContractListResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /contracts [get]
func (uc *ContractController) ListContracts(c *gin.Context) {
	lg := uc.logger.With("method", "ListContracts")
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

	contracts, total, err := uc.contractService.ListContracts(context.Background(), filters, sortField, sortOrder, page, pageSize)
	if err != nil {
		lg.Error("failed to list contracts", "error", err)
		response.InternalError(c)
		return
	}
	contractResponses := make([]dto.ContractResponse, len(contracts))

	for i, contract := range contracts {
		contractResponse := dto.ContractResponse{}
		contractResponse.FromEntity(&contract)
		contractResponses[i] = contractResponse
	}

	response.Ok(c, dto.ContractListResponse{
		Contracts: contractResponses,
		Total:     total,
	}, "")
}

// UpdateContract godoc
// @Summary      Update contract details
// @Description  Update an existing contract's details
// @Tags         contracts
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        contract  body      dto.ContractRequest  true  "Updated Contract Data"
// @Param        id  path      string  true  "id"
// @Success      200   {object}  response.Response[dto.ContractResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /contracts/{id} [put]
func (uc *ContractController) UpdateContract(c *gin.Context) {
	lg := uc.logger.With("method", "UpdateContract")
	var contractRequest dto.ContractRequest

	if err := c.ShouldBindJSON(&contractRequest); err != nil {
		lg.Error("failed to bind contract data", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	contract := contractRequest.ToEntity()
	contract.ID = c.Param("id")
	updatedContract, err := uc.contractService.UpdateContract(context.Background(), contract)
	if err != nil {
		lg.Error("failed to update contract", "error", err)
		response.InternalError(c)
		return
	}

	contractResponse := dto.ContractResponse{}
	contractResponse.FromEntity(updatedContract)

	response.Ok(c, contractResponse, "")
}

// DeleteContract godoc
// @Summary      Delete a contract
// @Description  Delete a contract by ID
// @Tags         contracts
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      string  true  "Contract ID"
// @Success      200   {object}  response.Response[swagger.EmptyObject]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Failure      404   {object}  response.Response[swagger.EmptyObject]
// @Router       /contracts/{id} [delete]
func (uc *ContractController) DeleteContract(c *gin.Context) {
	lg := uc.logger.With("method", "DeleteContract")
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		lg.Error("invalid contract ID", "contractID", id)
		response.BadRequest(c, "invalid request body")
		return
	}

	err := uc.contractService.DeleteContract(context.Background(), id)
	if err != nil {
		lg.Error("failed to delete contract", "error", err)
		response.InternalError(c)
		return
	}

	response.Ok(c, nil, "contract deleted successfully")
}

// GetContractDetail godoc
// @Summary      Get contract details
// @Description  Retrieve contract details by ID
// @Tags         contracts
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      string  true  "Contract ID"
// @Success      200   {object}  response.Response[dto.ContractResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      404   {object}  response.Response[swagger.EmptyObject]
// @Router       /contracts/{id} [get]
func (uc *ContractController) GetContractDetail(c *gin.Context) {
	lg := uc.logger.With("method", "GetContractDetail")
	id := c.Param("id")

	contract, err := uc.contractService.Detail(context.Background(), "id", id)
	if err != nil {
		lg.Error("failed to get contract details", "error", err)
		response.NotFound(c)
		return
	}
	var contractResponse dto.ContractResponse
	contractResponse.FromEntity(contract)

	response.Ok(c, contractResponse, "")
}
