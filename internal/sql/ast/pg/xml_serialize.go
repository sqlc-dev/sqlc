package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type XmlSerialize struct {
	Xmloption XmlOptionType
	Expr      ast.Node
	TypeName  *TypeName
	Location  int
}

func (n *XmlSerialize) Pos() int {
	return n.Location
}
