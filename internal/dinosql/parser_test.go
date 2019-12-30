package dinosql

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	pg "github.com/lfittl/pg_query_go"
	nodes "github.com/lfittl/pg_query_go/nodes"
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
			q, err := pluckQuery(pluck, n)
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
		line, column := lineno(lineColumn, raw.StmtLocation)
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
		query string
		count int
	}{
		{"SELECT * FROM venue WHERE slug = $1 AND city = $2", 2},
		{"SELECT * FROM venue WHERE slug = $1", 1},
		{"SELECT * FROM venue LIMIT $1", 1},
		{"SELECT * FROM venue OFFSET $1", 1},
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
			if len(refs) != q.count {
				t.Errorf("expected %d refs, got %d", q.count, len(refs))
			}
		}
	}
}

func cmpDirectory(t *testing.T, dir string, actual map[string]string) {
	t.Helper()

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]string{}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), ".go") {
			continue
		}
		if strings.HasSuffix(file.Name(), "_test.go") {
			continue
		}
		blob, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			t.Fatal(err)
		}
		expected[file.Name()] = string(blob)
	}

	if !cmp.Equal(expected, actual) {
		t.Errorf("%s contents differ", dir)
		for name, contents := range expected {
			if actual[name] == "" {
				t.Errorf("%s is empty", name)
				continue
			}
			if diff := cmp.Diff(contents, actual[name]); diff != "" {
				t.Errorf("%s differed (-want +got):\n%s", name, diff)
			}
		}
	}
}

func TestParseSchema(t *testing.T) {
	c, err := ParseCatalog(filepath.Join("testdata", "ondeck", "schema"))
	if err != nil {
		t.Fatal(err)
	}

	q, err := ParseQueries(c, PackageSettings{
		Name:         "db",
		Queries:      filepath.Join("testdata", "ondeck", "query"),
		EmitJSONTags: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("default", func(t *testing.T) {
		output, err := Generate(q, mockSettings)
		if err != nil {
			t.Fatal(err)
		}

		cmpDirectory(t, filepath.Join("testdata", "ondeck"), output)
	})

	t.Run("prepared", func(t *testing.T) {
		output, err := Generate(q, mockSettings)
		if err != nil {
			t.Fatal(err)
		}

		cmpDirectory(t, filepath.Join("testdata", "ondeck", "prepared"), output)
	})
}

func TestCompile(t *testing.T) {
	cmd := exec.Command("go", "build", "-mod", "readonly", "./...")
	cmd.Dir = filepath.Join("testdata", "ondeck")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("%s: %s", err, string(output))
	}
}

func TestParseMetadata(t *testing.T) {
	for _, query := range []string{
		`-- name: CreateFoo, :one`,
		`-- name: 9Foo_, :one`,
		`-- name: CreateFoo :two`,
		`-- name: CreateFoo`,
		`-- name: CreateFoo :one something`,
		`-- name: `,
	} {
		if _, _, err := parseMetadata(query); err == nil {
			t.Errorf("expected invalid metadata: %q", query)
		}
	}
}

func TestExpand(t *testing.T) {
	// pretend that foo has two columns, a and b
	raw := `SELECT *, *, foo.* FROM foo`
	expected := `SELECT a, b, a, b, foo.a, foo.b FROM foo`
	edits := []edit{
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
