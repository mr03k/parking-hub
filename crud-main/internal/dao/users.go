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

var _ UsersDao = (*usersDao)(nil)

// UsersDao defining the dao interface
type UsersDao interface {
	Create(ctx context.Context, table *model.Users) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.Users) error
	GetByID(ctx context.Context, id uint64) (*model.Users, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.Users, int64, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Users) (string, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Users) error
}

type usersDao struct {
	db    *gorm.DB
	cache cache.UsersCache    // if nil, the cache is not used.
	sfg   *singleflight.Group // if cache is nil, the sfg is not used.
}

// NewUsersDao creating the dao interface
func NewUsersDao(db *gorm.DB, xCache cache.UsersCache) UsersDao {
	if xCache == nil {
		return &usersDao{db: db}
	}
	return &usersDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *usersDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a record, insert the record and the id value is written back to the table
func (d *usersDao) Create(ctx context.Context, table *model.Users) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a record by id
func (d *usersDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Users{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a record by id
func (d *usersDao) UpdateByID(ctx context.Context, table *model.Users) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	//_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *usersDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.Users) error {
	if table.ID == "" {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.Username != "" {
		update["username"] = table.Username
	}
	if table.Password != "" {
		update["password"] = table.Password
	}
	if table.FirstName != "" {
		update["first_name"] = table.FirstName
	}
	if table.LastName != "" {
		update["last_name"] = table.LastName
	}
	if table.Email != "" {
		update["email"] = table.Email
	}
	if table.NumberPhone != "" {
		update["number_phone"] = table.NumberPhone
	}
	if table.NumberMobile != "" {
		update["number_mobile"] = table.NumberMobile
	}
	if table.IDNational != "" {
		update["id_national"] = table.IDNational
	}
	if table.CodePostal != "" {
		update["code_postal"] = table.CodePostal
	}
	if table.NameCompany != "" {
		update["name_company"] = table.NameCompany
	}
	if table.ImageProfile != "" {
		update["image_profile"] = table.ImageProfile
	}
	if table.Gender != "" {
		update["gender"] = table.Gender
	}
	if table.Address != "" {
		update["address"] = table.Address
	}
	if table.Status != "" {
		update["status"] = table.Status
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a record by id
func (d *usersDao) GetByID(ctx context.Context, id uint64) (*model.Users, error) {
	// no cache
	if d.cache == nil {
		record := &model.Users{}
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
			table := &model.Users{}
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
			err = d.cache.Set(ctx, id, table, cache.UsersExpireTime)
			if err != nil {
				return nil, fmt.Errorf("cache.Set error: %v, id=%d", err, id)
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.Users)
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
func (d *usersDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.Users, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.Users{}).Select([]string{"id"}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.Users{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *usersDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Users) (string, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *usersDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	err := tx.WithContext(ctx).Where("id = ?", id).Delete(&model.Users{}).Error
	if err != nil {
		return err
	}

	// delete cache
	//_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *usersDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Users) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	//_ = d.deleteCache(ctx, table.ID)

	return err
}
