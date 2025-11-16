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
	baseRatesCachePrefixKey = "baseRates:"
	// BaseRatesExpireTime expire time
	BaseRatesExpireTime = 5 * time.Minute
)

var _ BaseRatesCache = (*baseRatesCache)(nil)

// BaseRatesCache cache interface
type BaseRatesCache interface {
	Set(ctx context.Context, id uint64, data *model.BaseRates, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.BaseRates, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.BaseRates, error)
	MultiSet(ctx context.Context, data []*model.BaseRates, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// baseRatesCache define a cache struct
type baseRatesCache struct {
	cache cache.Cache
}

// NewBaseRatesCache new a cache
func NewBaseRatesCache(cacheType *model.CacheType) BaseRatesCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.BaseRates{}
		})
		return &baseRatesCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.BaseRates{}
		})
		return &baseRatesCache{cache: c}
	}

	return nil // no cache
}

// GetBaseRatesCacheKey cache key
func (c *baseRatesCache) GetBaseRatesCacheKey(id uint64) string {
	return baseRatesCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *baseRatesCache) Set(ctx context.Context, id uint64, data *model.BaseRates, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetBaseRatesCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *baseRatesCache) Get(ctx context.Context, id uint64) (*model.BaseRates, error) {
	var data *model.BaseRates
	cacheKey := c.GetBaseRatesCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *baseRatesCache) MultiSet(ctx context.Context, data []*model.BaseRates, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetBaseRatesCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *baseRatesCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.BaseRates, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetBaseRatesCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.BaseRates)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.BaseRates)
	for _, id := range ids {
		val, ok := itemMap[c.GetBaseRatesCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *baseRatesCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetBaseRatesCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *baseRatesCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetBaseRatesCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
