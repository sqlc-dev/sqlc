package postgresql

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
	"github.com/kyleconroy/sqlc/internal/sql/sqlerr"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestUpdate(t *testing.T) {
	p := NewParser()

	for _, tc := range []struct {
		name string
		stmt string
		s    *catalog.Schema
	}{
		{
			"create-enum",
			"CREATE TYPE status AS ENUM ('open', 'closed');",
			&catalog.Schema{
				Name: "public",
				Types: []catalog.Type{
					&catalog.Enum{
						Name: "status",
						Vals: []string{"open", "closed"},
					},
				},
			},
		},
		{
			"alter-type-rename-value",
			`
			CREATE TYPE status AS ENUM ('open', 'closed');
			ALTER TYPE status RENAME VALUE 'closed' TO 'shut';
			`,
			&catalog.Schema{
				Name: "public",
				Types: []catalog.Type{
					&catalog.Enum{
						Name: "status",
						Vals: []string{"open", "shut"},
					},
				},
			},
		},
		{
			"alter-type-add-value",
			`
			CREATE TYPE status AS ENUM ('open', 'closed');
			ALTER TYPE status ADD VALUE 'unknown';
			ALTER TYPE status ADD VALUE IF NOT EXISTS 'unknown';
			`,
			&catalog.Schema{
				Name: "public",
				Types: []catalog.Type{
					&catalog.Enum{
						Name: "status",
						Vals: []string{"open", "closed", "unknown"},
					},
				},
			},
		},
		{
			"create-table",
			"CREATE TABLE venues ();",
			&catalog.Schema{
				Name: "public",
				Tables: []*catalog.Table{
					{
						Rel: &ast.TableName{Name: "venues"},
					},
				},
			},
		},
		{
			"alter-table-drop-column",
			`
			CREATE TABLE foo ();
			ALTER TABLE foo ADD COLUMN bar text;
			ALTER TABLE foo DROP COLUMN bar;
			`,
			&catalog.Schema{
				Name: "public",
				Tables: []*catalog.Table{
					{
						Rel: &ast.TableName{Name: "foo"},
					},
				},
			},
		},
		{
			"alter-table-drop-column-if-exists",
			`
			CREATE TABLE foo ();
			ALTER TABLE foo DROP COLUMN IF EXISTS bar;
			`,
			&catalog.Schema{
				Name: "public",
				Tables: []*catalog.Table{
					{
						Rel: &ast.TableName{Name: "foo"},
					},
				},
			},
		},
		{
			"alter-table-set-not-null",
			`
			CREATE TABLE foo (bar text);
			ALTER TABLE foo ALTER bar SET NOT NULL;
			`,
			&catalog.Schema{
				Name: "public",
				Tables: []*catalog.Table{
					{
						Rel: &ast.TableName{Name: "foo"},
						Columns: []*catalog.Column{
							{
								Name:      "bar",
								Type:      ast.TypeName{Name: "text"},
								IsNotNull: true,
							},
						},
					},
				},
			},
		},
		/*
			{
				`
				CREATE TABLE foo (bar text[] not null);
				`,
				pg.Catalog{
					Schemas: map[string]pg.Schema{
						"public": {
							Tables: map[string]pg.Table{
								"foo": pg.Table{
									Name: "foo",
									Columns: []pg.Column{
										{Name: "bar", DataType: "text", IsArray: true, NotNull: true, Table: pg.FQN{Schema: "public", Rel: "foo"}},
									},
								},
							},
						},
					},
				},
			},
		*/
		{
			"alter-table-drop-not-null",
			`
			CREATE TABLE foo (bar text NOT NULL);
			ALTER TABLE foo ALTER bar DROP NOT NULL;
			`,
			&catalog.Schema{
				Name: "public",
				Tables: []*catalog.Table{
					{
						Rel: &ast.TableName{Name: "foo"},
						Columns: []*catalog.Column{
							{
								Name: "bar",
								Type: ast.TypeName{Name: "text"},
							},
						},
					},
				},
			},
		},
		{
			"alter-table-column-drop-not-null",
			`
			CREATE TABLE foo (bar text NOT NULL);
			ALTER TABLE foo ALTER COLUMN bar DROP NOT NULL;
			`,
			&catalog.Schema{
				Name: "public",
				Tables: []*catalog.Table{
					{
						Rel: &ast.TableName{Name: "foo"},
						Columns: []*catalog.Column{
							{
								Name: "bar",
								Type: ast.TypeName{Name: "text"},
							},
						},
					},
				},
			},
		},
		{
			"alter-table-rename-column",
			`
			CREATE TABLE foo (bar text);
			ALTER TABLE foo RENAME bar TO baz;
			`,
			&catalog.Schema{
				Name: "public",
				Tables: []*catalog.Table{
					{
						Rel: &ast.TableName{Name: "foo"},
						Columns: []*catalog.Column{
							{
								Name: "baz",
								Type: ast.TypeName{Name: "text"},
							},
						},
					},
				},
			},
		},
		{
			"alter-table-set-data-type",
			`
			CREATE TABLE foo (bar text);
			ALTER TABLE foo ALTER bar SET DATA TYPE bool;
			`,
			&catalog.Schema{
				Name: "public",
				Tables: []*catalog.Table{
					{
						Rel: &ast.TableName{Name: "foo"},
						Columns: []*catalog.Column{
							{
								Name: "bar",
								Type: ast.TypeName{Name: "bool"},
							},
						},
					},
				},
			},
		},
		{
			"alter-table-set-schema",
			`
			CREATE SCHEMA foo;
			CREATE TABLE bar ();
			ALTER TABLE bar SET SCHEMA foo;
			`,
			&catalog.Schema{
				Name: "foo",
				Tables: []*catalog.Table{
					{
						Rel: &ast.TableName{Name: "bar"},
					},
				},
			},
		},
		{
			"drop-table",
			`
			CREATE TABLE venues ();
			DROP TABLE venues;
			`,
			nil,
		},
		{
			"drop-table-if-exists",
			`
			CREATE TABLE venues ();
			DROP TABLE IF EXISTS venues;
			DROP TABLE IF EXISTS venues;
			`,
			nil,
		},
		{
			"alter-table-rename",
			`
			CREATE TABLE venues ();
			ALTER TABLE venues RENAME TO arenas;
			`,
			&catalog.Schema{
				Name: "public",
				Tables: []*catalog.Table{
					{
						Rel: &ast.TableName{Name: "arenas"},
					},
				},
			},
		},
		{
			"drop-type",
			`
			CREATE TYPE status AS ENUM ('open', 'closed');
			DROP TYPE status;
			`,
			nil,
		},
		{
			"drop-type-if-exists",
			`
			CREATE TYPE status AS ENUM ('open', 'closed');
			DROP TYPE IF EXISTS status;
			DROP TYPE IF EXISTS status;
			`,
			nil,
		},
		{
			"drop-table-in-schema",
			`
			CREATE TABLE venues ();
			DROP TABLE public.venues;
			`,
			nil,
		},
		{
			"drop-type-in-schema",
			`
			CREATE TYPE status AS ENUM ('open', 'closed');
			DROP TYPE public.status;
			`,
			nil,
		},
		{
			"drop-schema",
			`
			CREATE SCHEMA foo;
			DROP SCHEMA foo;
			`,
			nil,
		},
		{
			"drop-schema-if-exists",
			`
			DROP SCHEMA IF EXISTS foo;
			`,
			nil,
		},
		{
			"drop-function-if-exists",
			`
			DROP FUNCTION IF EXISTS bar;
			DROP FUNCTION IF EXISTS bar();
			`,
			nil,
		},
		{
			"alter-table-drop-constraint",
			`
			CREATE TABLE venues (id SERIAL PRIMARY KEY);
			ALTER TABLE venues DROP CONSTRAINT venues_id_pkey;
			`,
			&catalog.Schema{
				Name: "public",
				Tables: []*catalog.Table{
					{
						Rel: &ast.TableName{Name: "venues"},
						Columns: []*catalog.Column{
							{
								Name:      "id",
								Type:      ast.TypeName{Name: "serial"},
								IsNotNull: true,
							},
						},
					},
				},
			},
		},
		{
			"create-function",
			`
			CREATE FUNCTION foo(TEXT) RETURNS bool AS $$ SELECT true $$ LANGUAGE sql;
			`,
			&catalog.Schema{
				Name: "public",
				Funcs: []*catalog.Function{
					{
						Name: "foo",
						Args: []*catalog.Argument{
							{
								Type: &ast.TypeName{Name: "text"},
							},
						},
						ReturnType: &ast.TypeName{Name: "bool"},
					},
				},
			},
		},
		{
			"create-function-args",
			`
			CREATE FUNCTION foo(bar TEXT) RETURNS bool AS $$ SELECT true $$ LANGUAGE sql;
			CREATE FUNCTION foo(bar TEXT, baz TEXT) RETURNS TEXT AS $$ SELECT "baz" $$ LANGUAGE sql;
			`,
			&catalog.Schema{
				Name: "public",
				Funcs: []*catalog.Function{
					{
						Name: "foo",
						Args: []*catalog.Argument{
							{
								Name: "bar",
								Type: &ast.TypeName{Name: "text"},
							},
						},
						ReturnType: &ast.TypeName{Name: "bool"},
					},
					{
						Name: "foo",
						Args: []*catalog.Argument{
							{
								Name: "bar",
								Type: &ast.TypeName{Name: "text"},
							},
							{
								Name: "baz",
								Type: &ast.TypeName{Name: "text"},
							},
						},
						ReturnType: &ast.TypeName{Name: "text"},
					},
				},
			},
		},
		{
			"create-function-types",
			`
			CREATE FUNCTION foo(bar TEXT) RETURNS bool AS $$ SELECT true $$ LANGUAGE sql;
			CREATE FUNCTION foo(bar INTEGER) RETURNS TEXT AS $$ SELECT "baz" $$ LANGUAGE sql;
			`,
			&catalog.Schema{
				Name: "public",
				Funcs: []*catalog.Function{
					{
						Name: "foo",
						Args: []*catalog.Argument{
							{
								Name: "bar",
								Type: &ast.TypeName{Name: "text"},
							},
						},
						ReturnType: &ast.TypeName{Name: "bool"},
					},
					{
						Name: "foo",
						Args: []*catalog.Argument{
							{
								Name: "bar",
								Type: &ast.TypeName{Schema: "pg_catalog", Name: "int4"},
							},
						},
						ReturnType: &ast.TypeName{Name: "text"},
					},
				},
			},
		},
		{
			"create-function-return",
			`
			CREATE FUNCTION foo(bar TEXT, baz TEXT="baz") RETURNS bool AS $$ SELECT true $$ LANGUAGE sql;
			`,
			&catalog.Schema{
				Name: "public",
				Funcs: []*catalog.Function{
					{
						Name: "foo",
						Args: []*catalog.Argument{
							{
								Name: "bar",
								Type: &ast.TypeName{Name: "text"},
							},
							{
								Name:       "baz",
								Type:       &ast.TypeName{Name: "text"},
								HasDefault: true,
							},
						},
						ReturnType: &ast.TypeName{Name: "bool"},
					},
				},
			},
		},
		{
			"drop-function-args",
			`
			CREATE FUNCTION foo(bar text) RETURNS bool AS $$ SELECT true $$ LANGUAGE sql;
			DROP FUNCTION foo(text);
			`,
			nil,
		},
		{
			"drop-function",
			`
			CREATE FUNCTION foo(bar text) RETURNS bool AS $$ SELECT true $$ LANGUAGE sql;
			DROP FUNCTION foo;
			`,
			nil,
		},
		// CREATE FUNCTION foo(bar text) RETURNS bool AS $$ SELECT true $$ LANGUAGE sql;
		// CREATE FUNCTION foo() RETURNS bool AS $$ SELECT true $$ LANGUAGE sql;
		// DROP FUNCTION foo -- FAIL
		{
			"pg_temp",
			`
			CREATE TABLE pg_temp.migrate (val SERIAL);
			INSERT INTO pg_temp.migrate (val) SELECT val FROM old;
			INSERT INTO new (val) SELECT val FROM pg_temp.migrate;
			`,
			&catalog.Schema{
				Name: "pg_temp",
				Tables: []*catalog.Table{
					{
						Rel: &ast.TableName{Schema: "pg_temp", Name: "migrate"},
						Columns: []*catalog.Column{
							{
								Name: "val",
								Type: ast.TypeName{Name: "serial"},
							},
						},
					},
				},
			},
		},
		{
			"comment",
			`
			CREATE SCHEMA foo;
			CREATE TABLE foo.bar (baz text);
			CREATE TYPE foo.bat AS ENUM ('bat');
			COMMENT ON SCHEMA foo IS 'Schema comment';
			COMMENT ON TABLE foo.bar IS 'Table comment';
			COMMENT ON COLUMN foo.bar.baz IS 'Column comment';
			COMMENT ON TYPE foo.bat IS 'Enum comment';
			`,
			&catalog.Schema{
				Name:    "foo",
				Comment: "Schema comment",
				Tables: []*catalog.Table{
					{
						Rel:     &ast.TableName{Schema: "foo", Name: "bar"},
						Comment: "Table comment",
						Columns: []*catalog.Column{
							{
								Name:    "baz",
								Type:    ast.TypeName{Name: "text"},
								Comment: "Column comment",
							},
						},
					},
				},
				Types: []catalog.Type{
					&catalog.Enum{
						Name:    "bat",
						Vals:    []string{"bat"},
						Comment: "Enum comment",
					},
				},
			},
		},
	} {
		test := tc
		t.Run(test.name, func(t *testing.T) {
			stmts, err := p.Parse(strings.NewReader(test.stmt))
			if err != nil {
				t.Log(test.stmt)
				t.Fatal(err)
			}

			c := newTestCatalog()
			if err := c.Build(stmts); err != nil {
				t.Log(test.stmt)
				t.Fatal(err)
			}

			e := newTestCatalog()
			if test.s != nil {
				var replaced bool
				for i := range e.Schemas {
					if e.Schemas[i].Name == test.s.Name {
						e.Schemas[i] = test.s
						replaced = true
						break
					}
				}
				if !replaced {
					e.Schemas = append(e.Schemas, test.s)
				}
			}

			if diff := cmp.Diff(e.Schemas, c.Schemas, cmpopts.EquateEmpty()); diff != "" {
				t.Log(test.stmt)
				t.Errorf("catalog mismatch:\n%s", diff)
			}
		})
	}
}

