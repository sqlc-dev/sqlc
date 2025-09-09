package lib

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func AggregateFunctions() []*catalog.Function {
	var funcs []*catalog.Function

	// COUNT(*)
	funcs = append(funcs, &catalog.Function{
		Name:       "COUNT",
		Args:       []*catalog.Argument{},
		ReturnType: &ast.TypeName{Name: "Uint64"},
	})

	// COUNT(T) и COUNT(T?)
	for _, typ := range types {
		funcs = append(funcs, &catalog.Function{
			Name: "COUNT",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: typ}},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		})
		funcs = append(funcs, &catalog.Function{
			Name: "COUNT",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: typ}, Mode: ast.FuncParamVariadic},
			},
			ReturnType: &ast.TypeName{Name: "Uint64"},
		})
	}

	// MIN и MAX
	for _, typ := range types {
		funcs = append(funcs, &catalog.Function{
			Name: "MIN",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: typ}},
			},
			ReturnType:         &ast.TypeName{Name: typ},
			ReturnTypeNullable: true,
		})
		funcs = append(funcs, &catalog.Function{
			Name: "MAX",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: typ}},
			},
			ReturnType:         &ast.TypeName{Name: typ},
			ReturnTypeNullable: true,
		})
	}

	// SUM для unsigned типов
	for _, typ := range unsignedTypes {
		funcs = append(funcs, &catalog.Function{
			Name: "SUM",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: typ}},
			},
			ReturnType:         &ast.TypeName{Name: "Uint64"},
			ReturnTypeNullable: true,
		})
	}

	// SUM для signed типов
	for _, typ := range signedTypes {
		funcs = append(funcs, &catalog.Function{
			Name: "SUM",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: typ}},
			},
			ReturnType:         &ast.TypeName{Name: "Int64"},
			ReturnTypeNullable: true,
		})
	}

	// SUM для float/double
	for _, typ := range []string{"float", "double"} {
		funcs = append(funcs, &catalog.Function{
			Name: "SUM",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: typ}},
			},
			ReturnType:         &ast.TypeName{Name: typ},
			ReturnTypeNullable: true,
		})
	}

	// AVG для целочисленных типов
	for _, typ := range append(unsignedTypes, signedTypes...) {
		funcs = append(funcs, &catalog.Function{
			Name: "AVG",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: typ}},
			},
			ReturnType:         &ast.TypeName{Name: "Double"},
			ReturnTypeNullable: true,
		})
	}

	// AVG для float/double
	for _, typ := range []string{"float", "double"} {
		funcs = append(funcs, &catalog.Function{
			Name: "AVG",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: typ}},
			},
			ReturnType:         &ast.TypeName{Name: typ},
			ReturnTypeNullable: true,
		})
	}

	// COUNT_IF
	funcs = append(funcs, &catalog.Function{
		Name: "COUNT_IF",
		Args: []*catalog.Argument{
			{Type: &ast.TypeName{Name: "Bool"}},
		},
		ReturnType:         &ast.TypeName{Name: "Uint64"},
		ReturnTypeNullable: true,
	})

	// SUM_IF для unsigned
	for _, typ := range unsignedTypes {
		funcs = append(funcs, &catalog.Function{
			Name: "SUM_IF",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: typ}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType:         &ast.TypeName{Name: "Uint64"},
			ReturnTypeNullable: true,
		})
	}

	// SUM_IF для signed
	for _, typ := range signedTypes {
		funcs = append(funcs, &catalog.Function{
			Name: "SUM_IF",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: typ}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType:         &ast.TypeName{Name: "Int64"},
			ReturnTypeNullable: true,
		})
	}

	// SUM_IF для float/double
	for _, typ := range []string{"float", "double"} {
		funcs = append(funcs, &catalog.Function{
			Name: "SUM_IF",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: typ}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType:         &ast.TypeName{Name: typ},
			ReturnTypeNullable: true,
		})
	}

	// AVG_IF для целочисленных
	for _, typ := range append(unsignedTypes, signedTypes...) {
		funcs = append(funcs, &catalog.Function{
			Name: "AVG_IF",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: typ}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType:         &ast.TypeName{Name: "Double"},
			ReturnTypeNullable: true,
		})
	}

	// AVG_IF для float/double
	for _, typ := range []string{"float", "double"} {
		funcs = append(funcs, &catalog.Function{
			Name: "AVG_IF",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: typ}},
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType:         &ast.TypeName{Name: typ},
			ReturnTypeNullable: true,
		})
	}

	// SOME
	for _, typ := range types {
		funcs = append(funcs, &catalog.Function{
			Name: "SOME",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: typ}},
			},
			ReturnType:         &ast.TypeName{Name: typ},
			ReturnTypeNullable: true,
		})
	}

	// AGGREGATE_LIST и AGGREGATE_LIST_DISTINCT
	for _, typ := range types {
		funcs = append(funcs, &catalog.Function{
			Name: "AGGREGATE_LIST",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: typ}},
			},
			ReturnType: &ast.TypeName{Name: "List<" + typ + ">"},
		})
		funcs = append(funcs, &catalog.Function{
			Name: "AGGREGATE_LIST_DISTINCT",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: typ}},
			},
			ReturnType: &ast.TypeName{Name: "List<" + typ + ">"},
		})
	}

	// BOOL_AND, BOOL_OR, BOOL_XOR
	boolAggrs := []string{"BOOL_AND", "BOOL_OR", "BOOL_XOR"}
	for _, name := range boolAggrs {
		funcs = append(funcs, &catalog.Function{
			Name: name,
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Bool"}},
			},
			ReturnType:         &ast.TypeName{Name: "Bool"},
			ReturnTypeNullable: true,
		})
	}

	// BIT_AND, BIT_OR, BIT_XOR
	bitAggrs := []string{"BIT_AND", "BIT_OR", "BIT_XOR"}
	for _, typ := range append(unsignedTypes, signedTypes...) {
		for _, name := range bitAggrs {
			funcs = append(funcs, &catalog.Function{
				Name: name,
				Args: []*catalog.Argument{
					{Type: &ast.TypeName{Name: typ}},
				},
				ReturnType:         &ast.TypeName{Name: typ},
				ReturnTypeNullable: true,
			})
		}
	}

	// STDDEV и VARIANCE
	stdDevVariants := []struct {
		name   string
		returnType string
	}{
		{"STDDEV", "Double"},
		{"VARIANCE", "Double"},
		{"STDDEV_SAMPLE", "Double"},
		{"VARIANCE_SAMPLE", "Double"},
		{"STDDEV_POPULATION", "Double"},
		{"VARIANCE_POPULATION", "Double"},
	}
	for _, variant := range stdDevVariants {
		funcs = append(funcs, &catalog.Function{
			Name: variant.name,
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType:         &ast.TypeName{Name: variant.returnType},
			ReturnTypeNullable: true,
		})
	}

	// CORRELATION и COVARIANCE
	corrCovar := []string{"CORRELATION", "COVARIANCE", "COVARIANCE_SAMPLE", "COVARIANCE_POPULATION"}
	for _, name := range corrCovar {
		funcs = append(funcs, &catalog.Function{
			Name: name,
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Double"}},
				{Type: &ast.TypeName{Name: "Double"}},
			},
			ReturnType:         &ast.TypeName{Name: "Double"},
			ReturnTypeNullable: true,
		})
	}

	// HISTOGRAM
	funcs = append(funcs, &catalog.Function{
		Name: "HISTOGRAM",
		Args: []*catalog.Argument{
			{Type: &ast.TypeName{Name: "Double"}},
		},
		ReturnType:         &ast.TypeName{Name: "HistogramStruct"},
		ReturnTypeNullable: true,
	})

	// TOP и BOTTOM
	topBottom := []string{"TOP", "BOTTOM"}
	for _, name := range topBottom {
		for _, typ := range types {
			funcs = append(funcs, &catalog.Function{
				Name: name,
				Args: []*catalog.Argument{
					{Type: &ast.TypeName{Name: typ}},
					{Type: &ast.TypeName{Name: "Uint32"}},
				},
				ReturnType: &ast.TypeName{Name: "List<" + typ + ">"},
			})
		}
	}

	// MAX_BY и MIN_BY
	minMaxBy := []string{"MAX_BY", "MIN_BY"}
	for _, name := range minMaxBy {
		for _, typ := range types {
			funcs = append(funcs, &catalog.Function{
				Name: name,
				Args: []*catalog.Argument{
					{Type: &ast.TypeName{Name: typ}},
					{Type: &ast.TypeName{Name: "any"}},
				},
				ReturnType:         &ast.TypeName{Name: typ},
				ReturnTypeNullable: true,
			})
		}
	}

	// ... (добавьте другие агрегатные функции по аналогии)

	return funcs
}
