package compiler

import (
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func TestFindParametersSelectStmtUsesFromRangeVarForWhereParams(t *testing.T) {
	t.Parallel()

	tableName := "solar_commcard_mapping"
	refs, errs := findParameters(&ast.SelectStmt{
		FromClause: &ast.List{Items: []ast.Node{&ast.RangeVar{Relname: &tableName}}},
		WhereClause: &ast.A_Expr{
			Lexpr: &ast.ColumnRef{Fields: &ast.List{Items: []ast.Node{&ast.String{Str: "deviceId"}}}},
			Rexpr: &ast.ParamRef{Number: 1, Location: 1},
		},
	})
	if len(errs) > 0 {
		t.Fatalf("findParameters returned errors: %v", errs)
	}
	if len(refs) != 1 {
		t.Fatalf("expected 1 ref, got %d", len(refs))
	}
	if refs[0].rv == nil || refs[0].rv.Relname == nil {
		t.Fatal("expected ref to carry range var")
	}
	if got := *refs[0].rv.Relname; got != tableName {
		t.Fatalf("expected ref range var %q, got %q", tableName, got)
	}
}
