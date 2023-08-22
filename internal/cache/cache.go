package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/allegro/bigcache/v3"
)

var (
	ErrCacheCanNotCreated = fmt.Errorf("cannot create cache")
)

type Cache struct {
	cache *bigcache.BigCache
}

func New(ctx context.Context) (*Cache, error) {
	cache, err := bigcache.New(ctx, bigcache.Config{
		Shards:             1024,
		LifeWindow:         time.Hour * 24 * 7,
		MaxEntriesInWindow: 0,
		MaxEntrySize:       0,
		Verbose:            false,
		HardMaxCacheSize:   8192,
	})
	if err != nil {
		log.Printf("%v: %v", ErrCacheCanNotCreated, err)
		return nil, ErrCacheCanNotCreated
	}

	return &Cache{cache: cache}, nil
}

func (c *Cache) Close() error {
	return c.cache.Close()
}

func (c *Cache) Get(key string) ([]byte, error) {
	return c.cache.Get(key)
}

func (c *Cache) Set(key string, value []byte) error {
	return c.cache.Set(key, value)
}
