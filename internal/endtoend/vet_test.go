//go:build examples
// +build examples

package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/kyleconroy/sqlc/internal/cmd"
	"github.com/kyleconroy/sqlc/internal/sqltest"
)

func findSchema(t *testing.T, path string) string {
	t.Helper()
	schemaFile := filepath.Join(path, "postgresql", "schema.sql")
	if _, err := os.Stat(schemaFile); !os.IsNotExist(err) {
		return schemaFile
	}
	schemaDir := filepath.Join(path, "postgresql", "schema")
	if _, err := os.Stat(schemaDir); !os.IsNotExist(err) {
		return schemaDir
	}
	t.Fatalf("error: can't find schema files in %s", path)
	return ""
}

func TestExamplesVet(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

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
				sqltest.CreatePostgreSQLDatabase(t, tc, []string{
					findSchema(t, path),
				})
			}

			var stderr bytes.Buffer
			err := cmd.Vet(ctx, cmd.Env{}, path, "", &stderr)
			if err != nil {
				t.Fatalf("sqlc vet failed: %s %s", err, stderr.String())
			}
		})
	}
}
