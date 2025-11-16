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
	contractsCachePrefixKey = "contracts:"
	// ContractsExpireTime expire time
	ContractsExpireTime = 5 * time.Minute
)

var _ ContractsCache = (*contractsCache)(nil)

// ContractsCache cache interface
type ContractsCache interface {
	Set(ctx context.Context, id uint64, data *model.Contracts, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Contracts, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Contracts, error)
	MultiSet(ctx context.Context, data []*model.Contracts, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// contractsCache define a cache struct
type contractsCache struct {
	cache cache.Cache
}

// NewContractsCache new a cache
func NewContractsCache(cacheType *model.CacheType) ContractsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Contracts{}
		})
		return &contractsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Contracts{}
		})
		return &contractsCache{cache: c}
	}

	return nil // no cache
}

// GetContractsCacheKey cache key
func (c *contractsCache) GetContractsCacheKey(id uint64) string {
	return contractsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *contractsCache) Set(ctx context.Context, id uint64, data *model.Contracts, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetContractsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *contractsCache) Get(ctx context.Context, id uint64) (*model.Contracts, error) {
	var data *model.Contracts
	cacheKey := c.GetContractsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *contractsCache) MultiSet(ctx context.Context, data []*model.Contracts, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetContractsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *contractsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Contracts, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetContractsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Contracts)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Contracts)
	for _, id := range ids {
		val, ok := itemMap[c.GetContractsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *contractsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetContractsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *contractsCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetContractsCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
