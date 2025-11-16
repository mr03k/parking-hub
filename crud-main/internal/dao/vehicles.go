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

var _ VehiclesDao = (*vehiclesDao)(nil)

// VehiclesDao defining the dao interface
type VehiclesDao interface {
	Create(ctx context.Context, table *model.Vehicles) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.Vehicles) error
	GetByID(ctx context.Context, id uint64) (*model.Vehicles, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.Vehicles, int64, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Vehicles) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Vehicles) error
}

type vehiclesDao struct {
	db    *gorm.DB
	cache cache.VehiclesCache // if nil, the cache is not used.
	sfg   *singleflight.Group // if cache is nil, the sfg is not used.
}

// NewVehiclesDao creating the dao interface
func NewVehiclesDao(db *gorm.DB, xCache cache.VehiclesCache) VehiclesDao {
	if xCache == nil {
		return &vehiclesDao{db: db}
	}
	return &vehiclesDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *vehiclesDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a record, insert the record and the id value is written back to the table
func (d *vehiclesDao) Create(ctx context.Context, table *model.Vehicles) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a record by id
func (d *vehiclesDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Vehicles{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a record by id
func (d *vehiclesDao) UpdateByID(ctx context.Context, table *model.Vehicles) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *vehiclesDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.Vehicles) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.CodeVehicle != "" {
		update["code_vehicle"] = table.CodeVehicle
	}
	if table.Vin != "" {
		update["vin"] = table.Vin
	}
	if table.PlateLicense != "" {
		update["plate_license"] = table.PlateLicense
	}
	if table.TypeVehicle != "" {
		update["type_vehicle"] = table.TypeVehicle
	}
	if table.Brand != "" {
		update["brand"] = table.Brand
	}
	if table.Model != "" {
		update["model"] = table.Model
	}
	if table.Color != "" {
		update["color"] = table.Color
	}
	if table.ManufactureOfYear != 0 {
		update["manufacture_of_year"] = table.ManufactureOfYear
	}
	if table.KilometersInitial != 0 {
		update["kilometers_initial"] = table.KilometersInitial
	}
	if table.ExpiryInsurancePartyThird.IsZero() == false {
		update["expiry_insurance_party_third"] = table.ExpiryInsurancePartyThird
	}
	if table.ExpiryInsuranceBody.IsZero() == false {
		update["expiry_insurance_body"] = table.ExpiryInsuranceBody
	}
	if table.ImageDocumentVehicle != "" {
		update["image_document_vehicle"] = table.ImageDocumentVehicle
	}
	if table.ImageCardVehicle != "" {
		update["image_card_vehicle"] = table.ImageCardVehicle
	}
	if table.ThirdPartyInsuranceImage != "" {
		update["third_party_insurance_image"] = table.ThirdPartyInsuranceImage
	}
	if table.BodyInsuranceImage != "" {
		update["body_insurance_image"] = table.BodyInsuranceImage
	}
	if table.IDContractor != "" {
		update["id_contractor"] = table.IDContractor
	}
	if table.Status != "" {
		update["status"] = table.Status
	}
	if table.Description != "" {
		update["description"] = table.Description
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a record by id
func (d *vehiclesDao) GetByID(ctx context.Context, id uint64) (*model.Vehicles, error) {
	// no cache
	if d.cache == nil {
		record := &model.Vehicles{}
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
			table := &model.Vehicles{}
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
			err = d.cache.Set(ctx, id, table, cache.VehiclesExpireTime)
			if err != nil {
				return nil, fmt.Errorf("cache.Set error: %v, id=%d", err, id)
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.Vehicles)
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
func (d *vehiclesDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.Vehicles, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.Vehicles{}).Select([]string{"id"}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.Vehicles{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *vehiclesDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Vehicles) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *vehiclesDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	err := tx.WithContext(ctx).Where("id = ?", id).Delete(&model.Vehicles{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *vehiclesDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Vehicles) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
