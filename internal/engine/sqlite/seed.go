package sqlite

import (
	"embed"
	"strings"
	"sync"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/core/seed"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// The dialect directory describes SQLite's type system, and functions.jsonl is
// its standard library. Both the analysis core and the catalog the legacy
// compiler builds read them.
//
//go:embed dialect
var dialectFS embed.FS

// Dialect returns the catalog option that seeds SQLite's type system.
func Dialect() core.Option {
	return seed.Dialect(dialectFS, "dialect")
}

func init() {
	core.RegisterUserTypeBase("sqlite", affinity)
}

// affinity is the type a declared spelling SQLite has no name for stands
// on: the affinity its rule gives it, decided by the words in it. INT
// anywhere is INTEGER; CHAR, CLOB or TEXT is TEXT; BLOB is BLOB; REAL, FLOA
// or DOUB is REAL; anything else is NUMERIC. A column with no type at all
// has BLOB affinity, but sqlc reads one as any.
func affinity(name string) (base, category string) {
	upper := strings.ToUpper(name)
	switch {
	case strings.Contains(upper, "INT"):
		return "integer", "N"
	case strings.Contains(upper, "CHAR"), strings.Contains(upper, "CLOB"), strings.Contains(upper, "TEXT"):
		return "text", "S"
	case strings.Contains(upper, "BLOB"):
		return "blob", "U"
	case strings.Contains(upper, "REAL"), strings.Contains(upper, "FLOA"), strings.Contains(upper, "DOUB"):
		return "real", "N"
	}
	return "numeric", "N"
}

// stdlib is SQLite's functions in the form the catalog uses. They are embedded
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
