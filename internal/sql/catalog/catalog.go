package catalog

import (
	"errors"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

func Build(stmts []ast.Statement) (*Catalog, error) {
	c := &Catalog{
		DefaultSchema: "main", // TODO: Needs to be public for PostgreSQL
		Schemas: []*Schema{
			&Schema{Name: "main"},
		},
	}
	for i := range stmts {
		if stmts[i].Raw == nil {
			continue
		}
		var err error
		switch n := stmts[i].Raw.Stmt.(type) {
		case *ast.CreateTableStmt:
			err = c.createTable(n)
		case *ast.DropTableStmt:
			err = c.dropTable(n)
		}
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

var ErrRelationNotFound = errors.New("relation not found")
var ErrSchemaNotFound = errors.New("schema not found")

func (c *Catalog) getSchema(name string) (*Schema, error) {
	for i := range c.Schemas {
		if c.Schemas[i].Name == name {
			return c.Schemas[i], nil
		}
	}
	return nil, ErrSchemaNotFound
}

func (c *Catalog) createTable(stmt *ast.CreateTableStmt) error {
	ns := stmt.Name.Schema
	if ns == "" {
		ns = c.DefaultSchema
	}
	schema, err := c.getSchema(ns)
	if err != nil {
		return err
	}
	if _, _, err := schema.getTable(stmt.Name); err != nil {
		if !errors.Is(err, ErrRelationNotFound) {
			return err
		}
	} else if stmt.IfNotExists {
		return nil
	}
	tbl := Table{Rel: stmt.Name}
	for _, col := range stmt.Cols {
		tbl.Columns = append(tbl.Columns, &Column{
			Name: col.Colname,
			Type: *col.TypeName,
		})
	}
	schema.Tables = append(schema.Tables, &tbl)
	return nil
}

func (c *Catalog) dropTable(stmt *ast.DropTableStmt) error {
	for _, name := range stmt.Tables {
		ns := name.Schema
		if ns == "" {
			ns = c.DefaultSchema
		}
		schema, err := c.getSchema(ns)
		if errors.Is(err, ErrSchemaNotFound) && stmt.IfExists {
			continue
		} else if err != nil {
			return err
		}

		_, idx, err := schema.getTable(name)
		if errors.Is(err, ErrRelationNotFound) && stmt.IfExists {
			continue
		} else if err != nil {
			return err
		}

		schema.Tables = append(schema.Tables[:idx], schema.Tables[idx+1:]...)
	}
	return nil
}

type Catalog struct {
	Name    string
	Schemas []*Schema
	Comment string

	DefaultSchema string
}

type Schema struct {
	Name    string
	Tables  []*Table
	Comment string
}

func (s *Schema) getTable(rel *ast.TableName) (*Table, int, error) {
	for i := range s.Tables {
		if s.Tables[i].Rel.Name == rel.Name {
			return s.Tables[i], i, nil
		}
	}
	return nil, 0, ErrRelationNotFound
}

type Table struct {
	Rel     *ast.TableName
	Columns []*Column
	Comment string
}

type Column struct {
	Name string
	Type ast.TypeName
}
