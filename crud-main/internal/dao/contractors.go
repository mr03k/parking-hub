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

var _ ContractorsDao = (*contractorsDao)(nil)

// ContractorsDao defining the dao interface
type ContractorsDao interface {
	Create(ctx context.Context, table *model.Contractors) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.Contractors) error
	GetByID(ctx context.Context, id uint64) (*model.Contractors, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.Contractors, int64, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Contractors) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Contractors) error
}

type contractorsDao struct {
	db    *gorm.DB
	cache cache.ContractorsCache // if nil, the cache is not used.
	sfg   *singleflight.Group    // if cache is nil, the sfg is not used.
}

// NewContractorsDao creating the dao interface
func NewContractorsDao(db *gorm.DB, xCache cache.ContractorsCache) ContractorsDao {
	if xCache == nil {
		return &contractorsDao{db: db}
	}
	return &contractorsDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *contractorsDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a record, insert the record and the id value is written back to the table
func (d *contractorsDao) Create(ctx context.Context, table *model.Contractors) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a record by id
func (d *contractorsDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Contractors{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a record by id
func (d *contractorsDao) UpdateByID(ctx context.Context, table *model.Contractors) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *contractorsDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.Contractors) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.ContractorName != "" {
		update["contractor_name"] = table.ContractorName
	}
	if table.CodeContractor != "" {
		update["code_contractor"] = table.CodeContractor
	}
	if table.NumberRegistration != "" {
		update["number_registration"] = table.NumberRegistration
	}
	if table.PersonContact != "" {
		update["person_contact"] = table.PersonContact
	}
	if table.CeoName != "" {
		update["ceo_name"] = table.CeoName
	}
	if table.SignatoriesAuthorized != "" {
		update["signatories_authorized"] = table.SignatoriesAuthorized
	}
	if table.PhoneNumber != "" {
		update["phone_number"] = table.PhoneNumber
	}
	if table.Email != "" {
		update["email"] = table.Email
	}
	if table.Address != "" {
		update["address"] = table.Address
	}
	if table.TypeContract != "" {
		update["type_contract"] = table.TypeContract
	}
	if table.NumberAccountBank != "" {
		update["number_account_bank"] = table.NumberAccountBank
	}
	if table.Description != "" {
		update["description"] = table.Description
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a record by id
func (d *contractorsDao) GetByID(ctx context.Context, id uint64) (*model.Contractors, error) {
	// no cache
	if d.cache == nil {
		record := &model.Contractors{}
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
			table := &model.Contractors{}
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
			err = d.cache.Set(ctx, id, table, cache.ContractorsExpireTime)
			if err != nil {
				return nil, fmt.Errorf("cache.Set error: %v, id=%d", err, id)
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.Contractors)
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
func (d *contractorsDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.Contractors, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.Contractors{}).Select([]string{"id"}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.Contractors{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *contractorsDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Contractors) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *contractorsDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	err := tx.WithContext(ctx).Where("id = ?", id).Delete(&model.Contractors{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *contractorsDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Contractors) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
