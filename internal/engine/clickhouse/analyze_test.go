package clickhouse_test

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/core/analyzer"
	"github.com/sqlc-dev/sqlc/internal/engine/clickhouse"
)

func analyzeOne(t *testing.T, ddl, query string) core.PrepareResult {
	t.Helper()
	cat, err := core.New(clickhouse.Dialect())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { cat.Close() })

	if err := clickhouse.LoadSchema(cat, ddl); err != nil {
		t.Fatalf("load schema: %v", err)
	}

	stmts, err := clickhouse.NewParser().Parse(strings.NewReader(query))
	if err != nil {
		t.Fatalf("parse query: %v", err)
	}
	if len(stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(stmts))
	}
	res, err := analyzer.Prepare(cat, stmts[0].Raw)
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}
	return res
}

func colByName(res core.PrepareResult, name string) (core.Column, bool) {
	for _, c := range res.Columns {
		if c.Name == name {
			return c, true
		}
	}
	return core.Column{}, false
}

const eventsDDL = `
CREATE TABLE events (
	id UInt64,
	name String,
	tag Nullable(String),
	amount Decimal(18, 4)
) ENGINE = MergeTree ORDER BY id
`

func TestClickHouseSelectColumns(t *testing.T) {
	res := analyzeOne(t, eventsDDL, `SELECT id, name, tag FROM events`)

	if len(res.Columns) != 3 {
		t.Fatalf("got %d cols, want 3: %+v", len(res.Columns), res.Columns)
	}

	id, ok := colByName(res, "id")
	if !ok || id.DataType != "uint64" || !id.NotNull {
		t.Errorf("id: %+v (ok=%v)", id, ok)
	}
	name, ok := colByName(res, "name")
	if !ok || name.DataType != "string" || !name.NotNull {
		t.Errorf("name: %+v (ok=%v)", name, ok)
	}
	tag, ok := colByName(res, "tag")
	if !ok || tag.DataType != "string" || tag.NotNull {
		t.Errorf("tag: want string/nullable, got %+v (ok=%v)", tag, ok)
	}

	for _, c := range res.Columns {
		if c.SourceClassOID == 0 || c.SourceAttributeOID == 0 {
			t.Errorf("col %s missing source binding: %+v", c.Name, c)
		}
	}
}

func TestClickHouseSelectStar(t *testing.T) {
	res := analyzeOne(t, eventsDDL, `SELECT * FROM events`)
	if len(res.Columns) != 4 {
		t.Fatalf("got %d cols, want 4: %+v", len(res.Columns), res.Columns)
	}
	want := []string{"id", "name", "tag", "amount"}
	for i, w := range want {
		if res.Columns[i].Name != w {
			t.Errorf("col %d: got %q, want %q", i, res.Columns[i].Name, w)
		}
	}
}
