package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type FuncParamMode int

const (
	FuncParamIn FuncParamMode = iota
	FuncParamOut
	FuncParamInOut
	FuncParamVariadic
	FuncParamTable
	FuncParamDefault
)

type FuncParam struct {
	Name    *string
	Type    *TypeName
	DefExpr Node // Will always be &ast.TODO
	Mode    FuncParamMode
}

func (n *FuncParam) Pos() int {
	return 0
}

func (n *FuncParam) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	// Parameter mode prefix (OUT, INOUT, VARIADIC)
	switch n.Mode {
	case FuncParamOut:
		buf.WriteString("OUT ")
	case FuncParamInOut:
		buf.WriteString("INOUT ")
	case FuncParamVariadic:
		buf.WriteString("VARIADIC ")
	}
	// Parameter name (if present)
	if n.Name != nil {
		buf.WriteString(*n.Name)
		buf.WriteString(" ")
	}
	// Parameter type
	buf.astFormat(n.Type, d)
}
