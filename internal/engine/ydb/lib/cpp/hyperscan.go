package cpp

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func HyperscanFunctions() []*catalog.Function {
	var funcs []*catalog.Function

	funcs = append(funcs, hyperscanGrepFuncs()...)
	funcs = append(funcs, hyperscanMatchFuncs()...)
	funcs = append(funcs, hyperscanBacktrackingFuncs()...)
	funcs = append(funcs, hyperscanMultiFuncs()...)
	funcs = append(funcs, hyperscanCaptureFuncs()...)
	funcs = append(funcs, hyperscanReplaceFuncs()...)

	return funcs
}

func hyperscanGrepFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Hyperscan::Grep",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func hyperscanMatchFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Hyperscan::Match",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func hyperscanBacktrackingFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Hyperscan::BacktrackingGrep",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Hyperscan::BacktrackingMatch",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func hyperscanMultiFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Hyperscan::MultiGrep",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "Hyperscan::MultiMatch",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func hyperscanCaptureFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Hyperscan::Capture",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func hyperscanReplaceFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Hyperscan::Replace",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}
