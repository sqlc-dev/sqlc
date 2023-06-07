package catalog

import (
	"errors"
	"fmt"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/sqlerr"
)

// Table describes a relational database table
//
// A database table is a collection of related data held in a table format within a database.
// It consists of columns and rows.
type Table struct {
	Rel     *ast.TableName
	Columns []*Column
	Comment string
}

func (table *Table) isExistColumn(cmd *ast.AlterTableCmd) (int, error) {
	for i, c := range table.Columns {
		if c.Name == *cmd.Name {
			return i, nil
		}
	}
	if !cmd.MissingOk {
		return -1, sqlerr.ColumnNotFound(table.Rel.Name, *cmd.Name)
	}
	// Missing column is allowed
	return -1, nil
}

func (table *Table) addColumn(cmd *ast.AlterTableCmd) error {
	for _, c := range table.Columns {
		if c.Name == cmd.Def.Colname {
			if !cmd.MissingOk {
				return sqlerr.ColumnExists(table.Rel.Name, cmd.Def.Colname)
			}
			return nil
		}
	}

	table.Columns = append(table.Columns, &Column{
		Name:      cmd.Def.Colname,
		Type:      *cmd.Def.TypeName,
		IsNotNull: cmd.Def.IsNotNull,
		IsArray:   cmd.Def.IsArray,
		Length:    cmd.Def.Length,
	})
	return nil
}

func (table *Table) alterColumnType(cmd *ast.AlterTableCmd) error {
	index, err := table.isExistColumn(cmd)
	if err != nil {
		return err
	}
	if index >= 0 {
		table.Columns[index].Type = *cmd.Def.TypeName
		table.Columns[index].IsArray = cmd.Def.IsArray
	}
	return nil
}

func (table *Table) dropColumn(cmd *ast.AlterTableCmd) error {
	index, err := table.isExistColumn(cmd)
	if err != nil {
		return err
	}
	if index >= 0 {
		table.Columns = append(table.Columns[:index], table.Columns[index+1:]...)
	}
	return nil
}

func (table *Table) dropNotNull(cmd *ast.AlterTableCmd) error {
	index, err := table.isExistColumn(cmd)
	if err != nil {
		return err
	}
	if index >= 0 {
		table.Columns[index].IsNotNull = false
	}
	return nil
}

func (table *Table) setNotNull(cmd *ast.AlterTableCmd) error {
	index, err := table.isExistColumn(cmd)
	if err != nil {
		return err
	}
	if index >= 0 {
		table.Columns[index].IsNotNull = true
	}
	return nil
}

// Column describes a set of data values of a particular type in a relational database table
//
// TODO: Should this just be ast Nodes?
type Column struct {
	Name      string
	Type      ast.TypeName
	IsNotNull bool
	IsArray   bool
	Comment   string
	Length    *int
}

// An interface is used to resolve a circular import between the catalog and compiler packages.
// The createView function requires access to functions in the compiler package to parse the SELECT
// statement that defines the view.
type columnGenerator interface {
	OutputColumns(node ast.Node) ([]*Column, error)
}

func (c *Catalog) getTable(tableName *ast.TableName) (*Schema, *Table, error) {
	schemaName := tableName.Schema
	if schemaName == "" {
		schemaName = c.DefaultSchema
	}
	var schema *Schema
	for i := range c.Schemas {
		if c.Schemas[i].Name == schemaName {
			schema = c.Schemas[i]
			break
		}
	}
	if schema == nil {
		return nil, nil, sqlerr.SchemaNotFound(schemaName)
	}
	table, _, err := schema.getTable(tableName)
	if err != nil {
		return nil, nil, err
	}
	return schema, table, nil
}

func isStmtImplemented(stmt *ast.AlterTableStmt) bool {
	var implemented bool
	for _, item := range stmt.Cmds.Items {
		switch cmd := item.(type) {
		case *ast.AlterTableCmd:
			switch cmd.Subtype {
			case ast.AT_AddColumn:
				implemented = true
			case ast.AT_AlterColumnType:
				implemented = true
			case ast.AT_DropColumn:
				implemented = true
			case ast.AT_DropNotNull:
				implemented = true
			case ast.AT_SetNotNull:
				implemented = true
			}
		}
	}
	return implemented
}

