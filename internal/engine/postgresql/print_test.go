package postgresql

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func TestPrintAst(t *testing.T) {
	p := NewParser()

	queries := []string{
		`SELECT * FROM foo;`,
		`SELECT * FROM foo,bar;`,
		`SELECT * FROM foo WHERE EXISTS (SELECT * FROM foo);`,
		`WITH bar AS (SELECT * FROM foo), bat AS (SELECT 1) SELECT * FROM foo;`,
		`SELECT t.* FROM foo t;`,
		`SELECT *,*,foo.* FROM foo;`,
	}

	// Use astutils to look for select nodes
	// Search for the deepest select nodes

	for i, q := range queries {
		q := q
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			stmts, err := p.Parse(strings.NewReader(q))
			if err != nil {
				t.Fatal(err)
			}
			for _, stmt := range stmts {
				out := ast.Format(stmt.Raw)
				if diff := cmp.Diff(q, out); diff != "" {
					debug.Dump(stmt)
					t.Errorf("- %s", q)
					t.Errorf("+ %s", out)
				}
			}
		})
	}
}
