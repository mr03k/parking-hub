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
	devicesCachePrefixKey = "devices:"
	// DevicesExpireTime expire time
	DevicesExpireTime = 5 * time.Minute
)

var _ DevicesCache = (*devicesCache)(nil)

// DevicesCache cache interface
type DevicesCache interface {
	Set(ctx context.Context, id uint64, data *model.Devices, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Devices, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Devices, error)
	MultiSet(ctx context.Context, data []*model.Devices, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// devicesCache define a cache struct
type devicesCache struct {
	cache cache.Cache
}

// NewDevicesCache new a cache
func NewDevicesCache(cacheType *model.CacheType) DevicesCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Devices{}
		})
		return &devicesCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Devices{}
		})
		return &devicesCache{cache: c}
	}

	return nil // no cache
}

// GetDevicesCacheKey cache key
func (c *devicesCache) GetDevicesCacheKey(id uint64) string {
	return devicesCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *devicesCache) Set(ctx context.Context, id uint64, data *model.Devices, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetDevicesCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *devicesCache) Get(ctx context.Context, id uint64) (*model.Devices, error) {
	var data *model.Devices
	cacheKey := c.GetDevicesCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *devicesCache) MultiSet(ctx context.Context, data []*model.Devices, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetDevicesCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *devicesCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Devices, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetDevicesCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Devices)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Devices)
	for _, id := range ids {
		val, ok := itemMap[c.GetDevicesCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *devicesCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetDevicesCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *devicesCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetDevicesCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