func (c *Catalog) alterTable(stmt *ast.AlterTableStmt) error {
	if !isStmtImplemented(stmt) {
		return nil
	}
	_, table, err := c.getTable(stmt.Table)
	if err != nil {
		return err
	}
	for _, item := range stmt.Cmds.Items {
		switch cmd := item.(type) {
		case *ast.AlterTableCmd:
			switch cmd.Subtype {
			case ast.AT_AddColumn:
				if err := table.addColumn(cmd); err != nil {
					return err
				}
			case ast.AT_AlterColumnType:
				if err := table.alterColumnType(cmd); err != nil {
					return err
				}
			case ast.AT_DropColumn:
				if err := table.dropColumn(cmd); err != nil {
					return err
				}
			case ast.AT_DropNotNull:
				if err := table.dropNotNull(cmd); err != nil {
					return err
				}
			case ast.AT_SetNotNull:
				if err := table.setNotNull(cmd); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (c *Catalog) alterTableSetSchema(stmt *ast.AlterTableSetSchemaStmt) error {
	ns := stmt.Table.Schema
	if ns == "" {
		ns = c.DefaultSchema
	}
	oldSchema, err := c.getSchema(ns)
	if err != nil {
		return err
	}
	tbl, idx, err := oldSchema.getTable(stmt.Table)
	if err != nil {
		return err
	}
	tbl.Rel.Schema = *stmt.NewSchema
	newSchema, err := c.getSchema(*stmt.NewSchema)
	if err != nil {
		return err
	}
	if _, _, err := newSchema.getTable(stmt.Table); err == nil {
		return sqlerr.RelationExists(stmt.Table.Name)
	}
	oldSchema.Tables = append(oldSchema.Tables[:idx], oldSchema.Tables[idx+1:]...)
	newSchema.Tables = append(newSchema.Tables, tbl)
	return nil
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
	_, _, err = schema.getTable(stmt.Name)
	if err == nil && stmt.IfNotExists {
		return nil
	} else if err == nil {
		return sqlerr.RelationExists(stmt.Name.Name)
	}

	tbl := Table{Rel: stmt.Name, Comment: stmt.Comment}
	for _, inheritTable := range stmt.Inherits {
		t, _, err := schema.getTable(inheritTable)
		if err != nil {
			return err
		}
		tbl.Columns = append(tbl.Columns, t.Columns...)
	}

	if stmt.ReferTable != nil && len(stmt.Cols) != 0 {
		return errors.New("create table node cannot have both a ReferTable and Cols")
	}

	if stmt.ReferTable != nil {
		_, original, err := c.getTable(stmt.ReferTable)
		if err != nil {
			return err
		}
		for _, col := range original.Columns {
			newCol := *col // make a copy, so changes to the ReferTable don't propagate
			tbl.Columns = append(tbl.Columns, &newCol)
		}
	} else {
		for _, col := range stmt.Cols {
			tc := &Column{
				Name:      col.Colname,
				Type:      *col.TypeName,
				IsNotNull: col.IsNotNull,
				IsArray:   col.IsArray,
				Comment:   col.Comment,
				Length:    col.Length,
			}
			if col.Vals != nil {
				typeName := ast.TypeName{
					Name: fmt.Sprintf("%s_%s", stmt.Name.Name, col.Colname),
				}
				s := &ast.CreateEnumStmt{TypeName: &typeName, Vals: col.Vals}
				if err := c.createEnum(s); err != nil {
					return err
				}
				tc.Type = typeName
			}
			tbl.Columns = append(tbl.Columns, tc)
		}
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
		if errors.Is(err, sqlerr.NotFound) && stmt.IfExists {
			continue
		} else if err != nil {
			return err
		}

		_, idx, err := schema.getTable(name)
		if errors.Is(err, sqlerr.NotFound) && stmt.IfExists {
			continue
		} else if err != nil {
			return err
		}

		schema.Tables = append(schema.Tables[:idx], schema.Tables[idx+1:]...)
	}
	return nil
}

func (c *Catalog) renameColumn(stmt *ast.RenameColumnStmt) error {
	_, tbl, err := c.getTable(stmt.Table)
	if err != nil {
		return err
	}
	idx := -1
	for i := range tbl.Columns {
		if tbl.Columns[i].Name == stmt.Col.Name {
			idx = i
		}
		if tbl.Columns[i].Name == *stmt.NewName {
			return sqlerr.ColumnExists(tbl.Rel.Name, *stmt.NewName)
		}
	}
	if idx == -1 {
		return sqlerr.ColumnNotFound(tbl.Rel.Name, stmt.Col.Name)
	}
	tbl.Columns[idx].Name = *stmt.NewName
	return nil
}

func (c *Catalog) renameTable(stmt *ast.RenameTableStmt) error {
	sch, tbl, err := c.getTable(stmt.Table)
	if err != nil {
		return err
	}
	if _, _, err := sch.getTable(&ast.TableName{Name: *stmt.NewName}); err == nil {
		return sqlerr.RelationExists(*stmt.NewName)
	}
	if stmt.NewName != nil {
		tbl.Rel.Name = *stmt.NewName
	}
	return nil
}

func (c *Catalog) createTableAs(stmt *ast.CreateTableAsStmt, colGen columnGenerator) error {
	cols, err := colGen.OutputColumns(stmt.Query)
	if err != nil {
		return err
	}

	catName := ""
	if stmt.Into.Rel.Catalogname != nil {
		catName = *stmt.Into.Rel.Catalogname
	}
	schemaName := ""
	if stmt.Into.Rel.Schemaname != nil {
		schemaName = *stmt.Into.Rel.Schemaname
	}

	tbl := Table{
		Rel: &ast.TableName{
			Catalog: catName,
			Schema:  schemaName,
			Name:    *stmt.Into.Rel.Relname,
		},
		Columns: cols,
	}

	ns := tbl.Rel.Schema
	if ns == "" {
		ns = c.DefaultSchema
	}
	schema, err := c.getSchema(ns)
	if err != nil {
		return err
	}
	_, _, err = schema.getTable(tbl.Rel)
	if err == nil {
		return sqlerr.RelationExists(tbl.Rel.Name)
	}

	schema.Tables = append(schema.Tables, &tbl)

	return nil
}
