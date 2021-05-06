package cache_test

import (
	"github.com/mightymatth/earthquake-tools/eq-aggregator/cache"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInMemory_SetGet(t *testing.T) {
	c := cache.NewInMemory()
	key := "key"
	value := "value"

	set := c.Set(key, value)
	assert.True(t, set)

	v, found := c.Get(key)
	assert.True(t, found)
	assert.Equal(t, value, v)
}
