package cache

import (
	"context"
	"strings"
	"time"

	"git.abanppc.com/farin-project/crud/internal/model"
	"github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/encoding"
)

const (
	// cache prefix key, must end with a colon
	parkingsCachePrefixKey = "parkings:"
	// ParkingsExpireTime expire time
	ParkingsExpireTime = 5 * time.Minute
)

var _ ParkingsCache = (*parkingsCache)(nil)

// ParkingsCache cache interface
type ParkingsCache interface {
	Set(ctx context.Context, id string, data *model.Parkings, duration time.Duration) error
	Get(ctx context.Context, id string) (*model.Parkings, error)
	MultiGet(ctx context.Context, ids []string) (map[string]*model.Parkings, error)
	MultiSet(ctx context.Context, data []*model.Parkings, duration time.Duration) error
	Del(ctx context.Context, id string) error
	SetCacheWithNotFound(ctx context.Context, id string) error
}

// parkingsCache define a cache struct
type parkingsCache struct {
	cache cache.Cache
}

// NewParkingsCache new a cache
func NewParkingsCache(cacheType *model.CacheType) ParkingsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Parkings{}
		})
		return &parkingsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Parkings{}
		})
		return &parkingsCache{cache: c}
	}

	return nil // no cache
}

// GetParkingsCacheKey cache key
func (c *parkingsCache) GetParkingsCacheKey(id string) string {
	return parkingsCachePrefixKey + id
}

// Set write to cache
func (c *parkingsCache) Set(ctx context.Context, id string, data *model.Parkings, duration time.Duration) error {
	if data == nil || id == "" {
		return nil
	}
	cacheKey := c.GetParkingsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *parkingsCache) Get(ctx context.Context, id string) (*model.Parkings, error) {
	var data *model.Parkings
	cacheKey := c.GetParkingsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *parkingsCache) MultiSet(ctx context.Context, data []*model.Parkings, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetParkingsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *parkingsCache) MultiGet(ctx context.Context, ids []string) (map[string]*model.Parkings, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetParkingsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Parkings)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[string]*model.Parkings)
	for _, id := range ids {
		val, ok := itemMap[c.GetParkingsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *parkingsCache) Del(ctx context.Context, id string) error {
	cacheKey := c.GetParkingsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *parkingsCache) SetCacheWithNotFound(ctx context.Context, id string) error {
	cacheKey := c.GetParkingsCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
