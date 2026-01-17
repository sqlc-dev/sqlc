package clickhouse

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func defaultSchema(name string) *catalog.Schema {
	s := &catalog.Schema{Name: name}
	s.Funcs = []*catalog.Function{
		// Aggregate functions
		{
			Name: "count",
			Args: []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "UInt64"},
		},
		{
			Name: "sum",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "avg",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Float64"},
		},
		{
			Name: "min",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		{
			Name: "max",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		// Type conversion functions
		{
			Name: "toInt32",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Int32"},
		},
		{
			Name: "toInt64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Int64"},
		},
		{
			Name: "toUInt32",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "UInt32"},
		},
		{
			Name: "toUInt64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "UInt64"},
		},
		{
			Name: "toString",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "toFloat64",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Float64"},
		},
		// Date/time functions
		{
			Name: "now",
			Args: []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "DateTime"},
		},
		{
			Name: "today",
			Args: []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "Date"},
		},
		{
			Name: "toDate",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "Date"},
		},
		{
			Name: "toDateTime",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "DateTime"},
		},
		// String functions
		{
			Name: "concat",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "lower",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "upper",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "String"},
		},
		{
			Name: "length",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "String"}},
			},
			ReturnType: &ast.TypeName{Name: "UInt64"},
		},
		// Conditional functions
		{
			Name: "if",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "UInt8"}},
				{Type: &ast.TypeName{Name: "any"}},
				{Type: &ast.TypeName{Name: "any"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
		// Array functions
		{
			Name: "array",
			Args: []*catalog.Argument{},
			ReturnType: &ast.TypeName{Name: "Array"},
		},
		{
			Name: "arrayJoin",
			Args: []*catalog.Argument{
				{Type: &ast.TypeName{Name: "Array"}},
			},
			ReturnType: &ast.TypeName{Name: "any"},
		},
	}
	return s
}
