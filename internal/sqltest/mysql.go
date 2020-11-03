package sqltest

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/kyleconroy/sqlc/internal/sql/sqlpath"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	_ "github.com/go-sql-driver/mysql"
)

func MySQL(t *testing.T, migrations []string) *sql.DB {
	t.Helper()

	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_ROOT_PASSWORD")
	database := os.Getenv("MYSQL_DATABASE")

	if user == "" {
		user = "root"
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
		host := os.Getenv("MYSQL_HOST")
		port := os.Getenv("MYSQL_PORT")

		if host == "" {
			host = "127.0.0.1"
		}

		if port == "" {
			port = "3306"
		}

		source := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true", user, password, host, port, database)
		log.Printf("db: %s", source)

		var err error
		db, err = sql.Open("mysql", source)
		if err != nil {
			log.Fatal(err)
		}

		// For each test, pick a new database name at random.
		tempDatabase := "sqltest_mysql_" + id()
		if _, err := db.Exec("CREATE DATABASE " + tempDatabase); err != nil {
			t.Fatal(err)
		}

		if _, err := db.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON %s.* TO '%s'@'%%';", database, user)); err != nil {
			t.Fatal(err)
		}

		source = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true", user, password, host, port, tempDatabase)
		db, err = sql.Open("mysql", source)
		if err != nil {
			t.Fatal(err)
		}

		t.Cleanup(func() {
			// Drop the test db after test runs
			if _, err := db.Exec("DROP DATABASE " + tempDatabase); err != nil {
				t.Fatal(err)
			}
		})
	default:
		pool, err := dockertest.NewPool("")
		if err != nil {
			t.Fatalf("new pool: Could not connect to docker: %s", err)
		}

		resource, err := pool.RunWithOptions(&dockertest.RunOptions{
			Name:       containerName(t, "mysql"),
			Repository: "mysql",
			Tag:        "8",
			Env: []string{
				fmt.Sprintf("MYSQL_USER=%s", user),
				fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", password),
				fmt.Sprintf("MYSQL_DATABASE=%s", database),
				"MYSQL_INITDB_SKIP_TZINFO=yes",
			},
		}, func(c *docker.HostConfig) {
			c.Tmpfs = map[string]string{
				"/var/lib/mysql": "rw,exec",
			}
		})
		if err != nil {
			t.Fatalf("Could not start mysql: %s", err)
		}

		if err := pool.Retry(func() error {
			source := fmt.Sprintf("%s:%s@tcp(localhost:%s)/%s?multiStatements=true&parseTime=true", user, password, resource.GetPort("3306/tcp"), database)

			var err error
			db, err = sql.Open("mysql", source)
			if err != nil {
				return err
			}

			return db.Ping()
		}); err != nil {
			t.Fatalf("Could not connect to database: %s", err)
		}

		retain := os.Getenv("DOCKERTEST_RETAIN")
		if retain != "" {
			break
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
