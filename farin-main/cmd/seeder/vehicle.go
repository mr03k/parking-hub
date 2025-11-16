package seeder

import (
	"context"
	"farin/domain/entity"
	"farin/domain/repository"
	"fmt"
	"github.com/bxcodec/faker/v4"
	"math/rand"
	"time"
)

func SeedVehicles(
	contractors []*entity.Contractor,
	vr *repository.VehicleRepository,
) ([]*entity.Vehicle, error) {
	// Define vehicle types
	vehicleTypes := []string{
		"Sedan", "SUV", "Pickup", "Van", "Truck", "Minivan",
		"Hatchback", "Coupe", "Convertible",
	}

	// Define vehicle brands
	vehicleBrands := []string{
		"Toyota", "Honda", "Ford", "Chevrolet", "Mercedes", "BMW",
		"Volkswagen", "Nissan", "Hyundai", "Kia",
	}

	// Define vehicle colors
	vehicleColors := []string{
		"White", "Black", "Silver", "Gray", "Red", "Blue",
		"Green", "Yellow", "Bronze", "Navy",
	}

	// Define vehicle statuses
	vehicleStatuses := []string{
		"Active", "Inactive", "Maintenance", "Reserved", "In Use",
	}

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	var vehicles []*entity.Vehicle

	// Generate 50 sample vehicles
	for i := 0; i < 50; i++ {
		// Select a random contractor
		contractor := contractors[r.Intn(len(contractors))]

		// Generate a unique vehicle code

		// Generate a mock VIN (Vehicle Identification Number)
		vin := fmt.Sprintf("%s%d%06d", faker.UUIDHyphenated()[:8], r.Intn(10), r.Intn(1000000))
		codeVehicle := fmt.Sprintf("%s%d%06d", faker.UUIDHyphenated()[:8], r.Intn(10), r.Intn(1000000))

		// Generate a mock plate license
		plateLicense := fmt.Sprintf("%c%c %d %c%c%c",
			rune(r.Intn(26)+'A'),
			rune(r.Intn(26)+'A'),
			r.Intn(999),
			rune(r.Intn(26)+'A'),
			rune(r.Intn(26)+'A'),
			rune(r.Intn(26)+'A'),
		)
		// Current timestamp for insurance expiry
		now := time.Now()

		vehicle := &entity.Vehicle{
			CodeVehicle:               codeVehicle,
			VIN:                       vin,
			PlateLicense:              plateLicense,
			TypeVehicle:               vehicleTypes[r.Intn(len(vehicleTypes))],
			Brand:                     vehicleBrands[r.Intn(len(vehicleBrands))],
			Model:                     faker.Word(),
			Color:                     vehicleColors[r.Intn(len(vehicleColors))],
			ManufactureOfYear:         r.Intn(15) + 2010,           // Vehicles from 2010 to 2024
			KilometersInitial:         int64(r.Intn(200000)),       // Up to 200,000 km
			ExpiryInsurancePartyThird: now.AddDate(1, 0, 0).Unix(), // 1 year from now
			ExpiryInsuranceBody:       now.AddDate(1, 6, 0).Unix(), // 1.5 years from now
			ImageDocumentVehicle:      faker.URL(),
			ImageCardVehicle:          faker.URL(),
			ThirdPartyInsuranceImage:  faker.URL(),
			BodyInsuranceImage:        faker.URL(),
			ContractorID:              contractor.ID,
			Status:                    vehicleStatuses[r.Intn(len(vehicleStatuses))],
			Description:               faker.Paragraph(),
		}

		// Create the vehicle
		createdVehicle, err := vr.Create(context.Background(), vehicle)
		if err != nil {
			return nil, err
		}
		vehicles = append(vehicles, createdVehicle)
	}

	return vehicles, nil
}
