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
	citiesCachePrefixKey = "cities:"
	// CitiesExpireTime expire time
	CitiesExpireTime = 5 * time.Minute
)

var _ CitiesCache = (*citiesCache)(nil)

// CitiesCache cache interface
type CitiesCache interface {
	Set(ctx context.Context, id uint64, data *model.Cities, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Cities, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Cities, error)
	MultiSet(ctx context.Context, data []*model.Cities, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// citiesCache define a cache struct
type citiesCache struct {
	cache cache.Cache
}

// NewCitiesCache new a cache
func NewCitiesCache(cacheType *model.CacheType) CitiesCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Cities{}
		})
		return &citiesCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Cities{}
		})
		return &citiesCache{cache: c}
	}

	return nil // no cache
}

// GetCitiesCacheKey cache key
func (c *citiesCache) GetCitiesCacheKey(id uint64) string {
	return citiesCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *citiesCache) Set(ctx context.Context, id uint64, data *model.Cities, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetCitiesCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *citiesCache) Get(ctx context.Context, id uint64) (*model.Cities, error) {
	var data *model.Cities
	cacheKey := c.GetCitiesCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *citiesCache) MultiSet(ctx context.Context, data []*model.Cities, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetCitiesCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *citiesCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Cities, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetCitiesCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Cities)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Cities)
	for _, id := range ids {
		val, ok := itemMap[c.GetCitiesCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *citiesCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetCitiesCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *citiesCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetCitiesCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
