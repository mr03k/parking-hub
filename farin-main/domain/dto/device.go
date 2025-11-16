package dto

import (
	"farin/domain/entity"
	"farin/infrastructure/godotenv"
	"farin/util"
)

type DeviceRequest struct {
	CodeDevice          string `json:"codeDevice" binding:"required,min=1,max=120"`
	DeviceID            int64  `json:"deviceID" binding:"required,min=1,max=120"`
	NumberSerial        string `json:"numberSerial" binding:"required,min=1,max=50"`
	Model               string `json:"model,omitempty"`
	DateInstallation    int64  `json:"dateInstallation,omitempty"`
	DateExpiryWarranty  int64  `json:"dateExpiryWarranty,omitempty"`
	DateExpiryInsurance int64  `json:"dateExpiryInsurance,omitempty"`
	ClassDevice         string `json:"classDevice,omitempty"`
	ImageContract       []byte `json:"imageContract,omitempty" binding:"fileData=image/jpeg&image/png;4096000"`
	ImageInsurance      []byte `json:"imageInsurance,omitempty" binding:"fileData=image/jpeg&image/png;4096000"`
	ContractorID        string `json:"contractorId" binding:"required,fkGorm=contractors"`
	VehicleID           string `json:"vehicleId" binding:"required,fkGorm=vehicles"`
	Description         string `json:"description,omitempty"`
}

func (req *DeviceRequest) ToEntity() *entity.Device {
	return &entity.Device{
		CodeDevice:          req.CodeDevice,
		NumberSerial:        req.NumberSerial,
		Model:               req.Model,
		DateInstallation:    req.DateInstallation,
		DateExpiryWarranty:  req.DateExpiryWarranty,
		DateExpiryInsurance: req.DateExpiryInsurance,
		ClassDevice:         req.ClassDevice,
		ContractorID:        req.ContractorID,
		VehicleID:           req.VehicleID,
		Description:         req.Description,
		DeviceID:            req.DeviceID,
	}
}

type DeviceResponse struct {
	ID                  string `json:"id"`
	CodeDevice          string `json:"codeDevice"`
	DeviceID            int64  `json:"deviceID"`
	NumberSerial        string `json:"numberSerial"`
	Model               string `json:"model,omitempty"`
	DateInstallation    int64  `json:"dateInstallation,omitempty"`
	DateExpiryWarranty  int64  `json:"dateExpiryWarranty,omitempty"`
	DateExpiryInsurance int64  `json:"dateExpiryInsurance,omitempty"`
	ClassDevice         string `json:"classDevice,omitempty"`
	ImageContract       string `json:"imageContract,omitempty"`
	ImageInsurance      string `json:"imageInsurance,omitempty"`
	ContractorID        string `json:"contractorId"`
	VehicleID           string `json:"vehicleId"`
	Description         string `json:"description,omitempty"`
	CreatedAt           int64  `json:"createdAt"`
	UpdatedAt           int64  `json:"updatedAt"`
}

func (resp *DeviceResponse) FromEntity(env *godotenv.Env, device *entity.Device) {
	resp.ID = device.ID
	resp.DeviceID = device.DeviceID
	resp.CodeDevice = device.CodeDevice
	resp.NumberSerial = device.NumberSerial
	resp.Model = device.Model
	resp.DateInstallation = device.DateInstallation
	resp.DateExpiryWarranty = device.DateExpiryWarranty
	resp.DateExpiryInsurance = device.DateExpiryInsurance
	resp.ClassDevice = device.ClassDevice
	resp.ImageContract = util.GeneratePublicURL(device.ImageContract, env)
	resp.ImageInsurance = util.GeneratePublicURL(device.ImageInsurance, env)
	resp.ContractorID = device.ContractorID
	resp.Description = device.Description
	resp.CreatedAt = device.CreatedAt
	resp.UpdatedAt = device.UpdatedAt
}

type DeviceListResponse struct {
	Devices []DeviceResponse `json:"devices"`
	Total   int64            `json:"total"`
}
