package oracle

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// parseOne parses a single Oracle statement and returns its inner AST node.
func parseOne(t *testing.T, sql string) ast.Node {
	t.Helper()
	p := NewParser()
	stmts, err := p.Parse(strings.NewReader(sql))
	if err != nil {
		t.Fatalf("Parse(%q) returned error: %v", sql, err)
	}
	if len(stmts) != 1 {
		t.Fatalf("Parse(%q) returned %d statements, want 1", sql, len(stmts))
	}
	raw := stmts[0].Raw
	if raw == nil {
		t.Fatalf("Parse(%q) returned a statement with nil Raw", sql)
	}
	return raw.Stmt
}

func TestConvertCreateTable(t *testing.T) {
	stmt := parseOne(t, `CREATE TABLE employees (id NUMBER NOT NULL, name VARCHAR2(100), bio CLOB)`)

	ct, ok := stmt.(*ast.CreateTableStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateTableStmt, got %T", stmt)
	}
	if ct.Name == nil || ct.Name.Name != "employees" {
		t.Fatalf("expected table name 'employees', got %+v", ct.Name)
	}
	if len(ct.Cols) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(ct.Cols))
	}

	if ct.Cols[0].Colname != "id" {
		t.Errorf("col0 name = %q, want id", ct.Cols[0].Colname)
	}
	if !ct.Cols[0].IsNotNull {
		t.Errorf("col0 (id) should be NOT NULL")
	}
	if ct.Cols[0].TypeName == nil || ct.Cols[0].TypeName.Name != "number" {
		t.Errorf("col0 type = %+v, want number", ct.Cols[0].TypeName)
	}
	if ct.Cols[1].Colname != "name" {
		t.Errorf("col1 name = %q, want name", ct.Cols[1].Colname)
	}
	if ct.Cols[1].TypeName == nil || ct.Cols[1].TypeName.Name != "varchar2" {
		t.Errorf("col1 type = %+v, want varchar2", ct.Cols[1].TypeName)
	}
	if ct.Cols[2].TypeName == nil || ct.Cols[2].TypeName.Name != "clob" {
		t.Errorf("col2 type = %+v, want clob", ct.Cols[2].TypeName)
	}
}

func TestConvertSelect(t *testing.T) {
	stmt := parseOne(t, `SELECT id, name FROM employees WHERE id = :1`)

	sel, ok := stmt.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("expected *ast.SelectStmt, got %T", stmt)
	}
	if sel.TargetList == nil || len(sel.TargetList.Items) != 2 {
		t.Fatalf("expected 2 target columns, got %+v", sel.TargetList)
	}
	if sel.FromClause == nil || len(sel.FromClause.Items) != 1 {
		t.Fatalf("expected 1 FROM item, got %+v", sel.FromClause)
	}
	rv, ok := sel.FromClause.Items[0].(*ast.RangeVar)
	if !ok || rv.Relname == nil || *rv.Relname != "employees" {
		t.Fatalf("expected FROM employees RangeVar, got %+v", sel.FromClause.Items[0])
	}

	where, ok := sel.WhereClause.(*ast.A_Expr)
	if !ok {
		t.Fatalf("expected WHERE to be *ast.A_Expr, got %T", sel.WhereClause)
	}
	if _, ok := where.Rexpr.(*ast.ParamRef); !ok {
		t.Errorf("expected WHERE right side to be *ast.ParamRef, got %T", where.Rexpr)
	}
}

func TestConvertSelectStar(t *testing.T) {
	stmt := parseOne(t, `SELECT * FROM employees`)
	sel, ok := stmt.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("expected *ast.SelectStmt, got %T", stmt)
	}
	if sel.TargetList == nil || len(sel.TargetList.Items) != 1 {
		t.Fatalf("expected 1 target (star), got %+v", sel.TargetList)
	}
	rt, ok := sel.TargetList.Items[0].(*ast.ResTarget)
	if !ok {
		t.Fatalf("expected *ast.ResTarget, got %T", sel.TargetList.Items[0])
	}
	cr, ok := rt.Val.(*ast.ColumnRef)
	if !ok || cr.Fields == nil || len(cr.Fields.Items) != 1 {
		t.Fatalf("expected star ColumnRef, got %+v", rt.Val)
	}
	if _, ok := cr.Fields.Items[0].(*ast.A_Star); !ok {
		t.Errorf("expected A_Star, got %T", cr.Fields.Items[0])
	}
}

