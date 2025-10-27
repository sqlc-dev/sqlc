package cpp

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func PcreFunctions() []*catalog.Function {
	var funcs []*catalog.Function

	funcs = append(funcs, pcreGrepFuncs()...)
	funcs = append(funcs, pcreMatchFuncs()...)
	funcs = append(funcs, pcreBacktrackingFuncs()...)
	funcs = append(funcs, pcreMultiFuncs()...)
	funcs = append(funcs, pcreCaptureFuncs()...)
	funcs = append(funcs, pcreReplaceFuncs()...)

	return funcs
}

func pcreGrepFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "pcre_grep",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func pcreMatchFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "pcre_match",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func pcreBacktrackingFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "pcre_backtrackinggrep",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "pcre_backtrackingmatch",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func pcreMultiFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "pcre_multigrep",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "pcre_multimatch",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func pcreCaptureFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "pcre_capture",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func pcreReplaceFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "pcre_replace",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}
