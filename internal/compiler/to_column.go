package compiler

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast/pg"
	"github.com/kyleconroy/sqlc/internal/sql/astutils"
)

func isArray(n *pg.TypeName) bool {
	if n == nil {
		return false
	}
	return len(n.ArrayBounds.Items) > 0
}

func toColumn(n *pg.TypeName) *Column {
	if n == nil {
		panic("can't build column for nil type name")
	}
	return &Column{
		DataType: astutils.Join(n.Names, "."),
		NotNull:  true, // XXX: How do we know if this should be null?
		IsArray:  isArray(n),
	}
}
