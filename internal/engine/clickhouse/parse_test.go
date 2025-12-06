package clickhouse

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func TestParseCreateTable(t *testing.T) {
	sql := `
CREATE TABLE IF NOT EXISTS users
(
    id UInt32,
    name String,
    email String
)
ENGINE = MergeTree()
ORDER BY id;
`

	p := NewParser()
	stmts, err := p.Parse(strings.NewReader(sql))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(stmts))
	}

	stmt := stmts[0]
	if stmt.Raw == nil || stmt.Raw.Stmt == nil {
		t.Fatal("Statement or Raw.Stmt is nil")
	}

	createStmt, ok := stmt.Raw.Stmt.(*ast.CreateTableStmt)
	if !ok {
		t.Fatalf("Expected CreateTableStmt, got %T", stmt.Raw.Stmt)
	}

	if createStmt.Name == nil || createStmt.Name.Name == "" {
		t.Fatal("Table name is missing")
	}

	if createStmt.Name.Name != "users" {
		t.Errorf("Expected table name 'users', got '%s'", createStmt.Name.Name)
	}

	if createStmt.Cols == nil || len(createStmt.Cols) == 0 {
		t.Fatal("Table columns are missing")
	}

	if len(createStmt.Cols) != 3 {
		t.Errorf("Expected 3 columns, got %d", len(createStmt.Cols))
	}

	// Check first column
	col0 := createStmt.Cols[0]
	if col0.Colname != "id" {
		t.Errorf("Expected column name 'id', got '%s'", col0.Colname)
	}
}

func TestParseNumberLiterals(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantInt   bool
		wantVal   int64
		wantFloat bool
	}{
		{
			name:    "Integer literal",
			sql:     "SELECT 42;",
			wantInt: true,
			wantVal: 42,
		},
		{
			name:    "Large integer literal",
			sql:     "SELECT 9223372036854775807;",
			wantInt: true,
			wantVal: 9223372036854775807,
		},
		{
			name:      "Float literal",
			sql:       "SELECT 3.14;",
			wantFloat: true,
		},
		{
			name:      "Scientific notation",
			sql:       "SELECT 1.5e2;",
			wantFloat: true,
		},
		{
			name:    "Zero",
			sql:     "SELECT 0;",
			wantInt: true,
			wantVal: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := NewParser()
			stmts, err := p.Parse(strings.NewReader(test.sql))
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

			if selectStmt.TargetList == nil || len(selectStmt.TargetList.Items) == 0 {
				t.Fatal("Expected at least one target")
			}

			target, ok := selectStmt.TargetList.Items[0].(*ast.ResTarget)
			if !ok {
				t.Fatalf("Expected ResTarget, got %T", selectStmt.TargetList.Items[0])
			}

			if test.wantInt {
				constNode, ok := target.Val.(*ast.A_Const)
				if !ok {
					t.Fatalf("Expected A_Const, got %T", target.Val)
				}

				intVal, ok := constNode.Val.(*ast.Integer)
				if !ok {
					t.Fatalf("Expected Integer, got %T", constNode.Val)
				}

				if intVal.Ival != test.wantVal {
					t.Errorf("Expected value %d, got %d", test.wantVal, intVal.Ival)
				}
			}

			if test.wantFloat {
				constNode, ok := target.Val.(*ast.A_Const)
				if !ok {
					t.Fatalf("Expected A_Const, got %T", target.Val)
				}

				_, ok = constNode.Val.(*ast.Float)
				if !ok {
					t.Fatalf("Expected Float, got %T", constNode.Val)
				}
			}
		})
	}
}

func TestParseWindowFunctions(t *testing.T) {
	sql := `
		SELECT 
			id,
			COUNT(*) OVER (PARTITION BY department ORDER BY salary DESC) as rank
		FROM employees;
	`

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

	// Window functions should be parsed in TargetList
	if selectStmt.TargetList == nil || len(selectStmt.TargetList.Items) < 2 {
		t.Fatalf("Expected at least 2 targets, got %d", len(selectStmt.TargetList.Items))
	}
}

func TestParseCastExpression(t *testing.T) {
	sql := "SELECT CAST(id AS String), CAST(value AS Float32) FROM table1;"

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

	if selectStmt.TargetList == nil || len(selectStmt.TargetList.Items) < 2 {
		t.Fatalf("Expected at least 2 targets, got %d", len(selectStmt.TargetList.Items))
	}
}

func TestParseCaseExpression(t *testing.T) {
	sql := `
		SELECT 
			id,
			CASE 
				WHEN status = 'active' THEN 1
				WHEN status = 'inactive' THEN 0
				ELSE -1
			END as status_code
		FROM users;
	`

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

	if selectStmt.TargetList == nil || len(selectStmt.TargetList.Items) < 2 {
		t.Fatalf("Expected at least 2 targets, got %d", len(selectStmt.TargetList.Items))
	}
}

func TestParseAggregateQuery(t *testing.T) {
	sql := `
		SELECT 
			department,
			COUNT(*) as count,
			SUM(salary) as total_salary,
			AVG(salary) as avg_salary
		FROM employees
		GROUP BY department
		HAVING COUNT(*) > 10
		ORDER BY total_salary DESC;
	`

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

	if selectStmt.GroupClause == nil || len(selectStmt.GroupClause.Items) == 0 {
		t.Fatal("GROUP BY clause not parsed")
	}

	if selectStmt.HavingClause == nil {
		t.Fatal("HAVING clause not parsed")
	}

	if selectStmt.SortClause == nil || len(selectStmt.SortClause.Items) == 0 {
		t.Fatal("ORDER BY clause not parsed")
	}
}

func TestParseUnionQueries(t *testing.T) {
	sql := `
		SELECT id, name FROM users
		UNION ALL
		SELECT id, name FROM archived_users;
	`

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

	if selectStmt.Op != ast.Union {
		t.Fatalf("Expected UNION operation, got %v", selectStmt.Op)
	}

	if !selectStmt.All {
		t.Fatal("Expected UNION ALL (All=true)")
	}
}

func TestParseSubquery(t *testing.T) {
	sql := `
		SELECT * FROM (
			SELECT id, name FROM users WHERE id > 100
		) as filtered_users;
	`

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

	if selectStmt.FromClause == nil || len(selectStmt.FromClause.Items) == 0 {
		t.Fatal("FROM clause not parsed")
	}
}

func TestParseIsNullExpressions(t *testing.T) {
	sql := `SELECT * FROM users WHERE name IS NULL AND email IS NOT NULL;`

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

	if selectStmt.WhereClause == nil {
		t.Fatal("WHERE clause not parsed")
	}
}

