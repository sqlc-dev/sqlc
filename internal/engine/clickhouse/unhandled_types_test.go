package clickhouse

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// TestUnhandledNodeTypes tests queries that use ClickHouse AST node types
// that we haven't implemented converters for yet.
// These tests identify which unhandled types actually matter in practice.
func TestUnhandledNodeTypes(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name      string
		query     string
		wantErr   bool
		hasTodo   bool // whether we expect TODO nodes in output
		checkFunc func(*ast.SelectStmt) error
	}{
		{
			name:    "array indexing",
			query:   "SELECT arr[1] FROM table1",
			wantErr: false,
			hasTodo: true, // IndexOperation not handled yet
		},
		{
			name:    "array literal",
			query:   "SELECT [1, 2, 3] as arr FROM table1",
			wantErr: false,
			hasTodo: false, // ArrayParamList now handled
		},
		{
			name:    "tuple literal",
			query:   "SELECT (1, 2, 3) as t FROM table1",
			wantErr: false,
			hasTodo: false, // ParamExprList handles this
		},
		{
			name:    "table function",
			query:   "SELECT * FROM numbers(10)",
			wantErr: false,
			hasTodo: false, // TableFunctionExpr now handled
		},
		{
			name:    "array comparison",
			query:   "SELECT * FROM table1 WHERE x IN [1, 2, 3]",
			wantErr: false,
			hasTodo: false, // ArrayParamList now handled
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmts, err := parser.Parse(strings.NewReader(tt.query))
			if (err != nil) != tt.wantErr {
				t.Fatalf("Parse error: %v, wantErr %v", err, tt.wantErr)
			}

			if len(stmts) == 0 {
				t.Fatal("No statements parsed")
			}

			selectStmt, ok := stmts[0].Raw.Stmt.(*ast.SelectStmt)
			if !ok {
				t.Fatalf("Expected SelectStmt, got %T", stmts[0].Raw.Stmt)
			}

			hasTodoNode := containsTODO(selectStmt)
			if hasTodoNode != tt.hasTodo {
				t.Errorf("Expected hasTodo=%v, but got %v", tt.hasTodo, hasTodoNode)
			}

			t.Logf("Query parses: has TODO nodes=%v", hasTodoNode)
		})
	}
}

// containsTODO recursively checks if an AST node tree contains any TODO nodes
func containsTODO(node interface{}) bool {
	if node == nil {
		return false
	}

	switch n := node.(type) {
	case *ast.TODO:
		if n != nil {
			return true
		}
	case *ast.SelectStmt:
		if n != nil {
			if containsTODO(n.TargetList) {
				return true
			}
			if containsTODO(n.FromClause) {
				return true
			}
			if containsTODO(n.WhereClause) {
				return true
			}
			return containsTODO(n.Larg) || containsTODO(n.Rarg)
		}
	case *ast.List:
		if n != nil && n.Items != nil {
			for _, item := range n.Items {
				if item != nil && containsTODO(item) {
					return true
				}
			}
		}
	case *ast.ResTarget:
		if n != nil && containsTODO(n.Val) {
			return true
		}
	case *ast.RangeVar:
		// RangeVar doesn't have a Subquery field in this version
		return false
	case *ast.A_Expr:
		if n != nil {
			return containsTODO(n.Lexpr) || containsTODO(n.Rexpr)
		}
	case *ast.FuncCall:
		if n != nil && containsTODO(n.Args) {
			return true
		}
	case *ast.JoinExpr:
		if n != nil {
			return containsTODO(n.Larg) || containsTODO(n.Rarg) || containsTODO(n.Quals)
		}
	case *ast.ColumnRef:
		if n != nil {
			return containsTODO(n.Fields)
		}
	}

	return false
}
