package dinosql

import (
	"fmt"
	"testing"

	core "github.com/kyleconroy/dinosql/internal/pg"
	pg "github.com/lfittl/pg_query_go"

	"github.com/google/go-cmp/cmp"
)

func parseSQLTwo(in string) (QueryTwo, error) {
	tree, err := pg.Parse(in)
	if err != nil {
		return QueryTwo{}, err
	}
	c := core.NewCatalog()
	if err := updateCatalog(&c, tree); err != nil {
		return QueryTwo{}, err
	}

	q, _, err := parseQuery(c, tree.Statements[len(tree.Statements)-1], in)
	q.Stmt = nil
	return q, err
}

func TestQueries(t *testing.T) {
	for _, tc := range []struct {
		name  string
		stmt  string
		query QueryTwo
	}{
		{
			"alias",
			`
			CREATE TABLE bar (id serial not null);
			CREATE TABLE foo (id serial not null, bar serial references bar(id));
			
			DELETE FROM foo f USING bar b
			WHERE f.bar = b.id AND b.id = $1;
			`,
			QueryTwo{
				Params: []Parameter{{Number: 1, Name: "id", Type: "serial"}},
			},
		},
		{
			"star",
			`
			CREATE TABLE bar (bid serial not null);
			CREATE TABLE foo (fid serial not null);
			SELECT * FROM bar, foo;
			`,
			QueryTwo{
				Columns: []core.Column{
					{Name: "bid", DataType: "serial", NotNull: true},
					{Name: "fid", DataType: "serial", NotNull: true},
				},
			},
		},
		{
			"cte_count",
			`
			CREATE TABLE bar (ready bool not null);
			WITH all_count AS (
				SELECT count(*) FROM bar
			), ready_count AS (
				SELECT count(*) FROM bar WHERE ready = true
			)
			SELECT all_count.count, ready_count.count
			FROM all_count, ready_count;
			`,
			QueryTwo{
				Columns: []core.Column{
					{Name: "count", DataType: "integer", NotNull: false},
					{Name: "count", DataType: "integer", NotNull: false},
				},
			},
		},
		{
			"cte_filter",
			`
			CREATE TABLE bar (ready bool not null);
			WITH filter_count AS (
				SELECT count(*) FROM bar WHERE ready = $1
			)
			SELECT filter_count.count
			FROM filter_count;
			`,
			QueryTwo{
				Params: []Parameter{
					{Number: 1, Name: "ready", Type: "bool"},
				},
				Columns: []core.Column{
					{Name: "count", DataType: "integer", NotNull: false},
				},
			},
		},
		{
			"update_set",
			`
			CREATE TABLE foo (name text not null, slug text not null);
			UPDATE foo SET name = $2 WHERE slug = $1;
			`,
			QueryTwo{
				Params: []Parameter{
					{Number: 1, Name: "slug", Type: "text"},
					{Number: 2, Name: "name", Type: "text"},
				},
			},
		},
		{
			"insert_select",
			`
			CREATE TABLE bar (name text not null, ready bool not null);
			CREATE TABLE foo (name text not null, meta text not null);
			INSERT INTO foo (name, meta)
			SELECT name, $1
			FROM bar WHERE ready = $2;
			`,
			QueryTwo{
				Params: []Parameter{
					{Number: 1, Name: "meta", Type: "text"},
					{Number: 2, Name: "ready", Type: "bool"},
				},
			},
		},
	} {
		test := tc
		t.Run(test.name, func(t *testing.T) {
			q, err := parseSQLTwo(test.stmt)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(test.query, q); diff != "" {
				t.Errorf("query mismatch: \n%s", diff)
			}
		})
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
			q, err := parseSQLTwo(fmt.Sprintf(testComparisonSQL, o))
			if err != nil {
				t.Fatal(err)
			}
			expected := QueryTwo{
				Columns: []core.Column{
					{Name: "_", DataType: "bool", NotNull: true},
				},
			}
			if diff := cmp.Diff(expected, q); diff != "" {
				t.Errorf("query mismatch: \n%s", diff)
			}
		})
	}
}
