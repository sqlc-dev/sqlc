package postgresql

import (
	"github.com/sqlc-dev/sqlc/internal/core/seed"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func NewCatalog() *catalog.Catalog {
	c := catalog.New("public")
	c.Schemas = append(c.Schemas, pgTemp())
	c.Schemas = append(c.Schemas, genPGCatalog())
	c.Schemas = append(c.Schemas, genInformationSchema())
	c.SearchPath = []string{"pg_catalog"}
	c.LoadExtension = loadExtension
	return c
}

// genPGCatalog and genInformationSchema build the schemas sqlc knows
// PostgreSQL by. Both read the dialect directory, which sqlc-pg-gen writes,
// and which the analysis core seeds its catalog from as well.
func genPGCatalog() *catalog.Schema {
	return systemSchema("pg_catalog", pgCatalogFuncs())
}

func genInformationSchema() *catalog.Schema {
	return systemSchema("information_schema", nil)
}

// systemSchema builds a schema from the relations recorded for it. The tables
// are read fresh each time, because a catalog is free to change the schemas it
// was built with.
func systemSchema(name string, funcs []*catalog.Function) *catalog.Schema {
	tables, err := seed.Relations(dialectFS, "dialect", name)
	if err != nil {
		// The relations are embedded in the binary: a failure here means sqlc
		// was built from a broken tree, which no caller can do anything about.
		panic(err)
	}
	return &catalog.Schema{Name: name, Funcs: funcs, Tables: tables}
}
