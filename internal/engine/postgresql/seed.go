package postgresql

import (
	_ "embed"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/core/seed"
)

// dialect.json describes PostgreSQL's type system. Two things in it are worth
// knowing about: the serial types are column definition sugar for an integer
// plus a sequence, but they are what a schema declares, so the catalog carries
// them as types of their own; and the numeric, string and datetime categories
// are listed as mutually castable, which is what an unadorned
// "WHERE bigint_col = 1" depends on.
//
//go:embed dialect.json
var dialectJSON []byte

// Dialect returns the catalog option that seeds PostgreSQL's type system,
// together with the generated pg_catalog function list.
func Dialect() core.Option {
	return seed.Dialect(dialectJSON, genPGCatalog().Funcs)
}
