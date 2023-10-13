//go:build examples
// +build examples

package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/cmd"
	"github.com/sqlc-dev/sqlc/internal/sqltest"
)

func findSchema(t *testing.T, path string) (string, bool) {
	schemaFile := filepath.Join(path, "schema.sql")
	if _, err := os.Stat(schemaFile); !os.IsNotExist(err) {
		return schemaFile, true
	}
	schemaDir := filepath.Join(path, "schema")
	if _, err := os.Stat(schemaDir); !os.IsNotExist(err) {
		return schemaDir, true
	}
	return "", false
}

func TestExamplesVet(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	authToken := os.Getenv("SQLC_AUTH_TOKEN")
	if authToken == "" {
		t.Skip("missing auth token")
	}

	examples, err := filepath.Abs(filepath.Join("..", "..", "examples"))
	if err != nil {
		t.Fatal(err)
	}

	files, err := os.ReadDir(examples)
	if err != nil {
		t.Fatal(err)
	}

	for _, replay := range files {
		if !replay.IsDir() {
			continue
		}
		tc := replay.Name()
		t.Run(tc, func(t *testing.T) {
			t.Parallel()
			path := filepath.Join(examples, tc)

			if tc != "kotlin" && tc != "python" {
				if s, found := findSchema(t, filepath.Join(path, "mysql")); found {
					db, cleanup := sqltest.CreateMySQLDatabase(t, tc, []string{s})
					defer db.Close()
					defer cleanup()
				}
				if s, found := findSchema(t, filepath.Join(path, "sqlite")); found {
					dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", tc)
					db, cleanup := sqltest.CreateSQLiteDatabase(t, dsn, []string{s})
					defer db.Close()
					defer cleanup()
				}
			}

			var stderr bytes.Buffer
			opts := &cmd.Options{
				Stderr: &stderr,
				Env:    cmd.Env{},
			}
			err := cmd.Vet(ctx, path, "", opts)
			if err != nil {
				t.Fatalf("sqlc vet failed: %s %s", err, stderr.String())
			}
		})
	}
}
