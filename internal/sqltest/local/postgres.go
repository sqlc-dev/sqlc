package local

import (
	"context"
	"fmt"
	"hash/fnv"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"golang.org/x/sync/singleflight"

	migrate "github.com/sqlc-dev/sqlc/internal/migrations"
	"github.com/sqlc-dev/sqlc/internal/pgx/poolcache"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
	"github.com/sqlc-dev/sqlc/internal/sqltest/docker"
)

var flight singleflight.Group
var cache = poolcache.New()

func PostgreSQL(t *testing.T, migrations []string) string {
	return postgreSQL(t, migrations, true)
}

func ReadOnlyPostgreSQL(t *testing.T, migrations []string) string {
	return postgreSQL(t, migrations, false)
}

func postgreSQL(t *testing.T, migrations []string, rw bool) string {
	ctx := context.Background()
	t.Helper()

	dburi := os.Getenv("POSTGRESQL_SERVER_URI")
	if dburi == "" {
		if ierr := docker.Installed(); ierr == nil {
			u, err := docker.StartPostgreSQLServer(ctx)
			if err != nil {
				t.Fatal(err)
			}
			dburi = u
		} else {
			t.Skip("POSTGRESQL_SERVER_URI is empty")
		}
	}

	postgresPool, err := cache.Open(ctx, dburi)
	if err != nil {
		t.Fatalf("PostgreSQL pool creation failed: %s", err)
	}

	var seed []string
	files, err := sqlpath.Glob(migrations)
	if err != nil {
		t.Fatal(err)
	}

	h := fnv.New64()
	for _, f := range files {
		blob, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		h.Write(blob)
		seed = append(seed, migrate.RemoveRollbackStatements(string(blob)))
	}

	var name string
	if rw {
		name = fmt.Sprintf("sqlc_test_%s", id())
	} else {
		name = fmt.Sprintf("sqlc_test_%x", h.Sum(nil))
	}

	uri, err := url.Parse(dburi)
	if err != nil {
		t.Fatal(err)
	}
	uri.Path = name
	dropQuery := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s" WITH (FORCE)`, name)

	key := uri.String()

	_, err, _ = flight.Do(key, func() (interface{}, error) {
		row := postgresPool.QueryRow(ctx,
			fmt.Sprintf(`SELECT datname FROM pg_database WHERE datname = '%s'`, name))

		var datname string
		if err := row.Scan(&datname); err == nil {
			t.Logf("database exists: %s", name)
			return nil, nil
		}

		t.Logf("creating database: %s", name)
		if _, err := postgresPool.Exec(ctx, fmt.Sprintf(`CREATE DATABASE "%s"`, name)); err != nil {
			return nil, err
		}

		conn, err := pgx.Connect(ctx, uri.String())
		if err != nil {
			return nil, fmt.Errorf("connect %s: %s", name, err)
		}
		defer conn.Close(ctx)

		for _, q := range seed {
			if len(strings.TrimSpace(q)) == 0 {
				continue
			}
			if _, err := conn.Exec(ctx, q); err != nil {
				return nil, fmt.Errorf("%s: %s", q, err)
			}
		}
		return nil, nil
	})
	if rw || err != nil {
		t.Cleanup(func() {
			if _, err := postgresPool.Exec(ctx, dropQuery); err != nil {
				t.Fatalf("failed cleaning up: %s", err)
			}
		})
	}
	if err != nil {
		t.Fatalf("create db: %s", err)
	}
	return key
}
