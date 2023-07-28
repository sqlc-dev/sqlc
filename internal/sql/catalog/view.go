package catalog

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

func (c *Catalog) createView(stmt *ast.ViewStmt, colGen columnGenerator) error {
	cols, err := colGen.OutputColumns(stmt.Query)
	if err != nil {
		return err
	}

	catName := ""
	if stmt.View.Catalogname != nil {
		catName = *stmt.View.Catalogname
	}
	schemaName := ""
	if stmt.View.Schemaname != nil {
		schemaName = *stmt.View.Schemaname
	}

	tbl := Table{
		Rel: &ast.TableName{
			Catalog: catName,
			Schema:  schemaName,
			Name:    *stmt.View.Relname,
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
	_, existingIdx, err := schema.getTable(tbl.Rel)
	if err == nil && !stmt.Replace {
		return sqlerr.RelationExists(tbl.Rel.Name)
	}

	if stmt.Replace && err == nil {
		schema.Tables[existingIdx] = &tbl
	} else {
		schema.Tables = append(schema.Tables, &tbl)
	}

	return nil
}
