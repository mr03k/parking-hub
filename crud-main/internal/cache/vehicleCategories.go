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
	vehicleCategoriesCachePrefixKey = "vehicleCategories:"
	// VehicleCategoriesExpireTime expire time
	VehicleCategoriesExpireTime = 5 * time.Minute
)

var _ VehicleCategoriesCache = (*vehicleCategoriesCache)(nil)

// VehicleCategoriesCache cache interface
type VehicleCategoriesCache interface {
	Set(ctx context.Context, id uint64, data *model.VehicleCategories, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.VehicleCategories, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.VehicleCategories, error)
	MultiSet(ctx context.Context, data []*model.VehicleCategories, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// vehicleCategoriesCache define a cache struct
type vehicleCategoriesCache struct {
	cache cache.Cache
}

// NewVehicleCategoriesCache new a cache
func NewVehicleCategoriesCache(cacheType *model.CacheType) VehicleCategoriesCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.VehicleCategories{}
		})
		return &vehicleCategoriesCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.VehicleCategories{}
		})
		return &vehicleCategoriesCache{cache: c}
	}

	return nil // no cache
}

// GetVehicleCategoriesCacheKey cache key
func (c *vehicleCategoriesCache) GetVehicleCategoriesCacheKey(id uint64) string {
	return vehicleCategoriesCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *vehicleCategoriesCache) Set(ctx context.Context, id uint64, data *model.VehicleCategories, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetVehicleCategoriesCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *vehicleCategoriesCache) Get(ctx context.Context, id uint64) (*model.VehicleCategories, error) {
	var data *model.VehicleCategories
	cacheKey := c.GetVehicleCategoriesCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *vehicleCategoriesCache) MultiSet(ctx context.Context, data []*model.VehicleCategories, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetVehicleCategoriesCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *vehicleCategoriesCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.VehicleCategories, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetVehicleCategoriesCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.VehicleCategories)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.VehicleCategories)
	for _, id := range ids {
		val, ok := itemMap[c.GetVehicleCategoriesCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *vehicleCategoriesCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetVehicleCategoriesCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *vehicleCategoriesCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetVehicleCategoriesCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
