package mysql

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"vitess.io/vitess/go/vt/sqlparser"
)

func TestSyntaxErr(t *testing.T) {
	tokenizer := sqlparser.NewStringTokenizer("SELEC T id FROM users;")
	expectedLocation := 6
	expectedNear := "SELEC"

	_, parseErr := sqlparser.ParseNextStrictDDL(tokenizer)
	if parseErr == nil {
		t.Errorf("Tokenizer failed to error on invalid MySQL syntax")
	} else if posErr, ok := parseErr.(sqlparser.PositionedErr); ok {
		if posErr.Pos != expectedLocation {
			t.Errorf(cmp.Diff(posErr.Pos, expectedLocation))
		}
		if string(posErr.Near) != expectedNear {
			t.Errorf(cmp.Diff(string(posErr.Near), string(expectedNear)))
		}
	} else {
		t.Errorf("failed to return sqlparser.PositionedErr error for invalid mysql expression")
	}
}

func TestArgMessage(t *testing.T) {
	tcase := [...]struct {
		input  string
		output string
	}{
		{
			input:  "/* name: GetUser :one */\nselect id, first_name from users where id = sqlc.argh(target_id)",
			output: `invalid function call "sqlc.argh", did you mean "sqlc.arg"?`,
		},
		{
			input:  "/* name: GetUser :one */\nselect id, first_name from users where id = sqlc.arg(sqlc.arg(target_id))",
			output: `invalid custom argument value "sqlc.arg(sqlc.arg(target_id))"`,
		},
		{
			input:  "/* name: GetUser :one */\nselect id, first_name from users where id = sqlc.arg(?)",
			output: `invalid custom argument value "sqlc.arg(?)"`,
		},
	}

	for _, tc := range tcase {
		q, err := parseContents(mockFileName, tc.input, mockSchema, mockSettings)
		if err == nil && len(q) > 0 {
			t.Errorf("parse contents succeeded on an invalid query")
		}
		if diff := cmp.Diff(err.Error(), tc.output); diff != "" {
			t.Errorf(diff)
		}
	}
}

func TestPositionedErr(t *testing.T) {
	tcase := [...]struct {
		input  string
		output string
	}{
		{
			input:  "/* name: GetUser :one */\nselect id, first_name from users from where id = ?",
			output: `syntax error at or near 'from'`,
		},
		{
			input:  "/* name: GetUser :one */\nselectt id, first_name from users",
			output: `syntax error at or near 'selectt'`,
		},
		{
			input:  "/* name: GetUser :one */\nselect id from users where select id",
			output: `syntax error at or near 'select'`,
		},
	}

	for _, tc := range tcase {
		q, err := parseContents(mockFileName, tc.input, mockSchema, mockSettings)
		if err == nil && len(q) > 0 {
			t.Errorf("parse contents succeeded on an invalid query")
		}
		if diff := cmp.Diff(tc.output, err.Error()); diff != "" {
			t.Errorf(diff)
		}
	}
}
