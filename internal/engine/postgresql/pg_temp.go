package postgresql

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
)

func pgTemp() *catalog.Schema {
	return &catalog.Schema{Name: "pg_temp"}
}

func typeName(name string) *ast.TypeName {
	return &ast.TypeName{Name: name}
}

func argN(name string, n int) *catalog.Function {
	var args []*catalog.Argument
	for i := 0; i < n; i++ {
		args = append(args, &catalog.Argument{
			Type: &ast.TypeName{Name: "any"},
		})
	}
	return &catalog.Function{
		Name:       name,
		Args:       args,
		ReturnType: &ast.TypeName{Name: "any"},
	}
}
