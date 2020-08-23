package postgresql

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/kyleconroy/sqlc/internal/sql/sqlerr"

	"github.com/google/go-cmp/cmp"
)

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
