package compiler

import (
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// TestNarrowToInnermostScope covers the lexical-scope rule used by the
// resolver to disambiguate column references that live inside subqueries.
// See issue #4251.
func TestNarrowToInnermostScope(t *testing.T) {
	t.Parallel()

	str := func(s string) *string { return &s }
	tables := []*ast.TableName{
		{Name: "t1"},
		{Name: "t2"},
	}
	typeMap := map[string]map[string]map[string]*catalog.Column{
		"public": {
			"t1": {"id": {Name: "id"}},
			"t2": {"id": {Name: "id"}, "t1_id": {Name: "t1_id"}},
		},
	}

	t.Run("nil_rv_returns_nil", func(t *testing.T) {
		got := narrowToInnermostScope(tables, typeMap, "public", nil, "id")
		if got != nil {
			t.Fatalf("expected nil (no narrowing) got %v", got)
		}
	})

	t.Run("nil_relname_returns_nil", func(t *testing.T) {
		got := narrowToInnermostScope(tables, typeMap, "public", &ast.RangeVar{}, "id")
		if got != nil {
			t.Fatalf("expected nil (no narrowing) got %v", got)
		}
	})

	t.Run("column_in_inner_scope_narrows_to_inner_table", func(t *testing.T) {
		// The repro shape: ParamRef in inner SELECT (rv=t2) resolving column "id".
		// id exists in t2 -> narrow search to [t2] so the outer t1.id doesn't
		// trigger a spurious "ambiguous" error.
		rv := &ast.RangeVar{Relname: str("t2")}
		got := narrowToInnermostScope(tables, typeMap, "public", rv, "id")
		if len(got) != 1 || got[0].Name != "t2" {
			t.Fatalf("expected narrow to [t2], got %v", got)
		}
	})

	t.Run("column_absent_from_inner_falls_back_to_full_scope", func(t *testing.T) {
		// Correlated-subquery shape: inner SELECT (rv=t2) references an outer
		// column not present in t2. Returning nil tells the caller to keep the
		// full tables list, which lets the outer-scope match win.
		rv := &ast.RangeVar{Relname: str("t2")}
		got := narrowToInnermostScope(tables, typeMap, "public", rv, "not_a_t2_column")
		if got != nil {
			t.Fatalf("expected nil (fall back to all tables), got %v", got)
		}
	})

	t.Run("rv_points_to_unknown_table_falls_back", func(t *testing.T) {
		rv := &ast.RangeVar{Relname: str("nonexistent")}
		got := narrowToInnermostScope(tables, typeMap, "public", rv, "id")
		if got != nil {
			t.Fatalf("expected nil (fall back), got %v", got)
		}
	})
}
