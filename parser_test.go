package strongdb

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

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
	queries := []string{
		"SELECT * FROM venue WHERE slug = $1 AND city = $2",
		"SELECT * FROM venue WHERE slug = $1",
		"SELECT * FROM venue LIMIT $1",
		"SELECT * FROM venue OFFSET $1",
	}
	for _, q := range queries {
		tree, err := pg.Parse(q)
		if err != nil {
			t.Fatal(err)
		}
		for _, stmt := range tree.Statements {
			refs := extractArgs(stmt)
			if err != nil {
				t.Error(err)
			}
			t.Logf("refs: %#v", refs)
		}
	}
}

func TestParseSchema(t *testing.T) {
	s, err := ParseSchmea(filepath.Join("testdata", "ondeck", "schema"))
	if err != nil {
		t.Error(err)
	}

	q, err := ParseQueries(s, filepath.Join("testdata", "ondeck", "query"))
	if err != nil {
		t.Error(err)
	}

	source := generate(q, "ondeck")

	blob, err := ioutil.ReadFile(filepath.Join("testdata", "ondeck", "db.go"))
	if err != nil {
		log.Fatal(err)
	}

	if source != string(blob) {
		t.Errorf("output differs")
		fmt.Println(source)
	}
}
