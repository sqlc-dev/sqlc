package catalog

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	sqlerr "github.com/kyleconroy/sqlc/internal/sql/errors"
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
	Name    string
	Schemas []*Schema
	Comment string

	DefaultSchema string
}

func (c *Catalog) getSchema(name string) (*Schema, error) {
	for i := range c.Schemas {
		if c.Schemas[i].Name == name {
			return c.Schemas[i], nil
		}
	}
	return nil, sqlerr.SchemaNotFound(name)
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
		return nil, 0, err
	}
	return s.getType(rel)
}

type Schema struct {
	Name    string
	Tables  []*Table
	Types   []Type
	Comment string
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
	return nil, 0, sqlerr.TypeNotFound(rel.Name)
}

func (s *Schema) getTable(rel *ast.TableName) (*Table, int, error) {
	for i := range s.Tables {
		if s.Tables[i].Rel.Name == rel.Name {
			return s.Tables[i], i, nil
		}
	}
	return nil, 0, sqlerr.RelationNotFound(rel.Name)
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

func New(def string) *Catalog {
	return &Catalog{
		DefaultSchema: def,
		Schemas: []*Schema{
			&Schema{Name: def},
		},
	}
}

func (c *Catalog) Build(stmts []ast.Statement) error {
	for i := range stmts {
		if stmts[i].Raw == nil {
			continue
		}
		var err error
		switch n := stmts[i].Raw.Stmt.(type) {
		case *ast.AlterTableStmt:
			err = c.alterTable(n)
		case *ast.AlterTableSetSchemaStmt:
			err = c.alterTableSetSchema(n)
		case *ast.CommentOnColumnStmt:
			err = c.commentOnColumn(n)
		case *ast.CommentOnSchemaStmt:
			err = c.commentOnSchema(n)
		case *ast.CommentOnTableStmt:
			err = c.commentOnTable(n)
		case *ast.CommentOnTypeStmt:
			err = c.commentOnType(n)
		case *ast.CreateEnumStmt:
			err = c.createEnum(n)
		case *ast.CreateSchemaStmt:
			err = c.createSchema(n)
		case *ast.CreateTableStmt:
			err = c.createTable(n)
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
		if err != nil {
			return err
		}
	}
	return nil
}
