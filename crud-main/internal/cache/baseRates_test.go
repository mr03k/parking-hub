package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"git.abanppc.com/farin-project/crud/internal/model"
)

func newBaseRatesCache() *gotest.Cache {
	record1 := &model.BaseRates{}
	record1.ID = 1
	record2 := &model.BaseRates{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewBaseRatesCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_baseRatesCache_Set(t *testing.T) {
	c := newBaseRatesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.BaseRates)
	err := c.ICache.(BaseRatesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(BaseRatesCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_baseRatesCache_Get(t *testing.T) {
	c := newBaseRatesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.BaseRates)
	err := c.ICache.(BaseRatesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(BaseRatesCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(BaseRatesCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_baseRatesCache_MultiGet(t *testing.T) {
	c := newBaseRatesCache()
	defer c.Close()

	var testData []*model.BaseRates
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.BaseRates))
	}

	err := c.ICache.(BaseRatesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(BaseRatesCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.BaseRates))
	}
}

func Test_baseRatesCache_MultiSet(t *testing.T) {
	c := newBaseRatesCache()
	defer c.Close()

	var testData []*model.BaseRates
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.BaseRates))
	}

	err := c.ICache.(BaseRatesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_baseRatesCache_Del(t *testing.T) {
	c := newBaseRatesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.BaseRates)
	err := c.ICache.(BaseRatesCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_baseRatesCache_SetCacheWithNotFound(t *testing.T) {
	c := newBaseRatesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.BaseRates)
	err := c.ICache.(BaseRatesCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewBaseRatesCache(t *testing.T) {
	c := NewBaseRatesCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewBaseRatesCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewBaseRatesCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
