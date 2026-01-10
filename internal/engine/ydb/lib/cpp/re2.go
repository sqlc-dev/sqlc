package cpp

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func Re2Functions() []*catalog.Function {
	var funcs []*catalog.Function

	funcs = append(funcs, re2GrepFuncs()...)
	funcs = append(funcs, re2MatchFuncs()...)
	funcs = append(funcs, re2CaptureFuncs()...)
	funcs = append(funcs, re2FindAndConsumeFuncs()...)
	funcs = append(funcs, re2ReplaceFuncs()...)
	funcs = append(funcs, re2CountFuncs()...)
	funcs = append(funcs, re2OptionsFuncs()...)

	return funcs
}

func re2GrepFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "re2_grep",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "re2_grep",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func re2MatchFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "re2_match",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "re2_match",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func re2CaptureFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "re2_capture",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "re2_capture",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func re2FindAndConsumeFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "re2_findandconsume",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "re2_findandconsume",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func re2ReplaceFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "re2_replace",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "re2_replace",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func re2CountFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "re2_count",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "re2_count",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func re2OptionsFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "re2_options",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "re2_options",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "re2_options",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "re2_options",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "re2_options",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "re2_options",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "re2_options",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Uint64"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "re2_options",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Uint64"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "re2_options",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Uint64"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "re2_options",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Uint64"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "re2_options",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Uint64"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "re2_options",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Uint64"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "re2_options",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Uint64"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "re2_options",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Uint64"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}
