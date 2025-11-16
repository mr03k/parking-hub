package dto

import (
	"farin/domain/entity"
)

type CalenderRequest struct {
	ContractID     string `json:"contractId" binding:"required,fkGorm=contracts"`
	ShamsiDate     string `json:"shamsiDate" binding:"required,min=1,max=10"`
	WorkDate       int64  `json:"workDate" binding:"required"`
	Weekday        string `json:"weekday" binding:"required,oneof=Saturday Sunday Monday Tuesday Wednesday Thursday Friday"`
	Year           int    `json:"year" binding:"required"`
	IsHoliday      bool   `json:"isHoliday" binding:"required"`
	WorkShift      string `json:"workShift" binding:"required,oneof=Morning Afternoon Both"`
	Description    string `json:"description,omitempty"`
	WorkShiftStart int64  `json:"workShiftStart,omitempty" binding:"required"`
	WorkShiftEnd   int64  `json:"workShiftEnd,omitempty" binding:"required"`
}

func (req *CalenderRequest) ToEntity() *entity.Calender {
	return &entity.Calender{
		ContractID:     req.ContractID,
		ShamsiDate:     req.ShamsiDate,
		WorkDate:       req.WorkDate,
		Weekday:        entity.Weekday(req.Weekday),
		Year:           req.Year,
		IsHoliday:      req.IsHoliday,
		WorkShift:      entity.WorkShift(req.WorkShift),
		Description:    req.Description,
		WorkShiftStart: req.WorkShiftStart,
		WorkShiftEnd:   req.WorkShiftEnd,
	}
}

type CalenderResponse struct {
	ID             string `json:"id"`
	ContractID     string `json:"contractId"`
	ShamsiDate     string `json:"shamsiDate"`
	WorkDate       int64  `json:"workDate"`
	Weekday        string `json:"weekday"`
	Year           int    `json:"year"`
	IsHoliday      bool   `json:"isHoliday"`
	WorkShift      string `json:"workShift"`
	Description    string `json:"description,omitempty"`
	WorkShiftStart int64  `json:"workShiftStart,omitempty"`
	WorkShiftEnd   int64  `json:"workShiftEnd,omitempty"`
	CreatedAt      int64  `json:"createdAt"`
	UpdatedAt      int64  `json:"updatedAt"`
}

func (resp *CalenderResponse) FromEntity(calender *entity.Calender) {
	resp.ID = calender.ID
	resp.ContractID = calender.ContractID
	resp.ShamsiDate = calender.ShamsiDate
	resp.WorkDate = calender.WorkDate
	resp.Weekday = string(calender.Weekday)
	resp.Year = calender.Year
	resp.IsHoliday = calender.IsHoliday
	resp.WorkShift = string(calender.WorkShift)
	resp.Description = calender.Description
	resp.WorkShiftStart = calender.WorkShiftStart
	resp.WorkShiftEnd = calender.WorkShiftEnd
	resp.CreatedAt = calender.CreatedAt
	resp.UpdatedAt = calender.UpdatedAt
}

type CalenderListResponse struct {
	Calenders []CalenderResponse `json:"calenders"`
	Total     int64              `json:"total"`
}
