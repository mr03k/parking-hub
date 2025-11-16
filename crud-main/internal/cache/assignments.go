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
	assignmentsCachePrefixKey = "assignments:"
	// AssignmentsExpireTime expire time
	AssignmentsExpireTime = 5 * time.Minute
)

var _ AssignmentsCache = (*assignmentsCache)(nil)

// AssignmentsCache cache interface
type AssignmentsCache interface {
	Set(ctx context.Context, id uint64, data *model.Assignments, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Assignments, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Assignments, error)
	MultiSet(ctx context.Context, data []*model.Assignments, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// assignmentsCache define a cache struct
type assignmentsCache struct {
	cache cache.Cache
}

// NewAssignmentsCache new a cache
func NewAssignmentsCache(cacheType *model.CacheType) AssignmentsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Assignments{}
		})
		return &assignmentsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Assignments{}
		})
		return &assignmentsCache{cache: c}
	}

	return nil // no cache
}

// GetAssignmentsCacheKey cache key
func (c *assignmentsCache) GetAssignmentsCacheKey(id uint64) string {
	return assignmentsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *assignmentsCache) Set(ctx context.Context, id uint64, data *model.Assignments, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetAssignmentsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *assignmentsCache) Get(ctx context.Context, id uint64) (*model.Assignments, error) {
	var data *model.Assignments
	cacheKey := c.GetAssignmentsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *assignmentsCache) MultiSet(ctx context.Context, data []*model.Assignments, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetAssignmentsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *assignmentsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Assignments, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetAssignmentsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Assignments)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Assignments)
	for _, id := range ids {
		val, ok := itemMap[c.GetAssignmentsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *assignmentsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetAssignmentsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *assignmentsCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetAssignmentsCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
