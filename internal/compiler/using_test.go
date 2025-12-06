package compiler

import (
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func TestGetJoinUsingMap(t *testing.T) {
	// Create a mock JoinExpr with USING clause
	joinExpr := &ast.JoinExpr{
		UsingClause: &ast.List{
			Items: []ast.Node{
				&ast.String{Str: "order_id"},
			},
		},
		Rarg: &ast.RangeVar{
			Relname: strPtr("shipments"),
		},
	}

	selectStmt := &ast.SelectStmt{
		FromClause: &ast.List{
			Items: []ast.Node{
				joinExpr,
			},
		},
	}

	usingMap := getJoinUsingMap(selectStmt)

	if info, ok := usingMap["shipments"]; ok {
		if !info.HasColumn("order_id") {
			t.Errorf("Expected order_id to be in USING clause for shipments")
		}
	} else {
		t.Errorf("Expected shipments to be in using map")
	}
}

func TestNewUsingInfo(t *testing.T) {
	joinExpr := &ast.JoinExpr{
		UsingClause: &ast.List{
			Items: []ast.Node{
				&ast.String{Str: "id"},
				&ast.String{Str: "type_id"},
			},
		},
	}

	info := NewUsingInfo(joinExpr)

	if !info.HasColumn("id") {
		t.Errorf("Expected id to be in USING columns")
	}
	if !info.HasColumn("type_id") {
		t.Errorf("Expected type_id to be in USING columns")
	}
	if info.HasColumn("other_id") {
		t.Errorf("Did not expect other_id to be in USING columns")
	}
}

func strPtr(s string) *string {
	return &s
}