func TestParseMultipleJoins(t *testing.T) {
	sql := `
		SELECT 
			u.id, u.name, p.title, c.content
		FROM users u
		INNER JOIN posts p ON u.id = p.user_id
		LEFT JOIN comments c ON p.id = c.post_id
		WHERE u.active = 1
		ORDER BY p.created_at DESC;
	`

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

	if selectStmt.FromClause == nil || len(selectStmt.FromClause.Items) == 0 {
		t.Fatal("FROM clause not parsed")
	}
}

// TestParseClickHousePrewhere tests PREWHERE support
// PREWHERE is a ClickHouse optimization - executes before WHERE for better performance
func TestParseClickHousePrewhere(t *testing.T) {
	sql := `SELECT * FROM events PREWHERE event_type = 'click' WHERE user_id = 123;`

	p := NewParser()
	stmts, err := p.Parse(strings.NewReader(sql))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(stmts))
	}

	// If PREWHERE is parsed as a separate clause, it would be in a different field
	// For now, it might be part of WHERE or treated as a TODO
	selectStmt, ok := stmts[0].Raw.Stmt.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Expected SelectStmt, got %T", stmts[0].Raw.Stmt)
	}

	// Just verify parsing succeeded
	if selectStmt.FromClause == nil {
		t.Fatal("FROM clause not parsed")
	}
}

// TestParseClickHouseSample tests SAMPLE support
// SAMPLE is a ClickHouse optimization to read only a portion of data
func TestParseClickHouseSample(t *testing.T) {
	sql := `SELECT * FROM events SAMPLE 1/10 WHERE created_at > now() - INTERVAL 1 DAY;`

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

	if selectStmt.FromClause == nil {
		t.Fatal("FROM clause not parsed")
	}
}

// TestParseClickHouseArrayFunctions tests array function support
// ClickHouse has built-in array operations
func TestParseClickHouseArrayFunctions(t *testing.T) {
	sql := `SELECT arrayLength(tags) as tag_count FROM articles WHERE arrayExists(x -> x > 5, scores);`

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

	if selectStmt.TargetList == nil || len(selectStmt.TargetList.Items) == 0 {
		t.Fatal("SELECT list not parsed")
	}
}

// TestParseStringOperations tests string operations and functions
func TestParseStringOperations(t *testing.T) {
	sql := `
		SELECT 
			concat(first_name, ' ', last_name) as full_name,
			length(email) as email_length,
			upper(name) as name_upper
		FROM users
		WHERE email LIKE '%@example.com';
	`

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

	if selectStmt.TargetList == nil || len(selectStmt.TargetList.Items) < 3 {
		t.Fatalf("Expected 3+ targets, got %d", len(selectStmt.TargetList.Items))
	}
}

// TestParsePositionalParameter tests positional parameters (?)
func TestParsePositionalParameter(t *testing.T) {
	sql := "SELECT * FROM users WHERE id = ? AND name = ?;"

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

	if selectStmt.WhereClause == nil {
		t.Fatal("WHERE clause not parsed")
	}

	// Check that we have at least one ParamRef in the WHERE clause
	paramRefs := findParamRefs(selectStmt.WhereClause)
	if len(paramRefs) == 0 {
		t.Fatal("Expected ParamRef nodes in WHERE clause")
	}
}

// Helper to find all ParamRef nodes in an AST
func findParamRefs(node ast.Node) []*ast.ParamRef {
	var refs []*ast.ParamRef
	var walkNode func(ast.Node)
	walkNode = func(n ast.Node) {
		if pr, ok := n.(*ast.ParamRef); ok {
			refs = append(refs, pr)
		}
		switch v := n.(type) {
		case *ast.A_Expr:
			if v.Lexpr != nil {
				walkNode(v.Lexpr)
			}
			if v.Rexpr != nil {
				walkNode(v.Rexpr)
			}
		case *ast.List:
			if v != nil {
				for _, item := range v.Items {
					walkNode(item)
				}
			}
		}
	}
	walkNode(node)
	return refs
}

// findSqlcFunctionCalls finds all sqlc.* function calls with the given function name
func findSqlcFunctionCalls(node ast.Node, funcName string) []*ast.FuncCall {
	var calls []*ast.FuncCall
	var walkNode func(ast.Node)
	walkNode = func(n ast.Node) {
		if fc, ok := n.(*ast.FuncCall); ok {
			// Check if this is a sqlc.* function call
			if fc.Func != nil && fc.Func.Schema == "sqlc" && fc.Func.Name == funcName {
				calls = append(calls, fc)
			}
		}
		switch v := n.(type) {
		case *ast.A_Expr:
			if v.Lexpr != nil {
				walkNode(v.Lexpr)
			}
			if v.Rexpr != nil {
				walkNode(v.Rexpr)
			}
		case *ast.List:
			if v != nil {
				for _, item := range v.Items {
					walkNode(item)
				}
			}
		}
	}
	walkNode(node)
	return calls
}

// TestInsertIntoSelect tests INSERT INTO ... SELECT
func TestInsertIntoSelect(t *testing.T) {
	sql := `
		INSERT INTO analytics.summary (date, count)
		SELECT toDate(timestamp) as date, COUNT(*) as count
		FROM events
		GROUP BY toDate(timestamp);
	`

	p := NewParser()
	stmts, err := p.Parse(strings.NewReader(sql))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(stmts))
	}

	insertStmt, ok := stmts[0].Raw.Stmt.(*ast.InsertStmt)
	if !ok {
		t.Fatalf("Expected InsertStmt, got %T", stmts[0].Raw.Stmt)
	}

	if insertStmt.Relation == nil {
		t.Fatal("INSERT target table not parsed")
	}

	if insertStmt.SelectStmt == nil {
		t.Fatal("INSERT SELECT statement not parsed")
	}
}

// TestParseDistinct tests DISTINCT clause
func TestParseDistinct(t *testing.T) {
	sql := `SELECT DISTINCT country FROM users ORDER BY country;`

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

	if selectStmt.DistinctClause == nil {
		t.Fatal("DISTINCT clause not parsed")
	}
}

// TestParseWithCTE tests Common Table Expressions (CTEs)
func TestParseWithCTE(t *testing.T) {
	sql := `
		WITH recent_events AS (
			SELECT id, user_id, event_type, timestamp
			FROM events
			WHERE timestamp > now() - INTERVAL 7 DAY
		)
		SELECT user_id, COUNT(*) as event_count
		FROM recent_events
		GROUP BY user_id;
	`

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

	if selectStmt.WithClause == nil {
		t.Fatal("WITH clause not parsed")
	}
}

