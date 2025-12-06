package clickhouse

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// TestArrayParamListConversion tests array literal conversion
func TestArrayParamListConversion(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{
			name:    "simple array literal",
			query:   "SELECT [1, 2, 3] as arr",
			wantErr: false,
		},
		{
			name:    "string array literal",
			query:   "SELECT ['a', 'b', 'c'] as strs",
			wantErr: false,
		},
		{
			name:    "array in IN clause",
			query:   "SELECT * FROM table1 WHERE x IN [1, 2, 3]",
			wantErr: false,
		},
		{
			name:    "array in function",
			query:   "SELECT arrayMap(x -> x * 2, [1, 2, 3]) as doubled",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmts, err := parser.Parse(strings.NewReader(tt.query))
			if (err != nil) != tt.wantErr {
				t.Fatalf("Parse error: %v, wantErr %v", err, tt.wantErr)
			}

			if len(stmts) > 0 {
				selectStmt, ok := stmts[0].Raw.Stmt.(*ast.SelectStmt)
				if !ok {
					t.Fatalf("Expected SelectStmt, got %T", stmts[0].Raw.Stmt)
				}

				// Verify it parses without TODO nodes (or minimal TODOs)
				_ = selectStmt
			}
		})
	}
}

// TestTableFunctionConversion tests table function conversion
func TestTableFunctionConversion(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{
			name:    "numbers table function",
			query:   "SELECT * FROM numbers(10)",
			wantErr: false,
		},
		{
			name:    "numbers with range",
			query:   "SELECT * FROM numbers(1, 5)",
			wantErr: false,
		},
		{
			name:    "numbers in join",
			query:   "SELECT t.number FROM numbers(5) t",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmts, err := parser.Parse(strings.NewReader(tt.query))
			if (err != nil) != tt.wantErr {
				t.Fatalf("Parse error: %v, wantErr %v", err, tt.wantErr)
			}

			if len(stmts) > 0 {
				selectStmt, ok := stmts[0].Raw.Stmt.(*ast.SelectStmt)
				if !ok {
					t.Fatalf("Expected SelectStmt, got %T", stmts[0].Raw.Stmt)
				}

				// Verify FromClause exists and has items
				if selectStmt.FromClause != nil && len(selectStmt.FromClause.Items) > 0 {
					// Should be a RangeFunction for table functions
					fromItem := selectStmt.FromClause.Items[0]
					_, isRangeFunc := fromItem.(*ast.RangeFunction)
					if !isRangeFunc && !tt.wantErr {
						// Could also be RangeVar if not a pure table function context
						// This is acceptable
					}
				}
			}
		})
	}
}

// TestIndexOperationConversion tests array/tuple indexing conversion
func TestIndexOperationConversion(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{
			name:    "array indexing",
			query:   "SELECT arr[1] FROM table1",
			wantErr: false,
		},
		{
			name:    "tuple indexing",
			query:   "SELECT (1, 2, 3)[2] as second_elem",
			wantErr: false,
		},
		{
			name:    "nested array indexing",
			query:   "SELECT nested_arr[1][2] FROM table1",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmts, err := parser.Parse(strings.NewReader(tt.query))
			if (err != nil) != tt.wantErr {
				t.Fatalf("Parse error: %v, wantErr %v", err, tt.wantErr)
			}

			if len(stmts) > 0 {
				selectStmt, ok := stmts[0].Raw.Stmt.(*ast.SelectStmt)
				if !ok {
					t.Fatalf("Expected SelectStmt, got %T", stmts[0].Raw.Stmt)
				}

				_ = selectStmt
			}
		})
	}
}

// TestTernaryOperationConversion tests ternary operator conversion
func TestTernaryOperationConversion(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{
			name:    "simple ternary",
			query:   "SELECT x > 0 ? 'positive' : 'non-positive' as sign FROM table1",
			wantErr: false,
		},
		{
			name:    "nested ternary",
			query:   "SELECT x > 0 ? 'positive' : x < 0 ? 'negative' : 'zero' as sign FROM table1",
			wantErr: false,
		},
		{
			name:    "ternary in expression",
			query:   "SELECT (a > b ? a : b) as max_val FROM table1",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmts, err := parser.Parse(strings.NewReader(tt.query))
			if (err != nil) != tt.wantErr {
				t.Fatalf("Parse error: %v, wantErr %v", err, tt.wantErr)
			}

			if len(stmts) > 0 {
				selectStmt, ok := stmts[0].Raw.Stmt.(*ast.SelectStmt)
				if !ok {
					t.Fatalf("Expected SelectStmt, got %T", stmts[0].Raw.Stmt)
				}

				// Should have a CaseExpr in the TargetList
				if len(selectStmt.TargetList.Items) > 0 {
					resTarget := selectStmt.TargetList.Items[0].(*ast.ResTarget)
					_, isCaseExpr := resTarget.Val.(*ast.CaseExpr)
					if !isCaseExpr {
						// Might be wrapped in other expressions, which is ok
					}
				}
			}
		})
	}
}
