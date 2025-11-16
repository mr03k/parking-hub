package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"git.abanppc.com/farin-project/crud/internal/model"
)

func newPartsCache() *gotest.Cache {
	record1 := &model.Parts{}
	record1.ID = 1
	record2 := &model.Parts{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewPartsCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_partsCache_Set(t *testing.T) {
	c := newPartsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Parts)
	err := c.ICache.(PartsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(PartsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_partsCache_Get(t *testing.T) {
	c := newPartsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Parts)
	err := c.ICache.(PartsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(PartsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(PartsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_partsCache_MultiGet(t *testing.T) {
	c := newPartsCache()
	defer c.Close()

	var testData []*model.Parts
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Parts))
	}

	err := c.ICache.(PartsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(PartsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Parts))
	}
}

func Test_partsCache_MultiSet(t *testing.T) {
	c := newPartsCache()
	defer c.Close()

	var testData []*model.Parts
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Parts))
	}

	err := c.ICache.(PartsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_partsCache_Del(t *testing.T) {
	c := newPartsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Parts)
	err := c.ICache.(PartsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_partsCache_SetCacheWithNotFound(t *testing.T) {
	c := newPartsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Parts)
	err := c.ICache.(PartsCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewPartsCache(t *testing.T) {
	c := NewPartsCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewPartsCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewPartsCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
