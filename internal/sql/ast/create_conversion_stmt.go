package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateConversionStmt struct {
	ConversionName  *ast.List
	ForEncodingName *string
	ToEncodingName  *string
	FuncName        *ast.List
	Def             bool
}

func (n *CreateConversionStmt) Pos() int {
	return 0
}
