package clickhouse

import (
	"context"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/dialect"
)

// TestDialect verifies the committed ClickHouse dialect against what the
// pinned clickhouse binary reports. It skips unless the binary is installed.
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
