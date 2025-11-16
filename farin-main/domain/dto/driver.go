package dto

import (
	"farin/domain/entity"
)

type DriverRequest struct {
	ContractorID             string `json:"contractorId" binding:"required,uuid,fkGorm=contractors"`
	UserID                   string `json:"userId" binding:"required,uuid,fkGorm=users"`
	DriverType               string `json:"driverType" binding:"required,oneof=original reserve"`
	ShiftType                string `json:"shiftType" binding:"required,oneof=morning afternoon both"`
	EmploymentStatus         string `json:"employmentStatus" binding:"required"`
	EmploymentStartDate      int64  `json:"employmentStartDate" binding:"required"`
	EmploymentEndDate        *int64 `json:"employmentEndDate,omitempty"`
	Description              string `json:"description,omitempty"`
	DriverPhoto              []byte `binding:"fileData=image/jpeg&image/png;4096000"`
	IDCardImage              []byte `binding:"fileData=image/jpeg&image/png;4096000"`
	BirthCertificateImage    []byte `binding:"fileData=image/jpeg&image/png;4096000"`
	MilitaryServiceCardImage []byte `binding:"fileData=image/jpeg&image/png;4096000"`
	HealthCertificateImage   []byte `binding:"fileData=image/jpeg&image/png;4096000"`
	CriminalRecordImage      []byte `binding:"fileData=image/jpeg&image/png;4096000"`
}

func (req *DriverRequest) ToEntity() *entity.Driver {
	return &entity.Driver{
		ContractorID:        &req.ContractorID,
		UserID:              req.UserID,
		DriverType:          entity.DriverType(req.DriverType),
		ShiftType:           req.ShiftType,
		EmploymentStatus:    req.EmploymentStatus,
		EmploymentStartDate: req.EmploymentStartDate,
		EmploymentEndDate:   req.EmploymentEndDate,
		Description:         req.Description,
	}
}

type DriverResponse struct {
	ID                  string  `json:"id"`
	ContractorID        *string `json:"contractorId,omitempty"`
	UserID              string  `json:"userId"`
	DriverType          string  `json:"driverType"`
	ShiftType           string  `json:"shiftType"`
	EmploymentStatus    string  `json:"employmentStatus"`
	EmploymentStartDate int64   `json:"employmentStartDate"`
	EmploymentEndDate   *int64  `json:"employmentEndDate,omitempty"`
	Description         string  `json:"description,omitempty"`
	CreatedAt           int64   `json:"createdAt"`
	UpdatedAt           int64   `json:"updatedAt"`
}

func (resp *DriverResponse) FromEntity(driver *entity.Driver) {
	resp.ID = driver.ID
	resp.ContractorID = driver.ContractorID
	resp.UserID = driver.UserID
	resp.DriverType = string(driver.DriverType)
	resp.ShiftType = driver.ShiftType
	resp.EmploymentStatus = driver.EmploymentStatus
	resp.EmploymentStartDate = driver.EmploymentStartDate
	resp.EmploymentEndDate = driver.EmploymentEndDate
	resp.Description = driver.Description
	resp.CreatedAt = driver.CreatedAt
	resp.UpdatedAt = driver.UpdatedAt
}

type DriverListResponse struct {
	Drivers []DriverResponse `json:"drivers"`
	Total   int64            `json:"total"`
}
