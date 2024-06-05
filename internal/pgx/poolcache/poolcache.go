package poolcache

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Cache struct {
	lock   sync.RWMutex
	closed bool
	pools  map[string]*pgxpool.Pool
}

func New() *Cache {
	return &Cache{
		pools: map[string]*pgxpool.Pool{},
	}
}

func (c *Cache) Open(ctx context.Context, uri string) (*pgxpool.Pool, error) {
	if c.closed {
		return nil, fmt.Errorf("poolcache is closed")
	}

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

func (c *Cache) Close() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	var closeErr error
	for _, pool := range c.pools {
		pool.Close()
	}

	c.closed = true
	clear(c.pools)

	return closeErr
}
