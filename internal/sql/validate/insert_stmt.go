package validate

import (
	"strings"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

func InsertStmt(c *catalog.Catalog, fqn *ast.TableName, stmt *ast.InsertStmt) error {
	sel, ok := stmt.SelectStmt.(*ast.SelectStmt)
	if !ok {
		return nil
	}
	if sel.ValuesLists == nil {
		return nil
	}
	if len(sel.ValuesLists.Items) != 1 {
		return nil
	}
	sublist, ok := sel.ValuesLists.Items[0].(*ast.List)
	if !ok {
		return nil
	}

	colsLen := len(stmt.Cols.Items)
	valsLen := len(sublist.Items)
	switch {
	case colsLen > valsLen:
		return &sqlerr.Error{
			Code:    "42601",
			Message: "INSERT has more target columns than expressions",
		}
	case colsLen < valsLen:
		return &sqlerr.Error{
			Code:    "42601",
			Message: "INSERT has more expressions than target columns",
		}
	}
	return onConflictClause(c, fqn, stmt)
}

// onConflictClause validates an ON CONFLICT DO UPDATE clause against the target
// table. It checks:
//   - ON CONFLICT (col, ...) conflict target columns exist
//   - DO UPDATE SET col = ... assignment target columns exist
//   - EXCLUDED.col references exist
func onConflictClause(c *catalog.Catalog, fqn *ast.TableName, n *ast.InsertStmt) error {
	if n.OnConflictClause == nil || n.OnConflictClause.Action != ast.OnConflictActionUpdate {
		return nil
	}

	table, err := c.GetTable(fqn)
	if err != nil {
		return err
	}

	// Build set of column names for existence checks.
	colNames := make(map[string]struct{}, len(table.Columns))
	for _, col := range table.Columns {
		colNames[col.Name] = struct{}{}
	}

	// Validate ON CONFLICT (col, ...) conflict target columns.
	if n.OnConflictClause.Infer != nil && n.OnConflictClause.Infer.IndexElems != nil {
		for _, item := range n.OnConflictClause.Infer.IndexElems.Items {
			elem, ok := item.(*ast.IndexElem)
			if !ok || elem.Name == nil {
				continue
			}
			if _, exists := colNames[*elem.Name]; !exists {
				e := sqlerr.ColumnNotFound(table.Rel.Name, *elem.Name)
				e.Location = n.OnConflictClause.Infer.Location
				return e
			}
		}
	}

	// Validate DO UPDATE SET col = ... assignment target columns and EXCLUDED.col references.
	if n.OnConflictClause.TargetList == nil {
		return nil
	}
	for _, item := range n.OnConflictClause.TargetList.Items {
		target, ok := item.(*ast.ResTarget)
		if !ok || target.Name == nil {
			continue
		}
		if _, exists := colNames[*target.Name]; !exists {
			e := sqlerr.ColumnNotFound(table.Rel.Name, *target.Name)
			e.Location = target.Location
			return e
		}
		if ref, ok := target.Val.(*ast.ColumnRef); ok {
			if excludedCol, ok := excludedColumnRef(ref); ok {
				if _, exists := colNames[excludedCol]; !exists {
					e := sqlerr.ColumnNotFound(table.Rel.Name, excludedCol)
					e.Location = ref.Location
					return e
				}
			}
		}
	}
	return nil
}

// excludedColumnRef returns the column name if the ColumnRef is an EXCLUDED.col
// reference, and ok=true. Returns "", false otherwise.
func excludedColumnRef(ref *ast.ColumnRef) (string, bool) {
	if ref.Fields == nil || len(ref.Fields.Items) != 2 {
		return "", false
	}
	first, ok := ref.Fields.Items[0].(*ast.String)
	if !ok || !strings.EqualFold(first.Str, "excluded") {
		return "", false
	}
	second, ok := ref.Fields.Items[1].(*ast.String)
	if !ok {
		return "", false
	}
	return second.Str, true
}
