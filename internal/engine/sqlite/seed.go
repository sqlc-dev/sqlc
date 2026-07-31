package sqlite

import (
	_ "embed"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/core/seed"
)

// SQLite has five storage classes and assigns a column one of them from the
// declared type's affinity, but a schema may declare any name at all. The types
// in dialect.json are the spellings the affinity rules call out, grouped by the
// storage class they land in. They are all mutually castable, because SQLite
// compares values whatever their declared type.
//
//go:embed dialect.json
var dialectJSON []byte

// Dialect returns the catalog option that seeds SQLite's type system, together
// with the functions the engine's schema declares.
func Dialect() core.Option {
	return seed.Dialect(dialectJSON, defaultSchema("main").Funcs)
}
