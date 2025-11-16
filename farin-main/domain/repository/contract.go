package repository

import (
	"context"
	"errors"
	"farin/domain/entity"
	gormdb "farin/infrastructure/gorm"
	"fmt"
	"gorm.io/gorm"
)

var ErrContractNotFound = errors.New("contract not found")

type ContractRepository struct {
	DB *gormdb.GORMDB
}

func NewContractRepository(db *gormdb.GORMDB) *ContractRepository {
	return &ContractRepository{DB: db}
}

func (r *ContractRepository) Create(ctx context.Context, contract *entity.Contract) (*entity.Contract, error) {
	if err := r.DB.DB.WithContext(ctx).Create(contract).Error; err != nil {
		return nil, err
	}
	return contract, nil
}

func (r *ContractRepository) List(ctx context.Context, filters map[string]interface{}, sortField, sortOrder string, page, pageSize int) ([]entity.Contract, int64, error) {
	var contracts []entity.Contract
	var total int64

	query := r.DB.DB.WithContext(ctx).Model(&entity.Contract{})

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

	if err := query.Find(&contracts).Error; err != nil {
		return nil, 0, err
	}

	return contracts, total, nil
}

func (r *ContractRepository) Update(ctx context.Context, contract *entity.Contract) (*entity.Contract, error) {
	if err := r.DB.DB.WithContext(ctx).Save(contract).Error; err != nil {
		return nil, err
	}
	return contract, nil
}

func (r *ContractRepository) Delete(ctx context.Context, id string) error {
	if err := r.DB.DB.WithContext(ctx).Delete(&entity.Contract{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *ContractRepository) GetByField(ctx context.Context, field, value string) (*entity.Contract, error) {
	var contract entity.Contract
	if err := r.DB.DB.WithContext(ctx).Where(fmt.Sprintf("%s = ?", field), value).First(&contract).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrContractNotFound // Use custom error for contract not found
		}
		return nil, err
	}
	return &contract, nil
}
