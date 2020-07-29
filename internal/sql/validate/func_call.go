package validate

import (
	"errors"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/astutils"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
	"github.com/kyleconroy/sqlc/internal/sql/sqlerr"
)

type funcCallVisitor struct {
	catalog *catalog.Catalog
	err     error
}

func (v *funcCallVisitor) Visit(node ast.Node) astutils.Visitor {
	if v.err != nil {
		return nil
	}

	call, ok := node.(*ast.FuncCall)
	if !ok {
		return v
	}
	if call.Func == nil {
		return v
	}

	fun, err := v.catalog.ResolveFuncCall(call)
	if fun != nil || errors.Is(err, sqlerr.NotFound) {
		return v
	}

	v.err = err
	return nil
}

func FuncCall(c *catalog.Catalog, n ast.Node) error {
	visitor := funcCallVisitor{catalog: c}
	astutils.Walk(&visitor, n)
	return visitor.err
}
