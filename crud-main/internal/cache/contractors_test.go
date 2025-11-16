package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"git.abanppc.com/farin-project/crud/internal/model"
)

func newContractorsCache() *gotest.Cache {
	record1 := &model.Contractors{}
	record1.ID = 1
	record2 := &model.Contractors{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewContractorsCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_contractorsCache_Set(t *testing.T) {
	c := newContractorsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Contractors)
	err := c.ICache.(ContractorsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(ContractorsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_contractorsCache_Get(t *testing.T) {
	c := newContractorsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Contractors)
	err := c.ICache.(ContractorsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(ContractorsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(ContractorsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_contractorsCache_MultiGet(t *testing.T) {
	c := newContractorsCache()
	defer c.Close()

	var testData []*model.Contractors
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Contractors))
	}

	err := c.ICache.(ContractorsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(ContractorsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Contractors))
	}
}

func Test_contractorsCache_MultiSet(t *testing.T) {
	c := newContractorsCache()
	defer c.Close()

	var testData []*model.Contractors
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Contractors))
	}

	err := c.ICache.(ContractorsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_contractorsCache_Del(t *testing.T) {
	c := newContractorsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Contractors)
	err := c.ICache.(ContractorsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_contractorsCache_SetCacheWithNotFound(t *testing.T) {
	c := newContractorsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Contractors)
	err := c.ICache.(ContractorsCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewContractorsCache(t *testing.T) {
	c := NewContractorsCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewContractorsCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewContractorsCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
