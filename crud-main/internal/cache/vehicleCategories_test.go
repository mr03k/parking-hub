package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"git.abanppc.com/farin-project/crud/internal/model"
)

func newVehicleCategoriesCache() *gotest.Cache {
	record1 := &model.VehicleCategories{}
	record1.ID = 1
	record2 := &model.VehicleCategories{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewVehicleCategoriesCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_vehicleCategoriesCache_Set(t *testing.T) {
	c := newVehicleCategoriesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.VehicleCategories)
	err := c.ICache.(VehicleCategoriesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(VehicleCategoriesCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_vehicleCategoriesCache_Get(t *testing.T) {
	c := newVehicleCategoriesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.VehicleCategories)
	err := c.ICache.(VehicleCategoriesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(VehicleCategoriesCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(VehicleCategoriesCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_vehicleCategoriesCache_MultiGet(t *testing.T) {
	c := newVehicleCategoriesCache()
	defer c.Close()

	var testData []*model.VehicleCategories
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.VehicleCategories))
	}

	err := c.ICache.(VehicleCategoriesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(VehicleCategoriesCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.VehicleCategories))
	}
}

func Test_vehicleCategoriesCache_MultiSet(t *testing.T) {
	c := newVehicleCategoriesCache()
	defer c.Close()

	var testData []*model.VehicleCategories
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.VehicleCategories))
	}

	err := c.ICache.(VehicleCategoriesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_vehicleCategoriesCache_Del(t *testing.T) {
	c := newVehicleCategoriesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.VehicleCategories)
	err := c.ICache.(VehicleCategoriesCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_vehicleCategoriesCache_SetCacheWithNotFound(t *testing.T) {
	c := newVehicleCategoriesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.VehicleCategories)
	err := c.ICache.(VehicleCategoriesCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewVehicleCategoriesCache(t *testing.T) {
	c := NewVehicleCategoriesCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewVehicleCategoriesCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewVehicleCategoriesCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
