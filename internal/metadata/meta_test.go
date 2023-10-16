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
	} {
		if _, _, err := ParseQueryNameAndType(query, CommentSyntax{Dash: true}); err == nil {
			t.Errorf("expected invalid metadata: %q", query)
		}
	}

	for _, query := range []string{
		`-- some comment`,
		`-- name comment`,
		`--name comment`,
	} {
		if _, _, err := ParseQueryNameAndType(query, CommentSyntax{Dash: true}); err != nil {
			t.Errorf("expected valid comment: %q", query)
		}
	}

	for query, cs := range map[string]CommentSyntax{
		`-- name: CreateFoo :one`:    {Dash: true},
		`# name: CreateFoo :one`:     {Hash: true},
		`/* name: CreateFoo :one */`: {SlashStar: true},
	} {
		queryName, queryCmd, err := ParseQueryNameAndType(query, cs)
		if err != nil {
			t.Errorf("expected valid metadata: %q", query)
		}
		if queryName != "CreateFoo" {
			t.Errorf("incorrect queryName parsed: (%q) %q", queryName, query)
		}
		if queryCmd != CmdOne {
			t.Errorf("incorrect queryCmd parsed: (%q) %q", queryCmd, query)
		}
	}

}

func TestParseQueryParams(t *testing.T) {
	for _, comments := range [][]string{
		{
			" name: CreateFoo :one",
			" @param foo_id UUID",
		},
		{
			" name: CreateFoo :one ",
			" @param foo_id UUID ",
		},
		{
			" name: CreateFoo :one",
			"@param foo_id UUID",
			" invalid",
		},
		{
			" name: CreateFoo :one",
			" @invalid",
			" @param foo_id UUID",
		},
		{
			" name: GetFoos :many ",
			" @param foo_id UUID ",
			" @param @invalid UUID ",
		},
	} {
		params, _, err := ParseParamsAndFlags(comments)
		if err != nil {
			t.Errorf("expected comments to parse, got err: %s", err)
		}

		pt, ok := params["foo_id"]
		if !ok {
			t.Errorf("expected param not found")
		}

		if pt != "UUID" {
			t.Error("unexpected param metadata:", pt)
		}

		_, ok = params["invalid"]
		if ok {
			t.Errorf("unexpected param found")
		}
	}
}

func TestParseQueryFlags(t *testing.T) {
	for _, comments := range [][]string{
		{
			" name: CreateFoo :one",
			" @flag-foo",
		},
		{
			" name: CreateFoo :one ",
			"@flag-foo ",
		},
		{
			" name: CreateFoo :one",
			" @flag-foo @flag-bar",
		},
		{
			" name: GetFoos :many",
			" @param @flag-bar UUID",
			" @flag-foo",
		},
		{
			" name: GetFoos :many",
			" @flag-foo",
			" @param @flag-bar UUID",
		},
	} {
		_, flags, err := ParseParamsAndFlags(comments)
		if err != nil {
			t.Errorf("expected comments to parse, got err: %s", err)
		}

		if !flags["@flag-foo"] {
			t.Errorf("expected flag not found")
		}

		if flags["@flag-bar"] {
			t.Errorf("unexpected flag found")
		}
	}
}
