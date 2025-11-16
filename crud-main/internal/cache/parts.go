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
	partsCachePrefixKey = "parts:"
	// PartsExpireTime expire time
	PartsExpireTime = 5 * time.Minute
)

var _ PartsCache = (*partsCache)(nil)

// PartsCache cache interface
type PartsCache interface {
	Set(ctx context.Context, id uint64, data *model.Parts, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Parts, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Parts, error)
	MultiSet(ctx context.Context, data []*model.Parts, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// partsCache define a cache struct
type partsCache struct {
	cache cache.Cache
}

// NewPartsCache new a cache
func NewPartsCache(cacheType *model.CacheType) PartsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Parts{}
		})
		return &partsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Parts{}
		})
		return &partsCache{cache: c}
	}

	return nil // no cache
}

// GetPartsCacheKey cache key
func (c *partsCache) GetPartsCacheKey(id uint64) string {
	return partsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *partsCache) Set(ctx context.Context, id uint64, data *model.Parts, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetPartsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *partsCache) Get(ctx context.Context, id uint64) (*model.Parts, error) {
	var data *model.Parts
	cacheKey := c.GetPartsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *partsCache) MultiSet(ctx context.Context, data []*model.Parts, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetPartsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *partsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Parts, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetPartsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Parts)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Parts)
	for _, id := range ids {
		val, ok := itemMap[c.GetPartsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *partsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetPartsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *partsCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetPartsCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
