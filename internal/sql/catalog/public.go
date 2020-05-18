package catalog

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

// TODO: Decide on a real, exported interface
func (c *Catalog) ListFuncsByName(rel *ast.FuncName) ([]Function, error) {
	var funcs []Function

	ns := rel.Schema
	if ns == "" {
		ns = c.DefaultSchema
	}
	s, err := c.getSchema(ns)
	if err != nil {
		return nil, err
	}
	for i := range s.Funcs {
		if s.Funcs[i].Name == rel.Name {
			funcs = append(funcs, *s.Funcs[i])
		}
	}
	return funcs, nil
}
