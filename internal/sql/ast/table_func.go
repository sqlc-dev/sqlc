package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type TableFunc struct {
	NsUris        *ast.List
	NsNames       *ast.List
	Docexpr       ast.Node
	Rowexpr       ast.Node
	Colnames      *ast.List
	Coltypes      *ast.List
	Coltypmods    *ast.List
	Colcollations *ast.List
	Colexprs      *ast.List
	Coldefexprs   *ast.List
	Notnulls      []uint32
	Ordinalitycol int
	Location      int
}

func (n *TableFunc) Pos() int {
	return n.Location
}
