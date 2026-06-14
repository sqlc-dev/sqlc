package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/cmd"
)

func TestEmitIteratorsCodegenGlobal(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	root := filepath.Join("testdata", "emit_iterators", "stdlib")
	abs, err := filepath.Abs(root)
	if err != nil {
		t.Fatal(err)
	}

	var stderr strings.Builder
	_, err = cmd.Generate(ctx, abs, "", &cmd.Options{Stderr: &stderr})
	if err != nil {
		t.Fatalf("generate failed: %v\n%s", err, stderr.String())
	}

	queryFile := filepath.Join(abs, "go", "query.sql.go")
	body, err := os.ReadFile(queryFile)
	if err != nil {
		t.Fatal(err)
	}
	src := string(body)
	for _, want := range []string{
		"func (q *Queries) ListAuthors",
		"func (q *Queries) IterAuthors",
		"iter.Seq2[Author, error]",
		"defer rows.Close()",
	} {
		if !strings.Contains(src, want) {
			t.Fatalf("missing %q in generated code:\n%s", want, src)
		}
	}
	if strings.Contains(src, "sqlc:iterator-stream") {
		t.Fatal("internal stream marker leaked into generated code")
	}
}

func TestEmitIteratorsCodegenExplicitOnly(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	root := filepath.Join("testdata", "emit_iterators_explicit", "stdlib")
	abs, err := filepath.Abs(root)
	if err != nil {
		t.Fatal(err)
	}

	var stderr strings.Builder
	_, err = cmd.Generate(ctx, abs, "", &cmd.Options{Stderr: &stderr})
	if err != nil {
		t.Fatalf("generate failed: %v\n%s", err, stderr.String())
	}

	queryFile := filepath.Join(abs, "go", "query.sql.go")
	body, err := os.ReadFile(queryFile)
	if err != nil {
		t.Fatal(err)
	}
	src := string(body)
	if strings.Contains(src, "sqlc:iterator-stream") {
		t.Fatal("internal stream marker leaked into generated code")
	}
	if strings.Count(src, "func (q *Queries) IterAuthors") != 1 {
		t.Fatalf("expected exactly one IterAuthors, got:\n%s", src)
	}
	for _, want := range []string{
		"func (q *Queries) StreamAuthors",
		"func (q *Queries) IterAuthors",
		"q.db.QueryContext(ctx, streamAuthors)",
	} {
		if !strings.Contains(src, want) {
			t.Fatalf("missing %q in generated code:\n%s", want, src)
		}
	}
	listAuthorsBlock := src[strings.Index(src, "func (q *Queries) ListAuthors"):strings.Index(src, "const streamAuthors")]
	if strings.Contains(listAuthorsBlock, "IterAuthors") {
		t.Fatal("ListAuthors block must not contain iterator method")
	}
}
