package golang

import (
	"github.com/kyleconroy/sqlc/internal/plugin"
)

func sameTableName(tableID, f *plugin.Identifier, defaultSchema string) bool {
	if tableID == nil {
		return false
	}
	schema := tableID.Schema
	if tableID.Schema == "" {
		schema = defaultSchema
	}
	return tableID.Catalog == f.Catalog && schema == f.Schema && tableID.Name == f.Name
}
