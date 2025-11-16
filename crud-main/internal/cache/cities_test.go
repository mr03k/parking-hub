package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"git.abanppc.com/farin-project/crud/internal/model"
)

func newCitiesCache() *gotest.Cache {
	record1 := &model.Cities{}
	record1.ID = 1
	record2 := &model.Cities{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewCitiesCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_citiesCache_Set(t *testing.T) {
	c := newCitiesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Cities)
	err := c.ICache.(CitiesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(CitiesCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_citiesCache_Get(t *testing.T) {
	c := newCitiesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Cities)
	err := c.ICache.(CitiesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(CitiesCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(CitiesCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_citiesCache_MultiGet(t *testing.T) {
	c := newCitiesCache()
	defer c.Close()

	var testData []*model.Cities
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Cities))
	}

	err := c.ICache.(CitiesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(CitiesCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Cities))
	}
}

func Test_citiesCache_MultiSet(t *testing.T) {
	c := newCitiesCache()
	defer c.Close()

	var testData []*model.Cities
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Cities))
	}

	err := c.ICache.(CitiesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_citiesCache_Del(t *testing.T) {
	c := newCitiesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Cities)
	err := c.ICache.(CitiesCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_citiesCache_SetCacheWithNotFound(t *testing.T) {
	c := newCitiesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Cities)
	err := c.ICache.(CitiesCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewCitiesCache(t *testing.T) {
	c := NewCitiesCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewCitiesCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewCitiesCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
