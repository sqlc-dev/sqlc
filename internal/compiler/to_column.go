package compiler

import (
	"strings"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
)

func isArray(n *ast.TypeName) bool {
	if n == nil || n.ArrayBounds == nil {
		return false
	}
	return len(n.ArrayBounds.Items) > 0
}

func arrayDims(n *ast.TypeName) int {
	if n == nil || n.ArrayBounds == nil {
		return 0
	}
	return len(n.ArrayBounds.Items)
}

func toColumn(n *ast.TypeName) *Column {
	if n == nil {
		panic("can't build column for nil type name")
	}
	typ, err := ParseTypeName(n)
	if err != nil {
		panic("toColumn: " + err.Error())
	}
	return &Column{
		Type:      typ,
		DataType:  strings.TrimPrefix(astutils.Join(n.Names, "."), "."),
		NotNull:   true, // XXX: How do we know if this should be null?
		IsArray:   isArray(n),
		ArrayDims: arrayDims(n),
	}
}
