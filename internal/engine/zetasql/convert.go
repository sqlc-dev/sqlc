package zetasql

import (
	"log"

	zast "github.com/goccy/go-zetasql/ast"

	"github.com/kyleconroy/sqlc/internal/debug"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type cc struct {
}

func todo(n zast.Node) *ast.TODO {
	if debug.Active {
		log.Printf("zetasql.convert: Unknown node type %T\n", n)
	}
	return &ast.TODO{}
}

func (c *cc) convert(node zast.Node) ast.Node {
	switch n := node.(type) {
	case nil:
		return nil
	default:
		return todo(n)
	}
}
