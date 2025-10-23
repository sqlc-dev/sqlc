package sqlite

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestConvertComparison(t *testing.T) {
	p := NewParser()

	tests := []struct {
		name     string
		sql      string
		expected string
	}{
		// Basic comparison operators
		{
			name:     "less than",
			sql:      "SELECT * FROM users WHERE age < 18",
			expected: "<",
		},
		{
			name:     "greater than",
			sql:      "SELECT * FROM users WHERE age > 65",
			expected: ">",
		},
		{
			name:     "less than or equal",
			sql:      "SELECT * FROM users WHERE age <= 18",
			expected: "<=",
		},
		{
			name:     "greater than or equal",
			sql:      "SELECT * FROM users WHERE age >= 65",
			expected: ">=",
		},
		{
			name:     "equals",
			sql:      "SELECT * FROM users WHERE status = 'active'",
			expected: "=",
		},
		{
			name:     "not equals (!=)",
			sql:      "SELECT * FROM users WHERE status != 'inactive'",
			expected: "!=",
		},
		{
			name:     "not equals (<>)",
			sql:      "SELECT * FROM users WHERE status <> 'inactive'",
			expected: "<>",
		},
		// Bit operations
		{
			name:     "left shift",
			sql:      "SELECT * FROM users WHERE flags << 2",
			expected: "<<",
		},
		{
			name:     "right shift",
			sql:      "SELECT * FROM users WHERE flags >> 1",
			expected: ">>",
		},
		{
			name:     "bitwise and",
			sql:      "SELECT * FROM users WHERE flags & 4",
			expected: "&",
		},
		{
			name:     "bitwise or",
			sql:      "SELECT * FROM users WHERE flags | 8",
			expected: "|",
		},
		// IS operators
		{
			name:     "is null",
			sql:      "SELECT * FROM users WHERE email IS NULL",
			expected: "IS",
		},
		{
			name:     "is not null",
			sql:      "SELECT * FROM users WHERE email IS NOT NULL",
			expected: "IS NOT",
		},
		// LIKE operators
		{
			name:     "like",
			sql:      "SELECT * FROM users WHERE name LIKE 'John%'",
			expected: "LIKE",
		},
		{
			name:     "not like",
			sql:      "SELECT * FROM users WHERE name NOT LIKE 'Admin%'",
			expected: "NOT LIKE",
		},
		// GLOB operators
		{
			name:     "glob",
			sql:      "SELECT * FROM users WHERE name GLOB 'J*'",
			expected: "GLOB",
		},
		{
			name:     "not glob",
			sql:      "SELECT * FROM users WHERE name NOT GLOB 'A*'",
			expected: "NOT GLOB",
		},
		// MATCH operators
		{
			name:     "match",
			sql:      "SELECT * FROM users WHERE name MATCH 'pattern'",
			expected: "MATCH",
		},
		{
			name:     "not match",
			sql:      "SELECT * FROM users WHERE name NOT MATCH 'pattern'",
			expected: "NOT MATCH",
		},
		// REGEXP operators
		{
			name:     "regexp",
			sql:      "SELECT * FROM users WHERE email REGEXP '.*@example\\.com'",
			expected: "REGEXP",
		},
		{
			name:     "not regexp",
			sql:      "SELECT * FROM users WHERE email NOT REGEXP '.*@spam\\.com'",
			expected: "NOT REGEXP",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stmts, err := p.Parse(strings.NewReader(tc.sql))
			if err != nil {
				t.Fatalf("Failed to parse SQL: %v", err)
			}

			if len(stmts) != 1 {
				t.Fatalf("Expected 1 statement, got %d", len(stmts))
			}

			stmt := stmts[0].Raw.Stmt
			selectStmt, ok := stmt.(*ast.SelectStmt)
			if !ok {
				t.Fatalf("Expected SelectStmt, got %T", stmt)
			}

			// Find the comparison expression in the WHERE clause
			if selectStmt.WhereClause == nil {
				t.Fatal("Expected WHERE clause")
			}

			var foundOperator string
			astutils.Walk(astutils.VisitorFunc(func(node ast.Node) {
				if aExpr, ok := node.(*ast.A_Expr); ok {
					if aExpr.Name != nil && len(aExpr.Name.Items) > 0 {
						if str, ok := aExpr.Name.Items[0].(*ast.String); ok {
							foundOperator = str.Str
						}
					}
				}
			}), selectStmt.WhereClause)

			if foundOperator != tc.expected {
				t.Errorf("Expected operator %q, got %q", tc.expected, foundOperator)
			}
		})
	}
}

