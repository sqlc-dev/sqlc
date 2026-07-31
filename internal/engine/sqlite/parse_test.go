package sqlite

import (
	"errors"
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

// TestStatementExtent checks the byte ranges sqlc slices the source with. A
// statement runs from just after the previous terminator, so that the
// "-- name:" comment above it is inside its extent, up to the last token
// before its own terminator.
func TestStatementExtent(t *testing.T) {
	const query = `-- name: One :one
SELECT 1;

/* name: Two :many */
SELECT 2 ;
-- name: Three :exec
SELECT 3`

	stmts, err := NewParser().Parse(strings.NewReader(query))
	if err != nil {
		t.Fatal(err)
	}
	if len(stmts) != 3 {
		t.Fatalf("parsed %d statements, want 3", len(stmts))
	}

	want := []string{
		"-- name: One :one\nSELECT 1",
		"\n\n/* name: Two :many */\nSELECT 2",
		"\n-- name: Three :exec\nSELECT 3",
	}
	for i, stmt := range stmts {
		start := stmt.Raw.StmtLocation
		got := query[start : start+stmt.Raw.StmtLen]
		if got != want[i] {
			t.Errorf("statement %d extent = %q, want %q", i, got, want[i])
		}
	}
}

// TestSqlcFunctions checks that sqlc's own functions survive a parser that
// implements SQLite's grammar, which has no schema-qualified function call.
func TestSqlcFunctions(t *testing.T) {
	for _, tc := range []struct {
		query  string
		schema string
		name   string
	}{
		{"SELECT * FROM foo WHERE bar = sqlc.arg(bar)", "sqlc", "arg"},
		{"SELECT * FROM foo WHERE bar = sqlc.narg(bar)", "sqlc", "narg"},
		{"SELECT * FROM foo WHERE bar IN (sqlc.slice(bars))", "sqlc", "slice"},
		{"SELECT sqlc.embed(foo) FROM foo", "sqlc", "embed"},
		// An unknown one still parses; the name is rejected later, with a
		// message that has to be able to say "sqlc.argh".
		{"SELECT * FROM foo WHERE bar = sqlc.argh(bar)", "sqlc", "argh"},
		// A call written in one part is left alone.
		{"SELECT * FROM foo WHERE bar = abs(bar)", "", "abs"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stmts, err := NewParser().Parse(strings.NewReader(tc.query))
			if err != nil {
				t.Fatal(err)
			}
			var found *ast.FuncCall
			walk(stmts[0].Raw.Stmt, func(n ast.Node) {
				if call, ok := n.(*ast.FuncCall); ok && found == nil {
					found = call
				}
			})
			if found == nil {
				t.Fatalf("no function call found in %q", tc.query)
			}
			if found.Func.Schema != tc.schema || found.Func.Name != tc.name {
				t.Errorf("got %q.%q, want %q.%q",
					found.Func.Schema, found.Func.Name, tc.schema, tc.name)
			}
		})
	}
}

// TestSqlcFunctionInString checks that the rewrite behind TestSqlcFunctions
// works off the token stream and so cannot reach inside a string literal.
func TestSqlcFunctionInString(t *testing.T) {
	stmts, err := NewParser().Parse(strings.NewReader(`SELECT 'sqlc.arg(x)' FROM foo`))
	if err != nil {
		t.Fatal(err)
	}
	var found *ast.String
	walk(stmts[0].Raw.Stmt, func(n ast.Node) {
		if c, ok := n.(*ast.A_Const); ok {
			if s, ok := c.Val.(*ast.String); ok && found == nil {
				found = s
			}
		}
	})
	if found == nil {
		t.Fatal("no string constant found")
	}
	if found.Str != "sqlc.arg(x)" {
		t.Errorf("string constant = %q, want %q", found.Str, "sqlc.arg(x)")
	}
}

func TestSyntaxError(t *testing.T) {
	_, err := NewParser().Parse(strings.NewReader("SELECT 1;\nSELECT FROM foo;"))
	if err == nil {
		t.Fatal("expected a syntax error")
	}
	var serr *sqlerr.Error
	if !errors.As(err, &serr) {
		t.Fatalf("error is %T, want *sqlerr.Error", err)
	}
	if serr.Line != 2 || serr.Column != 8 {
		t.Errorf("error at %d:%d, want 2:8", serr.Line, serr.Column)
	}
	if !strings.Contains(serr.Message, `near "FROM"`) {
		t.Errorf("message = %q, want it to mention the offending token", serr.Message)
	}
}

// walk visits every node reachable from n through the fields the sqlite
// converter populates. It is a test helper, not a general traversal.
func walk(n ast.Node, fn func(ast.Node)) {
	if n == nil {
		return
	}
	fn(n)
	switch t := n.(type) {
	case *ast.RawStmt:
		walk(t.Stmt, fn)
	case *ast.List:
		for _, item := range t.Items {
			walk(item, fn)
		}
	case *ast.SelectStmt:
		walk(t.TargetList, fn)
		walk(t.FromClause, fn)
		walk(t.WhereClause, fn)
	case *ast.ResTarget:
		walk(t.Val, fn)
	case *ast.A_Expr:
		walk(t.Lexpr, fn)
		walk(t.Rexpr, fn)
	case *ast.In:
		walk(t.Expr, fn)
		for _, item := range t.List {
			walk(item, fn)
		}
	case *ast.A_Const:
		walk(t.Val, fn)
	case *ast.FuncCall:
		walk(t.Args, fn)
	}
}
