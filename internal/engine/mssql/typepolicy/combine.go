// Package typepolicy implements SQL Server's type-combination rules without
// depending on its parser or the shared analyzer.
package typepolicy

import (
	"slices"

	"github.com/sqlc-dev/sqlc/internal/core"
)

// Combine applies SQL Server's precedence and set-operation type modifiers.
// sameCategory says whether lengths/precisions of the two families are
// comparable. An unlisted family leaves type selection to the caller.
func Combine(rules *core.SetOperationRules, left, right *core.TypeExpr, sameCategory bool) (*core.TypeExpr, bool) {
	li, ri := slices.Index(rules.Precedence, left.Name), slices.Index(rules.Precedence, right.Name)
	if li < 0 || ri < 0 {
		return nil, false
	}
	winner, loser := left, right
	if ri < li {
		winner, loser = right, left
	}
	out := winner.WithNullable(false)
	if varying, ok := rules.Varying[out.Name]; ok {
		for _, name := range rules.Varying {
			if loser.Name == name {
				out.Name = varying
				break
			}
		}
	}
	switch {
	case out.Name == rules.Decimal && rules.Decimal != "":
		combineDecimal(rules, out, loser)
	case len(out.Args) == 1 && len(loser.Args) == 1 && sameCategory:
		out.Args[0] = widerArg(out.Args[0], loser.Args[0])
	}
	return out, true
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

// combineDecimal applies SQL Server's UNION/INTERSECT/EXCEPT formula:
// retain the larger scale and integer-digit count, reducing scale at the
// precision limit. This is not a cross-engine decimal-coercion rule.
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
