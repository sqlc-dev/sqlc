package mssql

import (
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// NewCatalog creates a new MSSQL catalog with the default schema set to 'dbo'
func NewCatalog() *catalog.Catalog {
	return &catalog.Catalog{
		DefaultSchema: "dbo",
		Schemas: []*catalog.Schema{
			defaultSchema("dbo"),
		},
		Extensions: map[string]struct{}{},
	}
}

func defaultSchema(name string) *catalog.Schema {
	return &catalog.Schema{
		Name: name,
	}
}
