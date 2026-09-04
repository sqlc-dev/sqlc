package catalog

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func (c *Catalog) createExtension(stmt *ast.CreateExtensionStmt) error {
	if stmt.Extname == nil {
		return nil
	}
	// TODO: Implement IF NOT EXISTS
	return c.loadExtension(*stmt.Extname)
}

// loadExtension adds what the engine knows of the named extension — or of
// the virtual table module that stands for one — to the default schema,
// once. An engine without extension data, or without data for this one,
// adds nothing.
func (c *Catalog) loadExtension(name string) error {
	if _, exists := c.Extensions[name]; exists {
		return nil
	}
	if c.LoadExtension == nil {
		return nil
	}
	ext := c.LoadExtension(name)
	if ext == nil {
		return nil
	}
	if c.Extensions == nil {
		c.Extensions = map[string]struct{}{}
	}
	c.Extensions[name] = struct{}{}
	s, err := c.getSchema(c.DefaultSchema)
	if err != nil {
		return err
	}
	// TODO: Error on duplicate functions
	s.Funcs = append(s.Funcs, ext.Funcs...)
	return nil
}
