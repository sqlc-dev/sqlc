package dolphin

import (
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// toPointer converts an int to a pointer without a temporary
// variable at the call-site, and is used by the generated schemas
func toPointer(x int) *int {
	return &x
}

func NewCatalog() *catalog.Catalog {
	def := "public" // TODO: What is the default database for MySQL?

	c := catalog.New(def)
	// New() creates an empty schema which we'll replace with MySQL stdlib functions
	c.Schemas[0] = defaultSchema(def)
	c.Schemas = append(c.Schemas, genInformationSchema())
	c.Schemas = append(c.Schemas, genPerformanceSchema())
	c.Schemas = append(c.Schemas, genSysSchema())
	c.Schemas = append(c.Schemas, genMysqlCatalog())

	return c
}
