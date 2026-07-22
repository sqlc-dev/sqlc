package compiler

import (
	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

type Result struct {
	Catalog *catalog.Catalog
	Queries []*Query

	// CoreCatalog, when set, is the xqlc-derived catalog that drives both
	// analysis and codegen for engines on the core (currently ClickHouse).
	// When it is non-nil, codegen projects the catalog from it instead of
	// from Catalog.
	CoreCatalog *core.Catalog
}
