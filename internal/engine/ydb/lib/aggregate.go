package lib

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func AggregateFunctions() []*catalog.Function {
	var funcs []*catalog.Function

	funcs = append(funcs, countFuncs()...)
	funcs = append(funcs, minMaxFuncs()...)
	funcs = append(funcs, sumFuncs()...)
	funcs = append(funcs, avgFuncs()...)
	funcs = append(funcs, countIfFuncs()...)
	funcs = append(funcs, sumIfFuncs()...)
	funcs = append(funcs, avgIfFuncs()...)
	funcs = append(funcs, someFuncs()...)
	funcs = append(funcs, countDistinctEstimateHLLFuncs()...)
	funcs = append(funcs, maxByMinByFuncs()...)
	funcs = append(funcs, stddevVarianceFuncs()...)
	funcs = append(funcs, correlationCovarianceFuncs()...)
	funcs = append(funcs, percentileMedianFuncs()...)
	funcs = append(funcs, boolAndOrXorFuncs()...)
	funcs = append(funcs, bitAndOrXorFuncs()...)

	// TODO: Aggregate_List, Top, Bottom, Top_By, Bottom_By, TopFreq, Mode,
	// Histogram LinearHistogram, LogarithmicHistogram, LogHistogram, CDF,
	// SessionStart, AGGREGATE_BY, MULTI_AGGREGATE_BY

	return funcs
}

func countFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "COUNT",
			Args: []*catalog.Argument{
				{
					Type: &ast.TypeName{Name: "any"},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		},
	}
}

func minMaxFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "MIN",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "MAX",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func sumFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "SUM",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
	}
}

func avgFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "AVG",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
	}
}

func countIfFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "COUNT_IF",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType:         &ast.TypeName{Name: "Uint64"},
			ReturnTypeNullable: true,
		},
	}
}

func sumIfFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "SUM_IF",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
	}
}

func avgIfFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "AVG_IF",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
	}
}

func someFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "SOME",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func countDistinctEstimateHLLFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "CountDistinctEstimate",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "Uint64"},
			ReturnTypeNullable: true,
		},
		{
			Name: "HyperLogLog",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "Uint64"},
			ReturnTypeNullable: true,
		},
		{
			Name: "HLL",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "Uint64"},
			ReturnTypeNullable: true,
		},
	}
}

func maxByMinByFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "MAX_BY",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "MIN_BY",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
	}
	// todo: min/max_by with third argument returning list
}

func stddevVarianceFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "STDDEV",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType:         &ast.TypeName{Name: "Double"},
			ReturnTypeNullable: true,
		},
		{
			Name: "STDDEV_POPULATION",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType:         &ast.TypeName{Name: "Double"},
			ReturnTypeNullable: true,
		},
		{
			Name: "POPULATION_STDDEV",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType:         &ast.TypeName{Name: "Double"},
			ReturnTypeNullable: true,
		},
		{
			Name: "STDDEV_SAMPLE",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType:         &ast.TypeName{Name: "Double"},
			ReturnTypeNullable: true,
		},
		{
			Name: "STDDEVSAMP",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType:         &ast.TypeName{Name: "Double"},
			ReturnTypeNullable: true,
		},
		{
			Name: "VARIANCE",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType:         &ast.TypeName{Name: "Double"},
			ReturnTypeNullable: true,
		},
		{
			Name: "VARIANCE_POPULATION",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType:         &ast.TypeName{Name: "Double"},
			ReturnTypeNullable: true,
		},
		{
			Name: "POPULATION_VARIANCE",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType:         &ast.TypeName{Name: "Double"},
			ReturnTypeNullable: true,
		},
		{
			Name: "VARPOP",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType:         &ast.TypeName{Name: "Double"},
			ReturnTypeNullable: true,
		},
		{
			Name: "VARIANCE_SAMPLE",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType:         &ast.TypeName{Name: "Double"},
			ReturnTypeNullable: true,
		},
	}
}

func correlationCovarianceFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "CORRELATION",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType:         &ast.TypeName{Name: "Double"},
			ReturnTypeNullable: true,
		},
		{
			Name: "COVARIANCE",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType:         &ast.TypeName{Name: "Double"},
			ReturnTypeNullable: true,
		},
		{
			Name: "COVARIANCE_SAMPLE",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType:         &ast.TypeName{Name: "Double"},
			ReturnTypeNullable: true,
		},
		{
			Name: "COVARIANCE_POPULATION",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType:         &ast.TypeName{Name: "Double"},
			ReturnTypeNullable: true,
		},
	}
}

func percentileMedianFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "PERCENTILE",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "MEDIAN",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "MEDIAN",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
}

func boolAndOrXorFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "BOOL_AND",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType:         &ast.TypeName{Name: "Bool"},
			ReturnTypeNullable: true,
		},
		{
			Name: "BOOL_OR",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType:         &ast.TypeName{Name: "Bool"},
			ReturnTypeNullable: true,
		},
		{
			Name: "BOOL_XOR",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType:         &ast.TypeName{Name: "Bool"},
			ReturnTypeNullable: true,
		},
	}
}

func bitAndOrXorFuncs() []*catalog.Function {
	return []*catalog.Function{
		{
			Name: "BIT_AND",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "BIT_OR",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
		{
			Name: "BIT_XOR",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType:         &ast.TypeName{Name: "any"},
			ReturnTypeNullable: true,
		},
	}
}
