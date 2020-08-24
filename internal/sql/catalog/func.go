package catalog

import (
	"errors"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/sqlerr"
)

func (c *Catalog) createFunction(stmt *ast.CreateFunctionStmt) error {
	ns := stmt.Func.Schema
	if ns == "" {
		ns = c.DefaultSchema
	}
	s, err := c.getSchema(ns)
	if err != nil {
		return err
	}
	fn := &Function{
		Name:       stmt.Func.Name,
		Args:       make([]*Argument, len(stmt.Params.Items)),
		ReturnType: stmt.ReturnType,
	}
	types := make([]*ast.TypeName, len(stmt.Params.Items))
	for i, item := range stmt.Params.Items {
		arg := item.(*ast.FuncParam)
		var name string
		if arg.Name != nil {
			name = *arg.Name
		}
		fn.Args[i] = &Argument{
			Name:       name,
			Type:       arg.Type,
			Mode:       arg.Mode,
			HasDefault: arg.DefExpr != nil,
		}
		types[i] = arg.Type
	}

	_, idx, err := s.getFunc(stmt.Func, types)
	if err == nil && !stmt.Replace {
		return sqlerr.RelationExists(stmt.Func.Name)
	}

	if idx >= 0 {
		s.Funcs[idx] = fn
	} else {
		s.Funcs = append(s.Funcs, fn)
	}
	return nil
}

func (c *Catalog) dropFunction(stmt *ast.DropFunctionStmt) error {
	for _, spec := range stmt.Funcs {
		ns := spec.Name.Schema
		if ns == "" {
			ns = c.DefaultSchema
		}
		s, err := c.getSchema(ns)
		if errors.Is(err, sqlerr.NotFound) && stmt.MissingOk {
			continue
		} else if err != nil {
			return err
		}
		var idx int
		if spec.HasArgs {
			_, idx, err = s.getFunc(spec.Name, spec.Args)
		} else {
			_, idx, err = s.getFuncByName(spec.Name)
		}
		if errors.Is(err, sqlerr.NotFound) && stmt.MissingOk {
			continue
		} else if err != nil {
			return err
		}
		s.Funcs = append(s.Funcs[:idx], s.Funcs[idx+1:]...)
	}
	return nil
}
