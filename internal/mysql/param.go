package mysql

import (
	"fmt"
	"strings"

	"github.com/kyleconroy/sqlc/internal/dinosql"
	"vitess.io/vitess/go/vt/sqlparser"
)

type Param struct {
	originalName string
	name         string
	typ          string
}

func paramsInLimitExpr(limit *sqlparser.Limit, s *Schema, settings dinosql.GenerateSettings) ([]*Param, error) {
	params := []*Param{}
	if limit == nil {
		return params, nil
	}

	parseLimitSubExp := func(node sqlparser.Expr) {
		switch v := node.(type) {
		case *sqlparser.SQLVal:
			if v.Type == sqlparser.ValArg {
				params = append(params, &Param{
					originalName: string(v.Val),
					name:         "limit",
					typ:          "uint32",
				})
			}
		}
	}

	parseLimitSubExp(limit.Offset)
	parseLimitSubExp(limit.Rowcount)

	return params, nil
}

func paramsInWhereExpr(e sqlparser.SQLNode, s *Schema, defaultTable string, settings dinosql.GenerateSettings) ([]*Param, error) {
	params := []*Param{}
	switch v := e.(type) {
	case *sqlparser.Where:
		if v == nil {
			return params, nil
		}
		return paramsInWhereExpr(v.Expr, s, defaultTable, settings)
	case *sqlparser.ComparisonExpr:
		p, found, err := paramInComparison(v, s, defaultTable, settings)
		if err != nil {
			return nil, err
		}
		if found {
			params = append(params, p)
		}
	case *sqlparser.AndExpr:
		left, err := paramsInWhereExpr(v.Left, s, defaultTable, settings)
		if err != nil {
			return nil, err
		}
		params = append(params, left...)
		right, err := paramsInWhereExpr(v.Right, s, defaultTable, settings)
		if err != nil {
			return nil, err
		}
		params = append(params, right...)
	case *sqlparser.OrExpr:
		left, err := paramsInWhereExpr(v.Left, s, defaultTable, settings)
		if err != nil {
			return nil, err
		}
		params = append(params, left...)
		right, err := paramsInWhereExpr(v.Right, s, defaultTable, settings)
		if err != nil {
			return nil, err
		}
		params = append(params, right...)
	default:
		panic(fmt.Sprintf("Failed to handle %T in where", v))
	}

	return params, nil
}

func paramInComparison(cond *sqlparser.ComparisonExpr, s *Schema, defaultTable string, settings dinosql.GenerateSettings) (*Param, bool, error) {
	p := &Param{}
	var colIdent sqlparser.ColIdent
	walker := func(node sqlparser.SQLNode) (bool, error) {
		switch v := node.(type) {
		case *sqlparser.ColName:
			colDfn, err := s.getColType(v, defaultTable)
			if err != nil {
				return false, err
			}
			p.typ = goTypeCol(colDfn, settings)
			colIdent = colDfn.Name

		case *sqlparser.SQLVal:
			if v.Type == sqlparser.ValArg {
				p.originalName = string(v.Val)
			}
		}
		return true, nil
	}
	err := sqlparser.Walk(walker, cond)
	if err != nil {
		return nil, false, err
	}
	if p.originalName != "" && p.typ != "" {
		p.name = paramName(colIdent, p.originalName)
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
