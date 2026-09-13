package mssql

import (
	"context"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/dialect"
	"github.com/sqlc-dev/sqlc/internal/goldeneye/endtoend"
)

// TestDialect verifies the committed SQL Server dialect against what the
// server reports. It skips unless MSSQL_SERVER_URI names a server.
func TestDialect(t *testing.T) {
	dsn, err := Locate()
	if err != nil {
		t.Skip(err)
	}
	ctx := context.Background()
	version, err := Version(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	files, err := Generate(ctx, dsn)
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

// TestAnalyzeCases verifies every SQL Server analyze case under
// internal/endtoend/testdata against what the server reports. It skips
// unless MSSQL_SERVER_URI names a server.
func TestAnalyzeCases(t *testing.T) {
	dsn, err := Locate()
	if err != nil {
		t.Skip(err)
	}
	cases, err := endtoend.Cases(Engine)
	if err != nil {
		t.Fatal(err)
	}
	if len(cases) == 0 {
		t.Fatal("no mssql analyze cases found")
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			diff, err := Check(context.Background(), dsn, c)
			if err != nil {
				t.Fatal(err)
			}
			if diff != "" {
				t.Errorf("%s does not match what SQL Server reports (-committed +mssql):\n%s", c.Output, diff)
			}
		})
	}
}
