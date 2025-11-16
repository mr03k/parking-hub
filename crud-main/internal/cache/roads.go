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
	roadsCachePrefixKey = "roads:"
	// RoadsExpireTime expire time
	RoadsExpireTime = 5 * time.Minute
)

var _ RoadsCache = (*roadsCache)(nil)

// RoadsCache cache interface
type RoadsCache interface {
	Set(ctx context.Context, id uint64, data *model.Roads, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Roads, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Roads, error)
	MultiSet(ctx context.Context, data []*model.Roads, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// roadsCache define a cache struct
type roadsCache struct {
	cache cache.Cache
}

// NewRoadsCache new a cache
func NewRoadsCache(cacheType *model.CacheType) RoadsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Roads{}
		})
		return &roadsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Roads{}
		})
		return &roadsCache{cache: c}
	}

	return nil // no cache
}

// GetRoadsCacheKey cache key
func (c *roadsCache) GetRoadsCacheKey(id uint64) string {
	return roadsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *roadsCache) Set(ctx context.Context, id uint64, data *model.Roads, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetRoadsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *roadsCache) Get(ctx context.Context, id uint64) (*model.Roads, error) {
	var data *model.Roads
	cacheKey := c.GetRoadsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *roadsCache) MultiSet(ctx context.Context, data []*model.Roads, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetRoadsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *roadsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Roads, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetRoadsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Roads)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Roads)
	for _, id := range ids {
		val, ok := itemMap[c.GetRoadsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *roadsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetRoadsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *roadsCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetRoadsCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
