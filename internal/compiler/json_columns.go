package compiler

import (
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
)

// outputJSONColumn types an sqlc.jsonb_build_object."Name"(key, value, ...)
// call as a synthesized struct: Name is the call's own function name, fields
// come from the key/value pairs. Keys must be string literals; values are
// typed like any other SELECT target.
func (c *Compiler) outputJSONColumn(qc *QueryCatalog, tables []*Table, res *ast.ResTarget, call *ast.FuncCall) (*Column, error) {
	var args []ast.Node
	if call.Args != nil {
		args = call.Args.Items
	}

	jsonName := call.Func.Name

	if len(args)%2 != 0 {
		return nil, fmt.Errorf("sqlc.jsonb_build_object.%q(...) requires an even number of key/value arguments", jsonName)
	}

	var fields []*Column
	for i := 0; i < len(args); i += 2 {
		key, ok := jsonStringLiteral(args[i])
		if !ok {
			return nil, fmt.Errorf("sqlc.jsonb_build_object.%q(...) argument %d must be a string literal key", jsonName, i+1)
		}

		valCols, err := c.outputColumn(qc, tables, &ast.ResTarget{Val: args[i+1]})
		if err != nil {
			return nil, err
		}
		if len(valCols) != 1 {
			return nil, fmt.Errorf("sqlc.jsonb_build_object.%q(...) value for key %q must resolve to a single column", jsonName, key)
		}

		val := valCols[0]
		val.Name = key
		fields = append(fields, val)
	}

	name := "json"
	if res.Name != nil {
		name = *res.Name
	}
	return &Column{
		Name:       name,
		DataType:   "any",
		NotNull:    true,
		JSONFields: fields,
		JSONName:   jsonName,
	}, nil
}

// outputEmbedJSONColumn types an sqlc.embed.jsonb(table) call as a struct
// mirroring the table's columns, decoded from a single to_jsonb(table) value.
// The name comes from the result alias; the array case leaves it empty for
// the enclosing ARRAY_SUBLINK to fill from its own alias.
func (c *Compiler) outputEmbedJSONColumn(qc *QueryCatalog, tables []*Table, res *ast.ResTarget, call *ast.FuncCall) (*Column, error) {
	if call.Args == nil || len(call.Args.Items) != 1 {
		return nil, fmt.Errorf("sqlc.embed.jsonb(...) takes a single table argument")
	}
	ref, ok := call.Args.Items[0].(*ast.ColumnRef)
	if !ok {
		return nil, fmt.Errorf("sqlc.embed.jsonb(...) argument must be a table reference")
	}
	target := astutils.Join(ref.Fields, ".")

	var table *Table
	for _, t := range tables {
		if t.Rel.Name == target {
			table = t
			break
		}
	}
	if table == nil {
		return nil, fmt.Errorf("sqlc.embed.jsonb(%s): table not found in the query's FROM clause", target)
	}

	var fields []*Column
	for _, col := range table.Columns {
		fields = append(fields, &Column{
			Name:      col.Name,
			DataType:  col.DataType,
			NotNull:   col.NotNull,
			Unsigned:  col.Unsigned,
			IsArray:   col.IsArray,
			ArrayDims: col.ArrayDims,
			Length:    col.Length,
			Type:      col.Type,
		})
	}

	name := "json"
	jsonName := ""
	if res.Name != nil {
		name = *res.Name
		jsonName = *res.Name
	}
	return &Column{
		Name:       name,
		DataType:   "any",
		NotNull:    true,
		JSONFields: fields,
		JSONName:   jsonName,
	}, nil
}

func jsonStringLiteral(node ast.Node) (string, bool) {
	aconst, ok := node.(*ast.A_Const)
	if !ok {
		return "", false
	}
	str, ok := aconst.Val.(*ast.String)
	if !ok {
		return "", false
	}
	return str.Str, true
}
