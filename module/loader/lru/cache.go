package lru

import (
	"context"

	"github.com/hashicorp/golang-lru"
	"gopkg.in/nicksrandall/dataloader.v4"
)

// Cache implements the dataloader.Cache interface
type Cache struct {
	*lru.ARCCache
	Prefix string
}

// Get gets an item from the cache
func (c *Cache) Get(_ context.Context, key interface{}) (dataloader.Thunk, bool) {
	v, ok := c.ARCCache.Get(c.Prefix + "::" + key.(string))
	if ok {
		return v.(dataloader.Thunk), ok
	}
	return nil, ok
}

// Set sets an item in the cache
func (c *Cache) Set(_ context.Context, key interface{}, value dataloader.Thunk) {
	c.ARCCache.Add(c.Prefix+"::"+key.(string), value)
}

// Delete deletes an item in the cache
func (c *Cache) Delete(_ context.Context, key interface{}) bool {
	if c.ARCCache.Contains(c.Prefix + "::" + key.(string)) {
		c.ARCCache.Remove(c.Prefix + "::" + key.(string))
		return true
	}
	return false
}

// Clear cleasrs the cache
func (c *Cache) Clear() {
	c.ARCCache.Purge()
}

// NewEngine for lru engine
func NewEngine(prefix string) *Cache {
	c, _ := lru.NewARC(100)

	return &Cache{
		ARCCache: c,
		Prefix:   prefix,
	}
}
