package catalog

import (
	"errors"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/sqlerr"
)

func (c *Catalog) alterTable(stmt *ast.AlterTableStmt) error {
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
	if !implemented {
		return nil
	}
	_, table, err := c.getTable(stmt.Table)
	if err != nil {
		return err
	}

	for _, cmd := range stmt.Cmds.Items {
		switch cmd := cmd.(type) {
		case *ast.AlterTableCmd:
			idx := -1

			// Lookup column names for column-related commands
			switch cmd.Subtype {
			case ast.AT_AlterColumnType,
				ast.AT_DropColumn,
				ast.AT_DropNotNull,
				ast.AT_SetNotNull:
				for i, c := range table.Columns {
					if c.Name == *cmd.Name {
						idx = i
						break
					}
				}
				if idx < 0 && !cmd.MissingOk {
					return sqlerr.ColumnNotFound(table.Rel.Name, *cmd.Name)
				}
				// If a missing column is allowed, skip this command
				if idx < 0 && cmd.MissingOk {
					continue
				}
			}

			switch cmd.Subtype {

			case ast.AT_AddColumn:
				for _, c := range table.Columns {
					if c.Name == cmd.Def.Colname {
						return sqlerr.ColumnExists(table.Rel.Name, c.Name)
					}
				}
				table.Columns = append(table.Columns, &Column{
					Name:      cmd.Def.Colname,
					Type:      *cmd.Def.TypeName,
					IsNotNull: cmd.Def.IsNotNull,
					IsArray:   cmd.Def.IsArray,
				})

			case ast.AT_AlterColumnType:
				table.Columns[idx].Type = *cmd.Def.TypeName
				// table.Columns[idx].IsArray = isArray(d.TypeName)

			case ast.AT_DropColumn:
				table.Columns = append(table.Columns[:idx], table.Columns[idx+1:]...)

			case ast.AT_DropNotNull:
				table.Columns[idx].IsNotNull = false

			case ast.AT_SetNotNull:
				table.Columns[idx].IsNotNull = true

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
			}
			if col.Vals != nil {
				typeName := ast.TypeName{
					Name: col.Colname,
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
