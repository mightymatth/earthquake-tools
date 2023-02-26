package cache

import "sync"

type Cacher interface {
	Get(key interface{}) (interface{}, bool)
	Set(key, value interface{}) bool
}

type InMemory struct {
	m sync.Map
}

func (c *InMemory) Get(key interface{}) (interface{}, bool) {
	return c.m.Load(key)
}

func (c *InMemory) Set(key, value interface{}) bool {
	c.m.Store(key, value)
	return true
}

func NewInMemory() *InMemory {
	return &InMemory{}
}
