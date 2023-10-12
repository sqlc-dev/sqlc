package metadata

import "testing"

func TestParseQueryNameAndType(t *testing.T) {

	for _, query := range []string{
		`-- name: CreateFoo, :one`,
		`-- name: 9Foo_, :one`,
		`-- name: CreateFoo :two`,
		`-- name: CreateFoo`,
		`-- name: CreateFoo :one something`,
		`-- name: `,
		`--name: CreateFoo :one`,
		`--name CreateFoo :one`,
		`--name: CreateFoo :two`,
		"-- name:CreateFoo",
		`--name:CreateFoo :two`,
		`--  name: CreateFoo :two`,
		`-- name:  CreateFoo :two`,
		`-- name: CreateFoo  :two`,
	} {
		if _, err := ParseQueryMetadata(query, CommentSyntax{Dash: true}); err == nil {
			t.Errorf("expected invalid metadata: %q", query)
		}
	}

	for _, query := range []string{
		`-- some comment`,
		`-- name comment`,
		`--name comment`,
	} {
		if _, err := ParseQueryMetadata(query, CommentSyntax{Dash: true}); err != nil {
			t.Errorf("expected valid comment: %q", query)
		}
	}

	for query, cs := range map[string]CommentSyntax{
		`-- name: CreateFoo :one`:    {Dash: true},
		`# name: CreateFoo :one`:     {Hash: true},
		`/* name: CreateFoo :one */`: {SlashStar: true},
	} {
		queryMetadata, err := ParseQueryMetadata(query, cs)
		if err != nil {
			t.Errorf("expected valid metadata: %q", query)
		}
		if queryMetadata.Name != "CreateFoo" {
			t.Errorf("incorrect queryName parsed: %q", query)
		}
		if queryMetadata.Cmd != CmdOne {
			t.Errorf("incorrect queryType parsed: %q", query)
		}
	}

}

func TestParseQueryFlags(t *testing.T) {
	for _, comments := range [][]string{
		{
			"-- name: CreateFoo :one",
			"-- @flag-foo",
		},
		{
			"-- name: CreateFoo :one",
			"-- @flag-foo @flag-bar",
		},
		{
			"-- name: GetFoos :many",
			"-- @param @flag-bar UUID",
			"-- @flag-foo",
		},
	} {
		_, flags := parseParamsAndFlags(comments)

		if !flags["@flag-foo"] {
			t.Errorf("expected flag not found")
		}

		if flags["@flag-bar"] {
			t.Errorf("unexpected flag found")
		}
	}
}
