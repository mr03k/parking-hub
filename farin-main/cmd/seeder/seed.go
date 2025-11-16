package seeder

import (
	"context"
	"errors"
	"farin/domain/entity"
	"farin/domain/repository"
	"farin/infrastructure/godotenv"
	gormdb "farin/infrastructure/gorm"
	"farin/util/encrypt"
	"log"
	"log/slog"
	"os"
)

func Seed(fakeData bool) error {
	env := godotenv.NewEnv()
	env.Load()
	ctx := context.Background()
	gorm := gormdb.NewGORMDB(env, slog.New(slog.NewTextHandler(os.Stdout, nil)))
	if err := gorm.Setup(ctx); err != nil {
		log.Fatalf("failed to setup gorm:%s", err)
	}
	userRepository := repository.NewUserRepository(gorm)
	contractorRepository := repository.NewContractorRepository(gorm)
	contractRepository := repository.NewContractRepository(gorm)
	calenderRepository := repository.NewCalenderRepository(gorm)
	vehicleRepository := repository.NewVehicleRepository(gorm)
	ringRepository := repository.NewRingRepository(gorm)
	deviceRepository := repository.NewDeviceRepository(gorm)
	driverRepository := repository.NewDriverRepository(gorm)
	role := repository.NewRoleRepository(gorm)
	driverAssignmentRepository := repository.NewDriverAssignmentRepository(gorm)

	_, err := role.GetByField(context.Background(), "title", "Driver")
	if err != nil {
		if !errors.Is(err, repository.ErrRoleNotFound) {
			_, err = role.Create(context.Background(), &entity.Role{
				Title: "Driver",
			})
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}

	adminRole, err := role.GetByField(context.Background(), "title", "Admin")
	if err != nil {
		if !errors.Is(err, repository.ErrRoleNotFound) {
			_, err = role.Create(context.Background(), &entity.Role{
				Title: "Admin",
			})
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}

	_, err = userRepository.Create(ctx, &entity.User{
		Email:       "admin@example.com",
		PhoneNumber: "+989120123456",
		Password:    encrypt.HashSHA256("password123"),
		RoleID:      adminRole.ID,
	})
	if fakeData {
		users, err := SeedUsers(userRepository)
		if err != nil {
			return err
		}
		// Seed contractors
		contractors, err := SeedContractors(contractorRepository)
		if err != nil {
			return err
		}

		// Seed contracts
		contracts, err := SeedContracts(contractors, contractRepository)
		if err != nil {
			return err
		}

		// Seed calenders
		calenders, err := SeedCalenders(contracts, calenderRepository)
		if err != nil {
			return err
		}

		// Seed vehicles
		vehicles, err := SeedVehicles(contractors, vehicleRepository)
		if err != nil {
			return err
		}

		rings, err := SeedRings(ringRepository)
		if err != nil {
			return err
		}

		_, err = SeedDevices(contractors, vehicles, deviceRepository)
		if err != nil {
			return err
		}

		drivers, err := SeedDrivers(contractors, users, driverRepository)
		if err != nil {
			return err
		}

		_, err = SeedDriverAssignments(drivers, vehicles, rings, calenders, driverAssignmentRepository)
		if err != nil {
			return err
		}
	}

	return err
}
