package repository

import (
	"context"
	"errors"
	"farin/domain/entity"
	gormdb "farin/infrastructure/gorm"
	"fmt"
	"gorm.io/gorm"
)

var ErrContractorNotFound = errors.New("contractor not found")

type ContractorRepository struct {
	DB *gormdb.GORMDB
}

func NewContractorRepository(db *gormdb.GORMDB) *ContractorRepository {
	return &ContractorRepository{DB: db}
}

func (r *ContractorRepository) Create(ctx context.Context, contractor *entity.Contractor) (*entity.Contractor, error) {
	if err := r.DB.DB.WithContext(ctx).Create(contractor).Error; err != nil {
		return nil, err
	}
	return contractor, nil
}

func (r *ContractorRepository) List(ctx context.Context, filters map[string]interface{}, sortField, sortOrder string, page, pageSize int) ([]entity.Contractor, int64, error) {
	var contractors []entity.Contractor
	var total int64

	query := r.DB.DB.WithContext(ctx).Model(&entity.Contractor{})

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

	if err := query.Find(&contractors).Error; err != nil {
		return nil, 0, err
	}

	return contractors, total, nil
}

func (r *ContractorRepository) Update(ctx context.Context, contractor *entity.Contractor) (*entity.Contractor, error) {
	if err := r.DB.DB.WithContext(ctx).Save(contractor).Error; err != nil {
		return nil, err
	}
	return contractor, nil
}

func (r *ContractorRepository) Delete(ctx context.Context, id string) error {
	if err := r.DB.DB.WithContext(ctx).Delete(&entity.Contractor{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *ContractorRepository) GetByField(ctx context.Context, field, value string) (*entity.Contractor, error) {
	var contractor entity.Contractor
	if err := r.DB.DB.WithContext(ctx).Where(fmt.Sprintf("%s = ?", field), value).First(&contractor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrContractorNotFound // Use custom error for contractor not found
		}
		return nil, err
	}
	return &contractor, nil
}
