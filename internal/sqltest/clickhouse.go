package sqltest

import (
	"context"
	"database/sql"
	"math/rand"
	"os"
	"testing"

	_ "github.com/ClickHouse/clickhouse-go/v2" // ClickHouse driver

	"github.com/sqlc-dev/sqlc/internal/sqltest/docker"
	"github.com/sqlc-dev/sqlc/internal/sqltest/native"
)

var clickhouseLetterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func clickhouseID() string {
	b := make([]rune, 10)
	for i := range b {
		b[i] = clickhouseLetterRunes[rand.Intn(len(clickhouseLetterRunes))]
	}
	return string(b)
}

func ClickHouse(t *testing.T, migrations []string) *sql.DB {
	ctx := context.Background()

	// Check environment variable first
	uri := os.Getenv("CLICKHOUSE_SERVER_URI")
	var err error

	// Try Docker if no URI provided
	if uri == "" {
		uri, err = docker.StartClickHouseServer(ctx)
		if err != nil {
			t.Log("docker clickhouse startup failed:", err)
		}
	}

	// Try native installation
	if uri == "" {
		uri, err = native.StartClickHouseServer(ctx)
		if err != nil {
			t.Log("native clickhouse startup failed:", err)
		}
	}

	if uri == "" {
		t.Skip("no clickhouse server available")
		return nil
	}

	db, err := sql.Open("clickhouse", uri)
	if err != nil {
		t.Fatalf("connect to clickhouse: %s", err)
	}

	// Create a unique test database
	dbName := "test_" + clickhouseID()

	// Create the test database
	_, err = db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbName)
	if err != nil {
		db.Close()
		t.Fatalf("create database: %s", err)
	}

	// Switch to the test database by reconnecting
	db.Close()
	testURI := uri
	// Append database name to URI
	if testURI[len(testURI)-1] != '/' {
		testURI = testURI[:len(testURI)-len("default")] + dbName
	}

	db, err = sql.Open("clickhouse", testURI)
	if err != nil {
		t.Fatalf("connect to test database: %s", err)
	}

	// Apply migrations
	for _, migration := range migrations {
		if _, err := db.ExecContext(ctx, migration); err != nil {
			db.Close()
			t.Fatalf("migration failed: %s: %s", migration, err)
		}
	}

	// Cleanup on test completion
	t.Cleanup(func() {
		db.Close()
		// Reconnect to default database to drop test database
		cleanupDB, err := sql.Open("clickhouse", uri)
		if err == nil {
			cleanupDB.ExecContext(ctx, "DROP DATABASE IF EXISTS "+dbName)
			cleanupDB.Close()
		}
	})

	return db
}
