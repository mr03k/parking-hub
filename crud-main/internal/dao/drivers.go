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

var _ DriversDao = (*driversDao)(nil)

// DriversDao defining the dao interface
type DriversDao interface {
	Create(ctx context.Context, table *model.Drivers) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.Drivers) error
	GetByID(ctx context.Context, id uint64) (*model.Drivers, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.Drivers, int64, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Drivers) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Drivers) error
}

type driversDao struct {
	db    *gorm.DB
	cache cache.DriversCache  // if nil, the cache is not used.
	sfg   *singleflight.Group // if cache is nil, the sfg is not used.
}

// NewDriversDao creating the dao interface
func NewDriversDao(db *gorm.DB, xCache cache.DriversCache) DriversDao {
	if xCache == nil {
		return &driversDao{db: db}
	}
	return &driversDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *driversDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a record, insert the record and the id value is written back to the table
func (d *driversDao) Create(ctx context.Context, table *model.Drivers) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a record by id
func (d *driversDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Drivers{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a record by id
func (d *driversDao) UpdateByID(ctx context.Context, table *model.Drivers) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *driversDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.Drivers) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.FirstName != "" {
		update["first_name"] = table.FirstName
	}
	if table.NameLast != "" {
		update["name_last"] = table.NameLast
	}
	if table.Gender != "" {
		update["gender"] = table.Gender
	}
	if table.CodeDriver != "" {
		update["code_driver"] = table.CodeDriver
	}
	if table.IDNational != "" {
		update["id_national"] = table.IDNational
	}
	if table.CodePostal != "" {
		update["code_postal"] = table.CodePostal
	}
	if table.NumberPhone != "" {
		update["number_phone"] = table.NumberPhone
	}
	if table.NumberMobile != "" {
		update["number_mobile"] = table.NumberMobile
	}
	if table.Email != "" {
		update["email"] = table.Email
	}
	if table.Address != "" {
		update["address"] = table.Address
	}
	if table.IDContractor != "" {
		update["id_contractor"] = table.IDContractor
	}
	if table.TypeDriver != "" {
		update["type_driver"] = table.TypeDriver
	}
	if table.TypeShift != "" {
		update["type_shift"] = table.TypeShift
	}
	if table.StatusEmployment != "" {
		update["status_employment"] = table.StatusEmployment
	}
	if table.DateStartEmployment.IsZero() == false {
		update["date_start_employment"] = table.DateStartEmployment
	}
	if table.DateEndEmployment.IsZero() == false {
		update["date_end_employment"] = table.DateEndEmployment
	}
	if table.DriverPhoto != "" {
		update["driver_photo"] = table.DriverPhoto
	}
	if table.ImageCardID != "" {
		update["image_card_id"] = table.ImageCardID
	}
	if table.BirthCertificateImage != "" {
		update["birth_certificate_image"] = table.BirthCertificateImage
	}
	if table.ImageCardServiceMilitary != "" {
		update["image_card_service_military"] = table.ImageCardServiceMilitary
	}
	if table.ImageCertificateHealth != "" {
		update["image_certificate_health"] = table.ImageCertificateHealth
	}
	if table.ImageRecordCriminal != "" {
		update["image_record_criminal"] = table.ImageRecordCriminal
	}
	if table.Description != "" {
		update["description"] = table.Description
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a record by id
func (d *driversDao) GetByID(ctx context.Context, id uint64) (*model.Drivers, error) {
	// no cache
	if d.cache == nil {
		record := &model.Drivers{}
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
			table := &model.Drivers{}
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
			err = d.cache.Set(ctx, id, table, cache.DriversExpireTime)
			if err != nil {
				return nil, fmt.Errorf("cache.Set error: %v, id=%d", err, id)
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.Drivers)
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
func (d *driversDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.Drivers, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.Drivers{}).Select([]string{"id"}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.Drivers{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *driversDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Drivers) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *driversDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	err := tx.WithContext(ctx).Where("id = ?", id).Delete(&model.Drivers{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *driversDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Drivers) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
