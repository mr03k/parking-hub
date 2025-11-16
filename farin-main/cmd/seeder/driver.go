package seeder

import (
	"context"
	"farin/domain/entity"
	"farin/domain/repository"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func SeedDrivers(
	contractors []*entity.Contractor,
	users []*entity.User,
	driverRepository *repository.DriverRepository,
) ([]*entity.Driver, error) {
	// Ensure we have enough users and contractors
	if len(users) < 2 || len(contractors) < 2 {
		return nil, fmt.Errorf("not enough users or contractors to seed drivers")
	}

	drivers := []*entity.Driver{
		{
			Base: entity.Base{
				ID: uuid.New().String(),
			},
			ContractorID:             &contractors[0].ID,
			UserID:                   users[0].ID,
			DriverType:               entity.DriverTypeOriginal,
			ShiftType:                string(entity.ShiftTypeMorning),
			EmploymentStatus:         "Active",
			EmploymentStartDate:      time.Now().AddDate(-1, 0, 0).Unix(), // 1 year ago
			EmploymentEndDate:        nil,                                 // Still employed
			DriverPhoto:              "/path/to/driver/photo1.jpg",
			IDCardImage:              "/path/to/id/card1.jpg",
			BirthCertificateImage:    "/path/to/birth/certificate1.jpg",
			MilitaryServiceCardImage: "/path/to/military/service1.jpg",
			HealthCertificateImage:   "/path/to/health/certificate1.jpg",
			CriminalRecordImage:      "/path/to/criminal/record1.jpg",
			Description:              "Senior driver with excellent driving record",
		},
		{
			Base: entity.Base{
				ID: uuid.New().String(),
			},
			ContractorID:             &contractors[1].ID,
			UserID:                   users[1].ID,
			DriverType:               entity.DriverTypeReserved,
			ShiftType:                string(entity.ShiftTypeAfterNoon),
			EmploymentStatus:         "Probation",
			EmploymentStartDate:      time.Now().AddDate(0, -6, 0).Unix(),     // 6 months ago
			EmploymentEndDate:        ptr(time.Now().AddDate(0, 6, 0).Unix()), // 6 months from now
			DriverPhoto:              "/path/to/driver/photo2.jpg",
			IDCardImage:              "/path/to/id/card2.jpg",
			BirthCertificateImage:    "/path/to/birth/certificate2.jpg",
			MilitaryServiceCardImage: "/path/to/military/service2.jpg",
			HealthCertificateImage:   "/path/to/health/certificate2.jpg",
			CriminalRecordImage:      "/path/to/criminal/record2.jpg",
			Description:              "Reserve driver in training",
		},
	}

	// Seed drivers
	for _, driver := range drivers {
		_, err := driverRepository.Create(context.Background(), driver)
		if err != nil {
			return nil, err
		}
	}

	return drivers, nil
}

// Helper function to create a pointer to a value
func ptr[T any](v T) *T {
	return &v
}
