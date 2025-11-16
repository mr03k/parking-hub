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
	roadCategoriesCachePrefixKey = "roadCategories:"
	// RoadCategoriesExpireTime expire time
	RoadCategoriesExpireTime = 5 * time.Minute
)

var _ RoadCategoriesCache = (*roadCategoriesCache)(nil)

// RoadCategoriesCache cache interface
type RoadCategoriesCache interface {
	Set(ctx context.Context, id uint64, data *model.RoadCategories, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.RoadCategories, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.RoadCategories, error)
	MultiSet(ctx context.Context, data []*model.RoadCategories, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// roadCategoriesCache define a cache struct
type roadCategoriesCache struct {
	cache cache.Cache
}

// NewRoadCategoriesCache new a cache
func NewRoadCategoriesCache(cacheType *model.CacheType) RoadCategoriesCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.RoadCategories{}
		})
		return &roadCategoriesCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.RoadCategories{}
		})
		return &roadCategoriesCache{cache: c}
	}

	return nil // no cache
}

// GetRoadCategoriesCacheKey cache key
func (c *roadCategoriesCache) GetRoadCategoriesCacheKey(id uint64) string {
	return roadCategoriesCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *roadCategoriesCache) Set(ctx context.Context, id uint64, data *model.RoadCategories, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetRoadCategoriesCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *roadCategoriesCache) Get(ctx context.Context, id uint64) (*model.RoadCategories, error) {
	var data *model.RoadCategories
	cacheKey := c.GetRoadCategoriesCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *roadCategoriesCache) MultiSet(ctx context.Context, data []*model.RoadCategories, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetRoadCategoriesCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *roadCategoriesCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.RoadCategories, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetRoadCategoriesCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.RoadCategories)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.RoadCategories)
	for _, id := range ids {
		val, ok := itemMap[c.GetRoadCategoriesCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *roadCategoriesCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetRoadCategoriesCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *roadCategoriesCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetRoadCategoriesCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
