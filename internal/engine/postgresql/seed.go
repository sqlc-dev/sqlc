package postgresql

import (
	"embed"
	"sync"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/core/seed"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// The dialect directory describes PostgreSQL's type system. Two things in it
// are worth knowing about: the serial types are column definition sugar for an
// integer plus a sequence, but they are what a schema declares, so the catalog
// carries them as types of their own; and the numeric, string and datetime
// categories are listed as mutually castable, which is what an unadorned
// "WHERE bigint_col = 1" depends on.
//
// functions.jsonl and relations.jsonl are pg_catalog's functions and the
// tables and views of pg_catalog and information_schema, all written by
// sqlc-pg-gen. Both the analysis core and the catalog the legacy compiler
// builds read them. Each contrib extension sqlc knows is a directory under
// extensions/ holding the functions CREATE EXTENSION adds.
//
//go:embed dialect
var dialectFS embed.FS

// Dialect returns the catalog option that seeds PostgreSQL's type system.
func Dialect() core.Option {
	return seed.Dialect(dialectFS, "dialect")
}

// pgCatalogFuncs is pg_catalog's functions in the form the catalog uses. The
// list runs to thousands of entries and never changes within a run, so it is
// read once.
var pgCatalogFuncs = sync.OnceValue(func() []*catalog.Function {
	funcs, err := seed.Functions(dialectFS, "dialect")
	if err != nil {
		// The list is embedded in the binary: a failure here means sqlc was
		// built from a broken tree, which no caller can do anything about.
		panic(err)
	}
	return funcs
})
