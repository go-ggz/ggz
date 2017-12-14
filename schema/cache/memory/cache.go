package memory

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
	"gopkg.in/nicksrandall/dataloader.v4"
)

// Cache implements the dataloader.Cache interface
type Cache struct {
	c *cache.Cache
}

// Get gets a value from the cache
func (c *Cache) Get(_ context.Context, key interface{}) (dataloader.Thunk, bool) {
	v, ok := c.c.Get(key.(string))
	if ok {
		return v.(dataloader.Thunk), ok
	}
	return nil, ok
}

// Set sets a value in the cache
func (c *Cache) Set(_ context.Context, key interface{}, value dataloader.Thunk) {
	c.c.Set(key.(string), value, 0)
}

// Delete deletes and item in the cache
func (c *Cache) Delete(_ context.Context, key interface{}) bool {
	if _, found := c.c.Get(key.(string)); found {
		c.c.Delete(key.(string))
		return true
	}
	return false
}

// Clear clears the cache
func (c *Cache) Clear() {
	c.c.Flush()
}

// NewEngine for memory engine
func NewEngine() *Cache {
	c := cache.New(15*time.Minute, 15*time.Minute)

	return &Cache{c}
}
