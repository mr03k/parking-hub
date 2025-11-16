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
	segmentsCachePrefixKey = "segments:"
	// SegmentsExpireTime expire time
	SegmentsExpireTime = 5 * time.Minute
)

var _ SegmentsCache = (*segmentsCache)(nil)

// SegmentsCache cache interface
type SegmentsCache interface {
	Set(ctx context.Context, id uint64, data *model.Segments, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Segments, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Segments, error)
	MultiSet(ctx context.Context, data []*model.Segments, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// segmentsCache define a cache struct
type segmentsCache struct {
	cache cache.Cache
}

// NewSegmentsCache new a cache
func NewSegmentsCache(cacheType *model.CacheType) SegmentsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Segments{}
		})
		return &segmentsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Segments{}
		})
		return &segmentsCache{cache: c}
	}

	return nil // no cache
}

// GetSegmentsCacheKey cache key
func (c *segmentsCache) GetSegmentsCacheKey(id uint64) string {
	return segmentsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *segmentsCache) Set(ctx context.Context, id uint64, data *model.Segments, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetSegmentsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *segmentsCache) Get(ctx context.Context, id uint64) (*model.Segments, error) {
	var data *model.Segments
	cacheKey := c.GetSegmentsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *segmentsCache) MultiSet(ctx context.Context, data []*model.Segments, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetSegmentsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *segmentsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Segments, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetSegmentsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Segments)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Segments)
	for _, id := range ids {
		val, ok := itemMap[c.GetSegmentsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *segmentsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetSegmentsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *segmentsCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetSegmentsCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
