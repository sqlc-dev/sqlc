package compiler

import (
	"errors"
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/rewrite"
)

// outputJSONColumn types an sqlc.jsonb_build_object."Name"(key, value, ...)
// call as a synthesized struct: Name is the call's own function name, fields
// come from the key/value pairs. Keys must be string literals; values are
// typed by jsonValueColumn.
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
		val, err := c.jsonValueColumn(qc, tables, args[i+1], key)
		if err != nil {
			return nil, err
		}
		val.Name = key
		fields = append(fields, val)
	}

	name := "json"
	if res.Name != nil {
		name = *res.Name
	}
	return &Column{Name: name, DataType: "any", NotNull: true, JSONFields: fields, JSONName: jsonName}, nil
}

// jsonValueColumn types a single sqlc.jsonb_build_object value: a column
// reference, a nested JSON call, or an ARRAY(...) of one of those.
func (c *Compiler) jsonValueColumn(qc *QueryCatalog, tables []*Table, node ast.Node, key string) (*Column, error) {
	switch v := node.(type) {
	case *ast.ColumnRef:
		cols, err := outputColumnRefs(&ast.ResTarget{Val: v}, tables, v)
		if err != nil {
			return nil, err
		}
		if len(cols) != 1 {
			return nil, fmt.Errorf("sqlc.jsonb_build_object value for key %q must resolve to a single column", key)
		}
		return cols[0], nil
	case *ast.SubLink:
		if v.SubLinkType == ast.ARRAY_SUBLINK {
			return c.arraySubLinkColumn(qc, tables, v)
		}
	case *ast.FuncCall:
		if rewrite.IsJSONCall(v) {
			return c.outputJSONColumn(qc, tables, &ast.ResTarget{}, v)
		}
	}
	return &Column{DataType: "any"}, nil
}

// arraySubLinkColumn types an ARRAY(subquery) expression: the subquery must
// yield exactly one column, which becomes the array element.
func (c *Compiler) arraySubLinkColumn(qc *QueryCatalog, tables []*Table, sublink *ast.SubLink) (*Column, error) {
	subcols, err := c.outputColumns(qc, sublink.Subselect)
	if err != nil {
		return nil, err
	}
	if len(subcols) != 1 {
		return nil, errors.New("ARRAY() subquery must return only one column")
	}
	first := subcols[0]
	if first.IsArray {
		first.ArrayDims++
	} else {
		first.IsArray = true
		first.ArrayDims = 1
	}
	first.NotNull = true
	return first, nil
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
