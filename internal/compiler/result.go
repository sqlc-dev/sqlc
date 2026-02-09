package compiler

import (
	"github.com/sqlc-dev/sqlc/internal/plugin"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

type Result struct {
	Catalog      *catalog.Catalog
	Queries      []*Query
	// PluginCatalog is set by plugin engines (e.g. sqlc-engine-ydb) from ParseResponse.Catalog.
	// When non-nil, it is passed to codegen instead of converting Catalog.
	PluginCatalog *plugin.Catalog
}
