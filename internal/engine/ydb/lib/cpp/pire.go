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
			Name: "pire_grep",
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
			Name: "pire_match",
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
			Name: "pire_multigrep",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "pire_multimatch",
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
			Name: "pire_capture",
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
			Name: "pire_replace",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}
