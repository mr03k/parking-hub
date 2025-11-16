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
	ringsCachePrefixKey = "rings:"
	// RingsExpireTime expire time
	RingsExpireTime = 5 * time.Minute
)

var _ RingsCache = (*ringsCache)(nil)

// RingsCache cache interface
type RingsCache interface {
	Set(ctx context.Context, id uint64, data *model.Rings, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Rings, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Rings, error)
	MultiSet(ctx context.Context, data []*model.Rings, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// ringsCache define a cache struct
type ringsCache struct {
	cache cache.Cache
}

// NewRingsCache new a cache
func NewRingsCache(cacheType *model.CacheType) RingsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Rings{}
		})
		return &ringsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Rings{}
		})
		return &ringsCache{cache: c}
	}

	return nil // no cache
}

// GetRingsCacheKey cache key
func (c *ringsCache) GetRingsCacheKey(id uint64) string {
	return ringsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *ringsCache) Set(ctx context.Context, id uint64, data *model.Rings, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetRingsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *ringsCache) Get(ctx context.Context, id uint64) (*model.Rings, error) {
	var data *model.Rings
	cacheKey := c.GetRingsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *ringsCache) MultiSet(ctx context.Context, data []*model.Rings, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetRingsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *ringsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Rings, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetRingsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Rings)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Rings)
	for _, id := range ids {
		val, ok := itemMap[c.GetRingsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *ringsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetRingsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *ringsCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetRingsCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
