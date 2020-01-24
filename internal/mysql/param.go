package mysql

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kyleconroy/sqlc/internal/dinosql"
	"vitess.io/vitess/go/vt/sqlparser"
)

// Param describes a runtime query parameter with its
// associated type. Example: "SELECT name FROM users id = ?"
type Param struct {
	OriginalName string
	Name         string
	Typ          string
}

func paramsInLimitExpr(limit *sqlparser.Limit, s *Schema, tableAliasMap FromTables, settings dinosql.CombinedSettings) ([]*Param, error) {
	params := []*Param{}
	if limit == nil {
		return params, nil
	}

	parseLimitSubExp := func(node sqlparser.Expr) error {
		switch v := node.(type) {
		case *sqlparser.SQLVal:
			if v.Type == sqlparser.ValArg {
				params = append(params, &Param{
					OriginalName: string(v.Val),
					Name:         "limit",
					Typ:          "uint32",
				})
			}
		case *sqlparser.FuncExpr:
			name, raw, err := matchFuncExpr(v)
			if err != nil {
				return err
			}
			if name != "" && raw != "" {
				params = append(params, &Param{
					OriginalName: raw,
					Name:         name,
					Typ:          "uint32",
				})
			}
		}
		return nil
	}

	err := parseLimitSubExp(limit.Offset)
	if err != nil {
		return nil, err
	}
	err = parseLimitSubExp(limit.Rowcount)
	if err != nil {
		return nil, err
	}

	return params, nil
}

func paramsInWhereExpr(e sqlparser.SQLNode, s *Schema, tableAliasMap FromTables, defaultTable string, settings dinosql.CombinedSettings) ([]*Param, error) {
	params := []*Param{}
	if e == nil {
		return params, nil
	} else if expr, ok := e.(*sqlparser.Where); ok {
		if expr == nil {
			return params, nil
		}
		e = expr.Expr
	}
	switch v := e.(type) {
	case *sqlparser.Where:
		if v == nil {
			return params, nil
		}
		return paramsInWhereExpr(v, s, tableAliasMap, defaultTable, settings)
	case *sqlparser.ComparisonExpr:
		p, found, err := paramInComparison(v, s, tableAliasMap, defaultTable, settings)
		if err != nil {
			return nil, err
		}
		if found {
			params = append(params, p)
		}
	case *sqlparser.AndExpr:
		left, err := paramsInWhereExpr(v.Left, s, tableAliasMap, defaultTable, settings)
		if err != nil {
			return nil, err
		}
		params = append(params, left...)
		right, err := paramsInWhereExpr(v.Right, s, tableAliasMap, defaultTable, settings)
		if err != nil {
			return nil, err
		}
		params = append(params, right...)
	case *sqlparser.OrExpr:
		left, err := paramsInWhereExpr(v.Left, s, tableAliasMap, defaultTable, settings)
		if err != nil {
			return nil, err
		}
		params = append(params, left...)
		right, err := paramsInWhereExpr(v.Right, s, tableAliasMap, defaultTable, settings)
		if err != nil {
			return nil, err
		}
		params = append(params, right...)
	case *sqlparser.IsExpr:
		// TODO: see if there is a use case for params in IS expressions
		return []*Param{}, nil
	default:
		panic(fmt.Sprintf("Failed to handle %T in where", v))
	}

	return params, nil
}

func paramInComparison(cond *sqlparser.ComparisonExpr, s *Schema, tableAliasMap FromTables, defaultTable string, settings dinosql.CombinedSettings) (*Param, bool, error) {
	p := &Param{}
	var colIdent sqlparser.ColIdent
	walker := func(node sqlparser.SQLNode) (bool, error) {
		switch v := node.(type) {
		case *sqlparser.ColName:
			colDfn, err := s.getColType(v, tableAliasMap, defaultTable)
			if err != nil {
				return false, err
			}
			p.Typ = goTypeCol(colDfn, settings)
			colIdent = colDfn.Name

		case *sqlparser.SQLVal:
			if v.Type == sqlparser.ValArg {
				p.OriginalName = string(v.Val)
			}
		case *sqlparser.FuncExpr:
			name, raw, err := matchFuncExpr(v)
			if err != nil {
				return false, err
			}
			if name != "" && raw != "" {
				p.OriginalName = raw
				p.Name = name
			}
			return false, nil
		}
		return true, nil
	}
	err := sqlparser.Walk(walker, cond)
	if err != nil {
		return nil, false, err
	}
	if p.Name != "" {
		return p, true, nil
	}
	if p.OriginalName != "" && p.Typ != "" {
		p.Name = paramName(colIdent, p.OriginalName)
		return p, true, nil
	}
	return nil, false, nil
}

func paramName(col sqlparser.ColIdent, originalName string) string {
	str := col.String()
	if !strings.HasPrefix(originalName, ":v") {
		return originalName[1:]
	}
	if str != "" {
		return str
	}
	num := originalName[2]
	return fmt.Sprintf("param%v", num)
}

func replaceParamStrs(query string, params []*Param) (string, error) {
	for _, p := range params {
		re, err := regexp.Compile(fmt.Sprintf("(%v)", regexp.QuoteMeta(p.OriginalName)))
		if err != nil {
			return "", err
		}
		query = re.ReplaceAllString(query, "?")
	}
	return query, nil
}

func matchFuncExpr(v *sqlparser.FuncExpr) (name string, raw string, err error) {
	namespace := "sqlc"
	fakeFunc := "arg"
	if v.Qualifier.String() == namespace {
		if v.Name.String() == fakeFunc {
			if expr, ok := v.Exprs[0].(*sqlparser.AliasedExpr); ok {
				if colName, ok := expr.Expr.(*sqlparser.ColName); ok {
					customName := colName.Name.String()
					return customName, fmt.Sprintf("%s.%s(%s)", namespace, fakeFunc, customName), nil
				}
				return "", "", fmt.Errorf("invalid custom argument value \"%s.%s(%s)\"", namespace, fakeFunc, replaceVParamExprs(sqlparser.String(v.Exprs[0])))
			}
			return "", "", fmt.Errorf("invalid custom argument value \"%s.%s(%s)\"", namespace, fakeFunc, replaceVParamExprs(sqlparser.String(v.Exprs[0])))
		}
		return "", "", fmt.Errorf("invalid function call \"%s.%s\", did you mean \"%s.%s\"?", namespace, v.Name.String(), namespace, fakeFunc)
	}
	return "", "", nil
}

func replaceVParamExprs(sql string) string {
	/*
	   the sqlparser replaces "?" with ":v1"
	   to display a helpful error message, these should be replaced back to "?"
	*/
	matcher := regexp.MustCompile(":v[0-9]*")
	return matcher.ReplaceAllString(sql, "?")
}
