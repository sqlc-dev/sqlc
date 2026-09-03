package clickhouse

import (
	"embed"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/core/seed"
)

// The dialect directory describes ClickHouse's type system. types.jsonl is
// generated from system.data_type_families of the pinned ClickHouse release
// by goldeneye (internal/goldeneye), which also checks it against one;
// dialect.json and functions.jsonl are authored by hand, since ClickHouse
// publishes no function signatures. Regenerate from internal/goldeneye with:
//
//	go run ./cmd/goldeneye install clickhouse
//	go run ./cmd/goldeneye generate clickhouse
//
//go:embed dialect
var dialectFS embed.FS

// Dialect returns the catalog option that seeds ClickHouse's type system.
func Dialect() core.Option {
	return seed.Dialect(dialectFS, "dialect")
}
