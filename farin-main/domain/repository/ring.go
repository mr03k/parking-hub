package repository

import (
	"context"
	"errors"
	"farin/domain/entity"
	gormdb "farin/infrastructure/gorm"
	"fmt"
	"gorm.io/gorm"
)

var ErrRingNotFound = errors.New("ring not found")

type RingRepository struct {
	DB *gormdb.GORMDB
}

func NewRingRepository(db *gormdb.GORMDB) *RingRepository {
	return &RingRepository{DB: db}
}

func (r *RingRepository) Create(ctx context.Context, ring *entity.Ring) (*entity.Ring, error) {
	if err := r.DB.DB.WithContext(ctx).Create(ring).Error; err != nil {
		return nil, err
	}
	return ring, nil
}

func (r *RingRepository) List(ctx context.Context, filters map[string]interface{}, sortField, sortOrder string, page, pageSize int) ([]entity.Ring, int64, error) {
	var rings []entity.Ring
	var total int64

	query := r.DB.DB.WithContext(ctx).Model(&entity.Ring{})

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

	if err := query.Find(&rings).Error; err != nil {
		return nil, 0, err
	}

	return rings, total, nil
}

func (r *RingRepository) Update(ctx context.Context, ring *entity.Ring) (*entity.Ring, error) {
	if err := r.DB.DB.WithContext(ctx).Save(ring).Error; err != nil {
		return nil, err
	}
	return ring, nil
}

func (r *RingRepository) Delete(ctx context.Context, id string) error {
	if err := r.DB.DB.WithContext(ctx).Delete(&entity.Ring{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *RingRepository) GetByField(ctx context.Context, field, value any) (*entity.Ring, error) {
	var ring entity.Ring
	if err := r.DB.DB.WithContext(ctx).Where(fmt.Sprintf("%s = ?", field), value).First(&ring).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRingNotFound // Use custom error for ring not found
		}
		return nil, err
	}
	return &ring, nil
}
