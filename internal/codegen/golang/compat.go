package golang

import (
	"github.com/kyleconroy/sqlc/internal/core"
	"github.com/kyleconroy/sqlc/internal/plugin"
)

func sameTableName(n *plugin.Identifier, f core.FQN, defaultSchema string) bool {
	if n == nil {
		return false
	}
	schema := n.Schema
	if n.Schema == "" {
		schema = defaultSchema
	}
	return n.Catalog == f.Catalog && schema == f.Schema && n.Name == f.Rel
}
