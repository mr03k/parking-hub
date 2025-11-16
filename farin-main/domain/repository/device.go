package repository

import (
	"context"
	"errors"
	"farin/domain/entity"
	gormdb "farin/infrastructure/gorm"
	"fmt"
	"gorm.io/gorm"
)

var ErrDeviceNotFound = errors.New("device not found")

type DeviceRepository struct {
	DB *gormdb.GORMDB
}

func NewDeviceRepository(db *gormdb.GORMDB) *DeviceRepository {
	return &DeviceRepository{DB: db}
}

func (r *DeviceRepository) Create(ctx context.Context, device *entity.Device) (*entity.Device, error) {
	if err := r.DB.DB.WithContext(ctx).Create(device).Error; err != nil {
		return nil, err
	}
	return device, nil
}

func (r *DeviceRepository) List(ctx context.Context, filters map[string]interface{}, sortField, sortOrder string, page, pageSize int) ([]entity.Device, int64, error) {
	var devices []entity.Device
	var total int64

	query := r.DB.DB.WithContext(ctx).Model(&entity.Device{})

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

	if err := query.Find(&devices).Error; err != nil {
		return nil, 0, err
	}

	return devices, total, nil
}

func (r *DeviceRepository) Update(ctx context.Context, device *entity.Device) (*entity.Device, error) {
	if err := r.DB.DB.WithContext(ctx).Save(device).Error; err != nil {
		return nil, err
	}
	return device, nil
}

func (r *DeviceRepository) Delete(ctx context.Context, id string) error {
	if err := r.DB.DB.WithContext(ctx).Delete(&entity.Device{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *DeviceRepository) GetByField(ctx context.Context, field, value string) (*entity.Device, error) {
	var device entity.Device
	if err := r.DB.DB.WithContext(ctx).Where(fmt.Sprintf("%s = ?", field), value).First(&device).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDeviceNotFound // Use custom error for device not found
		}
		return nil, err
	}
	return &device, nil
}

func (r *DeviceRepository) GetMultipleByField(ctx context.Context, field, value string) ([]*entity.Device, error) {
	var devices []*entity.Device
	if err := r.DB.DB.WithContext(ctx).Where(fmt.Sprintf("%s = ?", field), value).Find(&devices).Error; err != nil {
		if len(devices) < 1 {
			return nil, ErrDeviceNotFound // Use custom error for device not found
		}
		return nil, err
	}
	return devices, nil
}
