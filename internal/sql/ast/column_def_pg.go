package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type ColumnDef struct {
	Colname       *string
	TypeName      *TypeName
	Inhcount      int
	IsLocal       bool
	IsNotNull     bool
	IsFromType    bool
	IsFromParent  bool
	Storage       byte
	RawDefault    ast.Node
	CookedDefault ast.Node
	Identity      byte
	CollClause    *CollateClause
	CollOid       Oid
	Constraints   *ast.List
	Fdwoptions    *ast.List
	Location      int
}

func (n *ColumnDef) Pos() int {
	return n.Location
}
