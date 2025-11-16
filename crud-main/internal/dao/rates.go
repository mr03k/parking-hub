package dao

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	cacheBase "github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/ggorm/query"
	"github.com/zhufuyi/sponge/pkg/utils"

	"git.abanppc.com/farin-project/crud/internal/cache"
	"git.abanppc.com/farin-project/crud/internal/model"
)

var _ RatesDao = (*ratesDao)(nil)

// RatesDao defining the dao interface
type RatesDao interface {
	Create(ctx context.Context, table *model.Rates) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.Rates) error
	GetByID(ctx context.Context, id uint64) (*model.Rates, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.Rates, int64, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Rates) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Rates) error
}

type ratesDao struct {
	db    *gorm.DB
	cache cache.RatesCache    // if nil, the cache is not used.
	sfg   *singleflight.Group // if cache is nil, the sfg is not used.
}

// NewRatesDao creating the dao interface
func NewRatesDao(db *gorm.DB, xCache cache.RatesCache) RatesDao {
	if xCache == nil {
		return &ratesDao{db: db}
	}
	return &ratesDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *ratesDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a record, insert the record and the id value is written back to the table
func (d *ratesDao) Create(ctx context.Context, table *model.Rates) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a record by id
func (d *ratesDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Rates{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a record by id
func (d *ratesDao) UpdateByID(ctx context.Context, table *model.Rates) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *ratesDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.Rates) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.Code != "" {
		update["code"] = table.Code
	}
	if table.RoadCategoryID != "" {
		update["road_category_id"] = table.RoadCategoryID
	}
	if table.TimeCycleMinutes != 0 {
		update["time_cycle_minutes"] = table.TimeCycleMinutes
	}
	if table.RateMultiplier != "" {
		update["rate_multiplier"] = table.RateMultiplier
	}
	if table.PeakHourMultiplier != "" {
		update["peak_hour_multiplier"] = table.PeakHourMultiplier
	}
	if table.GoodPercentage != 0 {
		update["good_percentage"] = table.GoodPercentage
	}
	if table.NormalSettlementPeriod != 0 {
		update["normal_settlement_period"] = table.NormalSettlementPeriod
	}
	if table.LatePenalty != "" {
		update["late_penalty"] = table.LatePenalty
	}
	if table.LatePenaltyMax != "" {
		update["late_penalty_max"] = table.LatePenaltyMax
	}
	if table.ValidFrom.IsZero() == false {
		update["valid_from"] = table.ValidFrom
	}
	if table.ValidTo.IsZero() == false {
		update["valid_to"] = table.ValidTo
	}
	if table.Description != "" {
		update["description"] = table.Description
	}
	if table.StartTime != "" {
		update["start_time"] = table.StartTime
	}
	if table.EndTime != "" {
		update["end_time"] = table.EndTime
	}
	if table.CityID != "" {
		update["city_id"] = table.CityID
	}
	if table.ApprovalNumber != "" {
		update["approval_number"] = table.ApprovalNumber
	}
	if table.ApprovalDate.IsZero() == false {
		update["approval_date"] = table.ApprovalDate
	}
	if table.Year != 0 {
		update["year"] = table.Year
	}
	if table.BaseRateID != "" {
		update["base_rate_id"] = table.BaseRateID
	}
	if table.ExceptionsID != "" {
		update["exceptions_id"] = table.ExceptionsID
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a record by id
func (d *ratesDao) GetByID(ctx context.Context, id uint64) (*model.Rates, error) {
	// no cache
	if d.cache == nil {
		record := &model.Rates{}
		err := d.db.WithContext(ctx).Where("id = ?", id).First(record).Error
		return record, err
	}

	// get from cache or database
	record, err := d.cache.Get(ctx, id)
	if err == nil {
		return record, nil
	}

	if errors.Is(err, model.ErrCacheNotFound) {
		// for the same id, prevent high concurrent simultaneous access to database
		val, err, _ := d.sfg.Do(utils.Uint64ToStr(id), func() (interface{}, error) { //nolint
			table := &model.Rates{}
			err = d.db.WithContext(ctx).Where("id = ?", id).First(table).Error
			if err != nil {
				// if data is empty, set not found cache to prevent cache penetration, default expiration time 10 minutes
				if errors.Is(err, model.ErrRecordNotFound) {
					err = d.cache.SetCacheWithNotFound(ctx, id)
					if err != nil {
						return nil, err
					}
					return nil, model.ErrRecordNotFound
				}
				return nil, err
			}
			// set cache
			err = d.cache.Set(ctx, id, table, cache.RatesExpireTime)
			if err != nil {
				return nil, fmt.Errorf("cache.Set error: %v, id=%d", err, id)
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.Rates)
		if !ok {
			return nil, model.ErrRecordNotFound
		}
		return table, nil
	} else if errors.Is(err, cacheBase.ErrPlaceholder) {
		return nil, model.ErrRecordNotFound
	}

	// fail fast, if cache error return, don't request to db
	return nil, err
}

// GetByColumns get paging records by column information,
// Note: query performance degrades when table rows are very large because of the use of offset.
//
// params includes paging parameters and query parameters
// paging parameters (required):
//
//	page: page number, starting from 0
//	limit: lines per page
//	sort: sort fields, default is id backwards, you can add - sign before the field to indicate reverse order, no - sign to indicate ascending order, multiple fields separated by comma
//
// query parameters (not required):
//
//	name: column name
//	exp: expressions, which default is "=",  support =, !=, >, >=, <, <=, like, in
//	value: column value, if exp=in, multiple values are separated by commas
//	logic: logical type, defaults to and when value is null, only &(and), ||(or)
//
// example: search for a male over 20 years of age
//
//	params = &query.Params{
//	    Page: 0,
//	    Limit: 20,
//	    Columns: []query.Column{
//		{
//			Name:    "age",
//			Exp: ">",
//			Value:   20,
//		},
//		{
//			Name:  "gender",
//			Value: "male",
//		},
//	}
func (d *ratesDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.Rates, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.Rates{}).Select([]string{"id"}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.Rates{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *ratesDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Rates) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *ratesDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	err := tx.WithContext(ctx).Where("id = ?", id).Delete(&model.Rates{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *ratesDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Rates) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
