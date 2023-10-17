package catalog

import (
	"errors"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

// Function describes a database function
//
// A database function is a method written to perform a specific operation on data within the database.
type Function struct {
	Name               string
	Args               []*Argument
	ReturnType         *ast.TypeName
	Comment            string
	Desc               string
	ReturnTypeNullable bool
}

type Argument struct {
	Name       string
	Type       *ast.TypeName
	HasDefault bool
	Mode       ast.FuncParamMode
}

func (f *Function) InArgs() []*Argument {
	var args []*Argument
	for _, a := range f.Args {
		switch a.Mode {
		case ast.FuncParamTable, ast.FuncParamOut:
			continue
		default:
			args = append(args, a)
		}
	}
	return args
}

func (f *Function) OutArgs() []*Argument {
	var args []*Argument
	for _, a := range f.Args {
		switch a.Mode {
		case ast.FuncParamOut:
			args = append(args, a)
		}
	}
	return args
}

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