func TestConvertInOperation(t *testing.T) {
	p := NewParser()

	tests := []struct {
		name      string
		sql       string
		expectNot bool
	}{
		{
			name:      "in operation",
			sql:       "SELECT * FROM users WHERE status IN ('active', 'pending')",
			expectNot: false,
		},
		{
			name:      "not in operation",
			sql:       "SELECT * FROM users WHERE status NOT IN ('inactive', 'deleted')",
			expectNot: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stmts, err := p.Parse(strings.NewReader(tc.sql))
			if err != nil {
				t.Fatalf("Failed to parse SQL: %v", err)
			}

			if len(stmts) != 1 {
				t.Fatalf("Expected 1 statement, got %d", len(stmts))
			}

			stmt := stmts[0].Raw.Stmt
			selectStmt, ok := stmt.(*ast.SelectStmt)
			if !ok {
				t.Fatalf("Expected SelectStmt, got %T", stmt)
			}

			// Find the IN expression in the WHERE clause
			if selectStmt.WhereClause == nil {
				t.Fatal("Expected WHERE clause")
			}

			var foundIn *ast.In
			astutils.Walk(astutils.VisitorFunc(func(node ast.Node) {
				if inExpr, ok := node.(*ast.In); ok {
					foundIn = inExpr
				}
			}), selectStmt.WhereClause)

			if foundIn == nil {
				t.Fatal("Expected IN expression")
			}

			if foundIn.Not != tc.expectNot {
				t.Errorf("Expected NOT=%v, got NOT=%v", tc.expectNot, foundIn.Not)
			}
		})
	}
}

func TestConvertOrderBy(t *testing.T) {
	p := NewParser()

	tests := []struct {
		name          string
		sql           string
		expectedDirs  []ast.SortByDir
		expectedNulls []ast.SortByNulls
	}{
		{
			name:          "order by default",
			sql:           "SELECT * FROM users ORDER BY name",
			expectedDirs:  []ast.SortByDir{ast.SortByDirDefault},
			expectedNulls: []ast.SortByNulls{ast.SortByNullsDefault},
		},
		{
			name:          "order by asc",
			sql:           "SELECT * FROM users ORDER BY name ASC",
			expectedDirs:  []ast.SortByDir{ast.SortByDirAsc},
			expectedNulls: []ast.SortByNulls{ast.SortByNullsDefault},
		},
		{
			name:          "order by desc",
			sql:           "SELECT * FROM users ORDER BY age DESC",
			expectedDirs:  []ast.SortByDir{ast.SortByDirDesc},
			expectedNulls: []ast.SortByNulls{ast.SortByNullsDefault},
		},
		{
			name:          "order by nulls first",
			sql:           "SELECT * FROM users ORDER BY email NULLS FIRST",
			expectedDirs:  []ast.SortByDir{ast.SortByDirDefault},
			expectedNulls: []ast.SortByNulls{ast.SortByNullsFirst},
		},
		{
			name:          "order by nulls last",
			sql:           "SELECT * FROM users ORDER BY email NULLS LAST",
			expectedDirs:  []ast.SortByDir{ast.SortByDirDefault},
			expectedNulls: []ast.SortByNulls{ast.SortByNullsLast},
		},
		{
			name:          "order by desc nulls first",
			sql:           "SELECT * FROM users ORDER BY score DESC NULLS FIRST",
			expectedDirs:  []ast.SortByDir{ast.SortByDirDesc},
			expectedNulls: []ast.SortByNulls{ast.SortByNullsFirst},
		},
		{
			name:          "order by multiple columns",
			sql:           "SELECT * FROM users ORDER BY name ASC, age DESC",
			expectedDirs:  []ast.SortByDir{ast.SortByDirAsc, ast.SortByDirDesc},
			expectedNulls: []ast.SortByNulls{ast.SortByNullsDefault, ast.SortByNullsDefault},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stmts, err := p.Parse(strings.NewReader(tc.sql))
			if err != nil {
				t.Fatalf("Failed to parse SQL: %v", err)
			}

			if len(stmts) != 1 {
				t.Fatalf("Expected 1 statement, got %d", len(stmts))
			}

			stmt := stmts[0].Raw.Stmt
			selectStmt, ok := stmt.(*ast.SelectStmt)
			if !ok {
				t.Fatalf("Expected SelectStmt, got %T", stmt)
			}

			// Check if SortClause is properly set
			if selectStmt.SortClause == nil {
				t.Fatal("Expected SortClause to be set")
			}

			if len(selectStmt.SortClause.Items) != len(tc.expectedDirs) {
				t.Fatalf("Expected %d sort items, got %d", len(tc.expectedDirs), len(selectStmt.SortClause.Items))
			}

			// Check each sort item
			for i, item := range selectStmt.SortClause.Items {
				sortBy, ok := item.(*ast.SortBy)
				if !ok {
					t.Fatalf("Expected SortBy at index %d, got %T", i, item)
				}

				if sortBy.SortbyDir != tc.expectedDirs[i] {
					t.Errorf("Expected SortbyDir %v at index %d, got %v", tc.expectedDirs[i], i, sortBy.SortbyDir)
				}

				if sortBy.SortbyNulls != tc.expectedNulls[i] {
					t.Errorf("Expected SortbyNulls %v at index %d, got %v", tc.expectedNulls[i], i, sortBy.SortbyNulls)
				}
			}
		})
	}
}

