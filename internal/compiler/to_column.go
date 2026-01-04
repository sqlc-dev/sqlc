package compiler

import (
	"fmt"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
)

func arrayDims(n *ast.TypeName) int {
	if n == nil || n.ArrayBounds == nil {
		return 0
	}
	return len(n.ArrayBounds.Items)
}

func toColumn(n *ast.TypeName) (*Column, error) {
	if n == nil {
		return nil, fmt.Errorf("can't build column for nil type name")
	}
	typ, err := ParseTypeName(n)
	if err != nil {
		return nil, fmt.Errorf("toColumn: %w", err)
	}
	arrayDims := arrayDims(n)
	return &Column{
		Type:      typ,
		DataType:  strings.TrimPrefix(astutils.Join(n.Names, "."), "."),
		NotNull:   true, // XXX: How do we know if this should be null?
		IsArray:   arrayDims > 0,
		ArrayDims: arrayDims,
	}, nil
}
