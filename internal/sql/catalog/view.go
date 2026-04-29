package catalog

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
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

	// Extract table dependencies from the view's SELECT query
	tbl.DependsOnTables = extractTableDeps(stmt.Query)

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

// extractTableDeps walks the SELECT query AST and returns all table references (RangeVar nodes).
func extractTableDeps(node ast.Node) []*ast.TableName {
	var deps []*ast.TableName
	seen := make(map[string]bool)

	astutils.Walk(astutils.VisitorFunc(func(n ast.Node) {
		rv, ok := n.(*ast.RangeVar)
		if !ok || rv.Relname == nil {
			return
		}
		schema := ""
		if rv.Schemaname != nil {
			schema = *rv.Schemaname
		}
		key := schema + "." + *rv.Relname
		if seen[key] {
			return
		}
		seen[key] = true

		// Skip system catalogs and information schema
		if schema == "pg_catalog" || schema == "information_schema" {
			return
		}

		deps = append(deps, &ast.TableName{
			Schema: schema,
			Name:   *rv.Relname,
		})
	}), node)

	return deps
}
