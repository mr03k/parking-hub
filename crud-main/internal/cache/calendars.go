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
	calendarsCachePrefixKey = "calendars:"
	// CalendarsExpireTime expire time
	CalendarsExpireTime = 5 * time.Minute
)

var _ CalendarsCache = (*calendarsCache)(nil)

// CalendarsCache cache interface
type CalendarsCache interface {
	Set(ctx context.Context, id uint64, data *model.Calendars, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Calendars, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Calendars, error)
	MultiSet(ctx context.Context, data []*model.Calendars, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// calendarsCache define a cache struct
type calendarsCache struct {
	cache cache.Cache
}

// NewCalendarsCache new a cache
func NewCalendarsCache(cacheType *model.CacheType) CalendarsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Calendars{}
		})
		return &calendarsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Calendars{}
		})
		return &calendarsCache{cache: c}
	}

	return nil // no cache
}

// GetCalendarsCacheKey cache key
func (c *calendarsCache) GetCalendarsCacheKey(id uint64) string {
	return calendarsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *calendarsCache) Set(ctx context.Context, id uint64, data *model.Calendars, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetCalendarsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *calendarsCache) Get(ctx context.Context, id uint64) (*model.Calendars, error) {
	var data *model.Calendars
	cacheKey := c.GetCalendarsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *calendarsCache) MultiSet(ctx context.Context, data []*model.Calendars, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetCalendarsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *calendarsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Calendars, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetCalendarsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Calendars)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Calendars)
	for _, id := range ids {
		val, ok := itemMap[c.GetCalendarsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *calendarsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetCalendarsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *calendarsCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetCalendarsCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
