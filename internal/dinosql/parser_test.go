package dinosql

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/google/go-cmp/cmp"
	pg "github.com/lfittl/pg_query_go"
	nodes "github.com/lfittl/pg_query_go/nodes"

	"github.com/kyleconroy/sqlc/internal/source"
	"github.com/kyleconroy/sqlc/internal/sql/sqlpath"
)

const pluck = `
SELECT * FROM venue WHERE slug = $1 AND city = $2;
SELECT * FROM venue WHERE slug = $1;
SELECT * FROM venue LIMIT $1;
SELECT * FROM venue OFFSET $1;
`

func TestPluck(t *testing.T) {
	tree, err := pg.Parse(pluck)
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{
		"\nSELECT * FROM venue WHERE slug = $1 AND city = $2",
		"\nSELECT * FROM venue WHERE slug = $1",
		"\nSELECT * FROM venue LIMIT $1",
		"\nSELECT * FROM venue OFFSET $1",
	}

	for i, stmt := range tree.Statements {
		switch n := stmt.(type) {
		case nodes.RawStmt:
			q, err := source.Pluck(pluck, n.StmtLocation, n.StmtLen)
			if err != nil {
				t.Error(err)
				continue
			}
			if q != expected[i] {
				t.Errorf("expected %s, got %s", expected[i], q)
			}
		default:
			t.Fatalf("wrong type; %T", n)
		}
	}
}

const lineColumn = `SELECT 1; SELECT * FROM venue WHERE slug = $1 AND city = $2;

SELECT * FROM venue WHERE slug = $1;
  SELECT * 
FROM venue 
LIMIT $1;

-- comment here
SELECT * FROM venue
OFFSET $1; SELECT 1;
`

func TestLineColumn(t *testing.T) {
	tree, err := pg.Parse(lineColumn)
	if err != nil {
		t.Fatal(err)
	}

	for i, test := range []struct {
		node   nodes.Node
		line   int
		column int
	}{
		{tree.Statements[0], 1, 1},
		{tree.Statements[1], 1, 11},
		{tree.Statements[2], 3, 1},
		{tree.Statements[3], 4, 3},
		{tree.Statements[4], 9, 1},
		{tree.Statements[5], 10, 12},
	} {
		raw := test.node.(nodes.RawStmt)
		line, column := source.LineNumber(lineColumn, raw.StmtLocation)
		if line != test.line {
			t.Errorf("expected stmt %d to be on line %d, not %d", i, test.line, line)
		}
		if column != test.column {
			t.Errorf("expected stmt %d to be on column %d, not %d", i, test.column, column)
		}
	}
}

func TestExtractArgs(t *testing.T) {
	queries := []struct {
		query       string
		bindNumbers []int
	}{
		{"SELECT * FROM venue WHERE slug = $1 AND city = $2", []int{1, 2}},
		{"SELECT * FROM venue WHERE slug = $1 AND region = $2 AND city = $3 AND country = $2", []int{1, 2, 3, 2}},
		{"SELECT * FROM venue WHERE slug = $1", []int{1}},
		{"SELECT * FROM venue LIMIT $1", []int{1}},
		{"SELECT * FROM venue OFFSET $1", []int{1}},
	}
	for _, q := range queries {
		tree, err := pg.Parse(q.query)
		if err != nil {
			t.Fatal(err)
		}
		for _, stmt := range tree.Statements {
			refs := findParameters(stmt)
			if err != nil {
				t.Error(err)
			}
			nums := make([]int, len(refs))
			for i, n := range refs {
				nums[i] = n.ref.Number
			}
			if diff := cmp.Diff(q.bindNumbers, nums); diff != "" {
				t.Errorf("expected bindings %v, got %v", q.bindNumbers, nums)
			}
		}
	}
}

func TestRewriteParameters(t *testing.T) {
	queries := []struct {
		orig string
		new  string
	}{
		{"SELECT * FROM venue WHERE slug = $1 AND city = $3 AND bar = $2", "SELECT * FROM venue WHERE slug = ? AND city = ? AND bar = ?"},
		{"DELETE FROM venue WHERE slug = $1 AND slug = $1", "DELETE FROM venue WHERE slug = ? AND slug = ?"},
		{"SELECT * FROM venue LIMIT $1", "SELECT * FROM venue LIMIT ?"},
	}
	for _, q := range queries {
		tree, err := pg.Parse(q.orig)
		if err != nil {
			t.Fatal(err)
		}
		for _, stmt := range tree.Statements {
			refs := findParameters(stmt)
			if err != nil {
				t.Error(err)
			}
			edits, err := rewriteNumberedParameters(refs, stmt.(nodes.RawStmt), q.orig)
			if err != nil {
				t.Error(err)
			}
			rewritten, err := editQuery(q.orig, edits)
			if err != nil {
				t.Error(err)
			}
			if rewritten != q.new {
				t.Errorf("expected %q, got %q", q.new, rewritten)
			}
		}
	}
}

func TestExpand(t *testing.T) {
	// pretend that foo has two columns, a and b
	raw := `SELECT *, *, foo.* FROM foo`
	expected := `SELECT a, b, a, b, foo.a, foo.b FROM foo`
	edits := []source.Edit{
		{7, "*", "a, b"},
		{10, "*", "a, b"},
		{13, "foo.*", "foo.a, foo.b"},
	}
	actual, err := editQuery(raw, edits)
	if err != nil {
		t.Error(err)
	}
	if expected != actual {
		t.Errorf("mismatch:\nexpected: %s\n  acutal: %s", expected, actual)
	}
}

func TestReadFiles(t *testing.T) {
	dir, err := ioutil.TempDir("", "sqlc-test-read-files")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir)

	subdir1 := path.Join(dir, "subdir1")
	if err := os.Mkdir(subdir1, 0777); err != nil {
		t.Error(err)
	}

	subdir2 := path.Join(dir, "subdir2")
	if err := os.Mkdir(subdir2, 0777); err != nil {
		t.Error(err)
	}

	files := []string{
		path.Join(subdir1, "include-me.sql"),
		path.Join(subdir1, "include-me.up.sql"),
		path.Join(subdir1, "not-me.down.sql"),
		path.Join(subdir2, "include-me.sql"),
		path.Join(subdir2, "include-me.up.sql"),
		path.Join(subdir2, "not-me.down.sql"),
	}
	for _, filename := range files {
		fd, err := os.Create(filename)
		if err != nil {
			t.Error(err)
		}
		defer fd.Close()
	}

	input := []string{
		subdir1,
		path.Join(subdir2, "include-me.sql"),
		path.Join(subdir2, "include-me.up.sql"),
		path.Join(subdir2, "not-me.down.sql"),
	}

	expectedFiles := []string{
		path.Join(subdir1, "include-me.sql"),
		path.Join(subdir1, "include-me.up.sql"),
		path.Join(subdir2, "include-me.sql"),
		path.Join(subdir2, "include-me.up.sql"),
	}

	filesRead, err := sqlpath.Glob(input)
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(expectedFiles, filesRead) {
		t.Errorf("unexpected files: %s", cmp.Diff(expectedFiles, filesRead))
	}
}
