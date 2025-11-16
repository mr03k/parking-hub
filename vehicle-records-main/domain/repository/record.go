package repository

import (
	"context"
	"errors"
	"fmt"
	"git.abanppc.com/farin-project/vehicle-records/domain/entity"
	gormdb "git.abanppc.com/farin-project/vehicle-records/infrastructure/gorm"
	"gorm.io/gorm"
)

var ErrVehicleRecordNotFound = errors.New("vehicle record not found")
var ErrDuplicateRecord = errors.New("duplicate vehicle record")

type VehicleRecordRepository struct {
	DB *gormdb.GORMDB
}

func NewVehicleRecordRepository(db *gormdb.GORMDB) *VehicleRecordRepository {
	return &VehicleRecordRepository{DB: db}
}

// Create a new vehicle record
func (r *VehicleRecordRepository) Create(ctx context.Context, record *entity.VehicleRecord) (*entity.VehicleRecord, error) {
	var c int64
	err := r.DB.DB.WithContext(ctx).Model(record).Where("record_id = ?", record.RecordID).Count(&c).Error
	if c > 0 {
		return nil, ErrDuplicateRecord
	}
	if err != nil {
		return nil, fmt.Errorf("failed to check existing record: %w", err)
	}
	if err := r.DB.DB.WithContext(ctx).Create(record).Error; err != nil {
		return nil, fmt.Errorf("failed to create vehicle record: %w", err)
	}
	return record, nil
}
func (r *VehicleRecordRepository) GetNotSent(ctx context.Context, createAt int64, limit int) (
	[]entity.VehicleRecord, error) {
	var records []entity.VehicleRecord

	err := r.DB.DB.WithContext(ctx).Model(&entity.VehicleRecord{}).Preload("VehiclePhotos").
		Where("created_at > ?", createAt).Where("plate_detection_id=?", 0).Limit(limit).
		Find(&records).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get records: %w", err)
	}

	return records, nil
}

// List vehicle records with filtering, sorting, and pagination
func (r *VehicleRecordRepository) List(
	ctx context.Context,
	filters map[string]interface{},
	sortField,
	sortOrder string,
	page,
	pageSize int,
) ([]entity.VehicleRecord, int64, error) {
	var records []entity.VehicleRecord
	var total int64

	query := r.DB.DB.WithContext(ctx).Model(&entity.VehicleRecord{})

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
		return nil, 0, fmt.Errorf("failed to count vehicle records: %w", err)
	}

	// Fetch records
	if err := query.Preload("VehiclePhotos").Find(&records).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch vehicle records: %w", err)
	}

	return records, total, nil
}

// Update an existing vehicle record
func (r *VehicleRecordRepository) Update(ctx context.Context, record *entity.VehicleRecord) (*entity.VehicleRecord, error) {
	if err := r.DB.DB.WithContext(ctx).Omit("VehiclePhotos", "Retries").Updates(record).Error; err != nil {
		return nil, fmt.Errorf("failed to update vehicle record: %w", err)
	}
	return record, nil
}

func (r *VehicleRecordRepository) IncreaseRetry(ctx context.Context, record *entity.VehicleRecord) (*entity.VehicleRecord, error) {
	if err := r.DB.DB.WithContext(ctx).Exec("UPDATE vehicle_records SET retries=retries+1 WHERE record_id=?",
		record.RecordID).Error; err != nil {
		return nil, fmt.Errorf("failed to update vehicle record: %w", err)
	}
	return record, nil
}

// Delete a vehicle record by ID
func (r *VehicleRecordRepository) Delete(ctx context.Context, id string) error {
	result := r.DB.DB.WithContext(ctx).Delete(&entity.VehicleRecord{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete vehicle record: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrVehicleRecordNotFound
	}

	return nil
}

// Get a vehicle record by a specific field
func (r *VehicleRecordRepository) GetByField(ctx context.Context, field string, value interface{}) (*entity.VehicleRecord, error) {
	var record entity.VehicleRecord

	query := r.DB.DB.WithContext(ctx).
		Where(fmt.Sprintf("%s = ?", field), value).
		Preload("VehiclePhotos")

	if err := query.First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVehicleRecordNotFound
		}
		return nil, fmt.Errorf("failed to retrieve vehicle record: %w", err)
	}

	return &record, nil
}

// Get a vehicle record by ID
func (r *VehicleRecordRepository) GetByID(ctx context.Context, id string) (*entity.VehicleRecord, error) {
	return r.GetByField(ctx, "record_id", id)
}

// FindByPlateNumber finds vehicle records by plate number
func (r *VehicleRecordRepository) FindByPlateNumber(ctx context.Context, plateNumber string) ([]entity.VehicleRecord, error) {
	var records []entity.VehicleRecord

	if err := r.DB.DB.WithContext(ctx).
		Where("citizen_plate_number = ?", plateNumber).
		Preload("VehiclePhotos").
		Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to find vehicle records by plate number: %w", err)
	}

	return records, nil
}

// Bulk create vehicle records
func (r *VehicleRecordRepository) BulkCreate(ctx context.Context, records []*entity.VehicleRecord) error {
	if err := r.DB.DB.WithContext(ctx).Create(records).Error; err != nil {
		return fmt.Errorf("failed to bulk create vehicle records: %w", err)
	}
	return nil
}
