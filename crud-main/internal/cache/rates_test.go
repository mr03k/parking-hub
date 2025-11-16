package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"git.abanppc.com/farin-project/crud/internal/model"
)

func newRatesCache() *gotest.Cache {
	record1 := &model.Rates{}
	record1.ID = 1
	record2 := &model.Rates{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewRatesCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_ratesCache_Set(t *testing.T) {
	c := newRatesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Rates)
	err := c.ICache.(RatesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(RatesCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_ratesCache_Get(t *testing.T) {
	c := newRatesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Rates)
	err := c.ICache.(RatesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(RatesCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(RatesCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_ratesCache_MultiGet(t *testing.T) {
	c := newRatesCache()
	defer c.Close()

	var testData []*model.Rates
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Rates))
	}

	err := c.ICache.(RatesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(RatesCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Rates))
	}
}

func Test_ratesCache_MultiSet(t *testing.T) {
	c := newRatesCache()
	defer c.Close()

	var testData []*model.Rates
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Rates))
	}

	err := c.ICache.(RatesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_ratesCache_Del(t *testing.T) {
	c := newRatesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Rates)
	err := c.ICache.(RatesCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_ratesCache_SetCacheWithNotFound(t *testing.T) {
	c := newRatesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Rates)
	err := c.ICache.(RatesCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewRatesCache(t *testing.T) {
	c := NewRatesCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewRatesCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewRatesCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