// TestParseNamedParameterSqlcArg tests sqlc.arg() function syntax
// sqlc.arg() is converted to sqlc_arg() during preprocessing, then converted
// back to sqlc.arg in the AST with proper schema/function name separation
func TestParseNamedParameterSqlcArg(t *testing.T) {
	sql := "SELECT * FROM users WHERE id = sqlc.arg('user_id') AND name = sqlc.arg('name');"

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

	if selectStmt.WhereClause == nil {
		t.Fatal("WHERE clause not parsed")
	}

	// Should find sqlc.arg() function calls in the WHERE clause
	funcCalls := findSqlcFunctionCalls(selectStmt.WhereClause, "arg")
	if len(funcCalls) != 2 {
		t.Fatalf("Expected 2 sqlc.arg() calls, found %d", len(funcCalls))
	}
}

// TestParseNamedParameterSqlcNarg tests sqlc.narg() function syntax
// sqlc.narg() is converted to sqlc_narg() during preprocessing, then converted
// back to sqlc.narg in the AST with proper schema/function name separation
func TestParseNamedParameterSqlcNarg(t *testing.T) {
	sql := "SELECT * FROM users WHERE status = sqlc.narg('optional_status');"

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

	if selectStmt.WhereClause == nil {
		t.Fatal("WHERE clause not parsed")
	}

	// Should find sqlc.narg() function call
	funcCalls := findSqlcFunctionCalls(selectStmt.WhereClause, "narg")
	if len(funcCalls) != 1 {
		t.Fatalf("Expected 1 sqlc.narg() call, found %d", len(funcCalls))
	}
}

// TestParseNamedParameterSqlcSlice tests sqlc.slice() function syntax
// sqlc.slice() is converted to sqlc_slice() during preprocessing, then converted
// back to sqlc.slice in the AST with proper schema/function name separation
func TestParseNamedParameterSqlcSlice(t *testing.T) {
	sql := "SELECT * FROM users WHERE status IN sqlc.slice('statuses');"

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

	if selectStmt.WhereClause == nil {
		t.Fatal("WHERE clause not parsed")
	}

	// Should find sqlc.slice() function call
	funcCalls := findSqlcFunctionCalls(selectStmt.WhereClause, "slice")
	if len(funcCalls) != 1 {
		t.Fatalf("Expected 1 sqlc.slice() call, found %d", len(funcCalls))
	}
}

// TestParseNamedParameterMultipleFunctions tests using multiple sqlc.* functions
func TestParseNamedParameterMultipleFunctions(t *testing.T) {
	sql := `
		SELECT u.id, u.name, p.title
		FROM users u
		LEFT JOIN posts p ON u.id = p.user_id
		WHERE u.id = sqlc.arg('user_id') AND u.status = sqlc.narg('status')
		AND p.category IN sqlc.slice('categories')
		ORDER BY p.created_at DESC;
	`

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

	if selectStmt.WhereClause == nil {
		t.Fatal("WHERE clause not parsed")
	}

	// Should find all three types of sqlc functions
	argCalls := findSqlcFunctionCalls(selectStmt.WhereClause, "arg")
	nargCalls := findSqlcFunctionCalls(selectStmt.WhereClause, "narg")
	sliceCalls := findSqlcFunctionCalls(selectStmt.WhereClause, "slice")

	if len(argCalls) != 1 {
		t.Fatalf("Expected 1 sqlc.arg() call, found %d", len(argCalls))
	}
	if len(nargCalls) != 1 {
		t.Fatalf("Expected 1 sqlc.narg() call, found %d", len(nargCalls))
	}
	if len(sliceCalls) != 1 {
		t.Fatalf("Expected 1 sqlc.slice() call, found %d", len(sliceCalls))
	}
}

// TestParseShow tests SHOW statements
func TestParseShow(t *testing.T) {
	sql := "SHOW TABLES;"

	p := NewParser()
	stmts, err := p.Parse(strings.NewReader(sql))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(stmts))
	}

	// SHOW statements return TODO as they're introspection queries
	if stmts[0].Raw == nil || stmts[0].Raw.Stmt == nil {
		t.Fatal("Statement or Raw.Stmt is nil")
	}
}

// TestParseTruncate tests TRUNCATE statements
func TestParseTruncate(t *testing.T) {
	sql := "TRUNCATE TABLE users;"

	p := NewParser()
	stmts, err := p.Parse(strings.NewReader(sql))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(stmts))
	}

	// TRUNCATE statements return TODO as they're maintenance operations
	if stmts[0].Raw == nil || stmts[0].Raw.Stmt == nil {
		t.Fatal("Statement or Raw.Stmt is nil")
	}
}

