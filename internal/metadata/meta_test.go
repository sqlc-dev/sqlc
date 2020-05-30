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
		`/* name: CreateFoo, :one */`,
		`/* name: 9Foo_, :one */`,
		`/* name: CreateFoo :two */`,
		`/* name: CreateFoo */`,
		`/* name: CreateFoo :one something */`,
		`/* name: */`,
	} {
		if _, _, err := Parse(query); err == nil {
			t.Errorf("expected invalid metadata: %q", query)
		}
	}
}
