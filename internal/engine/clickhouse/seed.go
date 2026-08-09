package clickhouse

import (
	"embed"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/core/seed"
)

//go:embed dialect
var dialectFS embed.FS

// Dialect returns the catalog option that seeds ClickHouse's type system.
func Dialect() core.Option {
	return seed.Dialect(dialectFS, "dialect")
}
