package dolphin

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
                       CREATE TABLE foo (bar int);
                       CREATE TABLE foo (bar int);
                       `,
			sqlerr.RelationExists("foo"),
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

func TestSuccessfulUpdate(t *testing.T) {
	p := NewParser()
	for i, tc := range []struct {
		stmt string
	}{
		{
			`
                       CREATE TABLE authors (
                               id   INT PRIMARY KEY,
                               name text      NOT NULL,
                               bio  text      NOT NULL DEFAULT (bio_func())
                         );
                       `,
		},
		{
			`
			CREATE TABLE IF NOT EXISTS organizations
(
    id VARCHAR(36) DEFAULT (UUID()) NOT NULL PRIMARY KEY
);
			`,
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
			if err != nil {
				t.Log(test.stmt)
				t.Log(err)
				t.Fatal("err should have been nil")
			}
		})
	}
}
