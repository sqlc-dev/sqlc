package sqltest

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/kyleconroy/sqlc/internal/sql/sqlpath"

	_ "github.com/mattn/go-sqlite3"
)

func SQLite(t *testing.T, migrations []string) (*sql.DB, func()) {
	t.Helper()

	// For each test, pick a new database name at random.
	source, err := ioutil.TempFile("", "sqltest_sqlite_")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("open %s\n", source.Name())
	sdb, err := sql.Open("sqlite3", source.Name())
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
		os.Remove(source.Name())
	}
}
