package analyzer_test

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/core/analyzer"
	"github.com/sqlc-dev/sqlc/internal/engine/postgresql"
)

func seedUsers(t *testing.T) *core.Catalog {
	t.Helper()
	cat, err := core.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { cat.Close() })

	ns, err := cat.NamespaceOID("public")
	if err != nil {
		t.Fatal(err)
	}
	int4, err := cat.CreateType("int4", 4)
	if err != nil {
		t.Fatal(err)
	}
	text, err := cat.CreateType("text", -1)
	if err != nil {
		t.Fatal(err)
	}
	users, err := cat.CreateClass(ns, "users", "r")
	if err != nil {
		t.Fatal(err)
	}
	if err := cat.CreateAttribute(users, "id", int4, true, false, 1); err != nil {
		t.Fatal(err)
	}
	if err := cat.CreateAttribute(users, "name", text, true, false, 2); err != nil {
		t.Fatal(err)
	}
	return cat
}

func prepare(t *testing.T, cat *core.Catalog, query string) core.PrepareResult {
	t.Helper()
	stmts, err := postgresql.NewParser().Parse(strings.NewReader(query))
	if err != nil {
		t.Fatalf("parse: %v", err)
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

func TestPrepareSelectColumns(t *testing.T) {
	cat := seedUsers(t)
	res := prepare(t, cat, "SELECT id, name FROM users")

	if len(res.Columns) != 2 {
		t.Fatalf("got %d cols, want 2: %+v", len(res.Columns), res.Columns)
	}
	if res.Columns[0].Name != "id" || res.Columns[0].DataType != "int4" || !res.Columns[0].NotNull {
		t.Errorf("col 0: %+v", res.Columns[0])
	}
	if res.Columns[1].Name != "name" || res.Columns[1].DataType != "text" || !res.Columns[1].NotNull {
		t.Errorf("col 1: %+v", res.Columns[1])
	}
	for i, c := range res.Columns {
		if c.SourceClassOID == 0 || c.SourceAttributeOID == 0 {
			t.Errorf("col %d %s missing source binding: %+v", i, c.Name, c)
		}
	}
}

func TestPrepareSelectStar(t *testing.T) {
	cat := seedUsers(t)
	res := prepare(t, cat, "SELECT * FROM users")

	if len(res.Columns) != 2 {
		t.Fatalf("got %d cols, want 2: %+v", len(res.Columns), res.Columns)
	}
	if res.Columns[0].Name != "id" || res.Columns[1].Name != "name" {
		t.Errorf("got %q, %q; want id, name", res.Columns[0].Name, res.Columns[1].Name)
	}
}

func TestPrepareAlias(t *testing.T) {
	cat := seedUsers(t)
	res := prepare(t, cat, "SELECT id AS user_id FROM users u")

	if len(res.Columns) != 1 || res.Columns[0].Name != "user_id" {
		t.Fatalf("got %+v", res.Columns)
	}
}
