package validate

import (
	"fmt"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

func Cmd(n ast.Node, name, cmd string) error {
	// TODO: Convert cmd to an enum
	if !(cmd == ":many" || cmd == ":one") {
		return nil
	}
	var list *ast.List
	switch stmt := n.(type) {
	case *ast.SelectStmt:
		return nil
	case *ast.DeleteStmt:
		list = stmt.ReturningList
	case *ast.InsertStmt:
		list = stmt.ReturningList
	case *ast.UpdateStmt:
		list = stmt.ReturningList
	default:
		return nil
	}
	if list == nil || len(list.Items) == 0 {
		return fmt.Errorf("query %q specifies parameter %q without containing a RETURNING clause", name, cmd)
	}
	return nil
}
