package compiler

import (
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// returningTables builds virtual tables for the OLD and NEW aliases that
// PostgreSQL 18 makes available in the RETURNING clause of INSERT, UPDATE and
// DELETE statements. Each alias exposes the columns of the statement's target
// table. For INSERT there is usually no old row and for DELETE there is no
// new row, so every column reached through those aliases becomes nullable.
func (c *Compiler) returningTables(qc *QueryCatalog, node ast.Node) ([]*Table, error) {
	if c.conf.Engine != config.EnginePostgreSQL {
		return nil, nil
	}

	var rv *ast.RangeVar
	var returning *ast.List
	oldAlias, newAlias := "old", "new"
	var oldNullable, newNullable bool
	switch n := node.(type) {
	case *ast.DeleteStmt:
		rv = firstRangeVar(n.Relations)
		returning = n.ReturningList
		if n.ReturningOldAlias != "" {
			oldAlias = n.ReturningOldAlias
		}
		if n.ReturningNewAlias != "" {
			newAlias = n.ReturningNewAlias
		}
		// A deleted row has no new value
		newNullable = true
	case *ast.InsertStmt:
		rv = n.Relation
		returning = n.ReturningList
		if n.ReturningOldAlias != "" {
			oldAlias = n.ReturningOldAlias
		}
		if n.ReturningNewAlias != "" {
			newAlias = n.ReturningNewAlias
		}
		// An inserted row has no old value, except when an ON CONFLICT
		// clause updates an existing row instead
		oldNullable = true
	case *ast.UpdateStmt:
		rv = firstRangeVar(n.Relations)
		returning = n.ReturningList
		if n.ReturningOldAlias != "" {
			oldAlias = n.ReturningOldAlias
		}
		if n.ReturningNewAlias != "" {
			newAlias = n.ReturningNewAlias
		}
	default:
		return nil, nil
	}
	if rv == nil || returning == nil || len(returning.Items) == 0 {
		return nil, nil
	}

	fqn, err := ParseTableName(rv)
	if err != nil {
		return nil, err
	}

	build := func(alias string, nullable bool) *Table {
		table, err := qc.GetTable(fqn)
		if err != nil {
			// An unresolvable target table is reported by the regular
			// source table lookup, so ignore the error here
			return nil
		}
		table.Rel = &ast.TableName{Name: alias}
		if nullable {
			for _, col := range table.Columns {
				col.NotNull = false
			}
		}
		return table
	}

	var tables []*Table
	if t := build(oldAlias, oldNullable); t != nil {
		tables = append(tables, t)
	}
	if t := build(newAlias, newNullable); t != nil {
		tables = append(tables, t)
	}
	return tables, nil
}

func firstRangeVar(list *ast.List) *ast.RangeVar {
	if list == nil {
		return nil
	}
	for _, item := range list.Items {
		if rv, ok := item.(*ast.RangeVar); ok && rv != nil {
			return rv
		}
	}
	return nil
}

// returningTableForScope returns the OLD or NEW virtual table named by scope.
// A source table with the same name shadows the virtual table, matching
// PostgreSQL, where the implicit aliases are only available when no relation
// in the query is already known under that name.
func returningTableForScope(tables, rtables []*Table, scope string) *Table {
	if scope == "" {
		return nil
	}
	for _, t := range tables {
		if t.Rel.Name == scope {
			return nil
		}
	}
	for _, t := range rtables {
		if t.Rel.Name == scope {
			return t
		}
	}
	return nil
}

// tablesForRef resolves a column reference against the source tables,
// extended with the OLD or NEW virtual table when the reference is qualified
// with one of their names.
func tablesForRef(ref *ast.ColumnRef, tables, rtables []*Table) []*Table {
	parts := stringSlice(ref.Fields)
	if len(parts) != 2 {
		return tables
	}
	if vt := returningTableForScope(tables, rtables, parts[0]); vt != nil {
		return append(append([]*Table{}, tables...), vt)
	}
	return tables
}
