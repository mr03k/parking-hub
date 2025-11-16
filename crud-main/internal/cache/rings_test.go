package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"git.abanppc.com/farin-project/crud/internal/model"
)

func newRingsCache() *gotest.Cache {
	record1 := &model.Rings{}
	record1.ID = 1
	record2 := &model.Rings{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewRingsCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_ringsCache_Set(t *testing.T) {
	c := newRingsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Rings)
	err := c.ICache.(RingsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(RingsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_ringsCache_Get(t *testing.T) {
	c := newRingsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Rings)
	err := c.ICache.(RingsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(RingsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(RingsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_ringsCache_MultiGet(t *testing.T) {
	c := newRingsCache()
	defer c.Close()

	var testData []*model.Rings
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Rings))
	}

	err := c.ICache.(RingsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(RingsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Rings))
	}
}

func Test_ringsCache_MultiSet(t *testing.T) {
	c := newRingsCache()
	defer c.Close()

	var testData []*model.Rings
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Rings))
	}

	err := c.ICache.(RingsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_ringsCache_Del(t *testing.T) {
	c := newRingsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Rings)
	err := c.ICache.(RingsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_ringsCache_SetCacheWithNotFound(t *testing.T) {
	c := newRingsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Rings)
	err := c.ICache.(RingsCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewRingsCache(t *testing.T) {
	c := NewRingsCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewRingsCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewRingsCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
