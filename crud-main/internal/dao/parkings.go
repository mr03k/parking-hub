package dao

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	"git.abanppc.com/farin-project/crud/internal/cache"
	"git.abanppc.com/farin-project/crud/internal/model"
	cacheBase "github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ ParkingsDao = (*parkingsDao)(nil)

// ParkingsDao defining the dao interface
type ParkingsDao interface {
	Create(ctx context.Context, table *model.Parkings) error
	DeleteByID(ctx context.Context, id string) error
	UpdateByID(ctx context.Context, table *model.Parkings) error
	GetByID(ctx context.Context, id string) (*model.Parkings, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.Parkings, int64, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Parkings) (string, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id string) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Parkings) error
}

type parkingsDao struct {
	db    *gorm.DB
	cache cache.ParkingsCache // if nil, the cache is not used.
	sfg   *singleflight.Group // if cache is nil, the sfg is not used.
}

// NewParkingsDao creating the dao interface
func NewParkingsDao(db *gorm.DB, xCache cache.ParkingsCache) ParkingsDao {
	if xCache == nil {
		return &parkingsDao{db: db}
	}
	return &parkingsDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *parkingsDao) deleteCache(ctx context.Context, id string) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a record, insert the record and the id value is written back to the table
func (d *parkingsDao) Create(ctx context.Context, table *model.Parkings) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a record by id
func (d *parkingsDao) DeleteByID(ctx context.Context, id string) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Parkings{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a record by id
func (d *parkingsDao) UpdateByID(ctx context.Context, table *model.Parkings) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *parkingsDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.Parkings) error {
	if table.ID == "" {
		return errors.New("id cannot be empty")
	}

	update := map[string]interface{}{}

	if table.CodeParking != "" {
		update["code_parking"] = table.CodeParking
	}
	if table.IDSegment != "" {
		update["id_segment"] = table.IDSegment
	}
	if table.TypeParking != "" {
		update["type_parking"] = table.TypeParking
	}
	if table.BoundaryParking != "" {
		update["boundary_parking"] = table.BoundaryParking
	}
	if table.Position != "" {
		update["position"] = table.Position
	}
	if table.StatusAvailability != "" {
		update["status_availability"] = table.StatusAvailability
	}
	if table.Description != "" {
		update["description"] = table.Description
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a record by id
func (d *parkingsDao) GetByID(ctx context.Context, id string) (*model.Parkings, error) {
	// no cache
	if d.cache == nil {
		record := &model.Parkings{}
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
		val, err, _ := d.sfg.Do(id, func() (interface{}, error) { //nolint
			table := &model.Parkings{}
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
			err = d.cache.Set(ctx, id, table, cache.ParkingsExpireTime)
			if err != nil {
				return nil, fmt.Errorf("cache.Set error: %v, id=%d", err, id)
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.Parkings)
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
func (d *parkingsDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.Parkings, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.Parkings{}).Select([]string{"id"}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.Parkings{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *parkingsDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Parkings) (string, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *parkingsDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id string) error {
	err := tx.WithContext(ctx).Where("id = ?", id).Delete(&model.Parkings{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *parkingsDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Parkings) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