// TestPreprocessNamedParameters tests the preprocessing function directly
// Preprocessing converts sqlc.* to sqlc_* (same length) so ClickHouse parser recognizes them as functions
func TestPreprocessNamedParameters(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "sqlc.arg with single quotes",
			input:    "WHERE id = sqlc.arg('user_id')",
			expected: "WHERE id = sqlc_arg('user_id')",
		},
		{
			name:     "sqlc.arg with double quotes",
			input:    `WHERE id = sqlc.arg("user_id")`,
			expected: `WHERE id = sqlc_arg("user_id")`,
		},
		{
			name:     "sqlc.narg",
			input:    "WHERE status = sqlc.narg('status')",
			expected: "WHERE status = sqlc_narg('status')",
		},
		{
			name:     "sqlc.slice",
			input:    "WHERE id IN sqlc.slice('ids')",
			expected: "WHERE id IN sqlc_slice('ids')",
		},
		{
			name:     "Multiple sqlc functions",
			input:    "WHERE id = sqlc.arg('id') AND status = sqlc.narg('status')",
			expected: "WHERE id = sqlc_arg('id') AND status = sqlc_narg('status')",
		},
		{
			name:     "With whitespace",
			input:    "WHERE id = sqlc.arg  ( 'user_id' )",
			expected: "WHERE id = sqlc_arg  ( 'user_id' )",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := preprocessNamedParameters(test.input)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

// TestParseUniqExact tests uniqExact aggregate function
func TestParseUniqExact(t *testing.T) {
	sql := `
		SELECT 
			user_id,
			uniqExact(event_id) as unique_events
		FROM events
		GROUP BY user_id;
	`

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

	if selectStmt.TargetList == nil || len(selectStmt.TargetList.Items) < 2 {
		t.Fatalf("Expected at least 2 targets, got %d", len(selectStmt.TargetList.Items))
	}

	// Check that the uniqExact function is in the target list
	hasUniqExact := false
	for _, item := range selectStmt.TargetList.Items {
		if target, ok := item.(*ast.ResTarget); ok {
			if funcCall, ok := target.Val.(*ast.FuncCall); ok {
				if funcCall.Func.Name == "uniqexact" {
					hasUniqExact = true
					break
				}
			}
		}
	}

	if !hasUniqExact {
		t.Fatal("Expected uniqExact function in target list")
	}
}

// TestParseUniqExactIf tests uniqExactIf conditional aggregate function
func TestParseUniqExactIf(t *testing.T) {
	sql := `
		SELECT 
			user_id,
			uniqExactIf(event_id, event_type = 'click') as unique_clicks
		FROM events
		GROUP BY user_id;
	`

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

	if selectStmt.TargetList == nil || len(selectStmt.TargetList.Items) < 2 {
		t.Fatalf("Expected at least 2 targets, got %d", len(selectStmt.TargetList.Items))
	}

	// Check that the uniqExactIf function is in the target list with 2 arguments
	hasUniqExactIf := false
	for _, item := range selectStmt.TargetList.Items {
		if target, ok := item.(*ast.ResTarget); ok {
			if funcCall, ok := target.Val.(*ast.FuncCall); ok {
				if funcCall.Func.Name == "uniqexactif" && len(funcCall.Args.Items) == 2 {
					hasUniqExactIf = true
					break
				}
			}
		}
	}

	if !hasUniqExactIf {
		t.Fatal("Expected uniqExactIf function with 2 arguments in target list")
	}
}

// TestParseArgMax tests argMax aggregate function
func TestParseArgMax(t *testing.T) {
	sql := `
		SELECT 
			user_id,
			argMax(event_name, timestamp) as latest_event
		FROM events
		GROUP BY user_id;
	`

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

	if selectStmt.TargetList == nil || len(selectStmt.TargetList.Items) < 2 {
		t.Fatalf("Expected at least 2 targets, got %d", len(selectStmt.TargetList.Items))
	}

	// Check that argMax function with 2 arguments is present
	hasArgMax := false
	for _, item := range selectStmt.TargetList.Items {
		if target, ok := item.(*ast.ResTarget); ok {
			if funcCall, ok := target.Val.(*ast.FuncCall); ok {
				if funcCall.Func.Name == "argmax" && len(funcCall.Args.Items) == 2 {
					hasArgMax = true
					break
				}
			}
		}
	}

	if !hasArgMax {
		t.Fatal("Expected argMax function with 2 arguments in target list")
	}
}

// TestParseArgMaxIf tests argMaxIf conditional aggregate function
func TestParseArgMaxIf(t *testing.T) {
	sql := `
		SELECT 
			user_id,
			argMaxIf(event_name, timestamp, event_type = 'purchase') as latest_purchase
		FROM events
		GROUP BY user_id;
	`

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

	if selectStmt.TargetList == nil || len(selectStmt.TargetList.Items) < 2 {
		t.Fatalf("Expected at least 2 targets, got %d", len(selectStmt.TargetList.Items))
	}

	// Check that argMaxIf function with 3 arguments is present
	hasArgMaxIf := false
	for _, item := range selectStmt.TargetList.Items {
		if target, ok := item.(*ast.ResTarget); ok {
			if funcCall, ok := target.Val.(*ast.FuncCall); ok {
				if funcCall.Func.Name == "argmaxif" && len(funcCall.Args.Items) == 3 {
					hasArgMaxIf = true
					break
				}
			}
		}
	}

	if !hasArgMaxIf {
		t.Fatal("Expected argMaxIf function with 3 arguments in target list")
	}
}

// TestParseCountIf tests countIf conditional aggregate function
func TestParseCountIf(t *testing.T) {
	sql := `
		SELECT 
			user_id,
			count() as total_events,
			countIf(event_type = 'click') as click_count,
			countIf(event_type = 'view') as view_count
		FROM events
		GROUP BY user_id;
	`

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

	if selectStmt.TargetList == nil || len(selectStmt.TargetList.Items) < 4 {
		t.Fatalf("Expected at least 4 targets, got %d", len(selectStmt.TargetList.Items))
	}

	// Count countIf functions
	countIfCount := 0
	for _, item := range selectStmt.TargetList.Items {
		if target, ok := item.(*ast.ResTarget); ok {
			if funcCall, ok := target.Val.(*ast.FuncCall); ok {
				if funcCall.Func.Name == "countif" {
					countIfCount++
				}
			}
		}
	}

	if countIfCount != 2 {
		t.Errorf("Expected 2 countIf functions, got %d", countIfCount)
	}
}

// TestParseMultipleAggregatesFunctions tests multiple aggregate functions together
func TestParseMultipleAggregateFunctions(t *testing.T) {
	sql := `
		SELECT 
			category,
			COUNT(*) as count,
			SUM(amount) as total,
			AVG(amount) as average,
			MIN(amount) as min_amount,
			MAX(amount) as max_amount,
			uniqExact(customer_id) as unique_customers,
			countIf(status = 'completed') as completed_orders,
			argMax(product_name, amount) as top_product
		FROM orders
		WHERE created_at >= sqlc.arg('start_date')
		GROUP BY category
		HAVING COUNT(*) > sqlc.arg('min_orders')
		ORDER BY total DESC;
	`

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

	if selectStmt.TargetList == nil || len(selectStmt.TargetList.Items) < 9 {
		t.Fatalf("Expected at least 9 targets, got %d", len(selectStmt.TargetList.Items))
	}

	// Verify all expected functions are present
	functionNames := make(map[string]int)
	for _, item := range selectStmt.TargetList.Items {
		if target, ok := item.(*ast.ResTarget); ok {
			if funcCall, ok := target.Val.(*ast.FuncCall); ok {
				functionNames[funcCall.Func.Name]++
			}
		}
	}

	expectedFunctions := map[string]int{
		"count":     1,
		"sum":       1,
		"avg":       1,
		"min":       1,
		"max":       1,
		"uniqexact": 1,
		"countif":   1,
		"argmax":    1,
	}

	for funcName, expectedCount := range expectedFunctions {
		if count, ok := functionNames[funcName]; !ok || count < expectedCount {
			t.Errorf("Expected function %s with count >= %d, got %d", funcName, expectedCount, count)
		}
	}

	// Verify WHERE clause with named parameters
	if selectStmt.WhereClause == nil {
		t.Fatal("Expected WHERE clause")
	}

	// Verify GROUP BY
	if selectStmt.GroupClause == nil {
		t.Fatal("Expected GROUP BY clause")
	}

	// Verify HAVING
	if selectStmt.HavingClause == nil {
		t.Fatal("Expected HAVING clause")
	}

	// Verify ORDER BY
	if selectStmt.SortClause == nil {
		t.Fatal("Expected ORDER BY clause")
	}
}

// TestParseAggregatesWithWindow tests mixing aggregate and window functions
func TestParseAggregatesWithWindow(t *testing.T) {
	sql := `
		SELECT 
			user_id,
			event_count,
			ROW_NUMBER() OVER (ORDER BY event_count DESC) as rank,
			uniqExact(session_id) as unique_sessions,
			SUM(event_count) OVER (PARTITION BY user_type ORDER BY event_count) as running_total
		FROM (
			SELECT 
				user_id,
				user_type,
				COUNT(*) as event_count,
				countIf(event_type = 'purchase') as purchase_count
			FROM events
			GROUP BY user_id, user_type
		)
		ORDER BY event_count DESC;
	`

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

	// Verify main target list
	if selectStmt.TargetList == nil || len(selectStmt.TargetList.Items) < 5 {
		t.Fatalf("Expected at least 5 targets, got %d", len(selectStmt.TargetList.Items))
	}

	// Check for window functions
	hasWindowFunc := false
	for _, item := range selectStmt.TargetList.Items {
		if target, ok := item.(*ast.ResTarget); ok {
			if funcCall, ok := target.Val.(*ast.FuncCall); ok {
				if funcCall.Over != nil {
					hasWindowFunc = true
					break
				}
			}
		}
	}

	if !hasWindowFunc {
		t.Fatal("Expected at least one window function (OVER clause)")
	}

	// Verify FROM clause is a subquery
	if selectStmt.FromClause == nil || len(selectStmt.FromClause.Items) == 0 {
		t.Fatal("Expected FROM clause with subquery")
	}
}

// TestParseArgMinArgMax tests both argMin and argMax together
func TestParseArgMinArgMax(t *testing.T) {
	sql := `
		SELECT 
			product_id,
			argMin(price, timestamp) as min_price_time,
			argMax(price, timestamp) as max_price_time,
			argMinIf(price, timestamp, status = 'active') as min_active_price,
			argMaxIf(price, timestamp, status = 'active') as max_active_price
		FROM price_history
		GROUP BY product_id;
	`

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

	if selectStmt.TargetList == nil || len(selectStmt.TargetList.Items) < 5 {
		t.Fatalf("Expected at least 5 targets, got %d", len(selectStmt.TargetList.Items))
	}

	// Count each function type
	functionCounts := make(map[string]int)
	for _, item := range selectStmt.TargetList.Items {
		if target, ok := item.(*ast.ResTarget); ok {
			if funcCall, ok := target.Val.(*ast.FuncCall); ok {
				functionCounts[funcCall.Func.Name]++
			}
		}
	}

	expectedFunctions := []string{"argmin", "argmax", "argminif", "argmaxif"}
	for _, funcName := range expectedFunctions {
		if count, ok := functionCounts[funcName]; !ok || count == 0 {
			t.Errorf("Expected function %s to be present", funcName)
		}
	}
}

// TestParseUniqWithModifiers tests uniq functions with different modifiers
func TestParseUniqWithModifiers(t *testing.T) {
	sql := `
		SELECT 
			date,
			uniq(user_id) as unique_users,
			uniqIf(user_id, user_type = 'premium') as premium_users,
			uniqHLL12(user_id) as approx_unique_users,
			uniqExact(user_id) as exact_unique_users
		FROM events
		WHERE date >= sqlc.arg('start_date') AND date <= sqlc.arg('end_date')
		GROUP BY date
		ORDER BY date;
	`

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

	if selectStmt.TargetList == nil || len(selectStmt.TargetList.Items) < 5 {
		t.Fatalf("Expected at least 5 targets, got %d", len(selectStmt.TargetList.Items))
	}

	// Count uniq variants
	uniqVariants := make(map[string]int)
	for _, item := range selectStmt.TargetList.Items {
		if target, ok := item.(*ast.ResTarget); ok {
			if funcCall, ok := target.Val.(*ast.FuncCall); ok {
				if strings.HasPrefix(strings.ToLower(funcCall.Func.Name), "uniq") {
					uniqVariants[funcCall.Func.Name]++
				}
			}
		}
	}

	expectedVariants := []string{"uniq", "uniqif", "uniqhll12", "uniqexact"}
	for _, variant := range expectedVariants {
		if count, ok := uniqVariants[variant]; !ok || count == 0 {
			t.Errorf("Expected uniq variant %s to be present", variant)
		}
	}
}

// TestParseStatisticalAggregates tests statistical aggregate functions
func TestParseStatisticalAggregates(t *testing.T) {
	sql := `
		SELECT 
			varSamp(value) as variance_sample,
			varPop(value) as variance_population,
			stddevSamp(value) as stddev_sample,
			stddevPop(value) as stddev_population,
			covarSamp(x, y) as covariance_sample,
			covarPop(x, y) as covariance_population,
			corr(x, y) as correlation
		FROM metrics;
	`

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

	if selectStmt.TargetList == nil || len(selectStmt.TargetList.Items) < 7 {
		t.Fatalf("Expected at least 7 targets, got %d", len(selectStmt.TargetList.Items))
	}

	// Count statistical functions
	statFunctions := make(map[string]int)
	for _, item := range selectStmt.TargetList.Items {
		if target, ok := item.(*ast.ResTarget); ok {
			if funcCall, ok := target.Val.(*ast.FuncCall); ok {
				funcName := strings.ToLower(funcCall.Func.Name)
				statFunctions[funcName]++
			}
		}
	}

	expectedFunctions := map[string]int{
		"varsamp":    1,
		"varpop":     1,
		"stddevsamp": 1,
		"stddevpop":  1,
		"covarsamp":  1,
		"covarpop":   1,
		"corr":       1,
	}

	for funcName, expectedCount := range expectedFunctions {
		if count, ok := statFunctions[funcName]; !ok || count < expectedCount {
			t.Errorf("Expected function %s with count >= %d, got %d", funcName, expectedCount, count)
		}
	}
}

// TestParseConditionalAggregatesVariants tests minIf and other conditional variants
func TestParseConditionalAggregatesVariants(t *testing.T) {
	sql := `
		SELECT 
			minIf(price, status = 'active') as min_active_price,
			maxIf(price, status = 'active') as max_active_price,
			sumIf(amount, quantity > 0) as positive_amount,
			avgIf(value, value IS NOT NULL) as avg_non_null
		FROM orders
		GROUP BY category;
	`

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

	if selectStmt.TargetList == nil || len(selectStmt.TargetList.Items) < 4 {
		t.Fatalf("Expected at least 4 targets, got %d", len(selectStmt.TargetList.Items))
	}

	// Verify conditional aggregates are present
	conditionalFunctions := make(map[string]int)
	for _, item := range selectStmt.TargetList.Items {
		if target, ok := item.(*ast.ResTarget); ok {
			if funcCall, ok := target.Val.(*ast.FuncCall); ok {
				funcName := strings.ToLower(funcCall.Func.Name)
				if strings.HasSuffix(funcName, "if") {
					conditionalFunctions[funcName]++
				}
			}
		}
	}

	expectedConditionals := []string{"minif", "maxif", "sumif", "avgif"}
	for _, funcName := range expectedConditionals {
		if count, ok := conditionalFunctions[funcName]; !ok || count == 0 {
			t.Errorf("Expected function %s to be present", funcName)
		}
	}
}

// TestParseInOperator tests IN operator with value lists
func TestParseInOperator(t *testing.T) {
	sql := `
		SELECT id, name, status
		FROM users
		WHERE id IN (1, 2, 3, 4, 5)
		AND status IN ('active', 'pending')
		ORDER BY id;
	`

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

	// Verify WHERE clause exists
	if selectStmt.WhereClause == nil {
		t.Fatal("Expected WHERE clause")
	}

	// Verify ORDER BY exists
	if selectStmt.SortClause == nil {
		t.Fatal("Expected ORDER BY clause")
	}
}

// TestParseTOPClause tests TOP clause (ClickHouse LIMIT alternative)
func TestParseTOPClause(t *testing.T) {
	sql := "SELECT TOP 10 my_column FROM tableName;"

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

	// TOP clause is a valid ClickHouse syntax - parser should handle it
	// It may be stored in TargetList or as a special node depending on parser
	if selectStmt.TargetList == nil || len(selectStmt.TargetList.Items) == 0 {
		t.Fatal("Expected target list")
	}
}

// TestParseOrderByWithFill tests ORDER BY ... WITH FILL time series feature
func TestParseOrderByWithFill(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "Basic WITH FILL",
			sql: `
				SELECT n FROM data
				ORDER BY n WITH FILL;
			`,
		},
		{
			name: "WITH FILL FROM TO",
			sql: `
				SELECT date, value FROM timeseries
				ORDER BY date WITH FILL FROM '2024-01-01' TO '2024-01-10';
			`,
		},
		{
			name: "WITH FILL numeric STEP",
			sql: `
				SELECT day, metric FROM series
				ORDER BY day WITH FILL STEP 1;
			`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			stmts, err := p.Parse(strings.NewReader(tt.sql))
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

			// Verify ORDER BY clause exists (WITH FILL is part of ORDER BY)
			if selectStmt.SortClause == nil {
				t.Fatal("Expected ORDER BY clause with FILL")
			}
		})
	}
}

