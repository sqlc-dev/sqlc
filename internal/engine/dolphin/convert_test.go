package dolphin

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func TestConvertNotExistsSubquery(t *testing.T) {
	const query = `UPDATE target_table
SET active = 0
WHERE NOT (
    EXISTS (
        SELECT 1
        FROM source_table s
        WHERE s.updated_at >= sqlc.arg(updated_after)
    )
)`

	parser := NewParser()
	stmts, err := parser.Parse(strings.NewReader(query))
	if err != nil {
		t.Fatal(err)
	}
	if len(stmts) != 1 {
		t.Fatalf("expected one statement, got %d", len(stmts))
	}

	formatted := ast.Format(stmts[0].Raw, parser)
	if !strings.Contains(formatted, "NOT EXISTS") {
		t.Errorf("expected NOT EXISTS to be preserved, got:\n%s", formatted)
	}
	if !strings.Contains(formatted, "sqlc.arg(updated_after)") {
		t.Errorf("expected named parameter to be preserved, got:\n%s", formatted)
	}
}

func TestConvertNotExistsSubqueryWithoutParentheses(t *testing.T) {
	const query = `SELECT 1
WHERE NOT EXISTS (
    SELECT 1
    FROM source_table s
    WHERE s.updated_at >= sqlc.arg(updated_after)
)`

	parser := NewParser()
	stmts, err := parser.Parse(strings.NewReader(query))
	if err != nil {
		t.Fatal(err)
	}
	if len(stmts) != 1 {
		t.Fatalf("expected one statement, got %d", len(stmts))
	}

	formatted := ast.Format(stmts[0].Raw, parser)
	if !strings.Contains(formatted, "NOT EXISTS") {
		t.Errorf("expected NOT EXISTS to be preserved, got:\n%s", formatted)
	}
}
