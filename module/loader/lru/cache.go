package lru

import (
	"context"

	"github.com/hashicorp/golang-lru"
	"gopkg.in/nicksrandall/dataloader.v5"
)

// Cache implements the dataloader.Cache interface
type Cache struct {
	*lru.ARCCache
	Prefix string
}

// Get gets an item from the cache
func (c *Cache) Get(_ context.Context, key dataloader.Key) (dataloader.Thunk, bool) {
	v, ok := c.ARCCache.Get(c.Prefix + "::" + key.String())
	if ok {
		return v.(dataloader.Thunk), ok
	}
	return nil, ok
}

// Set sets an item in the cache
func (c *Cache) Set(_ context.Context, key dataloader.Key, value dataloader.Thunk) {
	c.ARCCache.Add(c.Prefix+"::"+key.String(), value)
}

// Delete deletes an item in the cache
func (c *Cache) Delete(_ context.Context, key dataloader.Key) bool {
	if c.ARCCache.Contains(c.Prefix + "::" + key.String()) {
		c.ARCCache.Remove(c.Prefix + "::" + key.String())
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
