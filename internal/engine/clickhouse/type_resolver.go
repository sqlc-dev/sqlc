package clickhouse

import (
	"strings"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func TypeResolver(call *ast.FuncCall, fun *catalog.Function, resolve func(n ast.Node) (*catalog.Column, error)) *ast.TypeName {
	funcName := strings.ToLower(fun.Name)

	switch funcName {
	case "arrayjoin":
		if len(call.Args.Items) != 1 {
			return nil
		}
		// arrayJoin(arr) -> returns element type of arr
		col, err := resolve(call.Args.Items[0])
		if err != nil {
			return nil
		}
		// If the argument is an array, return the element type
		// In sqlc, arrays are often represented by the IsArray flag on the column
		if col.IsArray {
			// Create a new type name based on the column's type
			// We need to "unwrap" the array.
			typeName := col.Type.Name
			if strings.HasSuffix(typeName, "[]") {
				typeName = strings.TrimSuffix(typeName, "[]")
			}
			return &ast.TypeName{
				Name: typeName,
			}
		}
		// If it's not marked as IsArray, it might still be an array type string (e.g. "Array(Int32)")
		// TODO: Parsing ClickHouse type strings might be needed here if DataType is "Array(T)"
		return nil

	case "argmin", "argmax", "any", "anylast", "anyheavy",
		"argminif", "argmaxif", "anyif", "anylastif", "anyheavyif",
		"min", "max", "sum", "minif", "maxif", "sumif":
		if len(call.Args.Items) < 1 {
			return nil
		}
		// These functions return the type of their first argument
		col, err := resolve(call.Args.Items[0])
		if err != nil {
			return nil
		}
		typeName := col.Type
		return &typeName

	case "count", "countif", "uniq", "uniqexact":
		// ClickHouse count returns UInt64
		return &ast.TypeName{Name: "uint64"}

	case "jsonextract":
		// JSONExtract(json, indices_or_keys..., return_type)
		// The last argument is usually the type
		if len(call.Args.Items) < 2 {
			return nil
		}
		lastArg := call.Args.Items[len(call.Args.Items)-1]
		// Check if it's a string literal
		if constVal, ok := lastArg.(*ast.A_Const); ok {
			if strVal, ok := constVal.Val.(*ast.String); ok {
				typeStr := strVal.Str
				// Map ClickHouse type string to sqlc type
				mappedType := mapClickHouseType(typeStr)
				// If it's an array type, we need to handle it
				if strings.HasSuffix(mappedType, "[]") {
					elemType := strings.TrimSuffix(mappedType, "[]")
					return &ast.TypeName{
						Name: elemType,
						ArrayBounds: &ast.List{
							Items: []ast.Node{&ast.A_Const{}},
						},
					}
				}
				return &ast.TypeName{Name: mappedType}
			}
		}
		return nil

	case "jsonextractkeysandvalues":
		// JSONExtractKeysAndValues(json, 'ValueType') -> Array(Tuple(String, ValueType))
		// We map this to just "any" or a complex type if possible.
		// In sqlc, we might represent Array(Tuple(String, T)) as... complex.
		// For now, let's return any, or maybe just handle the array part.
		return &ast.TypeName{Name: "any"}
	}

	return nil
}
