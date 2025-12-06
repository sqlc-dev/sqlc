package clickhouse

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func TestCaseSensitiveColumns(t *testing.T) {
	// ClickHouse is case-sensitive for identifiers
	// This test demonstrates the issue where columns with different cases
	// are incorrectly treated as the same column

	sql := `
CREATE TABLE test_table
(
    UserId UInt32,
    userName String,
    EMAIL String
)
ENGINE = MergeTree()
ORDER BY UserId;
`

	p := NewParser()
	stmts, err := p.Parse(strings.NewReader(sql))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(stmts))
	}

	createStmt, ok := stmts[0].Raw.Stmt.(*ast.CreateTableStmt)
	if !ok {
		t.Fatalf("Expected CreateTableStmt, got %T", stmts[0].Raw.Stmt)
	}

	// Check that column names preserve their case
	expectedColumns := map[string]bool{
		"UserId":   true,
		"userName": true,
		"EMAIL":    true,
	}

	actualColumns := make(map[string]bool)
	for _, col := range createStmt.Cols {
		actualColumns[col.Colname] = true
	}

	if len(actualColumns) != len(expectedColumns) {
		t.Errorf("Expected %d distinct columns, got %d", len(expectedColumns), len(actualColumns))
	}

	for expected := range expectedColumns {
		if !actualColumns[expected] {
			t.Errorf("Column '%s' not found. Found columns: %v", expected, actualColumns)
		}
	}
}

func TestCaseSensitiveColumnReference(t *testing.T) {
	// Test that column references preserve case in SELECT statements
	sql := "SELECT UserId, userName, EMAIL FROM test_table;"

	p := NewParser()
	stmts, err := p.Parse(strings.NewReader(sql))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(stmts))
	}

	selectStmt, ok := stmts[0].Raw.Stmt.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Expected SelectStmt, got %T", stmts[0].Raw.Stmt)
	}

	expectedColRefs := []string{"UserId", "userName", "EMAIL"}
	if len(selectStmt.TargetList.Items) != len(expectedColRefs) {
		t.Fatalf("Expected %d target items, got %d", len(expectedColRefs), len(selectStmt.TargetList.Items))
	}

	for i, expected := range expectedColRefs {
		target, ok := selectStmt.TargetList.Items[i].(*ast.ResTarget)
		if !ok {
			t.Fatalf("Item %d is not a ResTarget: %T", i, selectStmt.TargetList.Items[i])
		}

		// Check if Name is set (for aliased columns) or extract from ColumnRef
		var got string
		if target.Name != nil && *target.Name != "" {
			got = *target.Name
		} else if colRef, ok := target.Val.(*ast.ColumnRef); ok && colRef != nil && colRef.Fields != nil && len(colRef.Fields.Items) > 0 {
			// Extract the column name from the ColumnRef
			if s, ok := colRef.Fields.Items[len(colRef.Fields.Items)-1].(*ast.String); ok {
				got = s.Str
			}
		}

		if got != expected {
			t.Errorf("Column %d: expected '%s', got '%s'", i, expected, got)
		}
	}
}

func TestCaseSensitiveWhereClauses(t *testing.T) {
	// Test that WHERE clauses with case-sensitive column names work correctly
	sql := "SELECT * FROM users WHERE UserId = 123 AND userName = 'John';"

	p := NewParser()
	stmts, err := p.Parse(strings.NewReader(sql))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	selectStmt, ok := stmts[0].Raw.Stmt.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Expected SelectStmt, got %T", stmts[0].Raw.Stmt)
	}

	// Verify WHERE clause references preserve case
	if selectStmt.WhereClause == nil {
		t.Fatal("WHERE clause is nil")
	}

	// The WHERE clause should contain column references with preserved case
	// This is a simple check - we'd need to traverse the AST to verify
	// that column names in the WHERE clause preserve their case
	whereStr := astToString(selectStmt.WhereClause)

	// Check that the case is preserved in the where clause
	if !strings.Contains(whereStr, "UserId") || !strings.Contains(whereStr, "userName") {
		t.Errorf("WHERE clause should preserve column name case. Got: %s", whereStr)
	}
}

// astToString converts AST nodes to a string representation for testing
func astToString(node ast.Node) string {
	if node == nil {
		return ""
	}

	switch n := node.(type) {
	case *ast.A_Expr:
		left := astToString(n.Lexpr)
		right := astToString(n.Rexpr)
		return left + " " + right
	case *ast.ColumnRef:
		if n.Fields != nil && len(n.Fields.Items) > 0 {
			if s, ok := n.Fields.Items[len(n.Fields.Items)-1].(*ast.String); ok {
				return s.Str
			}
		}
	case *ast.A_Const:
		if s, ok := n.Val.(*ast.String); ok {
			return s.Str
		}
	}
	return ""
}
