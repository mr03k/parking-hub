package repository

import (
	"context"
	"errors"
	"farin/domain/entity"
	gormdb "farin/infrastructure/gorm"
	"fmt"
	"gorm.io/gorm"
)

var ErrCalenderNotFound = errors.New("calender not found")

type CalenderRepository struct {
	DB *gormdb.GORMDB
}

func NewCalenderRepository(db *gormdb.GORMDB) *CalenderRepository {
	return &CalenderRepository{DB: db}
}

func (r *CalenderRepository) Create(ctx context.Context, calender *entity.Calender) (*entity.Calender, error) {
	if err := r.DB.DB.WithContext(ctx).Create(calender).Error; err != nil {
		return nil, err
	}
	return calender, nil
}

func (r *CalenderRepository) List(ctx context.Context, filters map[string]interface{}, sortField, sortOrder string, page, pageSize int) ([]entity.Calender, int64, error) {
	var calenders []entity.Calender
	var total int64

	query := r.DB.DB.WithContext(ctx).Model(&entity.Calender{})

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

	if err := query.Find(&calenders).Error; err != nil {
		return nil, 0, err
	}

	return calenders, total, nil
}

func (r *CalenderRepository) Update(ctx context.Context, calender *entity.Calender) (*entity.Calender, error) {
	if err := r.DB.DB.WithContext(ctx).Save(calender).Error; err != nil {
		return nil, err
	}
	return calender, nil
}

func (r *CalenderRepository) Delete(ctx context.Context, id string) error {
	if err := r.DB.DB.WithContext(ctx).Delete(&entity.Calender{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *CalenderRepository) GetByField(ctx context.Context, field, value string) (*entity.Calender, error) {
	var calender entity.Calender
	if err := r.DB.DB.WithContext(ctx).Where(fmt.Sprintf("%s = ?", field), value).First(&calender).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCalenderNotFound // Use custom error for calender not found
		}
		return nil, err
	}
	return &calender, nil
}
