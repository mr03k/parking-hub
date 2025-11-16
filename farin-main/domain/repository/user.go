package repository

import (
	"context"
	"errors"
	"farin/domain/entity"
	gormdb "farin/infrastructure/gorm"
	"fmt"
	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	DB *gormdb.GORMDB
}

func NewUserRepository(db *gormdb.GORMDB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	if err := r.DB.DB.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) List(ctx context.Context, filters map[string]interface{}, sortField, sortOrder string, page, pageSize int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	query := r.DB.DB.WithContext(ctx).Model(&entity.User{})

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

	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	if err := r.DB.DB.WithContext(ctx).Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return r.DB.DB.WithContext(ctx).Where("id = ?", id).Delete(&entity.User{}).Error
}

func (r *UserRepository) GetByField(ctx context.Context, field, value string) (*entity.User, error) {
	var user entity.User
	if err := r.DB.DB.WithContext(ctx).Preload("Role").Where(fmt.Sprintf("%s = ?", field), value).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound // Use custom error for user not found
		}
		return nil, err
	}
	return &user, nil
}
