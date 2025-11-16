package device

import (
	"application/internal/entity/device"
	"context"
	"errors"
)

var (
	ErrVehicleNotFound = errors.New("device not found")
)

type VehicleRepositoryInterface interface {
	// CreateVehicle inserts a new vehicle into the repository
	CreateVehicle(ctx context.Context, vehicle *device.Vehicle) (*device.Vehicle, error)

	// GetVehicle retrieves a vehicle by its UserID from the repository
	GetVehicle(ctx context.Context, vehicleID string) (*device.Vehicle, error)

	// ListVehicles retrieves all vehicles from the repository
	ListVehicles(ctx context.Context) ([]device.Vehicle, error)

	// DeleteVehicle removes a vehicle by its UserID from the repository
	DeleteVehicle(ctx context.Context, vehicleID string) error
}

type VehicleServiceInterface interface {
	// CreateVehicle creates a new vehicle in the repository
	CreateVehicle(ctx context.Context, vehicle *device.Vehicle) (*device.Vehicle, error)

	// GetVehicle retrieves a vehicle by its UserID
	GetVehicle(ctx context.Context, vehicleID string) (*device.Vehicle, error)

	// ListVehicles retrieves all vehicles from the repository
	ListVehicles(ctx context.Context) ([]device.Vehicle, error)

	// DeleteVehicle removes a vehicle by its UserID from the repository
	DeleteVehicle(ctx context.Context, vehicleID string) error
}
