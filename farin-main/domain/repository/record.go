package repository

import (
	"context"
	"errors"
	"farin/domain/entity"
	gormdb "farin/infrastructure/gorm"
	"fmt"
	vent "git.abanppc.com/farin-project/vehicle-records/domain/entity"
	"gorm.io/gorm"
)

var ErrVehicleRecordNotFound = errors.New("vehicle record not found")

type VehicleRecordRepository struct {
	DB *gormdb.GORMDB
}

func NewVehicleRecordRepository(db *gormdb.GORMDB) *VehicleRecordRepository {
	return &VehicleRecordRepository{DB: db}
}

// Create a new vehicle record
func (r *VehicleRecordRepository) Create(ctx context.Context, record *vent.VehicleRecord) (*vent.VehicleRecord, error) {
	if err := r.DB.DB.WithContext(ctx).Create(record).Error; err != nil {
		return nil, fmt.Errorf("failed to create vehicle record: %w", err)
	}
	return record, nil
}

// List vehicle records with filtering, sorting, and pagination
func (r *VehicleRecordRepository) List(
	ctx context.Context,
	filters map[string]interface{},
	sortField,
	sortOrder string,
	page,
	pageSize int,
) ([]vent.VehicleRecord, int64, error) {
	var records []vent.VehicleRecord
	var total int64

	query := r.DB.DB.WithContext(ctx).Model(&vent.VehicleRecord{})

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
func (r *VehicleRecordRepository) Update(ctx context.Context, record *vent.VehicleRecord) (*vent.VehicleRecord, error) {
	if err := r.DB.DB.WithContext(ctx).Save(record).Error; err != nil {
		return nil, fmt.Errorf("failed to update vehicle record: %w", err)
	}
	return record, nil
}

// Delete a vehicle record by ID
func (r *VehicleRecordRepository) Delete(ctx context.Context, id string) error {
	result := r.DB.DB.WithContext(ctx).Delete(&vent.VehicleRecord{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete vehicle record: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrVehicleRecordNotFound
	}

	return nil
}

// Get a vehicle record by a specific field
func (r *VehicleRecordRepository) GetByField(ctx context.Context, field string, value interface{}) (*vent.VehicleRecord, error) {
	var record vent.VehicleRecord

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
func (r *VehicleRecordRepository) GetByID(ctx context.Context, id string) (*vent.VehicleRecord, error) {
	return r.GetByField(ctx, "record_id", id)
}

// FindByPlateNumber finds vehicle records by plate number
func (r *VehicleRecordRepository) FindByPlateNumber(ctx context.Context, plateNumber string) ([]vent.VehicleRecord, error) {
	var records []vent.VehicleRecord

	if err := r.DB.DB.WithContext(ctx).
		Where("citizen_plate_number = ?", plateNumber).
		Preload("VehiclePhotos").
		Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to find vehicle records by plate number: %w", err)
	}

	return records, nil
}

// Bulk create vehicle records
func (r *VehicleRecordRepository) BulkCreate(ctx context.Context, records []*vent.VehicleRecord) error {
	if err := r.DB.DB.WithContext(ctx).Create(records).Error; err != nil {
		return fmt.Errorf("failed to bulk create vehicle records: %w", err)
	}
	return nil
}

func (r *VehicleRecordRepository) FindRing(ctx context.Context, record *vent.VehicleRecord) ([]entity.Ring, error) {
	var rings []entity.Ring
	err := r.DB.DB.WithContext(ctx).Raw(
		`SELECT * FROM public.rings pgone WHERE ST_intersects(
            ST_Buffer(pgone.geom, 5),
            ST_Buffer(ST_Transform(ST_SetSRID(ST_MakePoint(?, ?), 4326), 32639), ?)
        )`,
		record.LPRVehicleGPSLongitude, record.LPRVehicleGPSLatitude, record.LPRVehicleGPSError).Scan(&rings).Error
	if err != nil {
		return nil, err
	}
	return rings, nil
}

func (r *VehicleRecordRepository) FindSegment(ctx context.Context, record *vent.VehicleRecord) ([]entity.Segment, error) {
	var segments []entity.Segment
	err := r.DB.DB.WithContext(ctx).Raw(
		`SELECT * FROM public.segments pgone WHERE ST_intersects(
            ST_Buffer(pgone.geom, 10),
            ST_Transform(ST_SetSRID(ST_MakePoint(?, ?), 4326), 32639)
        )`,
		record.LPRVehicleGPSLongitude, record.LPRVehicleGPSLatitude).Scan(&segments).Error
	if err != nil {
		return nil, err
	}
	return segments, nil
}

func (r *VehicleRecordRepository) FindRoad(ctx context.Context, record *vent.VehicleRecord) ([]entity.Road, error) {
	var roads []entity.Road
	err := r.DB.DB.WithContext(ctx).Raw(
		`SELECT * FROM public.roads pgone WHERE ST_intersects(
            ST_Buffer(pgone.geom, 5),
            ST_Buffer(ST_Transform(ST_SetSRID(ST_MakePoint(?, ?), 4326), 32639), ?)
        )`,
		record.LPRVehicleGPSLongitude, record.LPRVehicleGPSLatitude, record.LPRVehicleGPSError).Scan(&roads).Error
	if err != nil {
		return nil, err
	}
	return roads, nil
}

func (r *VehicleRecordRepository) FindParking(ctx context.Context, record *vent.VehicleRecord) ([]entity.Road, error) {
	var roads []entity.Road
	err := r.DB.DB.WithContext(ctx).Raw(
		`SELECT * FROM public.parkings pgone WHERE ST_intersects(
            ST_Buffer(pgone.geom, 5),
            ST_Buffer(ST_Transform(ST_SetSRID(ST_MakePoint(?, ?), 4326), 32639), ?)
        )`,
		record.LPRVehicleGPSLongitude, record.LPRVehicleGPSLatitude, record.LPRVehicleGPSError).Scan(&roads).Error
	if err != nil {
		return nil, err
	}
	return roads, nil
}
