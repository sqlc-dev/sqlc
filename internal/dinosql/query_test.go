package dinosql

import (
	"fmt"
	"strconv"
	"testing"

	core "github.com/kyleconroy/sqlc/internal/pg"

	"github.com/google/go-cmp/cmp"
	pg "github.com/lfittl/pg_query_go"
)

func parseSQL(in string) (Query, error) {
	tree, err := pg.Parse(in)
	if err != nil {
		return Query{}, err
	}
	c := core.NewCatalog()
	if err := updateCatalog(&c, tree); err != nil {
		return Query{}, err
	}

	q, err := parseQuery(c, tree.Statements[len(tree.Statements)-1], in)
	if q == nil {
		return Query{}, err
	}
	return *q, err
}

func public(rel string) core.FQN {
	return core.FQN{
		Catalog: "",
		Schema:  "public",
		Rel:     rel,
	}
}

const testComparisonSQL = `
CREATE TABLE bar (id serial not null);
SELECT count(*) %s 0 FROM bar;
`

func TestComparisonOperators(t *testing.T) {
	for _, op := range []string{">", "<", ">=", "<=", "<>", "!=", "="} {
		o := op
		t.Run(o, func(t *testing.T) {
			q, err := parseSQL(fmt.Sprintf(testComparisonSQL, o))
			if err != nil {
				t.Fatal(err)
			}
			expected := Query{
				SQL: q.SQL,
				Columns: []core.Column{
					{Name: "", DataType: "bool", NotNull: true},
				},
			}
			if diff := cmp.Diff(expected, q); diff != "" {
				t.Errorf("query mismatch: \n%s", diff)
			}
		})
	}
}

func TestUnknownFunctions(t *testing.T) {
	stmt := `
		CREATE TABLE foo (id text not null);
		-- name: ListFoos :one
		SELECT id FROM foo WHERE id = frobnicate($1);
		`
	_, err := parseSQL(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInvalidQueries(t *testing.T) {
	for i, tc := range []struct {
		stmt string
		msg  string
	}{
		{
			`
			CREATE TABLE foo (id text not null);
			-- name: ListFoos
			SELECT id FROM foo;
			`,
			"invalid query comment: -- name: ListFoos",
		},
		{
			`
			CREATE TABLE foo (id text not null);
			-- name: ListFoos :one :many
			SELECT id FROM foo;
			`,
			"invalid query comment: -- name: ListFoos :one :many",
		},
		{
			`
			CREATE TABLE foo (id text not null);
			-- name: ListFoos :two
			SELECT id FROM foo;
			`,
			"invalid query type: :two",
		},
		{
			`
			CREATE TABLE foo (id text not null);
			-- name: DeleteFoo :one
			DELETE FROM foo WHERE id = $1;
			`,
			`query "DeleteFoo" specifies parameter ":one" without containing a RETURNING clause`,
		},
		{
			`
			CREATE TABLE foo (id text not null);
			-- name: UpdateFoo :one
			UPDATE foo SET id = $2 WHERE id = $1;
			`,
			`query "UpdateFoo" specifies parameter ":one" without containing a RETURNING clause`,
		},
		{
			`
			CREATE TABLE foo (id text not null);
			-- name: InsertFoo :one
			INSERT INTO foo (id) VALUES ($1);
			`,
			`query "InsertFoo" specifies parameter ":one" without containing a RETURNING clause`,
		},
		{
			`
			CREATE TABLE foo (bar text not null, baz text not null);
			INSERT INTO foo (bar, baz) VALUES ($1);
			`,
			`INSERT has more target columns than expressions`,
		},
		{
			`
			CREATE TABLE foo (bar text not null, baz text not null);
			INSERT INTO foo (bar) VALUES ($1, $2);
			`,
			`INSERT has more expressions than target columns`,
		},
	} {
		test := tc
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			_, err := parseSQL(test.stmt)
			if err == nil {
				t.Fatalf("expected err, got nil")
			}
			if diff := cmp.Diff(test.msg, err.Error()); diff != "" {
				t.Errorf("error message differs: \n%s", diff)
			}
		})
	}
}
