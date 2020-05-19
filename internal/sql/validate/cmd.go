package validate

import (
	"fmt"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/ast/pg"
)

func Cmd(n ast.Node, name, cmd string) error {
	// TODO: Convert cmd to an enum
	if !(cmd == ":many" || cmd == ":one") {
		return nil
	}
	var list *ast.List
	switch stmt := n.(type) {
	case *pg.SelectStmt:
		return nil
	case *pg.DeleteStmt:
		list = stmt.ReturningList
	case *pg.InsertStmt:
		list = stmt.ReturningList
	case *pg.UpdateStmt:
		list = stmt.ReturningList
	default:
		return nil
	}
	if list == nil || len(list.Items) == 0 {
		return fmt.Errorf("query %q specifies parameter %q without containing a RETURNING clause", name, cmd)
	}
	return nil
}
