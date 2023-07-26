package sqltest

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"

	_ "github.com/lib/pq"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func id() string {
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func PostgreSQL(t *testing.T, migrations []string) (*sql.DB, func()) {
	t.Helper()

	// For each test, pick a new schema name at random.
	schema := "sqltest_postgresql_" + id()
	return CreatePostgreSQLDatabase(t, schema, true, migrations)
}

func CreatePostgreSQLDatabase(t *testing.T, name string, schema bool, migrations []string) (*sql.DB, func()) {
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

	db, err := sql.Open("postgres", source)
	if err != nil {
		t.Fatal(err)
	}

	// For each test, pick a new schema name at random.
	var newsource, dropQuery string
	if schema {
		if _, err := db.Exec("CREATE SCHEMA " + name); err != nil {
			t.Fatal(err)
		}
		newsource = source + "&search_path=" + name
		dropQuery = "DROP SCHEMA " + name + " CASCADE"
	} else {
		if _, err := db.Exec("CREATE DATABASE " + name); err != nil {
			t.Fatal(err)
		}
		newsource = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", pgUser, pgPass, pgHost, pgPort, name)
		dropQuery = "DROP DATABASE IF EXISTS " + name + " WITH (FORCE)"
	}

	sdb, err := sql.Open("postgres", newsource)
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
		if _, err := db.Exec(dropQuery); err != nil {
			t.Fatal(err)
		}
		db.Close()
	}
}
