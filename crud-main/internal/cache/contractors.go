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
	contractorsCachePrefixKey = "contractors:"
	// ContractorsExpireTime expire time
	ContractorsExpireTime = 5 * time.Minute
)

var _ ContractorsCache = (*contractorsCache)(nil)

// ContractorsCache cache interface
type ContractorsCache interface {
	Set(ctx context.Context, id uint64, data *model.Contractors, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Contractors, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Contractors, error)
	MultiSet(ctx context.Context, data []*model.Contractors, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// contractorsCache define a cache struct
type contractorsCache struct {
	cache cache.Cache
}

// NewContractorsCache new a cache
func NewContractorsCache(cacheType *model.CacheType) ContractorsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Contractors{}
		})
		return &contractorsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Contractors{}
		})
		return &contractorsCache{cache: c}
	}

	return nil // no cache
}

// GetContractorsCacheKey cache key
func (c *contractorsCache) GetContractorsCacheKey(id uint64) string {
	return contractorsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *contractorsCache) Set(ctx context.Context, id uint64, data *model.Contractors, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetContractorsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *contractorsCache) Get(ctx context.Context, id uint64) (*model.Contractors, error) {
	var data *model.Contractors
	cacheKey := c.GetContractorsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *contractorsCache) MultiSet(ctx context.Context, data []*model.Contractors, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetContractorsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *contractorsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Contractors, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetContractorsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Contractors)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Contractors)
	for _, id := range ids {
		val, ok := itemMap[c.GetContractorsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *contractorsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetContractorsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *contractorsCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetContractorsCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
