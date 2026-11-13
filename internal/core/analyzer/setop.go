package analyzer

import (
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/core"
	mssqltypes "github.com/sqlc-dev/sqlc/internal/engine/mssql/typepolicy"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// combineColumns types each column of a set operation from both of its
// branches by the dialect's rules. A column keeps the first branch's name
// and is read from no one table. It may be NULL in a UNION when either
// branch's may, in an INTERSECT when both may, and in an EXCEPT when the
// first branch's may, since every row it returns is one of those.
func (a *analyzer) combineColumns(rules *core.SetOperationRules, op ast.SetOperation, left, right []core.Column) ([]core.Column, error) {
	out := make([]core.Column, len(left))
	for i, l := range left {
		if i >= len(right) {
			out[i] = l
			continue
		}
		r := right[i]
		col := core.Column{Name: l.Name, NotNull: l.NotNull && r.NotNull}
		switch op {
		case ast.Intersect:
			col.NotNull = l.NotNull || r.NotNull
		case ast.Except:
			col.NotNull = l.NotNull
		}
		t, err := a.combineTypes(rules, columnExprType(l), columnExprType(r))
		if err != nil {
			return nil, err
		}
		t.nullable = !col.NotNull
		col.TypeOID = t.typeOID
		col.DataType, col.IsArray = a.typeNameOf(t)
		col.Type = a.typeExprOf(t)
		out[i] = col
	}
	return out, nil
}

// combineTypes delegates type selection and modifiers to the selected policy,
// then resolves its result against the catalog. No policy is enabled merely
// by listing a dialect's types.
func (a *analyzer) combineTypes(rules *core.SetOperationRules, l, r exprType) (exprType, error) {
	lt, rt := a.exprOf(l), a.exprOf(r)
	switch {
	case lt == nil:
		return r, nil
	case rt == nil:
		return l, nil
	}
	var out *core.TypeExpr
	var ok bool
	switch rules.Strategy {
	case "mssql":
		out, ok = mssqltypes.Combine(rules, lt, rt, a.sameCategory(lt.Name, rt.Name))
	default:
		return exprType{}, fmt.Errorf("unsupported type-combination strategy %q", rules.Strategy)
	}
	if !ok {
		return l, nil
	}
	// A type with no arguments is its family. Looked up as a declaration
	// it would take the dialect's defaults, as varchar is varchar(1).
	if len(out.Args) == 0 {
		if oid, err := a.cat.TypeOID(out.Name); err == nil {
			return exprType{typeOID: oid}, nil
		}
	}
	return a.lookupType(out), nil
}

// commonType uses the enabled policy for alternative-valued expressions.
// Currently only MSSQL combines these with its set-operation rules. Other
// dialects retain the first typed value; this is a compatibility fallback,
// not a claim that their databases use the same rules.
func (a *analyzer) commonType(types []exprType) (exprType, error) {
	var out exprType
	rules := a.cat.SetOperations()
	for i, t := range types {
		if t.typeOID == 0 && t.expr == nil {
			continue
		}
		switch {
		case out.typeOID == 0 && out.expr == nil:
			out = t
		case rules == nil:
			return out, nil
		default:
			var err error
			out, err = a.combineTypes(rules, out, t)
			if err != nil {
				return exprType{}, err
			}
		}
		if i > 0 {
			out.sourceClassOID, out.sourceAttributeOID, out.sourceTableAlias = 0, 0, ""
		}
	}
	return out, nil
}

// sameCategory reports whether two types are of one category, so that the
// length of one says something about the other: the 20 of varchar(20)
// widens an nvarchar(10), but not a datetime2(3).
func (a *analyzer) sameCategory(x, y string) bool {
	category := func(name string) string {
		oid, err := a.cat.TypeOID(name)
		if err != nil {
			return ""
		}
		t, err := a.cat.LookupType(a.familyOID(oid))
		if err != nil {
			return ""
		}
		return t.Category
	}
	cx := category(x)
	return cx != "" && cx == category(y)
}
