package repository

import (
	"context"
	"errors"
	"farin/domain/entity"
	gormdb "farin/infrastructure/gorm"
	"fmt"
	"gorm.io/gorm"
)

var ErrDriverNotFound = errors.New("driver not found")

type DriverRepository struct {
	DB *gormdb.GORMDB
}

func NewDriverRepository(db *gormdb.GORMDB) *DriverRepository {
	return &DriverRepository{DB: db}
}

func (r *DriverRepository) Create(ctx context.Context, driver *entity.Driver) (*entity.Driver, error) {
	if err := r.DB.DB.WithContext(ctx).Create(driver).Error; err != nil {
		return nil, err
	}
	return driver, nil
}

func (r *DriverRepository) List(ctx context.Context, filters map[string]interface{}, sortField, sortOrder string, page, pageSize int) ([]entity.Driver, int64, error) {
	var drivers []entity.Driver
	var total int64

	query := r.DB.DB.WithContext(ctx).Model(&entity.Driver{})

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

	if err := query.Find(&drivers).Error; err != nil {
		return nil, 0, err
	}

	return drivers, total, nil
}

func (r *DriverRepository) Update(ctx context.Context, driver *entity.Driver) (*entity.Driver, error) {
	if err := r.DB.DB.WithContext(ctx).Save(driver).Error; err != nil {
		return nil, err
	}
	return driver, nil
}

func (r *DriverRepository) Delete(ctx context.Context, id string) error {
	if err := r.DB.DB.WithContext(ctx).Delete(&entity.Driver{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *DriverRepository) GetByField(ctx context.Context, field, value string) (*entity.Driver, error) {
	var driver entity.Driver
	if err := r.DB.DB.WithContext(ctx).Where(fmt.Sprintf("%s = ?", field), value).First(&driver).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDriverNotFound // Use custom error for driver not found
		}
		return nil, err
	}
	return &driver, nil
}
