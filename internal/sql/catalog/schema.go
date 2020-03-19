package catalog

import (
	"fmt"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/sqlerr"
)

func (c *Catalog) createSchema(stmt *ast.CreateSchemaStmt) error {
	if stmt.Name == nil {
		return fmt.Errorf("create schema: empty name")
	}
	if _, err := c.getSchema(*stmt.Name); err == nil {
		if !stmt.IfNotExists {
			return sqlerr.SchemaExists(*stmt.Name)
		}
	}
	c.Schemas = append(c.Schemas, &Schema{Name: *stmt.Name})
	return nil
}

func (c *Catalog) dropSchema(stmt *ast.DropSchemaStmt) error {
	// TODO: n^2 in the worst-case
	for _, name := range stmt.Schemas {
		idx := -1
		for i := range c.Schemas {
			if c.Schemas[i].Name == name.Str {
				idx = i
			}
		}
		if idx == -1 {
			if stmt.MissingOk {
				continue
			}
			return sqlerr.SchemaNotFound(name.Str)
		}
		c.Schemas = append(c.Schemas[:idx], c.Schemas[idx+1:]...)
	}
	return nil
}
