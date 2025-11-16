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
	devicePartsCachePrefixKey = "deviceParts:"
	// DevicePartsExpireTime expire time
	DevicePartsExpireTime = 5 * time.Minute
)

var _ DevicePartsCache = (*devicePartsCache)(nil)

// DevicePartsCache cache interface
type DevicePartsCache interface {
	Set(ctx context.Context, id uint64, data *model.DeviceParts, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.DeviceParts, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.DeviceParts, error)
	MultiSet(ctx context.Context, data []*model.DeviceParts, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// devicePartsCache define a cache struct
type devicePartsCache struct {
	cache cache.Cache
}

// NewDevicePartsCache new a cache
func NewDevicePartsCache(cacheType *model.CacheType) DevicePartsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.DeviceParts{}
		})
		return &devicePartsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.DeviceParts{}
		})
		return &devicePartsCache{cache: c}
	}

	return nil // no cache
}

// GetDevicePartsCacheKey cache key
func (c *devicePartsCache) GetDevicePartsCacheKey(id uint64) string {
	return devicePartsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *devicePartsCache) Set(ctx context.Context, id uint64, data *model.DeviceParts, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetDevicePartsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *devicePartsCache) Get(ctx context.Context, id uint64) (*model.DeviceParts, error) {
	var data *model.DeviceParts
	cacheKey := c.GetDevicePartsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *devicePartsCache) MultiSet(ctx context.Context, data []*model.DeviceParts, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetDevicePartsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *devicePartsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.DeviceParts, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetDevicePartsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.DeviceParts)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.DeviceParts)
	for _, id := range ids {
		val, ok := itemMap[c.GetDevicePartsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *devicePartsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetDevicePartsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *devicePartsCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetDevicePartsCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
