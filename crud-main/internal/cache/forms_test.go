package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"git.abanppc.com/farin-project/crud/internal/model"
)

func newFormsCache() *gotest.Cache {
	record1 := &model.Forms{}
	record1.ID = 1
	record2 := &model.Forms{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewFormsCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_formsCache_Set(t *testing.T) {
	c := newFormsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Forms)
	err := c.ICache.(FormsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(FormsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_formsCache_Get(t *testing.T) {
	c := newFormsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Forms)
	err := c.ICache.(FormsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(FormsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(FormsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_formsCache_MultiGet(t *testing.T) {
	c := newFormsCache()
	defer c.Close()

	var testData []*model.Forms
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Forms))
	}

	err := c.ICache.(FormsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(FormsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Forms))
	}
}

func Test_formsCache_MultiSet(t *testing.T) {
	c := newFormsCache()
	defer c.Close()

	var testData []*model.Forms
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Forms))
	}

	err := c.ICache.(FormsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_formsCache_Del(t *testing.T) {
	c := newFormsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Forms)
	err := c.ICache.(FormsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_formsCache_SetCacheWithNotFound(t *testing.T) {
	c := newFormsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Forms)
	err := c.ICache.(FormsCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewFormsCache(t *testing.T) {
	c := NewFormsCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewFormsCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewFormsCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
