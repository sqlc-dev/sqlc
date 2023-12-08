package sqltest

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
)

func SQLite(t *testing.T, migrations []string) (*sql.DB, func()) {
	t.Helper()
	// For each test, pick a new database name at random.
	source, err := os.CreateTemp("", "sqltest_sqlite_")
	if err != nil {
		t.Fatal(err)
	}
	return CreateSQLiteDatabase(t, source.Name(), migrations)
}

func CreateSQLiteDatabase(t *testing.T, path string, migrations []string) (*sql.DB, func()) {
	t.Helper()

	t.Logf("open %s\n", path)
	sdb, err := sql.Open("sqlite", path)
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
		if _, err := sdb.Exec(string(blob)); err != nil {
			t.Fatalf("%s: %s", filepath.Base(f), err)
		}
	}

	return sdb, func() {
		if _, err := os.Stat(path); err == nil {
			os.Remove(path)
		}
	}
}
