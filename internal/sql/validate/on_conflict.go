package validate

import (
	"fmt"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

func OnConflictClause(cat *catalog.Catalog, stmt *ast.InsertStmt, tableName *ast.TableName) error {
	if stmt.OnConflictClause == nil {
		return nil
	}

	occ := stmt.OnConflictClause

	if occ.Action != ast.OnConflictActionUpdate {
		return nil
	}

	if tableName == nil {
		return nil
	}

	tbl, err := cat.GetTable(tableName)
	if err != nil {
		return err
	}

	relName := ""
	if tbl.Rel != nil {
		relName = tbl.Rel.Name
	}

	validCols := make(map[string]struct{}, len(tbl.Columns))
	for _, c := range tbl.Columns {
		validCols[strings.ToLower(c.Name)] = struct{}{}
	}

	if occ.TargetList == nil {
		return nil
	}

	for _, item := range occ.TargetList.Items {
		res, ok := item.(*ast.ResTarget)
		if !ok {
			continue
		}

		if res.Name != nil {
			colName := strings.ToLower(*res.Name)
			if _, exists := validCols[colName]; !exists {
				return &sqlerr.Error{
					Code:     "42703",
					Message:  fmt.Sprintf("column %q of relation %q does not exist", *res.Name, relName),
					Location: res.Location,
				}
			}
		}

		if res.Val != nil {
			if err := validateExcludedRefs(res.Val, validCols, relName); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateExcludedRefs(node ast.Node, validCols map[string]struct{}, tableName string) error {
	refs := astutils.Search(node, func(n ast.Node) bool {
		_, ok := n.(*ast.ColumnRef)
		return ok
	})

	for _, ref := range refs.Items {
		colRef, ok := ref.(*ast.ColumnRef)
		if !ok {
			continue
		}

		parts := make([]string, 0, len(colRef.Fields.Items))
		for _, field := range colRef.Fields.Items {
			if s, ok := field.(*ast.String); ok {
				parts = append(parts, s.Str)
			}
		}

		if len(parts) == 2 && strings.ToLower(parts[0]) == "excluded" {
			colName := strings.ToLower(parts[1])
			if _, exists := validCols[colName]; !exists {
				return &sqlerr.Error{
					Code:     "42703",
					Message:  fmt.Sprintf("column %q does not exist in relation %q (via EXCLUDED)", parts[1], tableName),
					Location: colRef.Location,
				}
			}
		}
	}

	return nil
}
