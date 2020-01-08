package dinosql

import (
	"path/filepath"
	"testing"

	"github.com/kyleconroy/sqlc/internal/pg"

	"github.com/google/go-cmp/cmp"
)

func TestFuncs(t *testing.T) {
	_, err := ParseQueries(
		pg.NewCatalog(),
		PackageSettings{
			Queries: filepath.Join("testdata", "funcs"),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

}

func TestParserErrors(t *testing.T) {
	for _, tc := range []struct {
		query string
		err   pg.Error
	}{
		{
			"SELECT foo FROM bar WHERE baz = $4;",
			pg.Error{Code: "42P18", Message: "could not determine data type of parameter $1"},
		},
		{
			"SELECT foo FROM bar WHERE baz = $1 AND baz = $3;",
			pg.Error{Code: "42P18", Message: "could not determine data type of parameter $2"},
		},
		{
			`
			CREATE TABLE bar (id serial not null);

			-- name: foo :one
			SELECT foo FROM bar;
			`,
			pg.Error{
				Code:     "42703",
				Message:  "column \"foo\" does not exist",
				Location: 75,
			},
		},
		{
			"SELECT random(1);",
			pg.Error{
				Code:     "42883",
				Message:  "function random(unknown) does not exist",
				Hint:     "No function matches the given name and argument types. You might need to add explicit type casts.",
				Location: 7,
			},
		},
		{
			"SELECT position()",
			pg.Error{
				Code:     "42883",
				Message:  "function position() does not exist",
				Hint:     "No function matches the given name and argument types. You might need to add explicit type casts.",
				Location: 7,
			},
		},
	} {
		test := tc
		t.Run(test.query, func(t *testing.T) {
			_, err := parseSQL(test.query)

			var actual pg.Error
			if err != nil {
				actual = err.(pg.Error)
			}

			if diff := cmp.Diff(test.err, actual); diff != "" {
				t.Errorf("error mismatch: \n%s", diff)
			}
		})
	}
}
