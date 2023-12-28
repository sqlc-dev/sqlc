package sqltest

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
)

func MySQL(t *testing.T, migrations []string) (*sql.DB, func()) {
	// For each test, pick a new database name at random.
	name := "sqltest_mysql_" + id()
	return CreateMySQLDatabase(t, name, migrations)
}

func CreateMySQLDatabase(t *testing.T, name string, migrations []string) (*sql.DB, func()) {
	t.Helper()

	data := os.Getenv("MYSQL_DATABASE")
	host := os.Getenv("MYSQL_HOST")
	pass := os.Getenv("MYSQL_ROOT_PASSWORD")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")

	if user == "" {
		user = "root"
	}

	if pass == "" {
		pass = "mysecretpassword"
	}

	if port == "" {
		port = "3306"
	}

	if host == "" {
		host = "127.0.0.1"
	}

	if data == "" {
		data = "dinotest"
	}

	source := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true", user, pass, host, port, data)
	t.Logf("db: %s", source)

	db, err := sql.Open("mysql", source)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := db.Exec("CREATE DATABASE " + name); err != nil {
		t.Fatal(err)
	}

	source = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true", user, pass, host, port, name)
	sdb, err := sql.Open("mysql", source)
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
		// Drop the test db after test runs
		if _, err := db.Exec("DROP DATABASE " + name); err != nil {
			t.Fatal(err)
		}
	}
}
