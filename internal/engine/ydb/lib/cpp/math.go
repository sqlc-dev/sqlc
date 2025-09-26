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
			Name:       "Math::Pi",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name:       "Math::E",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name:       "Math::Eps",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
	}
}

func mathCheckFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Math::IsInf",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "Math::IsNaN",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "Math::IsFinite",
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
			Name: "Math::Abs",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Acos",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Asin",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Asinh",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Atan",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Cbrt",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Ceil",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Cos",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Cosh",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Erf",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::ErfInv",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::ErfcInv",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Exp",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Exp2",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Fabs",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Floor",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Lgamma",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Rint",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Sigmoid",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Sin",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Sinh",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Sqrt",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Tan",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Tanh",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Tgamma",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Trunc",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Log",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Log2",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Log10",
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
			Name: "Math::Atan2",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Fmod",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Hypot",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Pow",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Remainder",
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
			Name: "Math::Ldexp",
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
			Name: "Math::Round",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Double"},
		},
		{
			Name: "Math::Round",
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
			Name: "Math::FuzzyEquals",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		},
		{
			Name: "Math::FuzzyEquals",
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
			Name: "Math::Mod",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Int64"}},
				{Type: &ast.TypeName{Name: "Int64"}},
			},
			ReturnType:         &ast.TypeName{Name: "Int64"},
			ReturnTypeNullable: true,
		},
		{
			Name: "Math::Rem",
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
			Name:       "Math::RoundDownward",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "Math::RoundToNearest",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "Math::RoundTowardZero",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name:       "Math::RoundUpward",
			Args:       []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func mathNearbyIntFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "Math::NearbyInt",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "Int64"},
			ReturnTypeNullable: true,
		},
	}
}
