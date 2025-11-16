package repository

import (
	"context"
	"errors"
	"farin/domain/entity"
	gormdb "farin/infrastructure/gorm"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

var ErrRoleNotFound = errors.New("Role not found")
var ErrRoleAlreadyExist = errors.New("Role already exists")

type RoleRepository struct {
	DB *gormdb.GORMDB
}

func NewRoleRepository(db *gormdb.GORMDB) *RoleRepository {
	return &RoleRepository{DB: db}
}

func (r *RoleRepository) Create(ctx context.Context, Role *entity.Role) (*entity.Role, error) {
	if err := r.DB.DB.WithContext(ctx).Create(Role).Error; err != nil {
		// Check if it's a duplicate key error
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique constraint") {
			return nil, ErrRoleAlreadyExist
		}
		return nil, err
	}
	return Role, nil
}

func (r *RoleRepository) List(ctx context.Context, filters map[string]interface{}, sortField, sortOrder string, page, pageSize int) ([]entity.Role, int64, error) {
	var Roles []entity.Role
	var total int64

	query := r.DB.DB.WithContext(ctx).Model(&entity.Role{})

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

	if err := query.Find(&Roles).Error; err != nil {
		return nil, 0, err
	}

	return Roles, total, nil
}

func (r *RoleRepository) Update(ctx context.Context, Role *entity.Role) (*entity.Role, error) {
	if err := r.DB.DB.WithContext(ctx).Save(Role).Error; err != nil {
		return nil, err
	}
	return Role, nil
}

func (r *RoleRepository) Delete(ctx context.Context, id string) error {
	return r.DB.DB.WithContext(ctx).Where("id = ?", id).Delete(&entity.Role{}).Error
}

func (r *RoleRepository) GetByField(ctx context.Context, field, value string) (*entity.Role, error) {
	var Role entity.Role
	tx := r.DB.DB.WithContext(ctx).Where(fmt.Sprintf("%s = ?", field), value)

	if err := tx.First(&Role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound // Use custom error for Role not found
		}
		return nil, err
	}
	return &Role, nil
}

func (r *RoleRepository) GetMultipleByField(ctx context.Context, field, value string) ([]*entity.Role, error) {
	var Roles []*entity.Role
	if err := r.DB.DB.WithContext(ctx).Where(fmt.Sprintf("%s = ?", field), value).Find(&Roles).Error; err != nil {
		if len(Roles) < 1 {
			return nil, ErrRoleNotFound // Use custom error for Role not found
		}
		return nil, err
	}
	return Roles, nil
}
