package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"git.abanppc.com/farin-project/crud/internal/model"
)

func newDevicesCache() *gotest.Cache {
	record1 := &model.Devices{}
	record1.ID = 1
	record2 := &model.Devices{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewDevicesCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_devicesCache_Set(t *testing.T) {
	c := newDevicesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Devices)
	err := c.ICache.(DevicesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(DevicesCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_devicesCache_Get(t *testing.T) {
	c := newDevicesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Devices)
	err := c.ICache.(DevicesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(DevicesCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(DevicesCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_devicesCache_MultiGet(t *testing.T) {
	c := newDevicesCache()
	defer c.Close()

	var testData []*model.Devices
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Devices))
	}

	err := c.ICache.(DevicesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(DevicesCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Devices))
	}
}

func Test_devicesCache_MultiSet(t *testing.T) {
	c := newDevicesCache()
	defer c.Close()

	var testData []*model.Devices
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Devices))
	}

	err := c.ICache.(DevicesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_devicesCache_Del(t *testing.T) {
	c := newDevicesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Devices)
	err := c.ICache.(DevicesCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_devicesCache_SetCacheWithNotFound(t *testing.T) {
	c := newDevicesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Devices)
	err := c.ICache.(DevicesCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewDevicesCache(t *testing.T) {
	c := NewDevicesCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewDevicesCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewDevicesCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
