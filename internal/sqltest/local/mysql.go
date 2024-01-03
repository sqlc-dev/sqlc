package local

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/go-sql-driver/mysql"

	migrate "github.com/sqlc-dev/sqlc/internal/migrations"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
)

var mysqlSync sync.Once
var mysqlPool *sql.DB

func MySQL(t *testing.T, migrations []string) string {
	ctx := context.Background()
	t.Helper()

	dburi := os.Getenv("MYSQL_SERVER_URI")
	if dburi == "" {
		t.Skip("MYSQL_SERVER_URI is empty")
	}

	mysqlSync.Do(func() {
		db, err := sql.Open("mysql", dburi)
		if err != nil {
			t.Fatal(err)
		}
		mysqlPool = db
	})

	if mysqlPool == nil {
		t.Fatalf("MySQL pool creation failed")
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

	cfg, err := mysql.ParseDSN(dburi)
	if err != nil {
		t.Fatal(err)
	}

	name := fmt.Sprintf("sqlc_test_%s", id())

	if _, err := mysqlPool.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE `%s`", name)); err != nil {
		t.Fatal(err)
	}

	cfg.DBName = name

	dropQuery := fmt.Sprintf("DROP DATABASE `%s`", name)

	t.Cleanup(func() {
		if _, err := mysqlPool.ExecContext(ctx, dropQuery); err != nil {
			t.Fatal(err)
		}
	})

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		t.Fatalf("connect %s: %s", name, err)
	}
	defer db.Close()

	for _, q := range seed {
		if len(strings.TrimSpace(q)) == 0 {
			continue
		}
		if _, err := db.ExecContext(ctx, q); err != nil {
			t.Fatalf("%s: %s", q, err)
		}
	}

	return cfg.FormatDSN()
}
