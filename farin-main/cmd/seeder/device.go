package seeder

import (
	"context"
	"farin/domain/entity"
	"farin/domain/repository"
	"github.com/google/uuid"
	"time"
)

func SeedDevices(
	contractors []*entity.Contractor,
	vehicles []*entity.Vehicle,
	deviceRepository *repository.DeviceRepository,
) ([]entity.Device, error) {
	devices := []entity.Device{
		{
			Base: entity.Base{
				ID: uuid.New().String(),
			},
			CodeDevice:          "DEV-001",
			NumberSerial:        "SN-12345-6789",
			Model:               "ModelX-2023",
			DateInstallation:    time.Now().Unix(),
			DateExpiryWarranty:  time.Now().AddDate(1, 0, 0).Unix(), // 1 year from now
			DateExpiryInsurance: time.Now().AddDate(2, 0, 0).Unix(), // 2 years from now
			ClassDevice:         "Tracking Device",
			ImageContract:       "/path/to/contract/image1.jpg",
			ImageInsurance:      "/path/to/insurance/image1.jpg",
			ContractorID:        contractors[0].ID,
			VehicleID:           vehicles[0].ID,
			Description:         "Primary tracking device for vehicle",
		},
		{
			Base: entity.Base{
				ID: uuid.New().String(),
			},
			CodeDevice:          "DEV-002",
			NumberSerial:        "SN-98765-4321",
			Model:               "ModelY-2023",
			DateInstallation:    time.Now().AddDate(0, -6, 0).Unix(), // 6 months ago
			DateExpiryWarranty:  time.Now().AddDate(0, 6, 0).Unix(),  // 6 months from now
			DateExpiryInsurance: time.Now().AddDate(1, 6, 0).Unix(),  // 1.5 years from now
			ClassDevice:         "Backup Tracking Device",
			ImageContract:       "/path/to/contract/image2.jpg",
			ImageInsurance:      "/path/to/insurance/image2.jpg",
			ContractorID:        contractors[1].ID,
			VehicleID:           vehicles[1].ID,
			Description:         "Secondary backup tracking device",
		},
	}

	// Seed devices
	for _, device := range devices {
		_, err := deviceRepository.Create(context.Background(), &device)
		if err != nil {
			return nil, err
		}
	}

	return devices, nil
}