func TestConvertInsert(t *testing.T) {
	stmt := parseOne(t, `INSERT INTO employees (id, name) VALUES (:1, :2)`)

	ins, ok := stmt.(*ast.InsertStmt)
	if !ok {
		t.Fatalf("expected *ast.InsertStmt, got %T", stmt)
	}
	if ins.Relation == nil || ins.Relation.Relname == nil || *ins.Relation.Relname != "employees" {
		t.Fatalf("expected INSERT INTO employees, got %+v", ins.Relation)
	}
	if ins.Cols == nil || len(ins.Cols.Items) != 2 {
		t.Fatalf("expected 2 insert columns, got %+v", ins.Cols)
	}
	sel, ok := ins.SelectStmt.(*ast.SelectStmt)
	if !ok || sel.ValuesLists == nil || len(sel.ValuesLists.Items) != 1 {
		t.Fatalf("expected a VALUES SelectStmt, got %+v", ins.SelectStmt)
	}
}

func TestConvertUpdate(t *testing.T) {
	stmt := parseOne(t, `UPDATE employees SET name = :1 WHERE id = :2`)

	upd, ok := stmt.(*ast.UpdateStmt)
	if !ok {
		t.Fatalf("expected *ast.UpdateStmt, got %T", stmt)
	}
	if upd.Relations == nil || len(upd.Relations.Items) != 1 {
		t.Fatalf("expected 1 relation, got %+v", upd.Relations)
	}
	if upd.TargetList == nil || len(upd.TargetList.Items) != 1 {
		t.Fatalf("expected 1 SET target, got %+v", upd.TargetList)
	}
	rt, ok := upd.TargetList.Items[0].(*ast.ResTarget)
	if !ok || rt.Name == nil || *rt.Name != "name" {
		t.Fatalf("expected SET name = ..., got %+v", upd.TargetList.Items[0])
	}
	if upd.WhereClause == nil {
		t.Errorf("expected a WHERE clause")
	}
}

func TestConvertDelete(t *testing.T) {
	stmt := parseOne(t, `DELETE FROM employees WHERE id = :1`)

	del, ok := stmt.(*ast.DeleteStmt)
	if !ok {
		t.Fatalf("expected *ast.DeleteStmt, got %T", stmt)
	}
	if del.Relations == nil || len(del.Relations.Items) != 1 {
		t.Fatalf("expected 1 relation, got %+v", del.Relations)
	}
	rv, ok := del.Relations.Items[0].(*ast.RangeVar)
	if !ok || rv.Relname == nil || *rv.Relname != "employees" {
		t.Fatalf("expected DELETE FROM employees, got %+v", del.Relations.Items[0])
	}
	if del.WhereClause == nil {
		t.Errorf("expected a WHERE clause")
	}
}

func TestConvertNamedBindVariable(t *testing.T) {
	stmt := parseOne(t, `SELECT id FROM employees WHERE name = :emp_name`)
	sel := stmt.(*ast.SelectStmt)
	where, ok := sel.WhereClause.(*ast.A_Expr)
	if !ok {
		t.Fatalf("expected WHERE *ast.A_Expr, got %T", sel.WhereClause)
	}
	// Named bind is represented like SQLite's @name: an A_Expr whose right side
	// is the bind name and whose Name is "@".
	inner, ok := where.Rexpr.(*ast.A_Expr)
	if !ok {
		t.Fatalf("expected named-param A_Expr on right side, got %T", where.Rexpr)
	}
	rname, ok := inner.Rexpr.(*ast.String)
	if !ok || rname.Str != "emp_name" {
		t.Errorf("expected named param 'emp_name', got %+v", inner.Rexpr)
	}
}

func TestParseMultipleStatements(t *testing.T) {
	sql := `CREATE TABLE t (id NUMBER);
SELECT id FROM t;`
	p := NewParser()
	stmts, err := p.Parse(strings.NewReader(sql))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(stmts) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(stmts))
	}
	if _, ok := stmts[0].Raw.Stmt.(*ast.CreateTableStmt); !ok {
		t.Errorf("stmt0 = %T, want *ast.CreateTableStmt", stmts[0].Raw.Stmt)
	}
	if _, ok := stmts[1].Raw.Stmt.(*ast.SelectStmt); !ok {
		t.Errorf("stmt1 = %T, want *ast.SelectStmt", stmts[1].Raw.Stmt)
	}
}
