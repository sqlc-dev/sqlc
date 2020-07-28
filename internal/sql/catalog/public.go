package catalog

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/sqlerr"
)

func (c *Catalog) schemasToSearch(ns string) []string {
	if ns == "" {
		ns = c.DefaultSchema
	}
	return append(c.SearchPath, ns)
}

func (c *Catalog) ListFuncsByName(rel *ast.FuncName) ([]Function, error) {
	var funcs []Function
	for _, ns := range c.schemasToSearch(rel.Schema) {
		s, err := c.getSchema(ns)
		if err != nil {
			return nil, err
		}
		for i := range s.Funcs {
			if s.Funcs[i].Name == rel.Name {
				funcs = append(funcs, *s.Funcs[i])
			}
		}
	}
	return funcs, nil
}

func (c *Catalog) GetFuncN(rel *ast.FuncName, n int) (Function, error) {
	for _, ns := range c.schemasToSearch(rel.Schema) {
		s, err := c.getSchema(ns)
		if err != nil {
			return Function{}, err
		}
		for i := range s.Funcs {
			if s.Funcs[i].Name != rel.Name {
				continue
			}
			args := s.Funcs[i].InArgs()
			if len(args) == n {
				return *s.Funcs[i], nil
			}
		}
	}
	return Function{}, sqlerr.RelationNotFound(rel.Name)
}

func (c *Catalog) GetTable(rel *ast.TableName) (Table, error) {
	_, table, err := c.getTable(rel)
	if table == nil {
		return Table{}, err
	} else {
		return *table, err
	}
}