func TestUpdateErrors(t *testing.T) {
	p := NewParser()

	for i, tc := range []struct {
		stmt string
		err  *sqlerr.Error
	}{
		{
			`
			CREATE TABLE foo ();
			CREATE TABLE foo ();
			`,
			sqlerr.RelationExists("foo"),
		},
		{
			`
			CREATE TYPE foo AS ENUM ('bar');
			CREATE TYPE foo AS ENUM ('bar');
			`,
			sqlerr.TypeExists("foo"),
		},
		{
			`
			DROP TABLE foo;
			`,
			sqlerr.RelationNotFound("foo"),
		},
		{
			`
			DROP TYPE foo;
			`,
			sqlerr.TypeNotFound("foo"),
		},
		{
			`
			CREATE TABLE foo ();
			CREATE TABLE bar ();
			ALTER TABLE foo RENAME TO bar;
			`,
			sqlerr.RelationExists("bar"),
		},
		{
			`
			ALTER TABLE foo RENAME TO bar;
			`,
			sqlerr.RelationNotFound("foo"),
		},
		{
			`
			CREATE TABLE foo ();
			ALTER TABLE foo ADD COLUMN bar text;
			ALTER TABLE foo ADD COLUMN bar text;
			`,
			sqlerr.ColumnExists("foo", "bar"),
		},
		{
			`
			CREATE TABLE foo ();
			ALTER TABLE foo DROP COLUMN bar;
			`,
			sqlerr.ColumnNotFound("foo", "bar"),
		},
		{
			`
			CREATE TABLE foo ();
			ALTER TABLE foo ALTER COLUMN bar SET NOT NULL;
			`,
			sqlerr.ColumnNotFound("foo", "bar"),
		},
		{
			`
			CREATE TABLE foo ();
			ALTER TABLE foo ALTER COLUMN bar DROP NOT NULL;
			`,
			sqlerr.ColumnNotFound("foo", "bar"),
		},
		{
			`
			CREATE SCHEMA foo;
			CREATE SCHEMA foo;
			`,
			sqlerr.SchemaExists("foo"),
		},
		{
			`
			ALTER TABLE foo.baz SET SCHEMA bar;
			`,
			sqlerr.SchemaNotFound("foo"),
		},
		{
			`
			CREATE SCHEMA foo;
			ALTER TABLE foo.baz SET SCHEMA bar;
			`,
			sqlerr.RelationNotFound("baz"),
		},
		{
			`
			CREATE SCHEMA foo;
			CREATE TABLE foo.baz ();
			ALTER TABLE foo.baz SET SCHEMA bar;
			`,
			sqlerr.SchemaNotFound("bar"),
		},
		{
			`
			DROP SCHEMA bar;
			`,
			sqlerr.SchemaNotFound("bar"),
		},
		{
			`
			ALTER TABLE foo RENAME bar TO baz;
			`,
			sqlerr.RelationNotFound("foo"),
		},
		{
			`
			CREATE TABLE foo ();
			ALTER TABLE foo RENAME bar TO baz;
			`,
			sqlerr.ColumnNotFound("foo", "bar"),
		},
		{
			`
			CREATE TABLE foo (bar text, baz text);
			ALTER TABLE foo RENAME bar TO baz;
			`,
			sqlerr.ColumnExists("foo", "baz"),
		},
	} {
		test := tc
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			stmts, err := p.Parse(strings.NewReader(test.stmt))
			if err != nil {
				t.Log(test.stmt)
				t.Fatal(err)
			}

			c := NewCatalog()
			err = c.Build(stmts)
			if err == nil {
				t.Log(test.stmt)
				t.Fatal("err was nil")
			}

			var actual *sqlerr.Error
			if !errors.As(err, &actual) {
				t.Fatalf("err is not *sqlerr.Error: %#v", err)
			}

			if diff := cmp.Diff(test.err.Error(), actual.Error()); diff != "" {
				t.Log(test.stmt)
				t.Errorf("error mismatch: \n%s", diff)
			}
		})
	}
}
