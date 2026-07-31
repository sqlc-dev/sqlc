package googlesql

import (
	_ "embed"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/core/seed"
)

//go:embed dialect.json
var dialectJSON []byte

// Dialect returns the catalog option that seeds GoogleSQL's type system.
func Dialect() core.Option {
	return seed.Dialect(dialectJSON, nil)
}