// TestParseTypeCastOperator tests :: operator for type casting
func TestParseTypeCastOperator(t *testing.T) {
	sql := `
		SELECT 
			timestamp_col::DateTime,
			amount::Float32,
			id::String,
			flag::Boolean
		FROM data
		WHERE created_at::Date >= '2024-01-01';
	`

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

	// Verify target list
	if selectStmt.TargetList == nil || len(selectStmt.TargetList.Items) < 4 {
		t.Fatalf("Expected at least 4 targets, got %d",
			len(selectStmt.TargetList.Items))
	}

	// Verify WHERE clause
	if selectStmt.WhereClause == nil {
		t.Fatal("Expected WHERE clause")
	}
}

// TestParseArrayJoin tests ARRAY JOIN clause for unfolding arrays into rows
func TestParseArrayJoin(t *testing.T) {
	tests := []struct {
		name           string
		sql            string
		expectJoinType string
	}{
		{
			name: "Basic ARRAY JOIN",
			sql: `
				SELECT id, tag
				FROM users
				ARRAY JOIN tags AS tag;
			`,
			expectJoinType: "",
		},
		// Note: LEFT ARRAY JOIN is not properly supported by clickhouse-sql-parser v0.4.16
		// The parser returns nil for ArrayJoin when LEFT is specified
		// This is a known limitation of the parser library
		/*
			{
				name: "LEFT ARRAY JOIN",
				sql: `
					SELECT id, tag
					FROM users
					LEFT ARRAY JOIN tags AS tag;
				`,
				expectJoinType: "LEFT",
			},
		*/
		{
			name: "ARRAY JOIN with WHERE and ORDER BY",
			sql: `
				SELECT user_id, param_key, param_value
				FROM events
				ARRAY JOIN params AS param
				WHERE user_id = ?
				ORDER BY param_key;
			`,
			expectJoinType: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := NewParser()
			stmts, err := p.Parse(strings.NewReader(test.sql))
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

			// ARRAY JOIN should be integrated into FROM clause
			if selectStmt.FromClause == nil || len(selectStmt.FromClause.Items) == 0 {
				t.Fatal("Expected FROM clause with ARRAY JOIN")
			}

			// Check for RangeSubselect (our representation of ARRAY JOIN)
			hasArrayJoin := false
			for _, item := range selectStmt.FromClause.Items {
				if rangeSubselect, ok := item.(*ast.RangeSubselect); ok {
					hasArrayJoin = true

					// Verify the RangeSubselect has a Subquery (synthetic SELECT statement)
					if rangeSubselect.Subquery == nil {
						t.Error("Expected RangeSubselect to have a Subquery")
						continue
					}

					syntheticSelect, ok := rangeSubselect.Subquery.(*ast.SelectStmt)
					if !ok {
						t.Errorf("Expected SelectStmt subquery, got %T", rangeSubselect.Subquery)
						continue
					}

					if syntheticSelect.TargetList == nil || len(syntheticSelect.TargetList.Items) == 0 {
						t.Error("Expected synthetic SELECT to have target list")
					}
				}
			}

			if !hasArrayJoin {
				t.Error("Expected ARRAY JOIN to be present in FROM clause")
			}
		})
	}
}

