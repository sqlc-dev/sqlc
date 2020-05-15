package validate

import (
	nodes "github.com/lfittl/pg_query_go/nodes"

	"github.com/kyleconroy/sqlc/internal/pg"
)

func InsertStmt(stmt nodes.InsertStmt) error {
	sel, ok := stmt.SelectStmt.(nodes.SelectStmt)
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
		return pg.Error{
			Code:    "42601",
			Message: "INSERT has more target columns than expressions",
		}
	case colsLen < valsLen:
		return pg.Error{
			Code:    "42601",
			Message: "INSERT has more expressions than target columns",
		}
	}
	return nil
}
