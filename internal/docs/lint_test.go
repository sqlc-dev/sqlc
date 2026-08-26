package docs

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestDocs validates the repository's docs/ directory against the content
// contract. This is a unit test rather than an endtoend case because the
// contract is a property of the documentation source tree, not of anything
// the sqlc CLI does: there is no command whose output could pin it down.
func TestDocs(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not locate source file")
	}
	root := filepath.Join(filepath.Dir(thisFile), "..", "..", "docs")

	problems, err := Lint(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range problems {
		t.Errorf("%s", p)
	}
}

// TestLintViolations pins down what the linter rejects, using a synthetic
// docs tree that violates each rule of the contract once.
func TestLintViolations(t *testing.T) {
	root := t.TempDir()
	write := func(name, content string) {
		t.Helper()
		p := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write("toc.yaml", `
index: index.md
sections:
  - title: Guides
    pages:
      - sub/page.md
      - sub/page.md
      - ghost.md
`)
	write("index.md", `# Index

A [broken link](missing.md) and a [bad anchor](sub/page.md#nope) and
a [bad self anchor](#nowhere) and an [absolute link](/sub/page.md) and
an [escaping link](../outside.md).

<div>raw html</div>

<!-- comments are fine -->

:::{note}
myst directive
:::

`+"```{warning}"+`
myst fence
`+"```"+`

# Second title
`)
	write("sub/page.md", "Not a heading first.\n\n# Title\n")
	write("orphan.md", "# Orphan\n")

	problems, err := Lint(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		`index.md:3: link target "missing.md" does not exist`,
		"index.md:3: anchor #nope not found in sub/page.md",
		"index.md:3: anchor #nowhere not found in this page",
		`index.md:3: absolute link "/sub/page.md"`,
		`index.md:3: link "../outside.md" points outside docs/`,
		"index.md:7: raw HTML is not allowed: <div>raw html</div>",
		"index.md:11: MyST directive syntax is not allowed",
		`index.md:16: MyST fenced directive "{warning}" is not allowed`,
		"index.md:19: page has more than one level-1 heading",
		"sub/page.md:1: page must start with a level-1 heading",
		"toc.yaml:1: sub/page.md is listed 2 times",
		"toc.yaml:1: ghost.md is listed but does not exist",
		"toc.yaml:1: orphan.md is not listed",
	}
	var got []string
	for _, p := range problems {
		got = append(got, p.String())
	}
	for _, w := range want {
		found := false
		for _, g := range got {
			if strings.Contains(g, w) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected a problem containing %q, got:\n  %s", w, strings.Join(got, "\n  "))
		}
	}
	if len(got) != len(want) {
		t.Errorf("got %d problems, want %d:\n  %s", len(got), len(want), strings.Join(got, "\n  "))
	}
}