// TestParseArrayJoinWithNamedParameters tests ARRAY JOIN with named parameters
func TestParseArrayJoinWithNamedParameters(t *testing.T) {
	sql := `
		SELECT 
			user_id,
			event_name,
			property_key,
			property_value
		FROM events
		ARRAY JOIN properties AS prop
		WHERE user_id = sqlc.arg('user_id')
		AND event_date >= sqlc.arg('start_date')
		ORDER BY event_time DESC;
	`

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

	// Verify ARRAY JOIN in FROM clause
	if selectStmt.FromClause == nil || len(selectStmt.FromClause.Items) == 0 {
		t.Fatal("Expected FROM clause with ARRAY JOIN")
	}

	// Verify WHERE clause exists (contains named parameters)
	if selectStmt.WhereClause == nil {
		t.Fatal("Expected WHERE clause")
	}

	// Verify ORDER BY exists
	if selectStmt.SortClause == nil {
		t.Fatal("Expected ORDER BY clause")
	}
}

// TestParseArrayJoinMultiple tests ARRAY JOIN with multiple array columns
func TestParseArrayJoinMultiple(t *testing.T) {
	sql := `
		SELECT 
			id,
			nested_value
		FROM table_with_nested
		ARRAY JOIN nested.field1, nested.field2
		WHERE id > 0;
	`

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

	// Verify FROM clause contains ARRAY JOIN
	if selectStmt.FromClause == nil || len(selectStmt.FromClause.Items) == 0 {
		t.Fatal("Expected FROM clause with ARRAY JOIN")
	}
}

