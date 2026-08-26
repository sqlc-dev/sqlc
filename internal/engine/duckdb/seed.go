package duckdb

import (
	"embed"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/core/seed"
)

// The dialect directory describes DuckDB's type system. types.jsonl,
// functions.jsonl and operators.jsonl are generated from a live DuckDB 2.0
// CLI by sqlc-duckdb-gen (internal/tools/sqlc-duckdb-gen); dialect.json is
// authored by hand. Regenerate with:
//
//	DUCKDB=/path/to/duckdb go run ./internal/tools/sqlc-duckdb-gen
//
//go:embed dialect
var dialectFS embed.FS

// Dialect returns the catalog option that seeds DuckDB's type system.
func Dialect() core.Option {
	return seed.Dialect(dialectFS, "dialect")
}
