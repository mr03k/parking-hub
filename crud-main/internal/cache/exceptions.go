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
	exceptionsCachePrefixKey = "exceptions:"
	// ExceptionsExpireTime expire time
	ExceptionsExpireTime = 5 * time.Minute
)

var _ ExceptionsCache = (*exceptionsCache)(nil)

// ExceptionsCache cache interface
type ExceptionsCache interface {
	Set(ctx context.Context, id uint64, data *model.Exceptions, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Exceptions, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Exceptions, error)
	MultiSet(ctx context.Context, data []*model.Exceptions, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// exceptionsCache define a cache struct
type exceptionsCache struct {
	cache cache.Cache
}

// NewExceptionsCache new a cache
func NewExceptionsCache(cacheType *model.CacheType) ExceptionsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Exceptions{}
		})
		return &exceptionsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Exceptions{}
		})
		return &exceptionsCache{cache: c}
	}

	return nil // no cache
}

// GetExceptionsCacheKey cache key
func (c *exceptionsCache) GetExceptionsCacheKey(id uint64) string {
	return exceptionsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *exceptionsCache) Set(ctx context.Context, id uint64, data *model.Exceptions, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetExceptionsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *exceptionsCache) Get(ctx context.Context, id uint64) (*model.Exceptions, error) {
	var data *model.Exceptions
	cacheKey := c.GetExceptionsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *exceptionsCache) MultiSet(ctx context.Context, data []*model.Exceptions, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetExceptionsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *exceptionsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Exceptions, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetExceptionsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Exceptions)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Exceptions)
	for _, id := range ids {
		val, ok := itemMap[c.GetExceptionsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *exceptionsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetExceptionsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *exceptionsCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetExceptionsCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
