package repository

import (
	"context"
	"errors"
	"fmt"
	"git.abanppc.com/farin-project/state/domain/entity"
	gormdb "git.abanppc.com/farin-project/state/infrastructure/gorm"
	"gorm.io/gorm"
)

var ErrStateNotFound = errors.New("state not found")

type StateRepository struct {
	DB *gormdb.GORMDB
}

func NewStateRepository(db *gormdb.GORMDB) *StateRepository {
	return &StateRepository{DB: db}
}

func (r *StateRepository) Create(ctx context.Context, state *entity.State) (*entity.State, error) {
	if err := r.DB.DB.WithContext(ctx).Create(state).Error; err != nil {
		return nil, err
	}
	return state, nil
}

func (r *StateRepository) List(ctx context.Context, filters map[string]interface{}, sortField, sortOrder string, page, pageSize int) ([]entity.State, int64, error) {
	var states []entity.State
	var total int64

	query := r.DB.DB.WithContext(ctx).Model(&entity.State{})

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

	if err := query.Find(&states).Error; err != nil {
		return nil, 0, err
	}

	return states, total, nil
}

func (r *StateRepository) Update(ctx context.Context, state *entity.State) (*entity.State, error) {
	if err := r.DB.DB.WithContext(ctx).Save(state).Error; err != nil {
		return nil, err
	}
	return state, nil
}

func (r *StateRepository) Delete(ctx context.Context, id string) error {
	if err := r.DB.DB.WithContext(ctx).Delete(&entity.State{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *StateRepository) GetByField(ctx context.Context, field, value string) (*entity.State, error) {
	var state entity.State
	if err := r.DB.DB.WithContext(ctx).Where(fmt.Sprintf("%s = ?", field), value).First(&state).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStateNotFound // Use custom error for state not found
		}
		return nil, err
	}
	return &state, nil
}

// GetByRecordID retrieves a state by its record ID
func (r *StateRepository) GetByRecordID(ctx context.Context, recordID string) (*entity.State, error) {
	return r.GetByField(ctx, "record_id", recordID)
}

// GetByLPRSystemID retrieves states by LPR system ID
func (r *StateRepository) GetByLPRSystemID(ctx context.Context, systemID string) ([]entity.State, error) {
	var states []entity.State
	if err := r.DB.DB.WithContext(ctx).Where("lpr_system_id = ?", systemID).Find(&states).Error; err != nil {
		return nil, err
	}
	return states, nil
}

// GetByLPRVehicleID retrieves states by LPR vehicle ID
func (r *StateRepository) GetByLPRVehicleID(ctx context.Context, vehicleID string) ([]entity.State, error) {
	var states []entity.State
	if err := r.DB.DB.WithContext(ctx).Where("lpr_vehicle_id = ?", vehicleID).Find(&states).Error; err != nil {
		return nil, err
	}
	return states, nil
}

// ListByTimeRange retrieves states within a specific time range
func (r *StateRepository) ListByTimeRange(ctx context.Context, startTime, endTime int64, page, pageSize int) ([]entity.State, int64, error) {
	var states []entity.State
	var total int64

	query := r.DB.DB.WithContext(ctx).Model(&entity.State{}).
		Where("record_store_time >= ? AND record_store_time <= ?", startTime, endTime)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page > 0 && pageSize > 0 {
		query = query.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	if err := query.Order("record_store_time DESC").Find(&states).Error; err != nil {
		return nil, 0, err
	}

	return states, total, nil
}

// ListAvailableServers retrieves states where the server is available
func (r *StateRepository) ListAvailableServers(ctx context.Context) ([]entity.State, error) {
	var states []entity.State
	if err := r.DB.DB.WithContext(ctx).Where("server_availability = ?", true).Find(&states).Error; err != nil {
		return nil, err
	}
	return states, nil
}
