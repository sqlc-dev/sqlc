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
			name: "array with name and paren in string literal",
			in:   `SELECT ARRAY(SELECT sqlc.jsonb_build_object."Child"('x', c.x, 'label', 'a(b)c') FROM c WHERE c.parent_id = p.id) AS children FROM p`,
			want: `SELECT ARRAY(SELECT jsonb_build_object('x', c.x, 'label', 'a(b)c') FROM c WHERE c.parent_id = p.id) AS children FROM p`,
		},
		{
			name: "two array embeds with names in one query",
			in:   `SELECT ARRAY(SELECT sqlc.jsonb_build_object."Child"('x', c.x) FROM c WHERE c.parent_id = p.id) AS children, ARRAY(SELECT sqlc.jsonb_build_object."Other"('y', d.y) FROM d WHERE d.parent_id = p.id) AS others FROM p`,
			want: `SELECT ARRAY(SELECT jsonb_build_object('x', c.x) FROM c WHERE c.parent_id = p.id) AS children, ARRAY(SELECT jsonb_build_object('y', d.y) FROM d WHERE d.parent_id = p.id) AS others FROM p`,
		},
		{
			name: "name split across multiple lines",
			in:   "SELECT ARRAY(\n  SELECT sqlc.jsonb_build_object.\"Child\"(\n    'x', c.x\n  )\n  FROM c\n) AS children FROM p",
			want: "SELECT ARRAY(\n  SELECT jsonb_build_object(\n    'x', c.x\n  )\n  FROM c\n) AS children FROM p",
		},
		{
			name: "unquoted name is lowercased by Postgres, still rewrites",
			in:   `SELECT sqlc.jsonb_build_object.Obj('x', a, 'y', b) AS obj FROM t`,
			want: `SELECT jsonb_build_object('x', a, 'y', b) AS obj FROM t`,
		},
		{
			name: "embed.jsonb scalar rewrites to to_jsonb",
			in:   `SELECT sqlc.embed.jsonb(t) AS obj FROM t`,
			want: `SELECT to_jsonb(t) AS obj FROM t`,
		},
		{
			name: "embed.jsonb inside ARRAY(...) leaves the wrapper untouched",
			in:   `SELECT ARRAY(SELECT sqlc.embed.jsonb(c) FROM c WHERE c.parent_id = p.id) AS children FROM p`,
			want: `SELECT ARRAY(SELECT to_jsonb(c) FROM c WHERE c.parent_id = p.id) AS children FROM p`,
		},
		{
			name: "mix of jsonb_build_object and embed.jsonb in one query",
			in:   `SELECT sqlc.jsonb_build_object."Obj"('x', a) AS obj, sqlc.embed.jsonb(t) AS row FROM t`,
			want: `SELECT jsonb_build_object('x', a) AS obj, to_jsonb(t) AS row FROM t`,
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
