package analyzer_test

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/core/analyzer"
	coreschema "github.com/sqlc-dev/sqlc/internal/core/schema"
	"github.com/sqlc-dev/sqlc/internal/engine/googlesql"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/rewrite"
	"github.com/sqlc-dev/sqlc/internal/sql/validate"
)

const benchSchema = `
CREATE TABLE users (
  id   INT64  NOT NULL,
  name STRING NOT NULL,
  bio  STRING,
) PRIMARY KEY (id);

CREATE TABLE posts (
  id      INT64     NOT NULL,
  user_id INT64     NOT NULL,
  title   STRING(255),
  created TIMESTAMP NOT NULL,
) PRIMARY KEY (id);
`

const benchQueries = `
SELECT * FROM users;
SELECT COUNT(*) AS total FROM users;
SELECT u.name, p.* FROM users u JOIN posts p ON p.user_id = u.id WHERE u.name = @name;
SELECT id, name, bio FROM users WHERE id = @id AND name = @name;
SELECT p.title, p.created FROM posts p WHERE p.user_id = @uid AND p.title = @t;
`

func parseAll(t testing.TB, src string) []ast.Node {
	p := googlesql.NewParser()
	stmts, err := p.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	out := make([]ast.Node, 0, len(stmts))
	for i := range stmts {
		out = append(out, stmts[i].Raw)
	}
	return out
}

// parseQueries mirrors the compiler pipeline: parse, then rewrite named
// parameters into ParamRef nodes before handing the statement to the analyzer.
func parseQueries(t testing.TB, src string) []ast.Node {
	out := make([]ast.Node, 0)
	for _, n := range parseAll(t, src) {
		raw, ok := n.(*ast.RawStmt)
		if !ok {
			t.Fatalf("not a raw statement: %T", n)
		}
		numbers, dollar, err := validate.ParamRef(raw)
		if err != nil {
			t.Fatal(err)
		}
		rewritten, _, _ := rewrite.NamedParameters(config.EngineGoogleSQL, raw, numbers, dollar)
		out = append(out, rewritten)
	}
	return out
}

func newBenchCatalog(t testing.TB) (*core.Catalog, []ast.Node) {
	cat, err := core.New(googlesql.Dialect())
	if err != nil {
		t.Fatal(err)
	}
	for _, n := range parseAll(t, benchSchema) {
		if err := coreschema.Apply(cat, n); err != nil {
			t.Fatal(err)
		}
	}
	return cat, parseQueries(t, benchQueries)
}

// BenchmarkCatalogNew measures per-compile catalog setup: open the SQLite
// catalog, install the schema, and run the dialect seed.
func BenchmarkCatalogNew(b *testing.B) {
	for b.Loop() {
		cat, err := core.New(googlesql.Dialect())
		if err != nil {
			b.Fatal(err)
		}
		cat.Close()
	}
}

// BenchmarkApplySchema measures loading user DDL into the catalog.
func BenchmarkApplySchema(b *testing.B) {
	stmts := parseAll(b, benchSchema)
	for b.Loop() {
		cat, err := core.New(googlesql.Dialect())
		if err != nil {
			b.Fatal(err)
		}
		for _, n := range stmts {
			if err := coreschema.Apply(cat, n); err != nil {
				b.Fatal(err)
			}
		}
		cat.Close()
	}
}

// BenchmarkPrepare measures steady-state query analysis against a warm catalog.
func BenchmarkPrepare(b *testing.B) {
	cat, queries := newBenchCatalog(b)
	defer cat.Close()
	for b.Loop() {
		for _, q := range queries {
			if _, err := analyzer.Prepare(cat, q); err != nil {
				b.Fatal(err)
			}
		}
	}
}
