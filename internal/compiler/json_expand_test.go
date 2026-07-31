package compiler

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/engine/postgresql"
	"github.com/sqlc-dev/sqlc/internal/source"
)

func TestExpandJSON(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "scalar with name",
			in:   `SELECT sqlc.jsonb_build_object."Obj"('x', a, 'y', b) AS obj FROM t`,
			want: `SELECT jsonb_build_object('x', a, 'y', b) AS obj FROM t`,
		},
		{
			// ARRAY(...) is left untouched; only the qualified call name
			// inside it is rewritten.
			name: "array with name leaves ARRAY(...) untouched",
			in:   `SELECT ARRAY(SELECT sqlc.jsonb_build_object."Child"('x', c.x, 'y', c.y) FROM c WHERE c.parent_id = p.id) AS children FROM p`,
			want: `SELECT ARRAY(SELECT jsonb_build_object('x', c.x, 'y', c.y) FROM c WHERE c.parent_id = p.id) AS children FROM p`,
		},
		{
			name: "name split across multiple lines",
			in:   "SELECT ARRAY(\n  SELECT sqlc.jsonb_build_object.\"Child\"(\n    'x', c.x\n  )\n  FROM c\n) AS children FROM p",
			want: "SELECT ARRAY(\n  SELECT jsonb_build_object(\n    'x', c.x\n  )\n  FROM c\n) AS children FROM p",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := postgresql.NewParser()
			stmts, err := parser.Parse(strings.NewReader(tt.in + ";"))
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			if len(stmts) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(stmts))
			}
			rawStmt := stmts[0].Raw

			c := &Compiler{}
			edits, err := c.expandJSON(rawStmt, tt.in)
			if err != nil {
				t.Fatalf("expandJSON error: %v", err)
			}
			got, err := source.Mutate(tt.in, edits)
			if err != nil {
				t.Fatalf("mutate error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got:\n  %s\nwant:\n  %s", got, tt.want)
			}
		})
	}
}
