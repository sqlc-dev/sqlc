package local

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"

	_ "github.com/ClickHouse/clickhouse-go/v2"

	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
)

var clickhouseSync sync.Once
var clickhousePool *sql.DB

// ClickHouse creates a database on the server CLICKHOUSE_SERVER_URI names,
// such as clickhouse://default:mysecretpassword@127.0.0.1:9000, runs the
// migrations in it
// and returns a URI that connects to it. The test is skipped when no
// server is named.
func ClickHouse(t *testing.T, migrations []string) string {
	ctx := context.Background()
	t.Helper()

	dburi := os.Getenv("CLICKHOUSE_SERVER_URI")
	if dburi == "" {
		t.Skip("CLICKHOUSE_SERVER_URI is empty")
	}

	clickhouseSync.Do(func() {
		db, err := sql.Open("clickhouse", dburi)
		if err != nil {
			t.Fatal(err)
		}
		clickhousePool = db
	})
	if clickhousePool == nil {
		t.Fatalf("ClickHouse pool creation failed")
	}

	seed, err := statements(migrations)
	if err != nil {
		t.Fatal(err)
	}

	name := fmt.Sprintf("sqlc_test_%s", id())
	if _, err := clickhousePool.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE `%s`", name)); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := clickhousePool.ExecContext(ctx, fmt.Sprintf("DROP DATABASE `%s`", name)); err != nil {
			t.Fatalf("failed cleaning up: %s", err)
		}
	})

	uri, err := url.Parse(dburi)
	if err != nil {
		t.Fatal(err)
	}
	uri.Path = name

	db, err := sql.Open("clickhouse", uri.String())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	for _, q := range seed {
		if _, err := db.ExecContext(ctx, q); err != nil {
			t.Fatalf("%s: %s", q, err)
		}
	}
	return uri.String()
}

// statements reads the migrations and splits them into the statements
// they hold, for a driver that runs one statement per call.
func statements(migrations []string) ([]string, error) {
	files, err := sqlpath.Glob(migrations)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, f := range files {
		blob, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		for _, q := range strings.Split(string(blob), ";") {
			if strings.TrimSpace(q) == "" {
				continue
			}
			out = append(out, q)
		}
	}
	return out, nil
}
