package dto

import (
	"farin/domain/entity"
)

type DriverAssignmentRequest struct {
	DriverID    string `json:"driverId" binding:"required,uuid,fkGorm=drivers"`
	CodeVehicle string `json:"codeVehicle" binding:"required"`
	RingID      int64  `json:"ringId" binding:"required,fkGorm=rings"`
	CalenderID  string `json:"calenderID" binding:"required,uuid,fkGorm=calenders"`
}

func (req *DriverAssignmentRequest) ToEntity() *entity.DriverAssignment {
	return &entity.DriverAssignment{
		DriverID:    req.DriverID,
		CodeVehicle: req.CodeVehicle,
		RingID:      req.RingID,
		CalenderID:  req.CalenderID,
	}
}

type DriverAssignmentResponse struct {
	ID          string `json:"id"`
	DriverID    string `json:"driverId"`
	CodeVehicle string `json:"codeVehicle"`
	RingID      int64  `json:"ringId"`
	CalenderID  string `json:"calendarId"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}

func (resp *DriverAssignmentResponse) FromEntity(driverAssignment *entity.DriverAssignment) {
	resp.ID = driverAssignment.ID
	resp.DriverID = driverAssignment.DriverID
	resp.CodeVehicle = driverAssignment.CodeVehicle
	resp.RingID = driverAssignment.RingID
	resp.CalenderID = driverAssignment.CalenderID
	resp.CreatedAt = driverAssignment.CreatedAt
	resp.UpdatedAt = driverAssignment.UpdatedAt
}

type DriverAssignmentListResponse struct {
	DriverAssignments []DriverAssignmentResponse `json:"driverAssignments"`
	Total             int64                      `json:"total"`
}
