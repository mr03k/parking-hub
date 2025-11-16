package seeder

import (
	"context"
	"farin/domain/entity"
	"farin/domain/repository"
	"fmt"

	"github.com/google/uuid"
)

func SeedDriverAssignments(
	drivers []*entity.Driver,
	vehicles []*entity.Vehicle,
	rings []*entity.Ring,
	calenders []*entity.Calender,
	driverAssignmentRepository *repository.DriverAssignmentRepository,
) ([]entity.DriverAssignment, error) {
	// Validate input data
	if len(drivers) == 0 || len(vehicles) == 0 || len(rings) == 0 || len(calenders) == 0 {
		return nil, fmt.Errorf("insufficient data to seed driver assignments")
	}

	// Ensure we have at least 2 of each entity type
	maxIndex := min(
		len(drivers)-1,
		len(vehicles)-1,
		len(rings)-1,
		len(calenders)-1,
	)

	driverAssignments := []entity.DriverAssignment{
		{
			Base: entity.Base{
				ID: uuid.New().String(),
			},
			DriverID:    drivers[0].ID,
			CodeVehicle: vehicles[0].ID,
			RingID:      rings[0].ID,
			CalenderID:  calenders[0].ID,
		},
		{
			Base: entity.Base{
				ID: uuid.New().String(),
			},
			DriverID:    drivers[min(1, maxIndex)].ID,
			CodeVehicle: vehicles[min(1, maxIndex)].ID,
			RingID:      rings[min(1, maxIndex)].ID,
			CalenderID:  calenders[min(1, maxIndex)].ID,
		},
	}

	// Seed driver assignments
	for _, assignment := range driverAssignments {
		_, err := driverAssignmentRepository.Create(context.Background(), &assignment)
		if err != nil {
			return nil, err
		}
	}

	return driverAssignments, nil
}

// Helper function to find the minimum of multiple integers
func min(vals ...int) int {
	if len(vals) == 0 {
		panic("min requires at least one argument")
	}
	minVal := vals[0]
	for _, val := range vals[1:] {
		if val < minVal {
			minVal = val
		}
	}
	return minVal
}
