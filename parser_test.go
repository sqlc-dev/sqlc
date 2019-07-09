package dinosql

import (
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
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
		"SELECT * FROM venue WHERE slug = $1 AND city = $2",
		"SELECT * FROM venue WHERE slug = $1",
		"SELECT * FROM venue LIMIT $1",
		"SELECT * FROM venue OFFSET $1",
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

func TestParseSchema(t *testing.T) {
	s, err := ParseSchmea(filepath.Join("testdata", "ondeck", "schema"), GenerateSettings{})
	if err != nil {
		t.Fatal(err)
	}

	q, err := ParseQueries(s, filepath.Join("testdata", "ondeck", "query"))
	if err != nil {
		t.Fatal(err)
	}

	t.Run("default", func(t *testing.T) {
		source := generate(q, GenerateSettings{
			Package: "ondeck",
		})

		blob, err := ioutil.ReadFile(filepath.Join("testdata", "ondeck", "db.go"))
		if err != nil {
			log.Fatal(err)
		}

		if diff := cmp.Diff(source, string(blob)); diff != "" {
			t.Errorf("genreated code differed (-want +got):\n%s", diff)
			t.Log(source)
		}
	})

	t.Run("prepared", func(t *testing.T) {
		source := generate(q, GenerateSettings{
			Package:             "prepared",
			EmitPreparedQueries: true,
		})

		blob, err := ioutil.ReadFile(filepath.Join("testdata", "ondeck", "prepared", "prepared.go"))
		if err != nil {
			log.Fatal(err)
		}

		if diff := cmp.Diff(source, string(blob)); diff != "" {
			t.Errorf("genreated code differed (-want +got):\n%s", diff)
			t.Log(source)
		}
	})
}

func TestCompile(t *testing.T) {
	files := []string{
		filepath.Join("testdata", "ondeck", "db.go"),
		filepath.Join("testdata", "ondeck", "prepared", "prepared.go"),
	}
	for _, filename := range files {
		f := filename
		t.Run(f, func(t *testing.T) {
			output, err := exec.Command("go", "build", f).CombinedOutput()
			if err != nil {
				t.Errorf("%s: %s:", err, string(output))
			}
		})
	}
}
