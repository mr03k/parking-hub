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
	peakHourMultipliersCachePrefixKey = "peakHourMultipliers:"
	// PeakHourMultipliersExpireTime expire time
	PeakHourMultipliersExpireTime = 5 * time.Minute
)

var _ PeakHourMultipliersCache = (*peakHourMultipliersCache)(nil)

// PeakHourMultipliersCache cache interface
type PeakHourMultipliersCache interface {
	Set(ctx context.Context, id uint64, data *model.PeakHourMultipliers, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.PeakHourMultipliers, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.PeakHourMultipliers, error)
	MultiSet(ctx context.Context, data []*model.PeakHourMultipliers, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// peakHourMultipliersCache define a cache struct
type peakHourMultipliersCache struct {
	cache cache.Cache
}

// NewPeakHourMultipliersCache new a cache
func NewPeakHourMultipliersCache(cacheType *model.CacheType) PeakHourMultipliersCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.PeakHourMultipliers{}
		})
		return &peakHourMultipliersCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.PeakHourMultipliers{}
		})
		return &peakHourMultipliersCache{cache: c}
	}

	return nil // no cache
}

// GetPeakHourMultipliersCacheKey cache key
func (c *peakHourMultipliersCache) GetPeakHourMultipliersCacheKey(id uint64) string {
	return peakHourMultipliersCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *peakHourMultipliersCache) Set(ctx context.Context, id uint64, data *model.PeakHourMultipliers, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetPeakHourMultipliersCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *peakHourMultipliersCache) Get(ctx context.Context, id uint64) (*model.PeakHourMultipliers, error) {
	var data *model.PeakHourMultipliers
	cacheKey := c.GetPeakHourMultipliersCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *peakHourMultipliersCache) MultiSet(ctx context.Context, data []*model.PeakHourMultipliers, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetPeakHourMultipliersCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *peakHourMultipliersCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.PeakHourMultipliers, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetPeakHourMultipliersCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.PeakHourMultipliers)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.PeakHourMultipliers)
	for _, id := range ids {
		val, ok := itemMap[c.GetPeakHourMultipliersCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *peakHourMultipliersCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetPeakHourMultipliersCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *peakHourMultipliersCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetPeakHourMultipliersCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
