package sqlite

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
)

// TODO: fill out sqlite functions from:
// 		 https://www.sqlite.org/lang_aggfunc.html
// 		 https://www.sqlite.org/lang_mathfunc.html
//		 https://www.sqlite.org/lang_corefunc.html

func defaultSchema(name string) *catalog.Schema {
	s := &catalog.Schema{Name: name}
	s.Funcs = []*catalog.Function{
		{
			Name:       "COUNT",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "bigint"},
		},
		{
			Name: "COUNT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bigint"},
		},
		{
			Name: "SUM",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "real"},
				},
			},
			ReturnType: &ast.TypeName{Name: "real"},
		},
		{
			Name: "SUM",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "integer"},
				},
			},
			ReturnType: &ast.TypeName{Name: "integer"},
		},
		{
			Name: "SUM",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
				},
			},
			ReturnType: &ast.TypeName{Name: "bigint"},
		},
	}
	return s
}
