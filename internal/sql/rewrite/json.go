package rewrite

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// IsJSONCall reports whether node is sqlc.jsonb_build_object."Name"(...). A
// 3-part qualified call parses as catalog.schema.name, so
// sqlc.jsonb_build_object."Book"(...) has Catalog="sqlc",
// Schema="jsonb_build_object", Name="Book".
func IsJSONCall(node ast.Node) bool {
	call, ok := node.(*ast.FuncCall)
	if !ok || call.Func == nil {
		return false
	}
	return call.Func.Catalog == "sqlc" && call.Func.Schema == "jsonb_build_object"
}
