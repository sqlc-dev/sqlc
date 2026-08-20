package ast_test

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/engine/postgresql"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func TestFormatSetClausePreservesMultiAssignments(t *testing.T) {
	const query = `UPDATE foo SET ("Value", other) = ($1, $2), simple = $3, (fourth, fifth) = ($4, $5)`

	parser := postgresql.NewParser()
	stmts, err := parser.Parse(strings.NewReader(query))
	if err != nil {
		t.Fatal(err)
	}
	formatted := ast.Format(stmts[0].Raw, parser)
	if !strings.Contains(formatted, `"Value"`) {
		t.Fatalf("formatted query lost quoted identifier: %s", formatted)
	}

	want, err := postgresql.Fingerprint(query)
	if err != nil {
		t.Fatal(err)
	}
	got, err := postgresql.Fingerprint(formatted)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("formatted query changed structure:\nwant fingerprint: %s\n got fingerprint: %s\nformatted query: %s", want, got, formatted)
	}
}
