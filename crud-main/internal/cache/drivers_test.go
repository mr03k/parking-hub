package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"git.abanppc.com/farin-project/crud/internal/model"
)

func newDriversCache() *gotest.Cache {
	record1 := &model.Drivers{}
	record1.ID = 1
	record2 := &model.Drivers{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewDriversCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_driversCache_Set(t *testing.T) {
	c := newDriversCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Drivers)
	err := c.ICache.(DriversCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(DriversCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_driversCache_Get(t *testing.T) {
	c := newDriversCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Drivers)
	err := c.ICache.(DriversCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(DriversCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(DriversCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_driversCache_MultiGet(t *testing.T) {
	c := newDriversCache()
	defer c.Close()

	var testData []*model.Drivers
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Drivers))
	}

	err := c.ICache.(DriversCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(DriversCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Drivers))
	}
}

func Test_driversCache_MultiSet(t *testing.T) {
	c := newDriversCache()
	defer c.Close()

	var testData []*model.Drivers
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Drivers))
	}

	err := c.ICache.(DriversCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_driversCache_Del(t *testing.T) {
	c := newDriversCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Drivers)
	err := c.ICache.(DriversCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_driversCache_SetCacheWithNotFound(t *testing.T) {
	c := newDriversCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Drivers)
	err := c.ICache.(DriversCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewDriversCache(t *testing.T) {
	c := NewDriversCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewDriversCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewDriversCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
