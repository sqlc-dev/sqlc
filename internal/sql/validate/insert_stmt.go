package validate

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast/pg"
	"github.com/kyleconroy/sqlc/internal/sql/sqlerr"
)

func InsertStmt(stmt *pg.InsertStmt) error {
	sel, ok := stmt.SelectStmt.(*pg.SelectStmt)
	if !ok {
		return nil
	}
	if len(sel.ValuesLists) != 1 {
		return nil
	}

	colsLen := len(stmt.Cols.Items)
	valsLen := len(sel.ValuesLists[0])
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