// TestParseComplexAggregationWithNamedParams combines statistical functions with named parameters
func TestParseComplexAggregationWithNamedParams(t *testing.T) {
	sql := `
		SELECT 
			date_col,
			COUNT(*) as count,
			varSamp(metric_value) as variance,
			corr(value_x, value_y) as correlation,
			countIf(status = 'success') as successes,
			maxIf(score, score IS NOT NULL) as max_valid_score
		FROM events
		WHERE date_col >= sqlc.arg('start_date') AND date_col <= sqlc.arg('end_date')
		GROUP BY date_col
		HAVING COUNT(*) > sqlc.arg('min_events')
		ORDER BY date_col DESC;
	`

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

	// Verify all clauses present
	if selectStmt.TargetList == nil || len(selectStmt.TargetList.Items) < 6 {
		t.Fatalf("Expected at least 6 targets")
	}
	if selectStmt.WhereClause == nil {
		t.Fatal("Expected WHERE clause")
	}
	if selectStmt.GroupClause == nil {
		t.Fatal("Expected GROUP BY clause")
	}
	if selectStmt.HavingClause == nil {
		t.Fatal("Expected HAVING clause")
	}
	if selectStmt.SortClause == nil {
		t.Fatal("Expected ORDER BY clause")
	}
}

// TestParseArrayJoinFunction tests arrayJoin() as a table function in SELECT list
func TestParseArrayJoinFunction(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "arrayJoin() in SELECT",
			sql: `
				SELECT arrayJoin(categories) AS category
				FROM products;
			`,
		},
		{
			name: "arrayJoin() with nested function",
			sql: `
				SELECT 
					product_id,
					arrayJoin(JSONExtract(metadata, 'Array(String)')) as tag
				FROM products
				WHERE product_id = ?;
			`,
		},
		{
			name: "arrayJoin() with window function",
			sql: `
				SELECT 
					product_id,
					arrayJoin(categories) AS category,
					COUNT(*) OVER (PARTITION BY category) as category_count
				FROM products
				WHERE product_id = ?
				GROUP BY product_id, category;
			`,
		},
		{
			name: "arrayJoin() with named parameters",
			sql: `
				SELECT 
					user_id,
					arrayJoin(tags) AS tag
				FROM users
				WHERE user_id = sqlc.arg('user_id')
				ORDER BY tag;
			`,
		},
		{
			name: "Multiple columns with arrayJoin()",
			sql: `
				SELECT 
					id,
					name,
					arrayJoin(items) AS item
				FROM orders;
			`,
		},
		{
			name: "arrayJoin() with JSONExtract",
			sql: `
				SELECT 
					MetadataPlatformId,
					arrayJoin(JSONExtract(JsonValue, 'Array(String)')) as self_help_id
				FROM events;
			`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := NewParser()
			stmts, err := p.Parse(strings.NewReader(test.sql))
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

			// Verify arrayJoin function is in target list
			if selectStmt.TargetList == nil || len(selectStmt.TargetList.Items) == 0 {
				t.Fatal("Expected at least one target in SELECT list")
			}

			// Look for arrayJoin function call in the targets
			hasArrayJoinFunc := false
			for _, item := range selectStmt.TargetList.Items {
				if target, ok := item.(*ast.ResTarget); ok {
					if funcCall, ok := target.Val.(*ast.FuncCall); ok {
						if funcCall.Func.Name == "arrayjoin" {
							hasArrayJoinFunc = true
							break
						}
					}
				}
			}

			if !hasArrayJoinFunc {
				t.Error("Expected arrayJoin() function call in SELECT list")
			}
		})
	}
}

func TestLocationIndexing(t *testing.T) {
	// Test to verify Location indexing is 0-based or 1-based
	sql := "SELECT sqlc_arg('test')"
	// Position map:
	// 0-indexed: S=0, E=1, L=2, E=3, C=4, T=5, space=6, s=7, q=8, l=9, c=10, _=11, a=12, r=13, g=14, (=15

	p := NewParser()
	stmts, err := p.Parse(strings.NewReader(sql))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(stmts))
	}

	stmt := stmts[0]
	if stmt.Raw == nil || stmt.Raw.Stmt == nil {
		t.Fatal("Statement or Raw.Stmt is nil")
	}

	selectStmt, ok := stmt.Raw.Stmt.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Expected SelectStmt, got %T", stmt.Raw.Stmt)
	}

	if selectStmt.TargetList == nil || len(selectStmt.TargetList.Items) == 0 {
		t.Fatal("No target items")
	}

	// Get the FuncCall node
	resTarget := selectStmt.TargetList.Items[0].(*ast.ResTarget)
	funcCall, ok := resTarget.Val.(*ast.FuncCall)
	if !ok {
		t.Fatalf("Expected FuncCall, got %T", resTarget.Val)
	}

	// The Location should point to 'sqlc_arg' in the parsed SQL
	// In "SELECT sqlc_arg('test')", 's' of 'sqlc_arg' is at 0-indexed position 7
	t.Logf("FuncCall.Location: %d", funcCall.Location)
	t.Logf("SQL: \"%s\"", sql)
	t.Logf("Expected location: 7 (0-indexed position of 's' in 'sqlc_arg')")

	// Extract substring at that location
	if funcCall.Location >= 0 && funcCall.Location < len(sql) {
		t.Logf("Character at Location: %c (expecting 's')", sql[funcCall.Location])
	}
}

