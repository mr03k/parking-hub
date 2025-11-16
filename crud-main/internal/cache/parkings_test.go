package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"git.abanppc.com/farin-project/crud/internal/model"
)

func newParkingsCache() *gotest.Cache {
	record1 := &model.Parkings{}
	record1.ID = 1
	record2 := &model.Parkings{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewParkingsCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_parkingsCache_Set(t *testing.T) {
	c := newParkingsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Parkings)
	err := c.ICache.(ParkingsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(ParkingsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_parkingsCache_Get(t *testing.T) {
	c := newParkingsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Parkings)
	err := c.ICache.(ParkingsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(ParkingsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(ParkingsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_parkingsCache_MultiGet(t *testing.T) {
	c := newParkingsCache()
	defer c.Close()

	var testData []*model.Parkings
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Parkings))
	}

	err := c.ICache.(ParkingsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(ParkingsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Parkings))
	}
}

func Test_parkingsCache_MultiSet(t *testing.T) {
	c := newParkingsCache()
	defer c.Close()

	var testData []*model.Parkings
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Parkings))
	}

	err := c.ICache.(ParkingsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_parkingsCache_Del(t *testing.T) {
	c := newParkingsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Parkings)
	err := c.ICache.(ParkingsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_parkingsCache_SetCacheWithNotFound(t *testing.T) {
	c := newParkingsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Parkings)
	err := c.ICache.(ParkingsCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewParkingsCache(t *testing.T) {
	c := NewParkingsCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewParkingsCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewParkingsCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
