package cpp

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func MathFunctions() []*catalog.Function {
	var funcs []*catalog.Function

	funcs = append(funcs, mathConstantsFuncs()...)
	funcs = append(funcs, mathCheckFuncs()...)
	funcs = append(funcs, mathUnaryFuncs()...)
	funcs = append(funcs, mathBinaryFuncs()...)
	funcs = append(funcs, mathLdexpFuncs()...)
	funcs = append(funcs, mathRoundFuncs()...)
	funcs = append(funcs, mathFuzzyEqualsFuncs()...)
	funcs = append(funcs, mathModRemFuncs()...)
	funcs = append(funcs, mathRoundingModeFuncs()...)
	funcs = append(funcs, mathNearbyIntFuncs()...)

	return funcs
}

func mathConstantsFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "math_pi",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_e",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_eps",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
	}
}

func mathCheckFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "math_isinf",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "math_isnan",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "math_isfinite",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
	}
}

func mathUnaryFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "math_abs",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_acos",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_asin",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_asinh",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_atan",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_cbrt",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_ceil",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_cos",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_cosh",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_erf",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_erfinv",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_erfcinv",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_exp",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_exp2",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_fabs",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_floor",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_lgamma",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_rint",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_sigmoid",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_sin",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_sinh",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_sqrt",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_tan",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_tanh",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_tgamma",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_trunc",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_log",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_log2",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_log10",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
	}
}

func mathBinaryFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "math_atan2",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_fmod",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_hypot",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_pow",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_remainder",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
	}
}

func mathLdexpFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "math_ldexp",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
				{Type: &ast.TypeName{Name: "Int32"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
	}
}

func mathRoundFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "math_round",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "math_round",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
				{Type: &ast.TypeName{Name: "Int32"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
	}
}

func mathFuzzyEqualsFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "math_fuzzyequals",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "math_fuzzyequals",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
				{Type: &ast.TypeName{Name: "Double"}},
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
	}
}

func mathModRemFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "math_mod",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int64"}},
				{Type: &ast.TypeName{Name: "Int64"}},
			},
			ReturnType:         &ast.TypeName{Name: "Int64"},
			ReturnTypeNullable: true,
		},
		{
			Name: "math_rem",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int64"}},
				{Type: &ast.TypeName{Name: "Int64"}},
			},
			ReturnType:         &ast.TypeName{Name: "Int64"},
			ReturnTypeNullable: true,
		},
	}
}

func mathRoundingModeFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "math_rounddownward",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "math_roundtonearest",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "math_roundtowardzero",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "math_roundupward",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func mathNearbyIntFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "math_nearbyint",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "Int64"},
			ReturnTypeNullable: true,
		},
	}
}
