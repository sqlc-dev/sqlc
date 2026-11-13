// Package typepolicy implements SQL Server's type-combination rules without
// depending on its parser or the shared analyzer.
package typepolicy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/sqlc-dev/sqlc/internal/core"
)

type rules struct {
	Precedence          []string            `json:"precedence"`
	Varying             map[string]string   `json:"varying,omitempty"`
	Decimal             string              `json:"decimal,omitempty"`
	DecimalMaxPrecision int64               `json:"decimal_max_precision,omitempty"`
	AsDecimal           map[string][2]int64 `json:"as_decimal,omitempty"`
}

type combiner struct {
	cat   *core.Catalog
	rules rules
}

// New binds SQL Server's policy to a freshly seeded or restored catalog.
func New(cat *core.Catalog) (core.TypeCombiner, error) {
	data, err := cat.TypeCombinationData()
	if err != nil {
		return nil, err
	}
	rules, err := decodeRules(data)
	if err != nil {
		return nil, fmt.Errorf("mssql: type-combination rules: %w", err)
	}
	return &combiner{cat: cat, rules: rules}, nil
}

func decodeRules(data []byte) (rules, error) {
	var out rules
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&out); err != nil {
		return rules{}, err
	}
	if len(out.Precedence) == 0 {
		return rules{}, fmt.Errorf("precedence must not be empty")
	}
	return out, nil
}

func (p *combiner) Combine(req core.TypeCombinationRequest) (core.TypeCombinationResult, error) {
	switch req.Context {
	case core.CombineUnion, core.CombineIntersect, core.CombineExcept,
		core.CombineCase, core.CombineCoalesce, core.CombineGreatest, core.CombineLeast:
	default:
		return core.TypeCombinationResult{}, nil
	}
	var out *core.TypeExpr
	for _, candidate := range req.Candidates {
		if candidate == nil {
			continue
		}
		if out == nil {
			out = candidate
		} else {
			out = p.combine(out, candidate)
		}
	}
	return core.TypeCombinationResult{Handled: true, Type: out.WithNullable(false)}, nil
}

// combine applies SQL Server's precedence and set-operation type modifiers.
// An unlisted family retains the left type, as the existing inference did.
func (p *combiner) combine(left, right *core.TypeExpr) *core.TypeExpr {
	rules := &p.rules
	li, ri := slices.Index(rules.Precedence, left.Name), slices.Index(rules.Precedence, right.Name)
	if li < 0 || ri < 0 {
		return left
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
	case len(out.Args) == 1 && len(loser.Args) == 1 && p.sameCategory(left.Name, right.Name):
		out.Args[0] = widerArg(out.Args[0], loser.Args[0])
	}
	return out
}

func (p *combiner) sameCategory(x, y string) bool {
	category := func(name string) string {
		oid, err := p.cat.TypeOID(name)
		if err != nil {
			return ""
		}
		chain := p.cat.ResolutionChain(oid)
		t, err := p.cat.LookupType(chain[len(chain)-1])
		if err != nil {
			return ""
		}
		return t.Category
	}
	cx := category(x)
	return cx != "" && cx == category(y)
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
func combineDecimal(rules *rules, out, other *core.TypeExpr) {
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

func decimalForm(rules *rules, t *core.TypeExpr) (int64, int64, bool) {
	if t.Name == rules.Decimal && len(t.Args) == 2 && t.Args[0].Int != nil && t.Args[1].Int != nil {
		return *t.Args[0].Int, *t.Args[1].Int, true
	}
	if form, ok := rules.AsDecimal[t.Name]; ok {
		return form[0], form[1], true
	}
	return 0, 0, false
}
