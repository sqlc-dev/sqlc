package postgresql

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// supplementalPGCatalogFuncs holds pg_catalog functions not covered by the
// generated catalog (pg_catalog.go). merge_action() is a PostgreSQL 17
// function that is only valid inside a MERGE ... RETURNING clause and returns
// the action taken ('INSERT', 'UPDATE', or 'DELETE').
var supplementalPGCatalogFuncs = []*catalog.Function{
	{
		Name:       "merge_action",
		Args:       []*catalog.Argument{},
		ReturnType: &ast.TypeName{Name: "text"},
	},
}

// toPointer converts an int to a pointer without a temporary
// variable at the call-site, and is used by the generated schemas
//
//go:fix inline
func toPointer(x int) *int {
	return new(x)
}

func NewCatalog() *catalog.Catalog {
	c := catalog.New("public")
	c.Schemas = append(c.Schemas, pgTemp())
	pgCatalog := genPGCatalog()
	// genPGCatalog aliases the package-global funcsgenPGCatalog slice; build a
	// fresh slice rather than appending into it so the global is never mutated.
	funcs := make([]*catalog.Function, 0, len(pgCatalog.Funcs)+len(supplementalPGCatalogFuncs))
	funcs = append(funcs, pgCatalog.Funcs...)
	funcs = append(funcs, supplementalPGCatalogFuncs...)
	pgCatalog.Funcs = funcs
	c.Schemas = append(c.Schemas, pgCatalog)
	c.Schemas = append(c.Schemas, genInformationSchema())
	c.SearchPath = []string{"pg_catalog"}
	c.LoadExtension = loadExtension
	return c
}
