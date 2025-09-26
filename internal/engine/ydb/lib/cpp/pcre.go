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
			Name: "Pcre::Grep",
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
			Name: "Pcre::Match",
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
			Name: "Pcre::BacktrackingGrep",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Pcre::BacktrackingMatch",
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
			Name: "Pcre::MultiGrep",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Pcre::MultiMatch",
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
			Name: "Pcre::Capture",
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
			Name: "Pcre::Replace",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}
