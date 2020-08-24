package compiler

import (
	"strings"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/astutils"
)

func isArray(n *ast.TypeName) bool {
	if n == nil {
		return false
	}
	return len(n.ArrayBounds.Items) > 0
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
		Type:     typ,
		DataType: strings.TrimPrefix(astutils.Join(n.Names, "."), "."),
		NotNull:  true, // XXX: How do we know if this should be null?
		IsArray:  isArray(n),
	}
}
