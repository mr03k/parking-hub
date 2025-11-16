package dto

import (
	"farin/domain/entity"
	"time"
)

type ContractRequest struct {
	ContractNumber  string `json:"contractNumber" binding:"required,min=1,max=50"`
	ContractDate    string `json:"contractDate" binding:"required"`
	StartDate       string `json:"startDate" binding:"required"`
	EndDate         string `json:"endDate" binding:"required"`
	ContractAmount  int64  `json:"contractAmount" binding:"required,min=1"`
	ContractType    string `json:"contractType" binding:"required,min=1,max=50"`
	ContractorID    string `json:"contractorId" binding:"required,fkGorm=contractors"`
	OperationPeriod int    `json:"operationPeriod" binding:"required,min=1"`
	EquipmentPeriod int    `json:"equipmentPeriod" binding:"required,min=1"`
	Description     string `json:"description,omitempty"`
}

func (req *ContractRequest) ToEntity() *entity.Contract {
	contractDate, _ := time.Parse("2006-01-02", req.ContractDate)
	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	return &entity.Contract{
		ContractNumber:  req.ContractNumber,
		ContractDate:    contractDate,
		StartDate:       startDate,
		EndDate:         endDate,
		ContractAmount:  req.ContractAmount,
		ContractType:    req.ContractType,
		ContractorID:    req.ContractorID,
		OperationPeriod: req.OperationPeriod,
		EquipmentPeriod: req.EquipmentPeriod,
		Description:     req.Description,
	}
}

type ContractResponse struct {
	ID              string `json:"id"`
	ContractNumber  string `json:"contractNumber"`
	ContractDate    string `json:"contractDate"`
	StartDate       string `json:"startDate"`
	EndDate         string `json:"endDate"`
	ContractAmount  int64  `json:"contractAmount"`
	ContractType    string `json:"contractType"`
	ContractorID    string `json:"contractorId"`
	OperationPeriod int    `json:"operationPeriod"`
	EquipmentPeriod int    `json:"equipmentPeriod"`
	Description     string `json:"description,omitempty"`
	CreatedAt       int64  `json:"createdAt"`
	UpdatedAt       int64  `json:"updatedAt"`
}

func (resp *ContractResponse) FromEntity(contract *entity.Contract) {
	resp.ID = contract.ID
	resp.ContractNumber = contract.ContractNumber
	resp.ContractDate = contract.ContractDate.Format("2006-01-02")
	resp.StartDate = contract.StartDate.Format("2006-01-02")
	resp.EndDate = contract.EndDate.Format("2006-01-02")
	resp.ContractAmount = contract.ContractAmount
	resp.ContractType = contract.ContractType
	resp.ContractorID = contract.ContractorID
	resp.OperationPeriod = contract.OperationPeriod
	resp.EquipmentPeriod = contract.EquipmentPeriod
	resp.Description = contract.Description
	resp.CreatedAt = contract.CreatedAt
	resp.UpdatedAt = contract.UpdatedAt
}

type ContractListResponse struct {
	Contracts []ContractResponse `json:"contracts"`
	Total     int64              `json:"total"`
}
