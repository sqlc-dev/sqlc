package dolphin

import (
	_ "embed"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/core/seed"
)

//go:embed dialect.json
var dialectJSON []byte

// Dialect returns the catalog option that seeds MySQL's type system, together
// with the functions the engine's schema declares.
func Dialect() core.Option {
	return seed.Dialect(dialectJSON, defaultSchema("public").Funcs)
}
