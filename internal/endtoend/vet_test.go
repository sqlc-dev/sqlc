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
)

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

		// Does this work?
		sqltest.PostgreSQL(t, tc)

		t.Run(tc, func(t *testing.T) {
			t.Parallel()
			path := filepath.Join(examples, tc)
			var stderr bytes.Buffer
			err := cmd.Vet(ctx, cmd.Env{}, path, "", &stderr)
			if err != nil {
				t.Fatalf("sqlc vet failed: %s %s", err, stderr.String())
			}
		})
	}
}
