package clickhouse

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// TestQualifiedColumnNameExtraction tests that qualified column references
// like u.id without explicit aliases get the column name extracted as the default name
func TestQualifiedColumnNameExtraction(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name          string
		query         string
		expectedNames []string
	}{
		{
			name:          "qualified columns without alias",
			query:         "SELECT u.id, u.name FROM users u",
			expectedNames: []string{"id", "name"},
		},
		{
			name:          "qualified columns with mixed aliases",
			query:         "SELECT u.id, u.name as user_name, u.email FROM users u",
			expectedNames: []string{"id", "user_name", "email"},
		},
		{
			name:          "join with qualified columns",
			query:         "SELECT u.id, u.name, p.title FROM users u LEFT JOIN posts p ON u.id = p.user_id",
			expectedNames: []string{"id", "name", "title"},
		},
		{
			name:          "nested qualified columns",
			query:         "SELECT t.a, t.b, t.c FROM table1 t",
			expectedNames: []string{"a", "b", "c"},
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

			if len(selectStmt.TargetList.Items) != len(tt.expectedNames) {
				t.Fatalf("Expected %d targets, got %d", len(tt.expectedNames), len(selectStmt.TargetList.Items))
			}

			for i, expected := range tt.expectedNames {
				resTarget, ok := selectStmt.TargetList.Items[i].(*ast.ResTarget)
				if !ok {
					t.Fatalf("Target %d is not ResTarget", i)
				}

				if resTarget.Name == nil {
					t.Errorf("Target %d has no name (expected %q)", i, expected)
				} else if *resTarget.Name != expected {
					t.Errorf("Target %d: expected name %q, got %q", i, expected, *resTarget.Name)
				}
			}
		})
	}
}
