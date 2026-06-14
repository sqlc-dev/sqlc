package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/cmd"
)

type emitIterCase struct {
	name      string
	root      string
	mustHave  []string
	mustNot   []string
	exactIter int // -1 = don't check count
}

func TestEmitIteratorsCodegenMatrix(t *testing.T) {
	cases := []emitIterCase{
		{
			name: "global stdlib",
			root: filepath.Join("testdata", "emit_iterators", "stdlib"),
			mustHave: []string{
				"func (q *Queries) ListAuthors",
				"func (q *Queries) IterAuthors",
				"iter.Seq2[Author, error]",
				"defer rows.Close()",
			},
			mustNot:   []string{"sqlc:iterator-stream"},
			exactIter: 1,
		},
		{
			name: "explicit_only stream query",
			root: filepath.Join("testdata", "emit_iterators_explicit", "stdlib"),
			mustHave: []string{
				"func (q *Queries) StreamAuthors",
				"func (q *Queries) IterAuthors",
				"q.db.QueryContext(ctx, streamAuthors)",
			},
			mustNot:   []string{"sqlc:iterator-stream"},
			exactIter: 1,
		},
		{
			name: "explicit_only many:stream annotation",
			root: filepath.Join("testdata", "emit_iterators_many_stream", "stdlib"),
			mustHave: []string{
				"func (q *Queries) IterAuthors",
				"func (q *Queries) ListAllAuthors",
			},
			mustNot: []string{
				"func (q *Queries) IterAllAuthors",
				"sqlc:iterator-stream",
			},
			exactIter: 1,
		},
		{
			name: "parameterized many",
			root: filepath.Join("testdata", "emit_iterators_params", "stdlib"),
			mustHave: []string{
				"func (q *Queries) ListAuthorsByMinID(ctx context.Context, id int64)",
				"func (q *Queries) IterAuthorsByMinID(ctx context.Context, id int64) iter.Seq2[Author, error]",
			},
			exactIter: 1,
		},
		{
			name: "emit_iterators disabled",
			root: filepath.Join("testdata", "emit_iterators_off", "stdlib"),
			mustHave: []string{
				"func (q *Queries) ListAuthors",
			},
			mustNot: []string{
				"iter.Seq2",
				"func (q *Queries) Iter",
			},
			exactIter: 0,
		},
		{
			name: "pgx v5",
			root: filepath.Join("testdata", "emit_iterators", "pgx", "v5"),
			mustHave: []string{
				"func (q *Queries) IterAuthors(ctx context.Context) iter.Seq2[Author, error]",
				"rows, err := q.db.Query(ctx, listAuthors",
				"defer rows.Close()",
			},
			exactIter: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			abs, err := filepath.Abs(tc.root)
			if err != nil {
				t.Fatal(err)
			}

			var stderr strings.Builder
			if _, err := cmd.Generate(ctx, abs, "", &cmd.Options{Stderr: &stderr}); err != nil {
				t.Fatalf("generate failed: %v\n%s", err, stderr.String())
			}

			queryFile := filepath.Join(abs, "go", "query.sql.go")
			body, err := os.ReadFile(queryFile)
			if err != nil {
				t.Fatal(err)
			}
			src := string(body)

			for _, want := range tc.mustHave {
				if !strings.Contains(src, want) {
					t.Fatalf("missing %q in:\n%s", want, src)
				}
			}
			for _, bad := range tc.mustNot {
				if strings.Contains(src, bad) {
					t.Fatalf("unexpected %q in:\n%s", bad, src)
				}
			}
			if tc.exactIter >= 0 {
				count := strings.Count(src, "func (q *Queries) Iter")
				if count != tc.exactIter {
					t.Fatalf("Iter method count = %d, want %d\n%s", count, tc.exactIter, src)
				}
			}
		})
	}
}

func TestEmitIteratorsExplicitOnlySkipsPlainMany(t *testing.T) {
	root := filepath.Join("testdata", "emit_iterators_many_stream", "stdlib")
	abs, err := filepath.Abs(root)
	if err != nil {
		t.Fatal(err)
	}
	var stderr strings.Builder
	if _, err := cmd.Generate(context.Background(), abs, "", &cmd.Options{Stderr: &stderr}); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(filepath.Join(abs, "go", "query.sql.go"))
	if err != nil {
		t.Fatal(err)
	}
	src := string(body)
	if strings.Contains(src, "IterAllAuthors") {
		t.Fatal("plain :many must not get iterator in explicit_only mode")
	}
}
