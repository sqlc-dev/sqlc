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

const testCTECount = `
CREATE TABLE bar (ready bool not null);

-- name: CountAllAndReady :one
WITH all_count AS (
	SELECT count(*) FROM bar
), ready_count AS (
	SELECT count(*) FROM bar WHERE ready = true
)
SELECT all_count.count, ready_count.count
FROM all_count, ready_count;
`

func TestCTECount(t *testing.T) {
	result, err := parseSQL(testCTECount)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Query{
		{
			Type:       ":one",
			MethodName: "CountAllAndReady",
			StmtName:   "countAllAndReady",
			QueryName:  "countAllAndReady",
			SQL:        "-- name: CountAllAndReady :one\nWITH all_count AS (\n\tSELECT count(*) FROM bar\n), ready_count AS (\n\tSELECT count(*) FROM bar WHERE ready = true\n)\nSELECT all_count.count, ready_count.count\nFROM all_count, ready_count",
			Args:       []Arg{},
			RowStruct:  true,
			ScanRecord: true,
			ReturnType: "CountAllAndReadyRow",
			Fields: []Field{
				{Name: "AllCountCount", Type: "int"},
				{Name: "ReadyCountCount", Type: "int"},
			},
		},
	}

	if diff := cmp.Diff(expected, result.Queries); diff != "" {
		t.Errorf("query mismatch: \n%s", diff)
	}
}

const testCTEFilter = `
CREATE TABLE bar (ready bool not null);

-- name: CountReady :one
WITH filter_count AS (
	SELECT count(*) FROM bar WHERE ready = $1
)
SELECT filter_count.count
FROM filter_count;
`

func TestCTEFilter(t *testing.T) {
	result, err := parseSQL(testCTEFilter)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Query{
		{
			Type:       ":one",
			MethodName: "CountReady",
			StmtName:   "countReady",
			QueryName:  "countReady",
			SQL:        "-- name: CountReady :one\nWITH filter_count AS (\n\tSELECT count(*) FROM bar WHERE ready = $1\n)\nSELECT filter_count.count\nFROM filter_count",
			Args:       []Arg{{Name: "ready", Type: "bool"}},
			ReturnType: "int",
		},
	}
	if diff := cmp.Diff(expected, result.Queries); diff != "" {
		t.Errorf("query mismatch: \n%s", diff)
	}
}

const testInsertSelect = `
CREATE TABLE bar (name text not null, ready bool not null);
CREATE TABLE foo (name text not null, meta text not null);

-- name: CreateFoo :exec
INSERT INTO foo (
	name,
	meta
)
SELECT name, $1
FROM bar
WHERE ready = $2;
`

func TestInsertSelect(t *testing.T) {
	result, err := parseSQL(testInsertSelect)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Query{
		{
			Type:       ":exec",
			MethodName: "CreateFoo",
			StmtName:   "createFoo",
			QueryName:  "createFoo",
			SQL:        "-- name: CreateFoo :exec\nINSERT INTO foo (\n\tname,\n\tmeta\n)\nSELECT name, $1\nFROM bar\nWHERE ready = $2",
			Args:       []Arg{{Name: "meta", Type: "string"}, {Name: "ready", Type: "bool"}},
			Table: postgres.Table{
				GoName: "Foo",
				Name:   "foo",
				Columns: []postgres.Column{
					{GoName: "Name", GoType: "string", Name: "name", Type: "text", NotNull: true},
					{GoName: "Meta", GoType: "string", Name: "meta", Type: "text", NotNull: true},
				},
			},
		},
	}
	if diff := cmp.Diff(expected, result.Queries); diff != "" {
		t.Errorf("query mismatch: \n%s", diff)
	}
}

const testUpdateSet = `
CREATE TABLE foo (name text not null, slug text not null);

-- name: UpdateFoo :exec
UPDATE foo
SET name = $2
WHERE slug = $1;
`

func TestUpdateSet(t *testing.T) {
	result, err := parseSQL(testUpdateSet)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Query{
		{
			Type:       ":exec",
			MethodName: "UpdateFoo",
			StmtName:   "updateFoo",
			QueryName:  "updateFoo",
			SQL:        "-- name: UpdateFoo :exec\nUPDATE foo\nSET name = $2\nWHERE slug = $1",
			Args:       []Arg{{Name: "slug", Type: "string"}, {Name: "name", Type: "string"}},
			Table: postgres.Table{
				GoName: "Foo",
				Name:   "foo",
				Columns: []postgres.Column{
					{GoName: "Name", GoType: "string", Name: "name", Type: "text", NotNull: true},
					{GoName: "Slug", GoType: "string", Name: "slug", Type: "text", NotNull: true},
				},
			},
		},
	}
	if diff := cmp.Diff(expected, result.Queries); diff != "" {
		t.Errorf("query mismatch: \n%s", diff)
	}
}
