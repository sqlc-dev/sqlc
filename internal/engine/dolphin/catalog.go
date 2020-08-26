package dolphin

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
)

func NewCatalog() *catalog.Catalog {
	def := "public" // TODO: What is the default database for MySQL?
	return &catalog.Catalog{
		DefaultSchema: def,
		Schemas: []*catalog.Schema{
			&catalog.Schema{
				Name: def,
				Funcs: []*catalog.Function{
					{
						Name: "count",
						Args: []*catalog.Argument{
							{
								Type: &ast.TypeName{Name: "any"},
							},
						},
						ReturnType: &ast.TypeName{Name: "bigint"},
					},
					{
						Name:       "count",
						Args:       []*catalog.Argument{},
						ReturnType: &ast.TypeName{Name: "bigint"},
					},
				},
			},
		},
		Extensions: map[string]struct{}{},
	}
}
