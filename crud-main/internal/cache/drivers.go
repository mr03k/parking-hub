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
	driversCachePrefixKey = "drivers:"
	// DriversExpireTime expire time
	DriversExpireTime = 5 * time.Minute
)

var _ DriversCache = (*driversCache)(nil)

// DriversCache cache interface
type DriversCache interface {
	Set(ctx context.Context, id uint64, data *model.Drivers, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Drivers, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Drivers, error)
	MultiSet(ctx context.Context, data []*model.Drivers, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// driversCache define a cache struct
type driversCache struct {
	cache cache.Cache
}

// NewDriversCache new a cache
func NewDriversCache(cacheType *model.CacheType) DriversCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Drivers{}
		})
		return &driversCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Drivers{}
		})
		return &driversCache{cache: c}
	}

	return nil // no cache
}

// GetDriversCacheKey cache key
func (c *driversCache) GetDriversCacheKey(id uint64) string {
	return driversCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *driversCache) Set(ctx context.Context, id uint64, data *model.Drivers, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetDriversCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *driversCache) Get(ctx context.Context, id uint64) (*model.Drivers, error) {
	var data *model.Drivers
	cacheKey := c.GetDriversCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *driversCache) MultiSet(ctx context.Context, data []*model.Drivers, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetDriversCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *driversCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Drivers, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetDriversCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Drivers)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Drivers)
	for _, id := range ids {
		val, ok := itemMap[c.GetDriversCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *driversCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetDriversCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *driversCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetDriversCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
