package duckdb

import (
	"embed"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/core/seed"
)

// The dialect directory describes DuckDB's type system. types.jsonl,
// functions.jsonl and operators.jsonl are generated from a live DuckDB 2.0
// CLI by goldeneye (internal/goldeneye), which also checks them against one;
// dialect.json is authored by hand. Regenerate from internal/goldeneye with:
//
//	DUCKDB=/path/to/duckdb go run ./cmd/goldeneye generate duckdb
//
//go:embed dialect
var dialectFS embed.FS

// Dialect returns the catalog option that seeds DuckDB's type system.
func Dialect() core.Option {
	return seed.Dialect(dialectFS, "dialect")
}

func init() {
	core.RegisterCanonicalizer("duckdb", canonicalize)
}

// canonicalize rewrites a type the way DuckDB stores it: a varchar's length
// is dropped, and a decimal declared without a precision is decimal(18,3).
func canonicalize(t *core.TypeExpr) *core.TypeExpr {
	switch t.Name {
	case "varchar":
		if len(t.Args) > 0 {
			return &core.TypeExpr{Name: t.Name, Nullable: t.Nullable}
		}
	case "decimal":
		if len(t.Args) == 0 {
			p, s := int64(18), int64(3)
			return &core.TypeExpr{Name: t.Name, Nullable: t.Nullable, Args: []core.TypeArg{{Int: &p}, {Int: &s}}}
		}
	}
	return t
}
