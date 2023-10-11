package postgresql

import (
	"fmt"
	"strings"
	"testing"

	pg_query "github.com/pganalyze/pg_query_go/v4"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func TestPrintAst(t *testing.T) {
	p := NewParser()

	queries := []string{
		`SELECT * FROM foo;`,
		`SELECT * 
FROM foo;`,
		`SELECT * FROM foo,bar;`,
		`SELECT * FROM foo WHERE EXISTS (SELECT * FROM foo);`,
		`WITH bar AS (SELECT * FROM foo), bat AS (SELECT 1) SELECT * FROM foo;`,
		`SELECT t.* FROM foo t;`,
		`SELECT *,*,foo.* FROM foo;`,
		`SELECT 'foo';`,
		`SELECT true;`,
		`SELECT 1.2;`,
		`SELECT "foo";`,
		`SELECT * FROM foo LIMIT 1;`,
		`SELECT * FROM foo OFFSET 1;`,
		`SELECT * FROM foo LIMIT 1 OFFSET 1;`,
		`SELECT * FROM foo ORDER BY name;`,
		`SELECT DISTINCT * FROM foo;`,
		`SELECT DISTINCT ON (location) location, time, report
		FROM weather_reports
		ORDER BY location, time DESC;`,
		`SELECT * FROM (SELECT * FROM mytable FOR SHARE) ss WHERE col1 = 5;`,
		`INSERT INTO myschema.foo (a, b) VALUES ($1, $2);`,
	}

	// Use astutils to look for select nodes
	// Search for the deepest select nodes

	for i, q := range queries {
		q := q
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			expected, err := pg_query.Fingerprint(q)
			if err != nil {
				t.Fatal(err)
			}
			stmts, err := p.Parse(strings.NewReader(q))
			if err != nil {
				t.Fatal(err)
			}
			if len(stmts) != 1 {
				t.Fatal("expected one statement")
			}
			out := ast.Format(stmts[0].Raw)
			actual, err := pg_query.Fingerprint(out)
			if err != nil {
				t.Error(err)
			}
			if expected != actual {
				t.Errorf("- %s", expected)
				t.Errorf("- %s", q)
				t.Errorf("+ %s", actual)
				t.Errorf("+ %s", out)
			}
		})
	}
}
