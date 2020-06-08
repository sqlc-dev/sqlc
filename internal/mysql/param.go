package mysql

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"vitess.io/vitess/go/vt/sqlparser"
)

// Param describes a runtime query parameter with its
// associated type. Example: "SELECT name FROM users id = ?"
type Param struct {
	OriginalName string
	Name         string
	Typ          string
}

func (pGen PackageGenerator) paramsInLimitExpr(limit *sqlparser.Limit, tableAliasMap FromTables) ([]*Param, error) {
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

func (pGen PackageGenerator) paramsInWhereExpr(e sqlparser.SQLNode, tableAliasMap FromTables, defaultTable string) ([]*Param, error) {
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
		return pGen.paramsInWhereExpr(v, tableAliasMap, defaultTable)
	case *sqlparser.ComparisonExpr:
		p, found, err := pGen.paramInComparison(v, tableAliasMap, defaultTable)
		if err != nil {
			return nil, err
		}
		if found {
			params = append(params, p)
		}
	case *sqlparser.AndExpr:
		left, err := pGen.paramsInWhereExpr(v.Left, tableAliasMap, defaultTable)
		if err != nil {
			return nil, err
		}
		params = append(params, left...)
		right, err := pGen.paramsInWhereExpr(v.Right, tableAliasMap, defaultTable)
		if err != nil {
			return nil, err
		}
		params = append(params, right...)
	case *sqlparser.OrExpr:
		left, err := pGen.paramsInWhereExpr(v.Left, tableAliasMap, defaultTable)
		if err != nil {
			return nil, err
		}
		params = append(params, left...)
		right, err := pGen.paramsInWhereExpr(v.Right, tableAliasMap, defaultTable)
		if err != nil {
			return nil, err
		}
		params = append(params, right...)
	case *sqlparser.IsExpr:
		// TODO: see if there is a use case for params in IS expressions
		return []*Param{}, nil
	case *sqlparser.ParenExpr:
		expr, err := pGen.paramsInWhereExpr(v.Expr, tableAliasMap, defaultTable)
		if err != nil {
			return nil, err
		}
		params = append(params, expr...)
	default:
		panic(fmt.Sprintf("Failed to handle %T in where", v))
	}

	return params, nil
}

func (pGen PackageGenerator) paramInComparison(cond *sqlparser.ComparisonExpr, tableAliasMap FromTables, defaultTable string) (*Param, bool, error) {
	param := &Param{}
	var colIdent sqlparser.ColIdent
	walker := func(node sqlparser.SQLNode) (bool, error) {
		switch v := node.(type) {
		case *sqlparser.ColName:
			col, err := pGen.getColType(v, tableAliasMap, defaultTable)
			if err != nil {
				return false, err
			}
			param.Typ = pGen.goTypeCol(*col)
			colIdent = col.Name

		case *sqlparser.SQLVal:
			if v.Type == sqlparser.ValArg {
				param.OriginalName = string(v.Val)
			}
		case *sqlparser.FuncExpr:
			name, raw, err := matchFuncExpr(v)
			if err != nil {
				return false, err
			}
			if name != "" && raw != "" {
				param.OriginalName = raw
				param.Name = name
			}
			return false, nil
		}
		return true, nil
	}
	err := sqlparser.Walk(walker, cond)
	if err != nil {
		return nil, false, err
	}
	if param.Name != "" {
		return param, true, nil
	}
	if param.OriginalName != "" && param.Typ != "" {
		param.Name = paramName(colIdent, param.OriginalName)
		return param, true, nil
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
	/*
		To ensure that ":v1" does not replace ":v12", we need to sort
		the params in decending order by length of the string.
		But, the original order of the params must be preserved.
	*/
	paramsCopy := make([]*Param, len(params))
	copy(paramsCopy, params)
	sort.Slice(paramsCopy, func(i, j int) bool {
		return len(paramsCopy[i].OriginalName) > len(paramsCopy[j].OriginalName)
	})

	for _, p := range paramsCopy {
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
