package sqlite

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/source"
)

// TestParseNonASCIIComment verifies that non-ASCII characters in SQL comments
// do not corrupt the plucked query text.
//
// ANTLR4 stores the input as []rune so all token positions are rune indices,
// not byte offsets. source.Pluck (and the rest of the pipeline) treats
// StmtLocation/StmtLen as byte offsets. For multi-byte UTF-8 characters the
// two differ, which previously caused the plucked query to be truncated by one
// byte per extra byte in each non-ASCII character.
func TestParseNonASCIIComment(t *testing.T) {
	p := NewParser()

	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "2-byte char (U+00DC Ü) in dash comment",
			sql:  "-- name: GetUser :one\n-- Ünïcode comment\nSELECT id FROM users WHERE id = ?",
		},
		{
			name: "3-byte char (U+2665 ♥) in dash comment",
			sql:  "-- name: GetUser :one\n-- ♥ love\nSELECT id FROM users WHERE id = ?",
		},
		{
			name: "4-byte char (U+1D11E 𝄞) in dash comment",
			sql:  "-- name: GetUser :one\n-- 𝄞 music\nSELECT id FROM users WHERE id = ?",
		},
		{
			name: "multiple non-ASCII chars in comment",
			sql:  "-- name: GetUser :one\n-- héllo wörld\nSELECT id FROM users WHERE id = ?",
		},
		{
			name: "non-ASCII only in first of two statements",
			sql:  "-- name: Q1 :one\n-- Ü\nSELECT 1;\n\n-- name: Q2 :one\nSELECT 2",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stmts, err := p.Parse(strings.NewReader(tc.sql))
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			if len(stmts) == 0 {
				t.Fatal("expected at least one statement")
			}

			// For every parsed statement, verify that the plucked text is a
			// valid substring of the original SQL (not truncated mid-character).
			for i, stmt := range stmts {
				raw := stmt.Raw
				plucked, err := source.Pluck(tc.sql, raw.StmtLocation, raw.StmtLen)
				if err != nil {
					t.Fatalf("stmt %d: Pluck error: %v", i, err)
				}
				if !strings.Contains(tc.sql, plucked) {
					t.Errorf("stmt %d: plucked text is not a substring of the input\ngot:  %q\ninput: %q", i, plucked, tc.sql)
				}
				if plucked == "" {
					t.Errorf("stmt %d: plucked text is empty", i)
				}
			}

			// For the single-statement cases the plucked text must equal the
			// full input, since there is exactly one statement and no trailing
			// semicolon to exclude.
			if len(stmts) == 1 {
				raw := stmts[0].Raw
				plucked, _ := source.Pluck(tc.sql, raw.StmtLocation, raw.StmtLen)
				if plucked != tc.sql {
					t.Errorf("single-statement pluck mismatch\ngot:  %q\nwant: %q", plucked, tc.sql)
				}
			}

			// For the two-statement case, verify each statement contains its
			// expected SELECT.
			if len(stmts) == 2 {
				for i, want := range []string{"SELECT 1", "SELECT 2"} {
					raw := stmts[i].Raw
					plucked, _ := source.Pluck(tc.sql, raw.StmtLocation, raw.StmtLen)
					if !strings.Contains(plucked, want) {
						t.Errorf("stmt %d: plucked text %q does not contain %q", i, plucked, want)
					}
				}
			}
		})
	}
}
