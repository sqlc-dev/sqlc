package rewrite

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// A 3-part qualified call parses as catalog.schema.name, so
// sqlc.jsonb_build_object."Book"(...) has Catalog="sqlc",
// Schema="jsonb_build_object", Name="Book", and sqlc.embed.jsonb(t) has
// Catalog="sqlc", Schema="embed", Name="jsonb".

// IsJSONCall reports whether node is sqlc.jsonb_build_object."Name"(...).
func IsJSONCall(node ast.Node) bool {
	call, ok := node.(*ast.FuncCall)
	if !ok || call.Func == nil {
		return false
	}
	return call.Func.Catalog == "sqlc" && call.Func.Schema == "jsonb_build_object"
}

// IsEmbedJSONCall reports whether node is sqlc.embed.jsonb(table).
func IsEmbedJSONCall(node ast.Node) bool {
	call, ok := node.(*ast.FuncCall)
	if !ok || call.Func == nil {
		return false
	}
	return call.Func.Catalog == "sqlc" && call.Func.Schema == "embed" && call.Func.Name == "jsonb"
}
