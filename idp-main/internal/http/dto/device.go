package dto

import "time"

// DeviceResponse represents a summary of device information.
type DeviceResponse struct {
	ID                  string    `json:"id"`
	CodeDevice          string    `json:"codeDevice"`
	NumberSerial        string    `json:"numberSerial"`
	Model               string    `json:"model"`
	DateInstallation    time.Time `json:"dateInstallation"`
	DateExpiryWarranty  time.Time `json:"dateExpiryWarranty"`
	DateExpiryInsurance time.Time `json:"dateExpiryInsurance"`
	ClassDevice         string    `json:"classDevice"`
	Status              string    `json:"status"` // Example status field
}

// DeviceDetailResponse represents detailed device information.
type DeviceDetailResponse struct {
	ID                  string    `json:"id"`
	CodeDevice          string    `json:"codeDevice"`
	NumberSerial        string    `json:"numberSerial"`
	Model               string    `json:"model"`
	DateInstallation    time.Time `json:"dateInstallation"`
	DateExpiryWarranty  time.Time `json:"dateExpiryWarranty"`
	DateExpiryInsurance time.Time `json:"dateExpiryInsurance"`
	ClassDevice         string    `json:"classDevice"`
	ImageContractURL    string    `json:"imageContractUrl,omitempty"`
	ImageInsuranceURL   string    `json:"imageInsuranceUrl,omitempty"`
	Description         string    `json:"description,omitempty"`
}

// DeviceListResponse represents the response for a list of devices.
type DeviceListResponse struct {
	Devices []DeviceResponse `json:"devices"`
}