func TestConvertComplexQueries(t *testing.T) {
	p := NewParser()

	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "complex where with multiple operators",
			sql:  "SELECT * FROM users WHERE age >= 18 AND status = 'active' AND email NOT LIKE '%@spam.com' ORDER BY name ASC, age DESC",
		},
		{
			name: "query with IN and ORDER BY",
			sql:  "SELECT * FROM products WHERE category_id IN (1, 2, 3) AND price > 100 ORDER BY price DESC NULLS LAST",
		},
		{
			name: "query with IS NOT and complex ordering",
			sql:  "SELECT * FROM orders WHERE processed_at IS NOT NULL AND total >= 50 ORDER BY created_at DESC, total ASC NULLS FIRST",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stmts, err := p.Parse(strings.NewReader(tc.sql))
			if err != nil {
				t.Fatalf("Failed to parse SQL: %v", err)
			}

			if len(stmts) != 1 {
				t.Fatalf("Expected 1 statement, got %d", len(stmts))
			}

			stmt := stmts[0].Raw.Stmt
			selectStmt, ok := stmt.(*ast.SelectStmt)
			if !ok {
				t.Fatalf("Expected SelectStmt, got %T", stmt)
			}

			// Basic checks to ensure parsing didn't fail
			if selectStmt.WhereClause == nil {
				t.Error("Expected WHERE clause")
			}

			if selectStmt.SortClause == nil {
				t.Error("Expected ORDER BY clause")
			}

			// Verify no TODO nodes were created (which would indicate parsing failures)
			var foundTodo bool
			astutils.Walk(astutils.VisitorFunc(func(node ast.Node) {
				if _, ok := node.(*ast.TODO); ok {
					foundTodo = true
				}
			}), selectStmt)

			if foundTodo {
				t.Error("Found TODO node, indicating incomplete parsing")
			}
		})
	}
}

// Helper function to extract all A_Expr operators from a WHERE clause
func extractOperators(whereClause ast.Node) []string {
	var operators []string
	astutils.Walk(astutils.VisitorFunc(func(node ast.Node) {
		if aExpr, ok := node.(*ast.A_Expr); ok {
			if aExpr.Name != nil && len(aExpr.Name.Items) > 0 {
				if str, ok := aExpr.Name.Items[0].(*ast.String); ok {
					operators = append(operators, str.Str)
				}
			}
		}
	}), whereClause)
	return operators
}

func TestExtractComparisonOperator(t *testing.T) {
	// Test that our helper can extract multiple operators from complex queries
	p := NewParser()

	sql := "SELECT * FROM users WHERE age >= 18 AND status != 'inactive' AND email LIKE '%@company.com'"
	stmts, err := p.Parse(strings.NewReader(sql))
	if err != nil {
		t.Fatalf("Failed to parse SQL: %v", err)
	}

	stmt := stmts[0].Raw.Stmt
	selectStmt, ok := stmt.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Expected SelectStmt, got %T", stmt)
	}

	operators := extractOperators(selectStmt.WhereClause)
	expectedOperators := []string{">=", "!=", "LIKE"}

	if diff := cmp.Diff(expectedOperators, operators, cmpopts.EquateEmpty()); diff != "" {
		t.Errorf("operators mismatch:\n%s", diff)
	}
}
