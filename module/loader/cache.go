package loader

import (
	"github.com/go-ggz/ggz/module/loader/lru"
	"github.com/go-ggz/ggz/module/loader/memory"

	"gopkg.in/nicksrandall/dataloader.v5"
)

// Cache for dataloader
var Cache dataloader.Cache

// NewEngine for initialize cache engine
func NewEngine(driver, prefix string, expire int) error {
	switch driver {
	case "lru":
		Cache = lru.NewEngine(prefix)
	case "memory":
		Cache = memory.NewEngine(prefix, expire)
	default:
		Cache = dataloader.NewCache()
	}

	return nil
}
