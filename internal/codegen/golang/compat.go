package golang

import (
	"github.com/kyleconroy/sqlc/internal/core"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

func sameTableName(n *ast.TableName, f core.FQN, defaultSchema string) bool {
	if n == nil {
		return false
	}
	schema := n.Schema
	if n.Schema == "" {
		schema = defaultSchema
	}
	return n.Catalog == f.Catalog && ((f.Schema != nil && f.Schema.Match(schema)) || schema == "") && ((f.Rel != nil && f.Rel.Match(n.Name)) || n.Name == "")
}
