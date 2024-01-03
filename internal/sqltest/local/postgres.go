package local

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	migrate "github.com/sqlc-dev/sqlc/internal/migrations"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
)

var postgresPool *pgxpool.Pool
var postgresSync sync.Once

func PostgreSQL(t *testing.T, migrations []string) string {
	ctx := context.Background()
	t.Helper()

	dburi := os.Getenv("POSTGRESQL_SERVER_URI")
	if dburi == "" {
		t.Skip("POSTGRESQL_SERVER_URI is empty")
	}

	postgresSync.Do(func() {
		pool, err := pgxpool.New(ctx, dburi)
		if err != nil {
			t.Fatal(err)
		}
		postgresPool = pool
	})

	if postgresPool == nil {
		t.Fatalf("PostgreSQL pool creation failed")
	}

	var seed []string
	files, err := sqlpath.Glob(migrations)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		blob, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		seed = append(seed, migrate.RemoveRollbackStatements(string(blob)))
	}

	uri, err := url.Parse(dburi)
	if err != nil {
		t.Fatal(err)
	}

	name := fmt.Sprintf("sqlc_test_%s", id())

	if _, err := postgresPool.Exec(ctx, fmt.Sprintf(`CREATE DATABASE "%s"`, name)); err != nil {
		t.Fatal(err)
	}

	uri.Path = name
	dropQuery := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s" WITH (FORCE)`, name)

	t.Cleanup(func() {
		if _, err := postgresPool.Exec(ctx, dropQuery); err != nil {
			t.Fatal(err)
		}
	})

	conn, err := pgx.Connect(ctx, uri.String())
	if err != nil {
		t.Fatalf("connect %s: %s", name, err)
	}
	defer conn.Close(ctx)

	for _, q := range seed {
		if len(strings.TrimSpace(q)) == 0 {
			continue
		}
		if _, err := conn.Exec(ctx, q); err != nil {
			t.Fatalf("%s: %s", q, err)
		}
	}

	return uri.String()
}
