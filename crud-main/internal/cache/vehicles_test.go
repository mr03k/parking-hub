package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"git.abanppc.com/farin-project/crud/internal/model"
)

func newVehiclesCache() *gotest.Cache {
	record1 := &model.Vehicles{}
	record1.ID = 1
	record2 := &model.Vehicles{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewVehiclesCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_vehiclesCache_Set(t *testing.T) {
	c := newVehiclesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Vehicles)
	err := c.ICache.(VehiclesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(VehiclesCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_vehiclesCache_Get(t *testing.T) {
	c := newVehiclesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Vehicles)
	err := c.ICache.(VehiclesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(VehiclesCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(VehiclesCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_vehiclesCache_MultiGet(t *testing.T) {
	c := newVehiclesCache()
	defer c.Close()

	var testData []*model.Vehicles
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Vehicles))
	}

	err := c.ICache.(VehiclesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(VehiclesCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Vehicles))
	}
}

func Test_vehiclesCache_MultiSet(t *testing.T) {
	c := newVehiclesCache()
	defer c.Close()

	var testData []*model.Vehicles
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Vehicles))
	}

	err := c.ICache.(VehiclesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_vehiclesCache_Del(t *testing.T) {
	c := newVehiclesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Vehicles)
	err := c.ICache.(VehiclesCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_vehiclesCache_SetCacheWithNotFound(t *testing.T) {
	c := newVehiclesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Vehicles)
	err := c.ICache.(VehiclesCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewVehiclesCache(t *testing.T) {
	c := NewVehiclesCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewVehiclesCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewVehiclesCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
