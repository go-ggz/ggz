package lru

import (
	"context"
	"fmt"

	"github.com/hashicorp/golang-lru"
	"gopkg.in/nicksrandall/dataloader.v4"
)

// Cache implements the dataloader.Cache interface
type Cache struct {
	*lru.ARCCache
}

// Get gets an item from the cache
func (c *Cache) Get(_ context.Context, key interface{}) (dataloader.Thunk, bool) {
	fmt.Println("Get key:", key)
	v, ok := c.ARCCache.Get(key)
	if ok {
		return v.(dataloader.Thunk), ok
	}
	return nil, ok
}

// Set sets an item in the cache
func (c *Cache) Set(_ context.Context, key interface{}, value dataloader.Thunk) {
	fmt.Println("Set key:", key)
	c.ARCCache.Add(key, value)
}

// Delete deletes an item in the cache
func (c *Cache) Delete(_ context.Context, key interface{}) bool {
	if c.ARCCache.Contains(key) {
		c.ARCCache.Remove(key)
		return true
	}
	return false
}

// Clear cleasrs the cache
func (c *Cache) Clear() {
	c.ARCCache.Purge()
}

// NewEngine for lru engine
func NewEngine() *Cache {
	c, _ := lru.NewARC(100)

	return &Cache{c}
}
