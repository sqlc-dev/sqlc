package catalog

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/sqlerr"
)

func stringSlice(list *ast.List) []string {
	items := []string{}
	for _, item := range list.Items {
		if n, ok := item.(*ast.String); ok {
			items = append(items, n.Str)
		}
	}
	return items
}

type Catalog struct {
	Comment       string
	DefaultSchema string
	Name          string
	Schemas       []*Schema
	SearchPath    []string
	LoadExtension func(string) *Schema

	// TODO: un-export
	Extensions map[string]struct{}
}

func (c *Catalog) getSchema(name string) (*Schema, error) {
	for i := range c.Schemas {
		if c.Schemas[i].Name == name {
			return c.Schemas[i], nil
		}
	}
	return nil, sqlerr.SchemaNotFound(name)
}

func (c *Catalog) getFunc(rel *ast.FuncName, tns []*ast.TypeName) (*Function, int, error) {
	ns := rel.Schema
	if ns == "" {
		ns = c.DefaultSchema
	}
	s, err := c.getSchema(ns)
	if err != nil {
		return nil, -1, err
	}
	return s.getFunc(rel, tns)
}

func (c *Catalog) getTable(name *ast.TableName) (*Schema, *Table, error) {
	ns := name.Schema
	if ns == "" {
		ns = c.DefaultSchema
	}
	var s *Schema
	for i := range c.Schemas {
		if c.Schemas[i].Name == ns {
			s = c.Schemas[i]
			break
		}
	}
	if s == nil {
		return nil, nil, sqlerr.SchemaNotFound(ns)
	}
	t, _, err := s.getTable(name)
	if err != nil {
		return nil, nil, err
	}
	return s, t, nil
}

func (c *Catalog) getType(rel *ast.TypeName) (Type, int, error) {
	ns := rel.Schema
	if ns == "" {
		ns = c.DefaultSchema
	}
	s, err := c.getSchema(ns)
	if err != nil {
		return nil, -1, err
	}
	return s.getType(rel)
}

type Schema struct {
	Name   string
	Tables []*Table
	Types  []Type
	Funcs  []*Function

	Comment string
}

func sameType(a, b *ast.TypeName) bool {
	if a.Catalog != b.Catalog {
		return false
	}
	// The pg_catalog schema is searched by default, so take that into
	// account when comparing schemas
	aSchema := a.Schema
	bSchema := b.Schema
	if aSchema == "pg_catalog" {
		aSchema = ""
	}
	if bSchema == "pg_catalog" {
		bSchema = ""
	}
	if aSchema != bSchema {
		return false
	}
	if a.Name != b.Name {
		return false
	}
	return true
}

func (s *Schema) getFunc(rel *ast.FuncName, tns []*ast.TypeName) (*Function, int, error) {
	for i := range s.Funcs {
		if s.Funcs[i].Name != rel.Name {
			continue
		}

		args := s.Funcs[i].InArgs()
		if len(args) != len(tns) {
			continue
		}
		found := true
		for j := range args {
			if !sameType(s.Funcs[i].Args[j].Type, tns[j]) {
				found = false
				break
			}
		}
		if !found {
			continue
		}
		return s.Funcs[i], i, nil
	}
	return nil, -1, sqlerr.RelationNotFound(rel.Name)
}

func (s *Schema) getFuncByName(rel *ast.FuncName) (*Function, int, error) {
	idx := -1
	for i := range s.Funcs {
		if s.Funcs[i].Name == rel.Name && idx >= 0 {
			return nil, -1, sqlerr.FunctionNotUnique(rel.Name)
		}
		if s.Funcs[i].Name == rel.Name {
			idx = i
		}
	}
	if idx < 0 {
		return nil, -1, sqlerr.RelationNotFound(rel.Name)
	}
	return s.Funcs[idx], idx, nil
}

func (s *Schema) getTable(rel *ast.TableName) (*Table, int, error) {
	for i := range s.Tables {
		if s.Tables[i].Rel.Name == rel.Name {
			return s.Tables[i], i, nil
		}
	}
	return nil, -1, sqlerr.RelationNotFound(rel.Name)
}

func (s *Schema) getType(rel *ast.TypeName) (Type, int, error) {
	for i := range s.Types {
		switch typ := s.Types[i].(type) {
		case *Enum:
			if typ.Name == rel.Name {
				return s.Types[i], i, nil
			}
		}
	}
	return nil, -1, sqlerr.TypeNotFound(rel.Name)
}

type Table struct {
	Rel     *ast.TableName
	Columns []*Column
	Comment string
}

// TODO: Should this just be ast Nodes?
type Column struct {
	Name      string
	Type      ast.TypeName
	IsNotNull bool
	IsArray   bool
	Comment   string
}

type Type interface {
	isType()

	SetComment(string)
}

type Enum struct {
	Name    string
	Vals    []string
	Comment string
}

func (e *Enum) SetComment(c string) {
	e.Comment = c
}

func (e *Enum) isType() {
}

type CompositeType struct {
	Name    string
	Comment string
}

func (ct *CompositeType) isType() {
}

func (ct *CompositeType) SetComment(c string) {
	ct.Comment = c
}

type Function struct {
	Name       string
	Args       []*Argument
	ReturnType *ast.TypeName
	Comment    string
	Desc       string
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

type Argument struct {
	Name       string
	Type       *ast.TypeName
	HasDefault bool
	Mode       ast.FuncParamMode
}

func New(def string) *Catalog {
	return &Catalog{
		DefaultSchema: def,
		Schemas: []*Schema{
			&Schema{Name: def},
		},
		Extensions: map[string]struct{}{},
	}
}

func (c *Catalog) Build(stmts []ast.Statement) error {
	for i := range stmts {
		if err := c.Update(stmts[i]); err != nil {
			return err
		}
	}
	return nil
}

func (c *Catalog) Update(stmt ast.Statement) error {
	if stmt.Raw == nil {
		return nil
	}
	var err error
	switch n := stmt.Raw.Stmt.(type) {
	case *ast.AlterTableStmt:
		err = c.alterTable(n)

	case *ast.AlterTableSetSchemaStmt:
		err = c.alterTableSetSchema(n)

	case *ast.AlterTypeAddValueStmt:
		err = c.alterTypeAddValue(n)

	case *ast.AlterTypeRenameValueStmt:
		err = c.alterTypeRenameValue(n)

	case *ast.CommentOnColumnStmt:
		err = c.commentOnColumn(n)

	case *ast.CommentOnSchemaStmt:
		err = c.commentOnSchema(n)

	case *ast.CommentOnTableStmt:
		err = c.commentOnTable(n)

	case *ast.CommentOnTypeStmt:
		err = c.commentOnType(n)

	case *ast.CompositeTypeStmt:
		err = c.createCompositeType(n)

	case *ast.CreateEnumStmt:
		err = c.createEnum(n)

	case *ast.CreateExtensionStmt:
		err = c.createExtension(n)

	case *ast.CreateFunctionStmt:
		err = c.createFunction(n)

	case *ast.CreateSchemaStmt:
		err = c.createSchema(n)

	case *ast.CreateTableStmt:
		err = c.createTable(n)

	case *ast.DropFunctionStmt:
		err = c.dropFunction(n)

	case *ast.DropSchemaStmt:
		err = c.dropSchema(n)

	case *ast.DropTableStmt:
		err = c.dropTable(n)

	case *ast.DropTypeStmt:
		err = c.dropType(n)

	case *ast.RenameColumnStmt:
		err = c.renameColumn(n)

	case *ast.RenameTableStmt:
		err = c.renameTable(n)

	}
	return err
}
