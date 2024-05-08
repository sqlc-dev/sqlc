package poolcache

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Cache struct {
	lock  sync.RWMutex
	pools map[string]*pgxpool.Pool
}

func New() *Cache {
	return &Cache{
		pools: map[string]*pgxpool.Pool{},
	}
}

// Should only be used in testing contexts
func (c *Cache) Open(ctx context.Context, uri string) (*pgxpool.Pool, error) {
	c.lock.RLock()
	existing, found := c.pools[uri]
	c.lock.RUnlock()

	if found {
		return existing, nil
	}

	pool, err := pgxpool.New(ctx, uri)
	if err != nil {
		return nil, err
	}

	c.lock.Lock()
	c.pools[uri] = pool
	c.lock.Unlock()

	return pool, nil
}
