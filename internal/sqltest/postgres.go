package sqltest

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kyleconroy/sqlc/internal/sql/sqlpath"
	"github.com/ory/dockertest/v3"

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

func PostgreSQL(t *testing.T, migrations []string) *sql.DB {
	t.Helper()

	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")
	database := os.Getenv("PG_DATABASE")

	if user == "" {
		user = "postgres"
	}

	if password == "" {
		password = "mysecretpassword"
	}

	if database == "" {
		database = "dinotest"
	}

	var db *sql.DB

	disabled := os.Getenv("DOCKERTEST_DISABLED")
	switch {
	case disabled != "":
		host := os.Getenv("PG_HOST")
		port := os.Getenv("PG_PORT")

		if host == "" {
			host = "127.0.0.1"
		}

		if port == "" {
			port = "5432"
		}

		source := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, database)

		var err error
		db, err = sql.Open("postgres", source)
		if err != nil {
			t.Fatalf("Could not connect to database: %s", err)
		}

		// For each test, pick a new schema name at random.
		schema := "sqltest_postgresql_" + id()
		if _, err := db.Exec("CREATE SCHEMA " + schema); err != nil {
			t.Fatal(err)
		}

		db, err = sql.Open("postgres", source+"&search_path="+schema)
		if err != nil {
			t.Fatal(err)
		}

		t.Cleanup(func() {
			if _, err := db.Exec("DROP SCHEMA " + schema + " CASCADE"); err != nil {
				t.Fatal(err)
			}
		})
	default:
		pool, err := dockertest.NewPool("")
		if err != nil {
			t.Fatalf("new pool: Could not connect to docker: %s", err)
		}

		resource, err := pool.Run("postgres", "13", []string{
			fmt.Sprintf("POSTGRES_USER=%s", user),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", password),
			fmt.Sprintf("POSTGRES_DB=%s", database),
		})
		if err != nil {
			t.Fatalf("Could not start postgres: %s", err)
		}

		if err := pool.Retry(func() error {
			source := fmt.Sprintf("postgres://%s:%s@:%s/%s?sslmode=disable", user, password, resource.GetPort("5432/tcp"), database)

			var err error
			db, err = sql.Open("postgres", source)
			if err != nil {
				return err
			}

			return db.Ping()
		}); err != nil {
			t.Fatalf("Could not connect to database: %s", err)
		}

		t.Cleanup(func() {
			if err := pool.Purge(resource); err != nil {
				t.Fatalf("Could not purge resource: %s", err)
			}
		})
	}

	files, err := sqlpath.Glob(migrations)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		blob, err := ioutil.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := db.Exec(string(blob)); err != nil {
			t.Fatalf("%s: %s", filepath.Base(f), err)
		}
	}

	return db
}
