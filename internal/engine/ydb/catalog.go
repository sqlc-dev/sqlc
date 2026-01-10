package ydb

import "github.com/sqlc-dev/sqlc/internal/sql/catalog"


func NewCatalog() *catalog.Catalog {
	def := "main"
	return &catalog.Catalog{
		DefaultSchema: def,
		Schemas: []*catalog.Schema{
			defaultSchema(def),
		},
		Extensions: map[string]struct{}{},
	}
}

func NewTestCatalog() *catalog.Catalog {
	return catalog.New("main")
}
