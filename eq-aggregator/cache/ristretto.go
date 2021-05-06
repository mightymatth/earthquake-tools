package cache

import (
	"github.com/dgraph-io/ristretto"
	"time"
)

type Ristretto struct {
	rc *ristretto.Cache
}

func (c *Ristretto) Get(key interface{}) (interface{}, bool) {
	value, found := c.rc.Get(key)
	if !found {
		// wait for value to pass through buffers
		time.Sleep(10 * time.Millisecond)

		value, found = c.rc.Get(key)
		if !found {
			return nil, false
		}
	}

	return value, true
}

func (c *Ristretto) Set(key, value interface{}) bool {
	return c.rc.Set(key, value, 0)
}

func NewRistretto() (*Ristretto, error) {
	c, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: (512 * 16) * 10, // number of keys to track frequency of (10M).
		MaxCost:     512 * 16,        // [max cost of single item] * [max items]
		BufferItems: 64,              // number of keys per Get buffer.
	})
	if err != nil {
		return nil, err
	}

	return &Ristretto{c}, nil
}
