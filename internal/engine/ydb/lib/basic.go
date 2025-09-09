package lib

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

var types = []string{
	"bool",
	"int8", "int16", "int32", "int64",
	"uint8", "uint16", "uint32", "uint64",
	"float", "double",
	"string", "utf8",
	"any",
}

var (
	unsignedTypes = []string{"uint8", "uint16", "uint32", "uint64"}
	signedTypes   = []string{"int8", "int16", "int32", "int64"}
	numericTypes  = append(append(unsignedTypes, signedTypes...), "float", "double")
)

func BasicFunctions() []*catalog.Function {
	var funcs []*catalog.Function

	for _, typ := range types {
		// COALESCE, NVL
		funcs = append(funcs, &catalog.Function{
			Name: "COALESCE",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: typ}},
				{Type: &ast.TypeName{Name: typ}},
				{
					Type: &ast.TypeName{Name: typ},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType:         &ast.TypeName{Name: typ},
			ReturnTypeNullable: false,
		})
		funcs = append(funcs, &catalog.Function{
			Name: "NVL",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: typ}},
				{Type: &ast.TypeName{Name: typ}},
				{
					Type: &ast.TypeName{Name: typ},
					Mode: ast.FuncParamVariadic,
				},
			},
			ReturnType:         &ast.TypeName{Name: typ},
			ReturnTypeNullable: false,
		})

		// IF(Bool, T, T) -> T
		funcs = append(funcs, &catalog.Function{
			Name: "IF",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Bool"}},
				{Type: &ast.TypeName{Name: typ}},
				{Type: &ast.TypeName{Name: typ}},
			},
			ReturnType:         &ast.TypeName{Name: typ},
			ReturnTypeNullable: false,
		})

		// LENGTH, LEN
		funcs = append(funcs, &catalog.Function{
			Name: "LENGTH",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: typ}},
			},
			ReturnType:         &ast.TypeName{Name: "Uint32"},
			ReturnTypeNullable: true,
		})
		funcs = append(funcs, &catalog.Function{
			Name: "LEN",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: typ}},
			},
			ReturnType:         &ast.TypeName{Name: "Uint32"},
			ReturnTypeNullable: true,
		})

		// StartsWith, EndsWith
		funcs = append(funcs, &catalog.Function{
			Name: "StartsWith",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: typ}},
				{Type: &ast.TypeName{Name: typ}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		})
		funcs = append(funcs, &catalog.Function{
			Name: "EndsWith",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: typ}},
				{Type: &ast.TypeName{Name: typ}},
			},
			ReturnType: &ast.TypeName{Name: "Bool"},
		})

		// ABS(T) -> T
	}

	// SUBSTRING
	funcs = append(funcs, &catalog.Function{
		Name: "Substring",
		Args: []*catalog.Argument{
			{Type: &ast.TypeName{Name: "String"}},
		},
		ReturnType: &ast.TypeName{Name: "String"},
	})
	funcs = append(funcs, &catalog.Function{
		Name: "Substring",
		Args: []*catalog.Argument{
			{Type: &ast.TypeName{Name: "String"}},
			{Type: &ast.TypeName{Name: "Uint32"}},
		},
		ReturnType: &ast.TypeName{Name: "String"},
	})
	funcs = append(funcs, &catalog.Function{
		Name: "Substring",
		Args: []*catalog.Argument{
			{Type: &ast.TypeName{Name: "String"}},
			{Type: &ast.TypeName{Name: "Uint32"}},
			{Type: &ast.TypeName{Name: "Uint32"}},
		},
		ReturnType: &ast.TypeName{Name: "String"},
	})

	// FIND / RFIND
	for _, name := range []string{"FIND", "RFIND"} {
		for _, typ := range []string{"String", "Utf8"} {
			funcs = append(funcs, &catalog.Function{
				Name: name,
				Args: []*catalog.Argument{
					{Type: &ast.TypeName{Name: typ}},
					{Type: &ast.TypeName{Name: typ}},
				},
				ReturnType: &ast.TypeName{Name: "Uint32"},
			})
			funcs = append(funcs, &catalog.Function{
				Name: name,
				Args: []*catalog.Argument{
					{Type: &ast.TypeName{Name: typ}},
					{Type: &ast.TypeName{Name: typ}},
					{Type: &ast.TypeName{Name: "Uint32"}},
				},
				ReturnType: &ast.TypeName{Name: "Uint32"},
			})
		}
	}

	for _, typ := range numericTypes {
		funcs = append(funcs, &catalog.Function{
			Name: "Abs",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: typ}},
			},
			ReturnType: &ast.TypeName{Name: typ},
		})
	}

	// NANVL
	funcs = append(funcs, &catalog.Function{
		Name: "NANVL",
		Args: []*catalog.Argument{
			{Type: &ast.TypeName{Name: "Float"}},
			{Type: &ast.TypeName{Name: "Float"}},
		},
		ReturnType: &ast.TypeName{Name: "Float"},
	})
	funcs = append(funcs, &catalog.Function{
		Name: "NANVL",
		Args: []*catalog.Argument{
			{Type: &ast.TypeName{Name: "Double"}},
			{Type: &ast.TypeName{Name: "Double"}},
		},
		ReturnType: &ast.TypeName{Name: "Double"},
	})

	// Random*
	funcs = append(funcs, &catalog.Function{
		Name:       "Random",
		Args:       []*catalog.Argument{},
		ReturnType: &ast.TypeName{Name: "Double"},
	})
	funcs = append(funcs, &catalog.Function{
		Name:       "RandomNumber",
		Args:       []*catalog.Argument{},
		ReturnType: &ast.TypeName{Name: "Uint64"},
	})
	funcs = append(funcs, &catalog.Function{
		Name:       "RandomUuid",
		Args:       []*catalog.Argument{},
		ReturnType: &ast.TypeName{Name: "Uuid"},
	})

	// todo: add all remain functions

	return funcs
}
