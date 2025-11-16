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
	modulesCachePrefixKey = "modules:"
	// ModulesExpireTime expire time
	ModulesExpireTime = 5 * time.Minute
)

var _ ModulesCache = (*modulesCache)(nil)

// ModulesCache cache interface
type ModulesCache interface {
	Set(ctx context.Context, id uint64, data *model.Modules, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Modules, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Modules, error)
	MultiSet(ctx context.Context, data []*model.Modules, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// modulesCache define a cache struct
type modulesCache struct {
	cache cache.Cache
}

// NewModulesCache new a cache
func NewModulesCache(cacheType *model.CacheType) ModulesCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Modules{}
		})
		return &modulesCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Modules{}
		})
		return &modulesCache{cache: c}
	}

	return nil // no cache
}

// GetModulesCacheKey cache key
func (c *modulesCache) GetModulesCacheKey(id uint64) string {
	return modulesCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *modulesCache) Set(ctx context.Context, id uint64, data *model.Modules, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetModulesCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *modulesCache) Get(ctx context.Context, id uint64) (*model.Modules, error) {
	var data *model.Modules
	cacheKey := c.GetModulesCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *modulesCache) MultiSet(ctx context.Context, data []*model.Modules, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetModulesCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *modulesCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Modules, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetModulesCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Modules)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Modules)
	for _, id := range ids {
		val, ok := itemMap[c.GetModulesCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *modulesCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetModulesCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *modulesCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetModulesCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
