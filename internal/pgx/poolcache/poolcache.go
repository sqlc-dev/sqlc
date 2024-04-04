package poolcache

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var lock sync.RWMutex
var pools = map[string]*pgxpool.Pool{}

func New(ctx context.Context, uri string) (*pgxpool.Pool, error) {
	lock.RLock()
	existing, found := pools[uri]
	lock.RUnlock()

	if found {
		return existing, nil
	}

	pool, err := pgxpool.New(ctx, uri)
	if err != nil {
		return nil, err
	}

	lock.Lock()
	pools[uri] = pool
	lock.Unlock()

	return pool, nil
}
