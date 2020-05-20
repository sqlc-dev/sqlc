package catalog

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/sqlerr"
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

func (c *Catalog) GetFuncN(rel *ast.FuncName, n int) (Function, error) {
	ns := rel.Schema
	if ns == "" {
		ns = c.DefaultSchema
	}
	s, err := c.getSchema(ns)
	if err != nil {
		return Function{}, err
	}
	for i := range s.Funcs {
		if s.Funcs[i].Name != rel.Name {
			continue
		}
		if len(s.Funcs[i].Args) == n {
			return *s.Funcs[i], nil
		}
	}
	return Function{}, sqlerr.RelationNotFound(rel.Name)
}

func (c *Catalog) GetTable(rel *ast.TableName) (Table, error) {
	_, table, err := c.getTable(rel)
	return *table, err
}
