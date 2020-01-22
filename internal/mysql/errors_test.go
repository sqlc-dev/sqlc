package mysql

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"vitess.io/vitess/go/vt/sqlparser"
)

func TestCustomArgErr(t *testing.T) {
	tests := [...]struct {
		input  string
		output sqlparser.PositionedErr
	}{
		{
			input: "/* name: GetUser :one */\nselect id, first_name from users where id = sqlc.argh(target_id)",
			output: sqlparser.PositionedErr{
				Err:  `invalid function call "sqlc.argh", did you mean "sqlc.arg"?`,
				Pos:  0,
				Near: nil,
			},
		},
		{
			input: "/* name: GetUser :one */\nselect id, first_name from users where id = sqlc.arg(sqlc.arg(target_id))",
			output: sqlparser.PositionedErr{
				Err:  `invalid custom argument value "sqlc.arg(sqlc.arg(target_id))"`,
				Pos:  0,
				Near: nil,
			},
		},
		{
			input: "/* name: GetUser :one */\nselect id, first_name from users where id = sqlc.arg(?)",
			output: sqlparser.PositionedErr{
				Err:  `invalid custom argument value "sqlc.arg(?)"`,
				Pos:  0,
				Near: nil,
			},
		},
	}

	for _, tcase := range tests {
		q, err := parseContents(mockFileName, tcase.input, mockSchema, mockSettings)
		if err == nil && len(q) > 0 {
			t.Errorf("parse contents succeeded on an invalid query")
		}
		if diff := cmp.Diff(tcase.output, err); diff != "" {
			t.Errorf(diff)
		}
	}
}

func TestPositionedErr(t *testing.T) {
	tests := [...]struct {
		input  string
		output sqlparser.PositionedErr
	}{
		{
			input: "/* name: GetUser :one */\nselect id, first_name from users from where id = ?",
			output: sqlparser.PositionedErr{
				Err:  `syntax error`,
				Pos:  63,
				Near: []byte("from"),
			},
		},
		{
			input: "/* name: GetUser :one */\nselectt id, first_name from users",
			output: sqlparser.PositionedErr{
				Err:  `syntax error`,
				Pos:  33,
				Near: []byte("selectt"),
			},
		},
		{
			input: "/* name: GetUser :one */\nselect id from users where select id",
			output: sqlparser.PositionedErr{
				Err:  `syntax error`,
				Pos:  59,
				Near: []byte("select"),
			},
		},
	}

	for _, tcase := range tests {
		q, err := parseContents(mockFileName, tcase.input, mockSchema, mockSettings)
		if err == nil && len(q) > 0 {
			t.Errorf("parse contents succeeded on an invalid query")
		}
		if diff := cmp.Diff(tcase.output, err); diff != "" {
			t.Errorf(diff)
		}
	}
}
