package rewrite

import (
	"fmt"

	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/source"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/astutils"
	"github.com/kyleconroy/sqlc/internal/sql/named"
)

// Given an AST node, return the string representation of names
func flatten(root ast.Node) (string, bool) {
	sw := &stringWalker{}
	astutils.Walk(sw, root)
	return sw.String, sw.IsConst
}

type stringWalker struct {
	String  string
	IsConst bool
}

func (s *stringWalker) Visit(node ast.Node) astutils.Visitor {
	if _, ok := node.(*ast.A_Const); ok {
		s.IsConst = true
	}
	if n, ok := node.(*ast.String); ok {
		s.String += n.Str
	}
	return s
}

func isNamedParamSignCast(node ast.Node) bool {
	expr, ok := node.(*ast.A_Expr)
	if !ok {
		return false
	}
	_, cast := expr.Rexpr.(*ast.TypeCast)
	return astutils.Join(expr.Name, ".") == "@" && cast
}

// paramFromName takes a user-defined parameter name, with an optional suffix of
// ? (nullable), or ! (non-null) and builds the appropiate parameter
func paramFromName(name string) named.Param {
	if len(name) == 0 {
		return named.NewUnspecifiedParam(name)
	}

	last := name[len(name)-1]
	if last == '!' {
		return named.NewUserDefinedParam(name[:len(name)-1], true)
	}

	if last == '?' {
		return named.NewUserDefinedParam(name[:len(name)-1], false)
	}

	return named.NewUnspecifiedParam(name)
}

func NamedParameters(engine config.Engine, raw *ast.RawStmt, numbs map[int]bool, dollar bool) (*ast.RawStmt, *named.ParamSet, []source.Edit) {
	foundFunc := astutils.Search(raw, named.IsParamFunc)
	foundSign := astutils.Search(raw, named.IsParamSign)
	hasNamedParameterSupport := engine != config.EngineMySQL
	allParams := named.NewParamSet(numbs, hasNamedParameterSupport)

	if len(foundFunc.Items)+len(foundSign.Items) == 0 {
		return raw, allParams, nil
	}

	var edits []source.Edit
	node := astutils.Apply(raw, func(cr *astutils.Cursor) bool {
		node := cr.Node()
		switch {
		case named.IsParamFunc(node):
			fun := node.(*ast.FuncCall)
			paramName, isConst := flatten(fun.Args)

			param := paramFromName(paramName)
			argn := allParams.Add(param)
			cr.Replace(&ast.ParamRef{
				Number:   argn,
				Location: fun.Location,
			})

			// TODO: This code assumes that sqlc.arg(name) is on a single line
			var old, replace string
			if isConst {
				old = fmt.Sprintf("sqlc.arg('%s')", paramName)
			} else {
				old = fmt.Sprintf("sqlc.arg(%s)", paramName)
			}
			if engine == config.EngineMySQL || !dollar {
				replace = "?"
			} else {
				replace = fmt.Sprintf("$%d", argn)
			}
			edits = append(edits, source.Edit{
				Location: fun.Location - raw.StmtLocation,
				Old:      old,
				New:      replace,
			})
			return false

		case isNamedParamSignCast(node):
			expr := node.(*ast.A_Expr)
			cast := expr.Rexpr.(*ast.TypeCast)
			paramName, _ := flatten(cast.Arg)
			param := paramFromName(paramName)

			argn := allParams.Add(param)
			cast.Arg = &ast.ParamRef{
				Number:   argn,
				Location: expr.Location,
			}
			cr.Replace(cast)

			// TODO: This code assumes that @foo::bool is on a single line
			var replace string
			if engine == config.EngineMySQL || !dollar {
				replace = "?"
			} else {
				replace = fmt.Sprintf("$%d", argn)
			}

			edits = append(edits, source.Edit{
				Location: expr.Location - raw.StmtLocation,
				Old:      fmt.Sprintf("@%s", paramName),
				New:      replace,
			})
			return false

		case named.IsParamSign(node):
			expr := node.(*ast.A_Expr)
			paramName, _ := flatten(expr.Rexpr)
			param := paramFromName(paramName)

			argn := allParams.Add(param)
			cr.Replace(&ast.ParamRef{
				Number:   argn,
				Location: expr.Location,
			})

			// TODO: This code assumes that @foo is on a single line
			var replace string
			if engine == config.EngineMySQL || !dollar {
				replace = "?"
			} else {
				replace = fmt.Sprintf("$%d", argn)
			}

			edits = append(edits, source.Edit{
				Location: expr.Location - raw.StmtLocation,
				Old:      fmt.Sprintf("@%s", paramName),
				New:      replace,
			})
			return false

		default:
			return true
		}
	}, nil)

	return node.(*ast.RawStmt), allParams, edits
}
