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
	ratesCachePrefixKey = "rates:"
	// RatesExpireTime expire time
	RatesExpireTime = 5 * time.Minute
)

var _ RatesCache = (*ratesCache)(nil)

// RatesCache cache interface
type RatesCache interface {
	Set(ctx context.Context, id uint64, data *model.Rates, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Rates, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Rates, error)
	MultiSet(ctx context.Context, data []*model.Rates, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// ratesCache define a cache struct
type ratesCache struct {
	cache cache.Cache
}

// NewRatesCache new a cache
func NewRatesCache(cacheType *model.CacheType) RatesCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Rates{}
		})
		return &ratesCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Rates{}
		})
		return &ratesCache{cache: c}
	}

	return nil // no cache
}

// GetRatesCacheKey cache key
func (c *ratesCache) GetRatesCacheKey(id uint64) string {
	return ratesCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *ratesCache) Set(ctx context.Context, id uint64, data *model.Rates, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetRatesCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *ratesCache) Get(ctx context.Context, id uint64) (*model.Rates, error) {
	var data *model.Rates
	cacheKey := c.GetRatesCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *ratesCache) MultiSet(ctx context.Context, data []*model.Rates, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetRatesCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *ratesCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Rates, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetRatesCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Rates)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Rates)
	for _, id := range ids {
		val, ok := itemMap[c.GetRatesCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *ratesCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetRatesCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *ratesCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetRatesCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
