package postgresql

import (
	"context"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/dialect"
)

// TestDialect verifies the committed PostgreSQL dialect against what the
// server reports. It skips unless POSTGRESQL_SERVER_URI names a server.
func TestDialect(t *testing.T) {
	url, err := Locate()
	if err != nil {
		t.Skip(err)
	}
	ctx := context.Background()
	version, err := Version(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	files, err := Generate(ctx, url)
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
