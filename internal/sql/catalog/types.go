package catalog

import (
	"errors"
	"fmt"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/sqlerr"
)

func (c *Catalog) createEnum(stmt *ast.CreateEnumStmt) error {
	ns := stmt.TypeName.Schema
	if ns == "" {
		ns = c.DefaultSchema
	}
	schema, err := c.getSchema(ns)
	if err != nil {
		return err
	}
	// Because tables have associated data types, the type name must also
	// be distinct from the name of any existing table in the same
	// schema.
	// https://www.postgresql.org/docs/current/sql-createtype.html
	tbl := &ast.TableName{
		Name: stmt.TypeName.Name,
	}
	if _, _, err := schema.getTable(tbl); err == nil {
		return sqlerr.RelationExists(tbl.Name)
	}
	if _, _, err := schema.getType(stmt.TypeName); err == nil {
		return sqlerr.TypeExists(tbl.Name)
	}
	schema.Types = append(schema.Types, &Enum{
		Name: stmt.TypeName.Name,
		Vals: stringSlice(stmt.Vals),
	})
	return nil
}

func (c *Catalog) createCompositeType(stmt *ast.CompositeTypeStmt) error {
	ns := stmt.TypeName.Schema
	if ns == "" {
		ns = c.DefaultSchema
	}
	schema, err := c.getSchema(ns)
	if err != nil {
		return err
	}
	// Because tables have associated data types, the type name must also
	// be distinct from the name of any existing table in the same
	// schema.
	// https://www.postgresql.org/docs/current/sql-createtype.html
	tbl := &ast.TableName{
		Name: stmt.TypeName.Name,
	}
	if _, _, err := schema.getTable(tbl); err == nil {
		return sqlerr.RelationExists(tbl.Name)
	}
	if _, _, err := schema.getType(stmt.TypeName); err == nil {
		return sqlerr.TypeExists(tbl.Name)
	}
	schema.Types = append(schema.Types, &CompositeType{
		Name: stmt.TypeName.Name,
	})
	return nil
}

func (c *Catalog) alterTypeRenameValue(stmt *ast.AlterTypeRenameValueStmt) error {
	ns := stmt.Type.Schema
	if ns == "" {
		ns = c.DefaultSchema
	}
	schema, err := c.getSchema(ns)
	if err != nil {
		return err
	}
	typ, _, err := schema.getType(stmt.Type)
	if err != nil {
		return err
	}
	enum, ok := typ.(*Enum)
	if !ok {
		return fmt.Errorf("type is not an enum: %T", stmt.Type)
	}

	oldIndex := -1
	newIndex := -1
	for i, val := range enum.Vals {
		if val == *stmt.OldValue {
			oldIndex = i
		}
		if val == *stmt.NewValue {
			newIndex = i
		}
	}
	if oldIndex < 0 {
		return fmt.Errorf("type %T does not have value %s", stmt.Type, *stmt.OldValue)
	}
	if newIndex >= 0 {
		return fmt.Errorf("type %T already has value %s", stmt.Type, *stmt.NewValue)
	}
	enum.Vals[oldIndex] = *stmt.NewValue
	return nil
}

func (c *Catalog) alterTypeAddValue(stmt *ast.AlterTypeAddValueStmt) error {
	ns := stmt.Type.Schema
	if ns == "" {
		ns = c.DefaultSchema
	}
	schema, err := c.getSchema(ns)
	if err != nil {
		return err
	}
	typ, _, err := schema.getType(stmt.Type)
	if err != nil {
		return err
	}
	enum, ok := typ.(*Enum)
	if !ok {
		return fmt.Errorf("type is not an enum: %T", stmt.Type)
	}

	newIndex := -1
	for i, val := range enum.Vals {
		if val == *stmt.NewValue {
			newIndex = i
		}
	}
	if newIndex >= 0 {
		if !stmt.SkipIfNewValExists {
			return fmt.Errorf("type %T already has value %s", stmt.Type, *stmt.NewValue)
		} else {
			return nil
		}
	}
	enum.Vals = append(enum.Vals, *stmt.NewValue)
	return nil
}

func (c *Catalog) dropType(stmt *ast.DropTypeStmt) error {
	for _, name := range stmt.Types {
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

		_, idx, err := schema.getType(name)
		if errors.Is(err, sqlerr.NotFound) && stmt.IfExists {
			continue
		} else if err != nil {
			return err
		}

		schema.Types = append(schema.Types[:idx], schema.Types[idx+1:]...)
	}
	return nil
}
