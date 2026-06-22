package sqlite

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/source"
)

// Regression test for a bug where non-ASCII characters in the source SQL
// (e.g. an em-dash in a comment) caused the generated query text to be
// truncated. antlr reports token offsets as rune indices, but sqlc slices the
// query text using byte offsets; the two diverge after any multi-byte rune.
func TestParseNonASCIIOffsets(t *testing.T) {
	const query = "-- name: GetItem :one\n" +
		"-- Returns an item — looks up by id\n" +
		"SELECT id, name FROM items WHERE id = ?;\n"

	p := NewParser()
	stmts, err := p.Parse(strings.NewReader(query))
	if err != nil {
		t.Fatal(err)
	}
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}

	raw := stmts[0].Raw
	plucked, err := source.Pluck(query, raw.StmtLocation, raw.StmtLen)
	if err != nil {
		t.Fatal(err)
	}

	// The plucked text must include the trailing "?" placeholder. Before the
	// fix the byte/rune mismatch dropped the last few bytes of the statement.
	if !strings.Contains(plucked, "WHERE id = ?") {
		t.Errorf("plucked query was truncated:\n%q", plucked)
	}
}