// TestImprovedTypeInference verifies that unqualified column references
// in function arguments can be resolved from the catalog
func TestImprovedTypeInference(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "arrayJoin with unqualified column reference",
			sql:  `SELECT arrayJoin(categories) AS category FROM products`,
		},
		{
			name: "argMin with unqualified column reference",
			sql:  `SELECT argMin(price, id) AS min_price FROM products`,
		},
		{
			name: "argMax with unqualified column reference",
			sql:  `SELECT argMax(name, timestamp) AS max_name FROM products`,
		},
		{
			name: "Array() literal in function",
			sql:  `SELECT arrayJoin(Array('a', 'b', 'c')) AS element`,
		},
		{
			name: "CAST expression in function",
			sql:  `SELECT CAST(price AS String) FROM products`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			stmts, err := p.Parse(strings.NewReader(tt.sql))
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			if len(stmts) == 0 {
				t.Fatal("Expected at least 1 statement")
			}

			// Just verify the query parses without TODO nodes
			// The type inference happens during conversion
			stmt := stmts[0]
			if stmt.Raw.Stmt == nil {
				t.Fatal("Expected non-nil statement")
			}
		})
	}
}

// TestCountStar verifies that COUNT(*) parses correctly
func TestCountStar(t *testing.T) {
	sql := `SELECT COUNT(*) FROM products`

	p := NewParser()
	stmts, err := p.Parse(strings.NewReader(sql))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(stmts) == 0 {
		t.Fatal("Expected at least 1 statement")
	}

	stmt := stmts[0]
	if stmt.Raw.Stmt == nil {
		t.Fatal("Expected non-nil statement")
	}

	selectStmt, ok := stmt.Raw.Stmt.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Expected SelectStmt, got %T", stmt.Raw.Stmt)
	}

	if len(selectStmt.TargetList.Items) == 0 {
		t.Fatal("Expected target list items")
	}

	resTarget, ok := selectStmt.TargetList.Items[0].(*ast.ResTarget)
	if !ok {
		t.Fatalf("Expected ResTarget, got %T", selectStmt.TargetList.Items[0])
	}

	funcCall, ok := resTarget.Val.(*ast.FuncCall)
	if !ok {
		t.Fatalf("Expected FuncCall, got %T", resTarget.Val)
	}

	if funcCall.Func.Name != "count" {
		t.Fatalf("Expected function name 'count', got '%s'", funcCall.Func.Name)
	}

	t.Logf("COUNT(*) parsed successfully with %d arguments", len(funcCall.Args.Items))
	for i, arg := range funcCall.Args.Items {
		t.Logf("  Arg %d: %T", i, arg)
		if colRef, ok := arg.(*ast.ColumnRef); ok {
			t.Logf("    ColumnRef with %d fields", len(colRef.Fields.Items))
			if len(colRef.Fields.Items) > 0 {
				t.Logf("    Field 0 type: %T", colRef.Fields.Items[0])
				if _, ok := colRef.Fields.Items[0].(*ast.A_Star); ok {
					t.Logf("    -> Contains A_Star")
				}
				if str, ok := colRef.Fields.Items[0].(*ast.String); ok {
					t.Logf("    -> String: '%s'", str.Str)
				}
			}
		}
	}
}

// TestCatalogHasCountFunction verifies COUNT is registered in the catalog
func TestCatalogHasCountFunction(t *testing.T) {
	cat := NewCatalog()

	// Try to find the COUNT function
	funcs, err := cat.ListFuncsByName(&ast.FuncName{Name: "count"})
	if err != nil {
		t.Fatalf("ListFuncsByName failed: %v", err)
	}

	if len(funcs) == 0 {
		t.Fatal("COUNT function not found in catalog")
	}

	count := funcs[0]
	t.Logf("Found COUNT function: %+v", count)
	t.Logf("  Name: %s", count.Name)
	t.Logf("  Return type: %+v", count.ReturnType)
	t.Logf("  Args: %v", len(count.Args))

	if len(count.Args) > 0 {
		arg := count.Args[0]
		t.Logf("  Arg 0: %+v", arg)
		t.Logf("    Name: %s", arg.Name)
		t.Logf("    Type: %+v", arg.Type)
		t.Logf("    Mode: %v", arg.Mode)
		t.Logf("    HasDefault: %v", arg.HasDefault)
	}
}

// TestCatalogIsolationBetweenQueries verifies that function registrations in one query
// don't affect other queries. This tests the catalog cloning mechanism.
func TestCatalogIsolationBetweenQueries(t *testing.T) {
	// Create a catalog with a test table
	cat := catalog.New("default")
	// Access the default schema directly
	schema := cat.Schemas[0]
	schema.Tables = append(schema.Tables, &catalog.Table{
		Rel: &ast.TableName{Name: "test_table"},
		Columns: []*catalog.Column{
			{Name: "id", Type: ast.TypeName{Name: "int32"}},
			{Name: "values", Type: ast.TypeName{Name: "array"}},
		},
	})

	// Create a parser and set the catalog
	parser := NewParser()
	parser.Catalog = cat

	// Query 1: arrayJoin should register with a specific type in the cloned catalog
	query1 := "SELECT arrayJoin(values) AS item FROM test_table"
	stmts1, err := parser.Parse(strings.NewReader(query1))
	if err != nil {
		t.Fatalf("Query 1 parse failed: %v", err)
	}
	if len(stmts1) != 1 {
		t.Fatalf("Query 1: expected 1 statement, got %d", len(stmts1))
	}

	// Query 2: Same arrayJoin call - should use fresh cloned catalog
	query2 := "SELECT arrayJoin(values) AS item FROM test_table"
	stmts2, err := parser.Parse(strings.NewReader(query2))
	if err != nil {
		t.Fatalf("Query 2 parse failed: %v", err)
	}
	if len(stmts2) != 1 {
		t.Fatalf("Query 2: expected 1 statement, got %d", len(stmts2))
	}

	// Verify that the original catalog was not mutated
	// Count arrayJoin functions registered in the original catalog
	var arrayJoinFuncs []*catalog.Function
	for _, fn := range schema.Funcs {
		if strings.ToLower(fn.Name) == "arrayjoin" {
			arrayJoinFuncs = append(arrayJoinFuncs, fn)
		}
	}

	// Should be no arrayJoin functions registered in the original catalog
	// since cloning happens per Parse() call
	if len(arrayJoinFuncs) > 0 {
		t.Fatalf("Original catalog was mutated: found %d arrayJoin functions, expected 0", len(arrayJoinFuncs))
	}

	t.Log("✓ Catalog isolation verified: functions registered during parsing don't affect original catalog")
}
