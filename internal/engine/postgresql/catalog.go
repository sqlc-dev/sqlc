package postgresql

import "github.com/kyleconroy/sqlc/internal/sql/catalog"

func NewCatalog() *catalog.Catalog {
	c := catalog.New("public")
	c.Schemas = append(c.Schemas, pgTemp())
	c.Schemas = append(c.Schemas, genPGCatalog())
	c.SearchPath = []string{"pg_catalog"}
	c.LoadExtension = loadExtension
	return c
}

// The generated pg_catalog is very slow to compare because it has so
// many entries. For testing, don't include it.
func newTestCatalog() *catalog.Catalog {
	c := catalog.New("public")
	c.Schemas = append(c.Schemas, pgTemp())
	c.LoadExtension = loadExtension
	return c
}

func loadExtension(name string) *catalog.Schema {
	switch name {
	case "citext":
		return genCIText()
	case "pg_trgm":
		return genPGTrigram()
	case "pgcrypto":
		return genPGCrypto()
	case "uuid-ossp":
		return genUUIDOSSP()
	}
	return nil
}
