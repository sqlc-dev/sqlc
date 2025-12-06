package clickhouse

import (
	"fmt"
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// TestConvertJoinExpr tests that JOIN expressions are properly converted
// to sqlc AST with RangeVar nodes instead of TODO nodes
func TestConvertJoinExpr(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name    string
		query   string
		wantErr bool
		check   func(*ast.SelectStmt) error
	}{
		{
			name:    "simple left join",
			query:   "SELECT u.id, u.name, p.id as post_id FROM users u LEFT JOIN posts p ON u.id = p.user_id",
			wantErr: false,
			check: func(selectStmt *ast.SelectStmt) error {
				// Check that FROM clause contains a JoinExpr
				if selectStmt.FromClause == nil {
					return errorf("FromClause is nil")
				}
				if len(selectStmt.FromClause.Items) == 0 {
					return errorf("FromClause items is empty")
				}

				// The first item in FromClause should be a JoinExpr
				fromItem := selectStmt.FromClause.Items[0]
				joinExpr, ok := fromItem.(*ast.JoinExpr)
				if !ok {
					return errorf("Expected JoinExpr, got %T", fromItem)
				}

				// Larg should be a RangeVar for the left table
				if joinExpr.Larg == nil {
					return errorf("JoinExpr.Larg is nil")
				}
				larg, ok := joinExpr.Larg.(*ast.RangeVar)
				if !ok {
					return errorf("Expected RangeVar for Larg, got %T", joinExpr.Larg)
				}
				if larg.Relname == nil || *larg.Relname != "users" {
					return errorf("Expected left table to be 'users', got %v", larg.Relname)
				}

				// Rarg should be a RangeVar for the right table (after normalization)
				// ClickHouse join structures are normalized to PostgreSQL style at conversion time
				if joinExpr.Rarg == nil {
					return errorf("JoinExpr.Rarg is nil")
				}
				rarg, ok := joinExpr.Rarg.(*ast.RangeVar)
				if !ok {
					return errorf("Expected RangeVar for Rarg (normalized from ClickHouse structure), got %T", joinExpr.Rarg)
				}
				if rarg.Relname == nil || *rarg.Relname != "posts" {
					return errorf("Expected right table to be 'posts', got %v", rarg.Relname)
				}

				return nil
			},
		},
		{
			name:    "join without aliases",
			query:   "SELECT * FROM users INNER JOIN posts ON users.id = posts.user_id",
			wantErr: false,
			check: func(selectStmt *ast.SelectStmt) error {
				if len(selectStmt.FromClause.Items) == 0 {
					return errorf("FromClause items is empty")
				}
				joinExpr, ok := selectStmt.FromClause.Items[0].(*ast.JoinExpr)
				if !ok {
					return errorf("Expected JoinExpr, got %T", selectStmt.FromClause.Items[0])
				}
				if joinExpr.Jointype != ast.JoinTypeInner {
					return errorf("Expected INNER join, got %v", joinExpr.Jointype)
				}
				return nil
			},
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

			if err := tt.check(selectStmt); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
