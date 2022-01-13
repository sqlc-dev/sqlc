package validate

import (
	"errors"
	"fmt"

	"github.com/kyleconroy/sqlc/internal/metadata"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

func validateCopyfrom(n ast.Node) error {
	stmt, ok := n.(*ast.InsertStmt)
	if !ok {
		return errors.New(":copyfrom requires an INSERT INTO statement")
	}
	if stmt.OnConflictClause != nil {
		return errors.New(":copyfrom is not compatible with ON CONFLICT")
	}
	if stmt.WithClause != nil {
		return errors.New(":copyfrom is not compatible with WITH clauses")
	}
	if stmt.ReturningList != nil && len(stmt.ReturningList.Items) > 0 {
		return errors.New(":copyfrom is not compatible with RETURNING")
	}
	sel, ok := stmt.SelectStmt.(*ast.SelectStmt)
	if !ok {
		return nil
	}
	if len(sel.FromClause.Items) > 0 {
		return errors.New(":copyfrom is not compatible with INSERT INTO ... SELECT")
	}
	if sel.ValuesLists == nil || len(sel.ValuesLists.Items) != 1 {
		return errors.New(":copyfrom requires exactly one example row to be inserted")
	}
	sublist, ok := sel.ValuesLists.Items[0].(*ast.List)
	if !ok {
		return nil
	}
	for _, v := range sublist.Items {
		if _, ok := v.(*ast.ParamRef); !ok {
			return errors.New(":copyfrom doesn't support non-parameter values")
		}
	}
	return nil
}

func Cmd(n ast.Node, name, cmd string) error {
	if cmd == metadata.CmdCopyFrom {
		return validateCopyfrom(n)
	}
	if !(cmd == metadata.CmdMany || cmd == metadata.CmdOne) {
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
