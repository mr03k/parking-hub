package repository

import (
	"context"
	"errors"
	"farin/domain/entity"
	gormdb "farin/infrastructure/gorm"
	"fmt"
	"gorm.io/gorm"
)

var ErrVehicleNotFound = errors.New("vehicle not found")

type VehicleRepository struct {
	DB *gormdb.GORMDB
}

func NewVehicleRepository(db *gormdb.GORMDB) *VehicleRepository {
	return &VehicleRepository{DB: db}
}

func (r *VehicleRepository) Create(ctx context.Context, vehicle *entity.Vehicle) (*entity.Vehicle, error) {
	if err := r.DB.DB.WithContext(ctx).Create(vehicle).Error; err != nil {
		return nil, err
	}
	return vehicle, nil
}

func (r *VehicleRepository) List(ctx context.Context, filters map[string]interface{}, sortField, sortOrder string, page, pageSize int) ([]entity.Vehicle, int64, error) {
	var vehicles []entity.Vehicle
	var total int64

	query := r.DB.DB.WithContext(ctx).Model(&entity.Vehicle{})

	for key, value := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	if sortField != "" && sortOrder != "" {
		query = query.Order(fmt.Sprintf("%s %s", sortField, sortOrder))
	}

	if page > 0 && pageSize > 0 {
		query = query.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Find(&vehicles).Error; err != nil {
		return nil, 0, err
	}

	return vehicles, total, nil
}

func (r *VehicleRepository) Update(ctx context.Context, vehicle *entity.Vehicle) (*entity.Vehicle, error) {
	if err := r.DB.DB.WithContext(ctx).Save(vehicle).Error; err != nil {
		return nil, err
	}
	return vehicle, nil
}

func (r *VehicleRepository) Delete(ctx context.Context, id string) error {
	if err := r.DB.DB.WithContext(ctx).Delete(&entity.Vehicle{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *VehicleRepository) GetByField(ctx context.Context, field, value string) (*entity.Vehicle, error) {
	var vehicle entity.Vehicle
	if err := r.DB.DB.WithContext(ctx).Where(fmt.Sprintf("%s = ?", field), value).First(&vehicle).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVehicleNotFound // Use custom error for vehicle not found
		}
		return nil, err
	}
	return &vehicle, nil
}
