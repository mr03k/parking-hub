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
	vehiclesCachePrefixKey = "vehicles:"
	// VehiclesExpireTime expire time
	VehiclesExpireTime = 5 * time.Minute
)

var _ VehiclesCache = (*vehiclesCache)(nil)

// VehiclesCache cache interface
type VehiclesCache interface {
	Set(ctx context.Context, id uint64, data *model.Vehicles, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Vehicles, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Vehicles, error)
	MultiSet(ctx context.Context, data []*model.Vehicles, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// vehiclesCache define a cache struct
type vehiclesCache struct {
	cache cache.Cache
}

// NewVehiclesCache new a cache
func NewVehiclesCache(cacheType *model.CacheType) VehiclesCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Vehicles{}
		})
		return &vehiclesCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Vehicles{}
		})
		return &vehiclesCache{cache: c}
	}

	return nil // no cache
}

// GetVehiclesCacheKey cache key
func (c *vehiclesCache) GetVehiclesCacheKey(id uint64) string {
	return vehiclesCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *vehiclesCache) Set(ctx context.Context, id uint64, data *model.Vehicles, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetVehiclesCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *vehiclesCache) Get(ctx context.Context, id uint64) (*model.Vehicles, error) {
	var data *model.Vehicles
	cacheKey := c.GetVehiclesCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *vehiclesCache) MultiSet(ctx context.Context, data []*model.Vehicles, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetVehiclesCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *vehiclesCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Vehicles, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetVehiclesCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Vehicles)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Vehicles)
	for _, id := range ids {
		val, ok := itemMap[c.GetVehiclesCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *vehiclesCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetVehiclesCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *vehiclesCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetVehiclesCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
