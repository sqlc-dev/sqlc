package cpp

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func PireFunctions() []*catalog.Function {
	var funcs []*catalog.Function

	funcs = append(funcs, pireGrepFuncs()...)
	funcs = append(funcs, pireMatchFuncs()...)
	funcs = append(funcs, pireMultiFuncs()...)
	funcs = append(funcs, pireCaptureFuncs()...)
	funcs = append(funcs, pireReplaceFuncs()...)

	return funcs
}

func pireGrepFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Pire::Grep",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func pireMatchFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Pire::Match",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func pireMultiFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Pire::MultiGrep",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Pire::MultiMatch",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func pireCaptureFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Pire::Capture",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func pireReplaceFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Pire::Replace",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}
