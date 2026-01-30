package lib

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func WindowFunctions() []*catalog.Function {
	var funcs []*catalog.Function

	funcs = append(funcs, rowNumberFuncs()...)
	funcs = append(funcs, lagLeadFuncs()...)
	funcs = append(funcs, firstLastValueFuncs()...)
	funcs = append(funcs, nthValueFuncs()...)
	funcs = append(funcs, rankFuncs()...)
	funcs = append(funcs, ntileFuncs()...)
	funcs = append(funcs, cumeDistFuncs()...)
	funcs = append(funcs, sessionStartFuncs()...)

	return funcs
}

func rowNumberFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name:       "ROW_NUMBER",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		},
	}
}

func lagLeadFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "LAG",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "LAG",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "Uint32"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "LEAD",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "LEAD",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "Uint32"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
	}
}

func firstLastValueFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "FIRST_VALUE",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "LAST_VALUE",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
	}
}

func nthValueFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "NTH_VALUE",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
	}
}

func rankFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "RANK",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		},
		{
			Name: "DENSE_RANK",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		},
		{
			Name: "PERCENT_RANK",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
	}
}

func ntileFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "NTILE",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Uint64"}},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		},
	}
}

func cumeDistFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name:       "CUME_DIST",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
	}
}

func sessionStartFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name:       "SESSION_START",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}
