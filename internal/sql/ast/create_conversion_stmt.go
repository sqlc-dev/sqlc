package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateConversionStmt struct {
	ConversionName  *List
	ForEncodingName *string
	ToEncodingName  *string
	FuncName        *List
	Def             bool
}

func (n *CreateConversionStmt) Pos() int {
	return 0
}
