package device

import (
	"context"
	"errors"
	"log/slog"

	device2 "application/internal/entity/device"
)

// VehicleService provides methods for vehicle-related business logic
type VehicleService struct {
	logger      *slog.Logger
	vehicleRepo VehicleRepositoryInterface
}

// NewVehicleService creates a new instance of VehicleService
func NewVehicleService(logger *slog.Logger, vehicleRepo VehicleRepositoryInterface) VehicleServiceInterface {
	return &VehicleService{
		logger:      logger.With("layer", "VehicleService"),
		vehicleRepo: vehicleRepo,
	}
}

// CreateVehicle creates a new vehicle and saves it to the repository
func (s *VehicleService) CreateVehicle(ctx context.Context, vehicle *device2.Vehicle) (*device2.Vehicle, error) {
	logger := s.logger.With("method", "CreateVehicle")
	logger.Debug("service CreateVehicle")

	if vehicle == nil {
		return nil, errors.New("vehicle entity cannot be nil")
	}

	// Here you could add additional business logic, e.g., validation
	createdVehicle, err := s.vehicleRepo.CreateVehicle(ctx, vehicle)
	if err != nil {
		logger.Error("error in CreateVehicle", "error", err)
		return nil, err
	}

	return createdVehicle, nil
}

// GetVehicle retrieves a vehicle by its UserID from the repository
func (s *VehicleService) GetVehicle(ctx context.Context, vehicleID string) (*device2.Vehicle, error) {
	logger := s.logger.With("method", "GetVehicle")
	logger.Debug("service GetVehicle")

	if vehicleID == "" {
		return nil, errors.New("vehicleID cannot be empty")
	}

	vehicle, err := s.vehicleRepo.GetVehicle(ctx, vehicleID)
	if err != nil {
		if errors.Is(err, ErrVehicleNotFound) {
			return nil, err
		}
		logger.Error("error in GetVehicle", "error", err)
		return nil, err
	}

	return vehicle, nil
}

// ListVehicles retrieves all vehicles from the repository
func (s *VehicleService) ListVehicles(ctx context.Context) ([]device2.Vehicle, error) {
	logger := s.logger.With("method", "ListVehicles")
	logger.Debug("service ListVehicles")

	vehicles, err := s.vehicleRepo.ListVehicles(ctx)
	if err != nil {
		logger.Error("error in ListVehicles", "error", err)
		return nil, err
	}

	return vehicles, nil
}

// DeleteVehicle deletes a vehicle by its UserID from the repository
func (s *VehicleService) DeleteVehicle(ctx context.Context, vehicleID string) error {
	logger := s.logger.With("method", "DeleteVehicle")
	logger.Debug("service DeleteVehicle")

	if vehicleID == "" {
		return errors.New("vehicleID cannot be empty")
	}

	err := s.vehicleRepo.DeleteVehicle(ctx, vehicleID)
	if err != nil {
		if errors.Is(err, ErrVehicleNotFound) {
			return err
		}
		logger.Error("error in DeleteVehicle", "error", err)
		return err
	}

	return nil
}
