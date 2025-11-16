package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"git.abanppc.com/farin-project/crud/internal/model"
)

func newRoadCategoriesCache() *gotest.Cache {
	record1 := &model.RoadCategories{}
	record1.ID = 1
	record2 := &model.RoadCategories{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewRoadCategoriesCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_roadCategoriesCache_Set(t *testing.T) {
	c := newRoadCategoriesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.RoadCategories)
	err := c.ICache.(RoadCategoriesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(RoadCategoriesCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_roadCategoriesCache_Get(t *testing.T) {
	c := newRoadCategoriesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.RoadCategories)
	err := c.ICache.(RoadCategoriesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(RoadCategoriesCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(RoadCategoriesCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_roadCategoriesCache_MultiGet(t *testing.T) {
	c := newRoadCategoriesCache()
	defer c.Close()

	var testData []*model.RoadCategories
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.RoadCategories))
	}

	err := c.ICache.(RoadCategoriesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(RoadCategoriesCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.RoadCategories))
	}
}

func Test_roadCategoriesCache_MultiSet(t *testing.T) {
	c := newRoadCategoriesCache()
	defer c.Close()

	var testData []*model.RoadCategories
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.RoadCategories))
	}

	err := c.ICache.(RoadCategoriesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_roadCategoriesCache_Del(t *testing.T) {
	c := newRoadCategoriesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.RoadCategories)
	err := c.ICache.(RoadCategoriesCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_roadCategoriesCache_SetCacheWithNotFound(t *testing.T) {
	c := newRoadCategoriesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.RoadCategories)
	err := c.ICache.(RoadCategoriesCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewRoadCategoriesCache(t *testing.T) {
	c := NewRoadCategoriesCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewRoadCategoriesCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewRoadCategoriesCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
