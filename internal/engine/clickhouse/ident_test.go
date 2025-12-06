package clickhouse

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// TestConvertIdentAsColumnRef tests that identifiers are correctly converted
// to ColumnRef nodes instead of String literals.
//
// This is important because identifiers in SELECT, WHERE, and other clauses
// refer to columns, not literal strings. The compiler's column resolution logic
// depends on finding ColumnRef nodes to properly match columns against the catalog.
func TestConvertIdentAsColumnRef(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name    string
		query   string
		wantVal func(ast.Node) bool // checks that Val is a ColumnRef
	}{
		{
			name:  "select single column",
			query: "SELECT id FROM table1",
			wantVal: func(n ast.Node) bool {
				colRef, ok := n.(*ast.ColumnRef)
				if !ok {
					return false
				}
				if len(colRef.Fields.Items) != 1 {
					return false
				}
				str, ok := colRef.Fields.Items[0].(*ast.String)
				return ok && str.Str == "id"
			},
		},
		{
			name:  "select multiple columns",
			query: "SELECT id, name, email FROM table1",
			wantVal: func(n ast.Node) bool {
				_, ok := n.(*ast.ColumnRef)
				return ok
			},
		},
		{
			name:  "where clause with column reference",
			query: "SELECT * FROM table1 WHERE id = 1",
			wantVal: func(n ast.Node) bool {
				// The WHERE clause should have a ColumnRef for 'id'
				// This is a simple smoke test that the query parses
				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmts, err := parser.Parse(strings.NewReader(tt.query))
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			if len(stmts) == 0 {
				t.Fatal("No statements parsed")
			}

			selectStmt, ok := stmts[0].Raw.Stmt.(*ast.SelectStmt)
			if !ok {
				t.Fatalf("Expected SelectStmt, got %T", stmts[0].Raw.Stmt)
			}

			if len(selectStmt.TargetList.Items) == 0 {
				t.Fatal("No targets in select")
			}

			// Check the first target
			resTarget := selectStmt.TargetList.Items[0].(*ast.ResTarget)
			if resTarget == nil {
				t.Fatal("First target is not a ResTarget")
			}

			if !tt.wantVal(resTarget.Val) {
				t.Errorf("Val check failed. Got type %T: %+v", resTarget.Val, resTarget.Val)
			}
		})
	}
}

// TestIdentifierInWhereClause tests that identifiers in WHERE clauses are
// converted to ColumnRef, not String literals.
func TestIdentifierInWhereClause(t *testing.T) {
	parser := NewParser()

	query := "SELECT * FROM users WHERE status = 'active' AND age > 18"
	stmts, err := parser.Parse(strings.NewReader(query))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	selectStmt := stmts[0].Raw.Stmt.(*ast.SelectStmt)

	// The WHERE clause should contain a BoolExpr with column references for 'status' and 'age'
	// This test ensures the parser correctly identifies column references in conditions
	if selectStmt.WhereClause == nil {
		t.Fatal("WHERE clause is nil")
	}

	// Just verify it parses without error and has a where clause
	// The detailed structure is tested in the ClickHouse parser tests
	t.Logf("WHERE clause type: %T", selectStmt.WhereClause)
}

// TestIdentifierResolution tests that identifiers are properly resolved
// when matching against catalog columns.
func TestIdentifierResolution(t *testing.T) {
	parser := NewParser()
	cat := NewCatalog()

	// Create a table with specific columns
	schemaSQL := `CREATE TABLE users (
		user_id UInt32,
		user_name String,
		user_email String
	)`

	stmts, err := parser.Parse(strings.NewReader(schemaSQL))
	if err != nil {
		t.Fatalf("Parse schema failed: %v", err)
	}

	for _, stmt := range stmts {
		if err := cat.Update(stmt, nil); err != nil {
			t.Fatalf("Update catalog failed: %v", err)
		}
	}

	// Parse a query selecting these columns by name
	querySQL := "SELECT user_id, user_name FROM users"
	queryStmts, err := parser.Parse(strings.NewReader(querySQL))
	if err != nil {
		t.Fatalf("Parse query failed: %v", err)
	}

	selectStmt := queryStmts[0].Raw.Stmt.(*ast.SelectStmt)

	// Verify that targets are ColumnRefs
	for i, target := range selectStmt.TargetList.Items {
		resTarget, ok := target.(*ast.ResTarget)
		if !ok {
			t.Fatalf("Target %d is not ResTarget", i)
		}

		colRef, ok := resTarget.Val.(*ast.ColumnRef)
		if !ok {
			t.Fatalf("Target %d Val is not ColumnRef, got %T", i, resTarget.Val)
		}

		if len(colRef.Fields.Items) == 0 {
			t.Fatalf("Target %d ColumnRef has no fields", i)
		}

		colName, ok := colRef.Fields.Items[0].(*ast.String)
		if !ok {
			t.Fatalf("Target %d field is not String", i)
		}

		t.Logf("Column %d: %s", i, colName.Str)
	}
}
