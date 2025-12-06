package clickhouse

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// TestClickHouseUsingMapStructure verifies the structure of the converted JoinExpr
func TestClickHouseUsingMapStructure(t *testing.T) {
	sql := "SELECT * FROM orders LEFT JOIN shipments USING (order_id)"
	parser := NewParser()
	
	stmts, err := parser.Parse(strings.NewReader(sql))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	selectStmt := stmts[0].Raw.Stmt.(*ast.SelectStmt)
	if selectStmt.FromClause == nil || len(selectStmt.FromClause.Items) == 0 {
		t.Fatal("No FROM clause")
	}

	fromItem := selectStmt.FromClause.Items[0]
	topJoin, ok := fromItem.(*ast.JoinExpr)
	if !ok {
		t.Fatalf("Expected JoinExpr, got %T", fromItem)
	}

	t.Logf("Top JoinExpr:")
	t.Logf("  Larg: %T", topJoin.Larg)
	t.Logf("  Rarg: %T", topJoin.Rarg)
	t.Logf("  UsingClause: %v", topJoin.UsingClause)

	// Check nested join
	if nestedJoin, ok := topJoin.Rarg.(*ast.JoinExpr); ok {
		t.Logf("Nested JoinExpr (Rarg):")
		t.Logf("  Larg: %T", nestedJoin.Larg)
		if rvar, ok := nestedJoin.Larg.(*ast.RangeVar); ok {
			t.Logf("    Relname: %v", rvar.Relname)
		}
		t.Logf("  Rarg: %T", nestedJoin.Rarg)
		t.Logf("  UsingClause: %v", nestedJoin.UsingClause)
		if nestedJoin.UsingClause != nil {
			t.Logf("    Items: %d", len(nestedJoin.UsingClause.Items))
			for i, item := range nestedJoin.UsingClause.Items {
				if str, ok := item.(*ast.String); ok {
					t.Logf("      %d: %s", i, str.Str)
				}
			}
		}
	}
}
