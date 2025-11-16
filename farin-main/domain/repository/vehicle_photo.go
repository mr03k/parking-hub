package repository

import (
	"context"
	"errors"
	gormdb "farin/infrastructure/gorm"
	"fmt"
	"git.abanppc.com/farin-project/vehicle-records/domain/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrCitizenVehiclePhotoNotFound = errors.New("citizen vehicle photo not found")

type CitizenVehiclePhotoRepository struct {
	DB *gormdb.GORMDB
}

func NewCitizenVehiclePhotoRepository(db *gormdb.GORMDB) *CitizenVehiclePhotoRepository {
	return &CitizenVehiclePhotoRepository{DB: db}
}

// Create a new citizen vehicle photo
func (r *CitizenVehiclePhotoRepository) Create(ctx context.Context, photo *entity.CitizenVehiclePhoto) (*entity.CitizenVehiclePhoto, error) {
	if err := r.DB.DB.WithContext(ctx).Create(photo).Error; err != nil {
		return nil, fmt.Errorf("failed to create citizen vehicle photo: %w", err)
	}
	return photo, nil
}

// List citizen vehicle photos with filtering, sorting, and pagination
func (r *CitizenVehiclePhotoRepository) List(
	ctx context.Context,
	filters map[string]interface{},
	sortField,
	sortOrder string,
	page,
	pageSize int,
) ([]entity.CitizenVehiclePhoto, int64, error) {
	var photos []entity.CitizenVehiclePhoto
	var total int64

	query := r.DB.DB.WithContext(ctx).Model(&entity.CitizenVehiclePhoto{})

	// Apply filters
	for key, value := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// Apply sorting
	if sortField != "" && sortOrder != "" {
		query = query.Order(fmt.Sprintf("%s %s", sortField, sortOrder))
	}

	// Apply pagination
	if page > 0 && pageSize > 0 {
		query = query.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count citizen vehicle photos: %w", err)
	}

	// Fetch records with related vehicle record
	if err := query.Preload("VehicleRecord").Find(&photos).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch citizen vehicle photos: %w", err)
	}

	return photos, total, nil
}

// Update an existing citizen vehicle photo
func (r *CitizenVehiclePhotoRepository) Update(ctx context.Context, photo *entity.CitizenVehiclePhoto) (*entity.CitizenVehiclePhoto, error) {
	if err := r.DB.DB.WithContext(ctx).Save(photo).Error; err != nil {
		return nil, fmt.Errorf("failed to update citizen vehicle photo: %w", err)
	}
	return photo, nil
}

// Delete a citizen vehicle photo by ID
func (r *CitizenVehiclePhotoRepository) Delete(ctx context.Context, id uint) error {
	result := r.DB.DB.WithContext(ctx).Delete(&entity.CitizenVehiclePhoto{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete citizen vehicle photo: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrCitizenVehiclePhotoNotFound
	}

	return nil
}

// Get a citizen vehicle photo by a specific field
func (r *CitizenVehiclePhotoRepository) GetByField(ctx context.Context, field string, value interface{}) (*entity.CitizenVehiclePhoto, error) {
	var photo entity.CitizenVehiclePhoto

	query := r.DB.DB.WithContext(ctx).
		Where(fmt.Sprintf("%s = ?", field), value).
		Preload("VehicleRecord")

	if err := query.First(&photo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCitizenVehiclePhotoNotFound
		}
		return nil, fmt.Errorf("failed to retrieve citizen vehicle photo: %w", err)
	}

	return &photo, nil
}

// Get a citizen vehicle photo by ID
func (r *CitizenVehiclePhotoRepository) GetByID(ctx context.Context, id uint) (*entity.CitizenVehiclePhoto, error) {
	return r.GetByField(ctx, "id", id)
}

// Find photos by record ID
func (r *CitizenVehiclePhotoRepository) FindByRecordID(ctx context.Context, recordID uuid.UUID) ([]entity.CitizenVehiclePhoto, error) {
	var photos []entity.CitizenVehiclePhoto

	if err := r.DB.DB.WithContext(ctx).
		Where("record_id = ?", recordID).
		Preload("VehicleRecord").
		Find(&photos).Error; err != nil {
		return nil, fmt.Errorf("failed to find citizen vehicle photos by record ID: %w", err)
	}

	return photos, nil
}

// Bulk create citizen vehicle photos
func (r *CitizenVehiclePhotoRepository) BulkCreate(ctx context.Context, photos []*entity.CitizenVehiclePhoto) error {
	if err := r.DB.DB.WithContext(ctx).Create(photos).Error; err != nil {
		return fmt.Errorf("failed to bulk create citizen vehicle photos: %w", err)
	}
	return nil
}

// Find photos by camera ID
func (r *CitizenVehiclePhotoRepository) FindByCameraID(ctx context.Context, cameraID int) ([]entity.CitizenVehiclePhoto, error) {
	var photos []entity.CitizenVehiclePhoto

	if err := r.DB.DB.WithContext(ctx).
		Where("lpr_vehicle_camera_id = ?", cameraID).
		Preload("VehicleRecord").
		Find(&photos).Error; err != nil {
		return nil, fmt.Errorf("failed to find citizen vehicle photos by camera ID: %w", err)
	}

	return photos, nil
}
