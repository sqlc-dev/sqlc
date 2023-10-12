package validate

import (
	"errors"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

type funcCallVisitor struct {
	catalog  *catalog.Catalog
	settings config.CombinedSettings
	err      error
}

func (v *funcCallVisitor) Visit(node ast.Node) astutils.Visitor {
	if v.err != nil {
		return nil
	}

	call, ok := node.(*ast.FuncCall)
	if !ok {
		return v
	}
	fn := call.Func
	if fn == nil {
		return v
	}

	if fn.Schema == "sqlc" {
		return nil
	}

	fun, err := v.catalog.ResolveFuncCall(call)
	if fun != nil {
		return v
	}
	if errors.Is(err, sqlerr.NotFound) && !v.settings.Package.StrictFunctionChecks {
		return v
	}
	v.err = err
	return nil
}

func FuncCall(c *catalog.Catalog, cs config.CombinedSettings, n ast.Node) error {
	visitor := funcCallVisitor{catalog: c, settings: cs}
	astutils.Walk(&visitor, n)
	return visitor.err
}
