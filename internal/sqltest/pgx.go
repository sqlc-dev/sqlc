package sqltest

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"

	pgxv4 "github.com/jackc/pgx/v4"
	pgxv5 "github.com/jackc/pgx/v5"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func PostgreSQLPgxV4(t *testing.T, migrations []string) (*pgxv4.Conn, func()) {
	t.Helper()

	ctx := context.Background()

	pgUser := os.Getenv("PG_USER")
	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgPass := os.Getenv("PG_PASSWORD")
	pgDB := os.Getenv("PG_DATABASE")

	if pgUser == "" {
		pgUser = "postgres"
	}

	if pgPass == "" {
		pgPass = "mysecretpassword"
	}

	if pgPort == "" {
		pgPort = "5432"
	}

	if pgHost == "" {
		pgHost = "127.0.0.1"
	}

	if pgDB == "" {
		pgDB = "dinotest"
	}

	source := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", pgUser, pgPass, pgHost, pgPort, pgDB)
	t.Logf("db: %s", source)

	db, err := pgxv4.Connect(ctx, source)
	if err != nil {
		t.Fatal(err)
	}

	// For each test, pick a new schema name at random.
	schema := "sqltest_postgresql_" + id()
	if _, err := db.Exec(ctx, "CREATE SCHEMA "+schema); err != nil {
		t.Fatal(err)
	}

	sdb, err := pgxv4.Connect(ctx, source+"&search_path="+schema)
	if err != nil {
		t.Fatal(err)
	}

	files, err := sqlpath.Glob(migrations)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		blob, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := sdb.Exec(ctx, string(blob)); err != nil {
			t.Fatalf("%s: %s", filepath.Base(f), err)
		}
	}

	return sdb, func() {
		if _, err := db.Exec(ctx, "DROP SCHEMA "+schema+" CASCADE"); err != nil {
			t.Fatal(err)
		}
	}
}

func PostgreSQLPgxV5(t *testing.T, migrations []string) (*pgxv5.Conn, func()) {
	t.Helper()

	pgUser := os.Getenv("PG_USER")
	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgPass := os.Getenv("PG_PASSWORD")
	pgDB := os.Getenv("PG_DATABASE")

	if pgUser == "" {
		pgUser = "postgres"
	}

	if pgPass == "" {
		pgPass = "mysecretpassword"
	}

	if pgPort == "" {
		pgPort = "5432"
	}

	if pgHost == "" {
		pgHost = "127.0.0.1"
	}

	if pgDB == "" {
		pgDB = "dinotest"
	}

	source := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", pgUser, pgPass, pgHost, pgPort, pgDB)
	t.Logf("db: %s", source)

	db, err := pgxv5.Connect(context.Background(), source)
	if err != nil {
		t.Fatal(err)
	}

	// For each test, pick a new schema name at random.
	schema := "sqltest_postgresql_" + id()
	if _, err := db.Exec(context.Background(), "CREATE SCHEMA "+schema); err != nil {
		t.Fatal(err)
	}

	sdb, err := pgxv5.Connect(context.Background(), source+"&search_path="+schema)
	if err != nil {
		t.Fatal(err)
	}

	files, err := sqlpath.Glob(migrations)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		blob, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := sdb.Exec(context.Background(), string(blob)); err != nil {
			t.Fatalf("%s: %s", filepath.Base(f), err)
		}
	}

	return sdb, func() {
		if _, err := db.Exec(context.Background(), "DROP SCHEMA "+schema+" CASCADE"); err != nil {
			t.Fatal(err)
		}
	}
}
