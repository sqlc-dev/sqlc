package sqltest

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/kyleconroy/sqlc/internal/dinosql"

	_ "github.com/lib/pq"
)

func id() string {
	bytes := make([]byte, 10)
	for i := 0; i < 10; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) // A=65 and Z = 65+25
	}
	return string(bytes)
}

func PostgreSQL(t *testing.T, migrations string) (*sql.DB, func()) {
	t.Helper()

	pgUser := os.Getenv("PG_USER")
	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgPass := os.Getenv("PG_PASSWORD")
	pgDB := os.Getenv("PG_DATABASE")

	if pgUser == "" {
		pgUser = "postgres"
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

	db, err := sql.Open("postgres", source)
	if err != nil {
		t.Fatal(err)
	}

	schema := "dinotest_" + id()

	// For each test, pick a new schema name at random.
	// `foo` is used here only as an example
	if _, err := db.Exec("CREATE SCHEMA " + schema); err != nil {
		t.Fatal(err)
	}

	sdb, err := sql.Open("postgres", source+"&search_path="+schema)
	if err != nil {
		t.Fatal(err)
	}

	files, err := dinosql.ReadSQLFiles(migrations)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		blob, err := ioutil.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := sdb.Exec(string(blob)); err != nil {
			t.Fatalf("%s: %s", filepath.Base(f), err)
		}
	}

	return sdb, func() {
		if _, err := db.Exec("DROP SCHEMA " + schema + " CASCADE"); err != nil {
			t.Fatal(err)
		}
	}
}
