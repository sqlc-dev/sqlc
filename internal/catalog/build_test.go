package catalog

import (
	"strconv"
	"testing"

	"github.com/kyleconroy/dinosql/internal/pg"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	query "github.com/lfittl/pg_query_go"
)

func buildCatalog(stmt string) (pg.Catalog, error) {
	c := pg.NewCatalog()
	tree, err := query.Parse(stmt)
	if err != nil {
		return c, err
	}
	for _, stmt := range tree.Statements {
		if err := Update(&c, stmt); err != nil {
			return c, err
		}
	}
	return c, nil
}

func TestUpdate(t *testing.T) {
	for i, tc := range []struct {
		stmt string
		c    pg.Catalog
	}{
		{
			"CREATE TYPE status AS ENUM ('open', 'closed');",
			pg.Catalog{
				Schemas: map[string]pg.Schema{
					"public": {
						Enums: map[string]pg.Enum{
							"status": pg.Enum{
								Name: "status",
								Vals: []string{"open", "closed"},
							},
						},
						Tables: map[string]pg.Table{},
					},
				},
			},
		},
		{
			"CREATE TABLE venues ();",
			pg.Catalog{
				Schemas: map[string]pg.Schema{
					"public": {
						Enums: map[string]pg.Enum{},
						Tables: map[string]pg.Table{
							"venues": pg.Table{
								Name: "venues",
							},
						},
					},
				},
			},
		},
		{
			`
			CREATE TABLE foo ();
			ALTER TABLE foo ADD COLUMN bar text;
			ALTER TABLE foo DROP COLUMN bar;
			`,
			pg.Catalog{
				Schemas: map[string]pg.Schema{
					"public": {
						Enums: map[string]pg.Enum{},
						Tables: map[string]pg.Table{
							"foo": pg.Table{
								Name: "foo",
							},
						},
					},
				},
			},
		},
		{
			`
			CREATE TABLE venues ();
			DROP TABLE venues;
			`,
			pg.NewCatalog(),
		},
		{
			`
			CREATE TYPE status AS ENUM ('open', 'closed');
			DROP TYPE status;
			`,
			pg.NewCatalog(),
		},
		{
			`
			CREATE TABLE venues ();
			DROP TABLE IF EXISTS venues;
			DROP TABLE IF EXISTS venues;
			`,
			pg.NewCatalog(),
		},
		{
			`
			CREATE TYPE status AS ENUM ('open', 'closed');
			DROP TYPE IF EXISTS status;
			DROP TYPE IF EXISTS status;
			`,
			pg.NewCatalog(),
		},
		{
			`
			CREATE TABLE venues ();
			DROP TABLE public.venues;
			`,
			pg.NewCatalog(),
		},
		{
			`
			CREATE TYPE status AS ENUM ('open', 'closed');
			DROP TYPE public.status;
			`,
			pg.NewCatalog(),
		},
		{
			`
			CREATE TABLE venues ();
			ALTER TABLE venues RENAME TO arenas;
			`,
			pg.Catalog{
				Schemas: map[string]pg.Schema{
					"public": {
						Enums: map[string]pg.Enum{},
						Tables: map[string]pg.Table{
							"arenas": pg.Table{
								Name: "arenas",
							},
						},
					},
				},
			},
		},
	} {
		test := tc
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(test.stmt)

			c, err := buildCatalog(test.stmt)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(test.c, c, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("catalog mismatch:\n%s", diff)
			}
		})
	}
}

func TestUpdateErrors(t *testing.T) {
	for i, tc := range []struct {
		stmt string
		err  pg.Error
	}{
		{
			`
			CREATE TABLE foo ();
			CREATE TABLE foo ();
			`,
			pg.Error{Code: "42P07", Message: "relation \"foo\" already exists"},
		},
		{
			`
			CREATE TYPE foo AS ENUM ('bar');
			CREATE TYPE foo AS ENUM ('bar');
			`,
			pg.Error{Code: "42710", Message: "type \"foo\" already exists"},
		},
		{
			`
			DROP TABLE foo;
			`,
			pg.Error{Code: "42P01", Message: "relation \"foo\" does not exist"},
		},
		{
			`
			DROP TYPE foo;
			`,
			pg.Error{Code: "42704", Message: "type \"foo\" does not exist"},
		},
		{
			`
			CREATE TABLE foo ();
			CREATE TABLE bar ();
			ALTER TABLE foo RENAME TO bar;
			`,
			pg.Error{Code: "42P07", Message: "relation \"bar\" already exists"},
		},
		{
			`
			ALTER TABLE foo RENAME TO bar;
			`,
			pg.Error{Code: "42P01", Message: "relation \"foo\" does not exist"},
		},
		{
			`
			CREATE TABLE foo ();
			ALTER TABLE foo ADD COLUMN bar text;
			ALTER TABLE foo ADD COLUMN bar text;
			`,
			pg.Error{Code: "42701", Message: "column \"bar\" of relation \"foo\" already exists"},
		},
		{
			`
			CREATE TABLE foo ();
			ALTER TABLE foo DROP COLUMN bar;
			`,
			pg.Error{Code: "42703", Message: "column \"bar\" of relation \"foo\" does not exist"},
		},
	} {
		test := tc
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(test.stmt)
			_, err := buildCatalog(test.stmt)
			if err == nil {
				t.Fatal("err was nil")
			}

			var actual pg.Error
			if err != nil {
				actual = err.(pg.Error)
			}

			if diff := cmp.Diff(test.err, actual, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("error mismatch: \n%s", diff)
			}
		})
	}
}
