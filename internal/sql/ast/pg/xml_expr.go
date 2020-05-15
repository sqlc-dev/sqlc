package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type XmlExpr struct {
	Xpr       ast.Node
	Op        XmlExprOp
	Name      *string
	NamedArgs *ast.List
	ArgNames  *ast.List
	Args      *ast.List
	Xmloption XmlOptionType
	Type      Oid
	Typmod    int32
	Location  int
}

func (n *XmlExpr) Pos() int {
	return n.Location
}
