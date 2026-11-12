package analyzer

import (
	"slices"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// combineColumns types each column of a set operation from both of its
// branches by the dialect's rules. A column keeps the first branch's name
// and is read from no one table. It may be NULL in a UNION when either
// branch's may, in an INTERSECT when both may, and in an EXCEPT when the
// first branch's may, since every row it returns is one of those.
func (a *analyzer) combineColumns(rules *core.SetOperationRules, op ast.SetOperation, left, right []core.Column) []core.Column {
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
		t := a.combineTypes(rules, l, r)
		t.nullable = !col.NotNull
		col.TypeOID = t.typeOID
		col.DataType, col.IsArray = a.typeNameOf(t)
		col.Type = a.typeExprOf(t)
		out[i] = col
	}
	return out
}

// combineTypes is the type a column takes from two branches: the one of
// higher precedence, varying when either branch varies, with the length or
// precision that holds both.
func (a *analyzer) combineTypes(rules *core.SetOperationRules, l, r core.Column) exprType {
	lt, rt := a.exprOf(columnExprType(l)), a.exprOf(columnExprType(r))
	switch {
	case lt == nil:
		return columnExprType(r)
	case rt == nil:
		return columnExprType(l)
	}
	li, ri := slices.Index(rules.Precedence, lt.Name), slices.Index(rules.Precedence, rt.Name)
	if li < 0 || ri < 0 {
		return columnExprType(l)
	}
	winner, loser := lt, rt
	if ri < li {
		winner, loser = rt, lt
	}
	out := winner.WithNullable(false)
	if varying, ok := rules.Varying[out.Name]; ok && slices.Contains(mapValues(rules.Varying), loser.Name) {
		out.Name = varying
	}
	switch {
	case out.Name == rules.Decimal && rules.Decimal != "":
		combineDecimal(rules, out, loser)
	case len(out.Args) == 1 && len(loser.Args) == 1 && a.sameCategory(lt.Name, rt.Name):
		out.Args[0] = widerArg(out.Args[0], loser.Args[0])
	}
	return a.lookupType(out)
}

func mapValues(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	return out
}

// widerArg is the larger of two lengths or precisions, MAX above any.
func widerArg(a, b core.TypeArg) core.TypeArg {
	switch {
	case a.Ident != nil:
		return a
	case b.Ident != nil:
		return b
	case a.Int != nil && b.Int != nil && *b.Int > *a.Int:
		return b
	}
	return a
}

// combineDecimal gives the exact numeric out the precision and scale that
// hold both it and the other branch's type: the larger scale, and the
// larger count of integer digits, with scale giving way past the maximum
// precision. A family with no decimal form leaves out as it is.
func combineDecimal(rules *core.SetOperationRules, out, other *core.TypeExpr) {
	p1, s1, ok1 := decimalForm(rules, out)
	p2, s2, ok2 := decimalForm(rules, other)
	if !ok1 || !ok2 {
		return
	}
	scale := max(s1, s2)
	precision := max(p1-s1, p2-s2) + scale
	if limit := rules.DecimalMaxPrecision; limit > 0 && precision > limit {
		scale = max(scale-(precision-limit), 0)
		precision = limit
	}
	out.Args = []core.TypeArg{{Int: &precision}, {Int: &scale}}
}

func decimalForm(rules *core.SetOperationRules, t *core.TypeExpr) (int64, int64, bool) {
	if t.Name == rules.Decimal && len(t.Args) == 2 && t.Args[0].Int != nil && t.Args[1].Int != nil {
		return *t.Args[0].Int, *t.Args[1].Int, true
	}
	if form, ok := rules.AsDecimal[t.Name]; ok {
		return form[0], form[1], true
	}
	return 0, 0, false
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
