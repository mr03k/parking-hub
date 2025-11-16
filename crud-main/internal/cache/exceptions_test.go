package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"git.abanppc.com/farin-project/crud/internal/model"
)

func newExceptionsCache() *gotest.Cache {
	record1 := &model.Exceptions{}
	record1.ID = 1
	record2 := &model.Exceptions{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewExceptionsCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_exceptionsCache_Set(t *testing.T) {
	c := newExceptionsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Exceptions)
	err := c.ICache.(ExceptionsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(ExceptionsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_exceptionsCache_Get(t *testing.T) {
	c := newExceptionsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Exceptions)
	err := c.ICache.(ExceptionsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(ExceptionsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(ExceptionsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_exceptionsCache_MultiGet(t *testing.T) {
	c := newExceptionsCache()
	defer c.Close()

	var testData []*model.Exceptions
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Exceptions))
	}

	err := c.ICache.(ExceptionsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(ExceptionsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Exceptions))
	}
}

func Test_exceptionsCache_MultiSet(t *testing.T) {
	c := newExceptionsCache()
	defer c.Close()

	var testData []*model.Exceptions
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Exceptions))
	}

	err := c.ICache.(ExceptionsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_exceptionsCache_Del(t *testing.T) {
	c := newExceptionsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Exceptions)
	err := c.ICache.(ExceptionsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_exceptionsCache_SetCacheWithNotFound(t *testing.T) {
	c := newExceptionsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Exceptions)
	err := c.ICache.(ExceptionsCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewExceptionsCache(t *testing.T) {
	c := NewExceptionsCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewExceptionsCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewExceptionsCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
