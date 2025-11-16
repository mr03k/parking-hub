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
	formsCachePrefixKey = "forms:"
	// FormsExpireTime expire time
	FormsExpireTime = 5 * time.Minute
)

var _ FormsCache = (*formsCache)(nil)

// FormsCache cache interface
type FormsCache interface {
	Set(ctx context.Context, id uint64, data *model.Forms, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Forms, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Forms, error)
	MultiSet(ctx context.Context, data []*model.Forms, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// formsCache define a cache struct
type formsCache struct {
	cache cache.Cache
}

// NewFormsCache new a cache
func NewFormsCache(cacheType *model.CacheType) FormsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Forms{}
		})
		return &formsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Forms{}
		})
		return &formsCache{cache: c}
	}

	return nil // no cache
}

// GetFormsCacheKey cache key
func (c *formsCache) GetFormsCacheKey(id uint64) string {
	return formsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *formsCache) Set(ctx context.Context, id uint64, data *model.Forms, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetFormsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *formsCache) Get(ctx context.Context, id uint64) (*model.Forms, error) {
	var data *model.Forms
	cacheKey := c.GetFormsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *formsCache) MultiSet(ctx context.Context, data []*model.Forms, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetFormsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *formsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Forms, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetFormsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Forms)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Forms)
	for _, id := range ids {
		val, ok := itemMap[c.GetFormsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *formsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetFormsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *formsCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetFormsCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
