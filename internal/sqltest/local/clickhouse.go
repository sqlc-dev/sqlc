package local

import (
	"context"
	"database/sql"
	"fmt"
	"hash/fnv"
	"os"
	"strings"
	"testing"

	migrate "github.com/sqlc-dev/sqlc/internal/migrations"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
	"github.com/sqlc-dev/sqlc/internal/sqltest/docker"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

func ClickHouse(t *testing.T, migrations []string) string {
	ctx := context.Background()
	t.Helper()

	dburi := os.Getenv("CLICKHOUSE_SERVER_URI")
	if dburi == "" {
		if ierr := docker.Installed(); ierr == nil {
			u, err := docker.StartClickHouseServer(ctx)
			if err != nil {
				t.Fatal(err)
			}
			dburi = u
		} else {
			t.Skip("CLICKHOUSE_SERVER_URI is empty")
		}
	}

	// Open connection to ClickHouse
	db, err := sql.Open("clickhouse", dburi)
	if err != nil {
		t.Fatalf("ClickHouse connection failed: %s", err)
	}
	defer db.Close()

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

	// Create unique database name
	name := fmt.Sprintf("sqlc_test_%x", h.Sum(nil))

	// Drop database if it exists (ClickHouse style)
	dropQuery := fmt.Sprintf(`DROP DATABASE IF EXISTS %s`, name)
	if _, err := db.ExecContext(ctx, dropQuery); err != nil {
		t.Logf("could not drop database (may not exist): %s", err)
	}

	// Create new database
	createQuery := fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS %s`, name)
	if _, err := db.ExecContext(ctx, createQuery); err != nil {
		t.Fatalf("failed to create database: %s", err)
	}

	// Execute migration scripts
	dbWithDatabase := fmt.Sprintf("%s?database=%s", dburi, name)
	dbConn, err := sql.Open("clickhouse", dbWithDatabase)
	if err != nil {
		t.Fatalf("ClickHouse connection to new database failed: %s", err)
	}
	defer dbConn.Close()

	for _, q := range seed {
		if len(strings.TrimSpace(q)) == 0 {
			continue
		}
		if _, err := dbConn.ExecContext(ctx, q); err != nil {
			t.Fatalf("migration failed: %s: %s", q, err)
		}
	}

	// Register cleanup
	t.Cleanup(func() {
		if _, err := db.ExecContext(ctx, dropQuery); err != nil {
			t.Logf("failed cleaning up database: %s", err)
		}
	})

	return dbWithDatabase
}
