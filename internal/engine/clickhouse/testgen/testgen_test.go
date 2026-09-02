package main

import (
	"bytes"
	"context"
	"flag"
	"os"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "rewrite the expected stdout.txt of every case")

// TestAnalyze runs the CLI over each directory under testdata, which holds
// the same files a sqlc analyze case does plus a fixture, and compares the
// output with the committed stdout.txt. It needs the clickhouse binary and
// skips when none is installed.
func TestAnalyze(t *testing.T) {
	if _, err := Locate(); err != nil {
		t.Skip(err)
	}
	dirs, err := filepath.Glob("testdata/*")
	if err != nil {
		t.Fatal(err)
	}
	for _, dir := range dirs {
		t.Run(filepath.Base(dir), func(t *testing.T) {
			args := []string{"analyze", "--schema", filepath.Join(dir, "schema.sql")}
			if _, err := os.Stat(filepath.Join(dir, "fixture.sql")); err == nil {
				args = append(args, "--fixture", filepath.Join(dir, "fixture.sql"))
			}
			args = append(args, filepath.Join(dir, "query.sql"))

			var stdout, stderr bytes.Buffer
			if err := run(context.Background(), args, &stdout, &stderr); err != nil {
				t.Fatalf("%v\n%s", err, stderr.String())
			}

			golden := filepath.Join(dir, "stdout.txt")
			if *update {
				if err := os.WriteFile(golden, stdout.Bytes(), 0o644); err != nil {
					t.Fatal(err)
				}
				return
			}
			want, err := os.ReadFile(golden)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(want, stdout.Bytes()) {
				t.Errorf("output differs from %s (run with -update to rewrite)\n--- want\n%s\n--- got\n%s", golden, want, stdout.Bytes())
			}
		})
	}
}
