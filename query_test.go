package dinosql

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kyleconroy/dinosql/postgres"
)

const testAliasSQL = `
CREATE TABLE bar (id serial not null);
CREATE TABLE foo (id serial not null, bar serial references bar(id));

-- name: DeleteFoo :exec
DELETE FROM foo f
USING bar b
WHERE f.bar = b.id AND b.id = $1;
`

func TestAlias(t *testing.T) {
	result, err := parseSQL(testAliasSQL)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Query{
		{
			Type:       ":exec",
			MethodName: "DeleteFoo",
			StmtName:   "deleteFoo",
			QueryName:  "deleteFoo",
			SQL:        "-- name: DeleteFoo :exec\nDELETE FROM foo f\nUSING bar b\nWHERE f.bar = b.id AND b.id = $1",
			Args:       []Arg{{Name: "id", Type: "int"}},
			Table: postgres.Table{
				GoName: "Foo",
				Name:   "foo",
				Columns: []postgres.Column{
					{GoName: "ID", GoType: "int", Name: "id", Type: "serial", NotNull: true},
					{GoName: "Bar", GoType: "int", Name: "bar", Type: "serial"},
				},
			},
		},
	}

	if diff := cmp.Diff(expected, result.Queries); diff != "" {
		t.Errorf("query mismatch: \n%s", diff)
	}
}

const testComparisonSQL = `
CREATE TABLE bar (id serial not null);

-- name: HasBar :one
SELECT count(*) %s 0
FROM bar;
`

func TestComparisonOperators(t *testing.T) {
	sqlTmpl := "-- name: HasBar :one\nSELECT count(*) %s 0\nFROM bar"

	for _, op := range []string{">", "<", ">=", "<=", "<>", "!=", "="} {
		o := op
		t.Run(o, func(t *testing.T) {
			result, err := parseSQL(fmt.Sprintf(testComparisonSQL, o))
			if err != nil {
				t.Fatal(err)
			}

			expected := []Query{
				{
					Type:       ":one",
					MethodName: "HasBar",
					StmtName:   "hasBar",
					QueryName:  "hasBar",
					SQL:        fmt.Sprintf(sqlTmpl, o),
					Args:       []Arg{},
					Table: postgres.Table{
						GoName:  "Bar",
						Name:    "bar",
						Columns: []postgres.Column{{GoName: "ID", GoType: "int", Name: "id", Type: "serial", NotNull: true}},
					},
					ReturnType: "bool",
				},
			}
			if diff := cmp.Diff(expected, result.Queries); diff != "" {
				t.Errorf("query mismatch: \n%s", diff)
			}
		})
	}
}
