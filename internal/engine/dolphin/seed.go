package dolphin

import (
	"embed"
	"sync"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/core/seed"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// The dialect directory describes MySQL's type system, and functions.jsonl is
// its standard library. Both the analysis core and the catalog the legacy
// compiler builds read them.
//
//go:embed dialect
var dialectFS embed.FS

// Dialect returns the catalog option that seeds MySQL's type system.
func Dialect() core.Option {
	return seed.Dialect(dialectFS, "dialect")
}

// stdlib is MySQL's functions in the form the catalog uses. They are embedded
// in the binary and never change within a run, so they are read once.
var stdlib = sync.OnceValue(func() []*catalog.Function {
	funcs, err := seed.Functions(dialectFS, "dialect")
	if err != nil {
		// A failure here means sqlc was built from a broken tree, which no
		// caller can do anything about.
		panic(err)
	}
	return funcs
})
