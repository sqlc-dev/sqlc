package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateOpClassItem struct {
	Itemtype    int
	Name        *ObjectWithArgs
	Number      int
	OrderFamily *List
	ClassArgs   *List
	Storedtype  *TypeName
}

func (n *CreateOpClassItem) Pos() int {
	return 0
}
