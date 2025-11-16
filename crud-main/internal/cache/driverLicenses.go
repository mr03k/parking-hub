package cache

import (
	"context"
	"strings"
	"time"

	"github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/encoding"
	"github.com/zhufuyi/sponge/pkg/utils"

	"git.abanppc.com/farin-project/crud/internal/model"
)

const (
	// cache prefix key, must end with a colon
	driverLicensesCachePrefixKey = "driverLicenses:"
	// DriverLicensesExpireTime expire time
	DriverLicensesExpireTime = 5 * time.Minute
)

var _ DriverLicensesCache = (*driverLicensesCache)(nil)

// DriverLicensesCache cache interface
type DriverLicensesCache interface {
	Set(ctx context.Context, id uint64, data *model.DriverLicenses, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.DriverLicenses, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.DriverLicenses, error)
	MultiSet(ctx context.Context, data []*model.DriverLicenses, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// driverLicensesCache define a cache struct
type driverLicensesCache struct {
	cache cache.Cache
}

// NewDriverLicensesCache new a cache
func NewDriverLicensesCache(cacheType *model.CacheType) DriverLicensesCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.DriverLicenses{}
		})
		return &driverLicensesCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.DriverLicenses{}
		})
		return &driverLicensesCache{cache: c}
	}

	return nil // no cache
}

// GetDriverLicensesCacheKey cache key
func (c *driverLicensesCache) GetDriverLicensesCacheKey(id uint64) string {
	return driverLicensesCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *driverLicensesCache) Set(ctx context.Context, id uint64, data *model.DriverLicenses, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetDriverLicensesCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *driverLicensesCache) Get(ctx context.Context, id uint64) (*model.DriverLicenses, error) {
	var data *model.DriverLicenses
	cacheKey := c.GetDriverLicensesCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *driverLicensesCache) MultiSet(ctx context.Context, data []*model.DriverLicenses, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetDriverLicensesCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *driverLicensesCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.DriverLicenses, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetDriverLicensesCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.DriverLicenses)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.DriverLicenses)
	for _, id := range ids {
		val, ok := itemMap[c.GetDriverLicensesCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *driverLicensesCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetDriverLicensesCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *driverLicensesCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetDriverLicensesCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
