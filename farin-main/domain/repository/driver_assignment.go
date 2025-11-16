package repository

import (
	"context"
	"errors"
	"farin/domain/entity"
	gormdb "farin/infrastructure/gorm"
	"fmt"
	"gorm.io/gorm"
)

var ErrDriverAssignmentNotFound = errors.New("driverAssignment not found")

type DriverAssignmentRepository struct {
	DB *gormdb.GORMDB
}

func NewDriverAssignmentRepository(db *gormdb.GORMDB) *DriverAssignmentRepository {
	return &DriverAssignmentRepository{DB: db}
}

func (r *DriverAssignmentRepository) Create(ctx context.Context, driverAssignment *entity.DriverAssignment) (*entity.DriverAssignment, error) {
	if err := r.DB.DB.WithContext(ctx).Create(driverAssignment).Error; err != nil {
		return nil, err
	}
	return driverAssignment, nil
}

func (r *DriverAssignmentRepository) List(ctx context.Context, filters map[string]interface{}, sortField, sortOrder string, page, pageSize int) ([]entity.DriverAssignment, int64, error) {
	var driverAssignments []entity.DriverAssignment
	var total int64

	query := r.DB.DB.WithContext(ctx).Model(&entity.DriverAssignment{})

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

	if err := query.Find(&driverAssignments).Error; err != nil {
		return nil, 0, err
	}

	return driverAssignments, total, nil
}

func (r *DriverAssignmentRepository) Update(ctx context.Context, driverAssignment *entity.DriverAssignment) (*entity.DriverAssignment, error) {
	if err := r.DB.DB.WithContext(ctx).Save(driverAssignment).Error; err != nil {
		return nil, err
	}
	return driverAssignment, nil
}

func (r *DriverAssignmentRepository) Delete(ctx context.Context, id string) error {
	if err := r.DB.DB.WithContext(ctx).Delete(&entity.DriverAssignment{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *DriverAssignmentRepository) GetByField(ctx context.Context, field, value string) (*entity.DriverAssignment, error) {
	var driverAssignment entity.DriverAssignment
	if err := r.DB.DB.WithContext(ctx).Where(fmt.Sprintf("%s = ?", field), value).First(&driverAssignment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDriverAssignmentNotFound // Use custom error for driverAssignment not found
		}
		return nil, err
	}
	return &driverAssignment, nil
}

func (r *DriverAssignmentRepository) GetByFields(ctx context.Context, driverID, vehicleCode string) (*entity.DriverAssignment, error) {
	var driverAssignment entity.DriverAssignment
	if err := r.DB.DB.WithContext(ctx).Preload("Driver").Preload("Vehicle").Preload("Ring").
		Preload("Calender").Where("driver_id=?", driverID).Where("code_vehicle=?", vehicleCode).
		First(&driverAssignment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDriverAssignmentNotFound // Use custom error for driverAssignment not found
		}
		return nil, err
	}
	return &driverAssignment, nil
}
