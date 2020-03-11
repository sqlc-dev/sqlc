package catalog

import (
	"errors"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	sqlerr "github.com/kyleconroy/sqlc/internal/sql/errors"
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
