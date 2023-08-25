package catalog

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func (c *Catalog) createExtension(stmt *ast.CreateExtensionStmt) error {
	if stmt.Extname == nil {
		return nil
	}
	// TODO: Implement IF NOT EXISTS
	if _, exists := c.Extensions[*stmt.Extname]; exists {
		return nil
	}
	if c.LoadExtension == nil {
		return nil
	}
	ext := c.LoadExtension(*stmt.Extname)
	if ext == nil {
		return nil
	}
	s, err := c.getSchema(c.DefaultSchema)
	if err != nil {
		return err
	}
	// TODO: Error on duplicate functions
	s.Funcs = append(s.Funcs, ext.Funcs...)
	return nil
}
