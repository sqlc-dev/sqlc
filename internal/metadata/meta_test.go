package metadata

import "testing"

func TestParseMetadata(t *testing.T) {

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
		if _, _, err := Parse(query, CommentSyntax{Dash: true}); err == nil {
			t.Errorf("expected invalid metadata: %q", query)
		}
	}

	for _, query := range []string{
		`-- some comment`,
		`-- name comment`,
		`--name comment`,
	} {
		if _, _, err := Parse(query, CommentSyntax{Dash: true}); err != nil {
			t.Errorf("expected valid comment: %q", query)
		}
	}

	query := `-- name: CreateFoo :one`
	queryName, queryType, err := Parse(query, CommentSyntax{Dash: true})
	if err != nil {
		t.Errorf("expected valid metadata: %q", query)
	}
	if queryName != "CreateFoo" {
		t.Errorf("incorrect queryName parsed: %q", query)
	}
	if queryType != CmdOne {
		t.Errorf("incorrect queryType parsed: %q", query)
	}

}
