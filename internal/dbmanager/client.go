package dbmanager

import (
	"context"
	"fmt"
	"hash/fnv"
	"io"
	"net/url"
	"strings"

	"github.com/jackc/pgx/v5"
	"golang.org/x/sync/singleflight"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/pgx/poolcache"
	"github.com/sqlc-dev/sqlc/internal/shfmt"
)

type CreateDatabaseRequest struct {
	Engine     string
	Migrations []string
	Prefix     string
}

type CreateDatabaseResponse struct {
	Uri string
}

type Client interface {
	CreateDatabase(context.Context, *CreateDatabaseRequest) (*CreateDatabaseResponse, error)
	Close(context.Context)
}

var flight singleflight.Group

type ManagedClient struct {
	cache    *poolcache.Cache
	replacer *shfmt.Replacer
	servers  []config.Server
}

func dbid(migrations []string) string {
	h := fnv.New64()
	for _, query := range migrations {
		io.WriteString(h, query)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (m *ManagedClient) CreateDatabase(ctx context.Context, req *CreateDatabaseRequest) (*CreateDatabaseResponse, error) {
	hash := dbid(req.Migrations)
	prefix := req.Prefix
	if prefix == "" {
		prefix = "sqlc_managed"
	}
	name := fmt.Sprintf("%s_%s", prefix, hash)

	engine := config.Engine(req.Engine)
	switch engine {
	case config.EngineMySQL:
		// pass
	case config.EnginePostgreSQL:
		// pass
	default:
		return nil, fmt.Errorf("unsupported engine: %s", engine)
	}

	var base string
	for _, server := range m.servers {
		if server.Engine == engine {
			base = server.URI
			break
		}
	}

	if strings.TrimSpace(base) == "" {
		return nil, fmt.Errorf("no PostgreSQL database server found")
	}

	serverUri := m.replacer.Replace(base)
	pool, err := m.cache.Open(ctx, serverUri)
	if err != nil {
		return nil, err
	}

	uri, err := url.Parse(serverUri)
	if err != nil {
		return nil, err
	}
	uri.Path = "/" + name

	key := uri.String()
	_, err, _ = flight.Do(key, func() (interface{}, error) {
		// TODO: Use a parameterized query
		row := pool.QueryRow(ctx,
			fmt.Sprintf(`SELECT datname FROM pg_database WHERE datname = '%s'`, name))

		var datname string
		if err := row.Scan(&datname); err == nil {
			return nil, nil
		}

		if _, err := pool.Exec(ctx, fmt.Sprintf(`CREATE DATABASE "%s"`, name)); err != nil {
			return nil, err
		}

		conn, err := pgx.Connect(ctx, uri.String())
		if err != nil {
			pool.Exec(ctx, fmt.Sprintf(`DROP DATABASE IF EXISTS "%s" WITH (FORCE)`, name))
			return nil, fmt.Errorf("connect %s: %s", name, err)
		}
		defer conn.Close(ctx)

		var migrationErr error
		for _, q := range req.Migrations {
			if len(strings.TrimSpace(q)) == 0 {
				continue
			}
			if _, err := conn.Exec(ctx, q); err != nil {
				migrationErr = fmt.Errorf("%s: %s", q, err)
				break
			}
		}

		if migrationErr != nil {
			pool.Exec(ctx, fmt.Sprintf(`DROP DATABASE IF EXISTS "%s" WITH (FORCE)`, name))
			return nil, migrationErr
		}

		return nil, nil
	})

	if err != nil {
		return nil, err
	}

	return &CreateDatabaseResponse{Uri: key}, err
}

func (m *ManagedClient) Close(ctx context.Context) {
	m.cache.Close()
}

func NewClient(servers []config.Server) *ManagedClient {
	return &ManagedClient{
		cache:    poolcache.New(),
		servers:  servers,
		replacer: shfmt.NewReplacer(nil),
	}
}
