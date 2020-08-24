package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type XmlSerialize struct {
	Xmloption XmlOptionType
	Expr      Node
	TypeName  *TypeName
	Location  int
}

func (n *XmlSerialize) Pos() int {
	return n.Location
}
