package duckdb

import (
	"context"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/dialect"
	"github.com/sqlc-dev/sqlc/internal/goldeneye/endtoend"
)

// TestDialect verifies the committed DuckDB dialect against what the DuckDB
// CLI reports. It skips unless a CLI is found.
func TestDialect(t *testing.T) {
	binary, err := Locate()
	if err != nil {
		t.Skip(err)
	}
	ctx := context.Background()
	version, err := Version(ctx, binary)
	if err != nil {
		t.Fatal(err)
	}
	files, err := Generate(ctx, binary)
	if err != nil {
		t.Fatal(err)
	}
	dir, err := dialect.Dir(Engine)
	if err != nil {
		t.Fatal(err)
	}
	report, err := dialect.Check(dir, files)
	if err != nil {
		t.Fatal(err)
	}
	if report != "" {
		t.Errorf("%s does not match what %s reports:\n%s", dir, version, report)
	}
}

// TestAnalyzeCases verifies every DuckDB analyze case under
// internal/endtoend/testdata against what the CLI reports. It skips
// unless a CLI is found.
func TestAnalyzeCases(t *testing.T) {
	binary, err := Locate()
	if err != nil {
		t.Skip(err)
	}
	cases, err := endtoend.Cases(Engine)
	if err != nil {
		t.Fatal(err)
	}
	if len(cases) == 0 {
		t.Fatal("no duckdb analyze cases found")
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			diff, err := Check(context.Background(), binary, c)
			if err != nil {
				t.Fatal(err)
			}
			if diff != "" {
				t.Errorf("%s does not match what DuckDB reports (-committed +duckdb):\n%s", c.Output, diff)
			}
		})
	}
}
