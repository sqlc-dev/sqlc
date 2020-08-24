package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateTableSpaceStmt struct {
	Tablespacename *string
	Owner          *RoleSpec
	Location       *string
	Options        *List
}

func (n *CreateTableSpaceStmt) Pos() int {
	return 0
}
