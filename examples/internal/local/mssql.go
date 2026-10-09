package local

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"sync"
	"testing"

	_ "github.com/microsoft/go-mssqldb"

	migrate "github.com/sqlc-dev/sqlc/internal/migrations"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
)

var mssqlSync sync.Once
var mssqlPool *sql.DB

// MSSQL creates a database on the server MSSQL_SERVER_URI names, such as
// sqlserver://sa:Mysecretpassword1!@127.0.0.1:1433?encrypt=disable, runs
// the migrations in it and returns a URI that connects to it. The test is
// skipped when no server is named.
func MSSQL(t *testing.T, migrations []string) string {
	ctx := context.Background()
	t.Helper()

	dburi := os.Getenv("MSSQL_SERVER_URI")
	if dburi == "" {
		t.Skip("MSSQL_SERVER_URI is empty")
	}

	mssqlSync.Do(func() {
		db, err := sql.Open("sqlserver", dburi)
		if err != nil {
			t.Fatal(err)
		}
		mssqlPool = db
	})
	if mssqlPool == nil {
		t.Fatalf("SQL Server pool creation failed")
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

	name := fmt.Sprintf("sqlc_test_%s", id())
	if _, err := mssqlPool.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE [%s]", name)); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		// Sessions still in the database, the test's own included, keep
		// it from being dropped until they are rolled back.
		drop := fmt.Sprintf("ALTER DATABASE [%s] SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE [%s]", name, name)
		if _, err := mssqlPool.ExecContext(ctx, drop); err != nil {
			t.Fatalf("failed cleaning up: %s", err)
		}
	})

	uri, err := url.Parse(dburi)
	if err != nil {
		t.Fatal(err)
	}
	query := uri.Query()
	query.Set("database", name)
	uri.RawQuery = query.Encode()

	db, err := sql.Open("sqlserver", uri.String())
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
