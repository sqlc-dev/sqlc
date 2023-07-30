package validate

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

func InsertStmt(stmt *ast.InsertStmt) error {
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
	return nil
}
