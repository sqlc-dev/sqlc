package core

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/sqlc-dev/sqlc/internal/core/catalogdb"
)

// A dialect describes, as data, what it does to a type before storing it.
// The seed loads three kinds of it from dialect.json: rewrites, which turn
// one expression into another (float(24) into real, a bare decimal into
// decimal(18, 0)); identifier words and positions, which say where a bare
// word is a word rather than a type (the max of nvarchar(max), the
// function of SimpleAggregateFunction(sum, UInt64)); and an affinity rule,
// which says what a family the schema declares and the seed does not list
// stands on. None of it is code, so an engine adds a dialect by writing
// files.

// FlagIdents holds the words that are identifiers wherever they stand as a
// type argument, comma-separated.
const FlagIdents = "types.idents"

// FlagIdentArgs holds the argument positions that are identifiers in a
// family, as "family:1,2;family:1".
const FlagIdentArgs = "types.ident_args"

// typeRewrite is one rewrite, parsed: a pattern whose $n arguments bind
// whatever stands there, a template the bindings are substituted into, and
// a condition on a binding.
type typeRewrite struct {
	pattern  *TypeExpr
	template *TypeExpr
	cond     rewriteCond
}

// rewriteCond bounds an integer binding: "$1 <= 24".
type rewriteCond struct {
	binding string
	op      string
	value   int64
}

// rules is what the catalog knows of its dialect's rewriting, read from the
// tables once and kept. The seed invalidates it as it adds to them.
type rules struct {
	mu        sync.Mutex
	loaded    bool
	rewrites  []typeRewrite
	idents    map[string]bool
	identArgs map[string][]int
}

func (c *Catalog) invalidateRules() {
	c.rules.mu.Lock()
	c.rules.loaded = false
	c.rules.mu.Unlock()
}

// loadRules reads the dialect's rewrites and identifier settings.
func (c *Catalog) loadRules() (*rules, error) {
	r := &c.rules
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.loaded {
		return r, nil
	}
	r.rewrites, r.idents, r.identArgs = nil, map[string]bool{}, map[string][]int{}
	if c.dialectOID == 0 {
		r.loaded = true
		return r, nil
	}
	rows, err := c.q.ListTypeRewrites(context.Background(), c.dialectOID)
	if err != nil {
		return nil, fmt.Errorf("type rewrites: %w", err)
	}
	for _, row := range rows {
		rw := typeRewrite{pattern: ParseTypeExpr(row.Pattern), template: ParseTypeExpr(row.Template)}
		if row.Cond != "" {
			cond, err := parseRewriteCond(row.Cond)
			if err != nil {
				return nil, fmt.Errorf("type rewrite %q: %w", row.Pattern, err)
			}
			rw.cond = cond
		}
		r.rewrites = append(r.rewrites, rw)
	}
	if idents, _ := c.DialectFlag(c.dialectOID, FlagIdents); idents != "" {
		for _, w := range strings.Split(idents, ",") {
			r.idents[strings.ToLower(strings.TrimSpace(w))] = true
		}
	}
	if identArgs, _ := c.DialectFlag(c.dialectOID, FlagIdentArgs); identArgs != "" {
		for _, entry := range strings.Split(identArgs, ";") {
			family, positions, ok := strings.Cut(entry, ":")
			if !ok {
				continue
			}
			for _, p := range strings.Split(positions, ",") {
				if n, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
					family = strings.ToLower(strings.TrimSpace(family))
					r.identArgs[family] = append(r.identArgs[family], n)
				}
			}
		}
	}
	r.loaded = true
	return r, nil
}

// parseRewriteCond reads "$1 <= 24".
func parseRewriteCond(s string) (rewriteCond, error) {
	fields := strings.Fields(s)
	if len(fields) != 3 || !strings.HasPrefix(fields[0], "$") {
		return rewriteCond{}, fmt.Errorf("condition %q: want \"$n op value\"", s)
	}
	switch fields[1] {
	case "<", "<=", "=", ">=", ">", "!=":
	default:
		return rewriteCond{}, fmt.Errorf("condition %q: unknown operator", s)
	}
	v, err := strconv.ParseInt(fields[2], 10, 64)
	if err != nil {
		return rewriteCond{}, fmt.Errorf("condition %q: %w", s, err)
	}
	return rewriteCond{binding: fields[0], op: fields[1], value: v}, nil
}

func (cond rewriteCond) holds(bindings map[string]TypeArg) bool {
	if cond.binding == "" {
		return true
	}
	a, ok := bindings[cond.binding]
	if !ok || a.Int == nil {
		return false
	}
	switch cond.op {
	case "<":
		return *a.Int < cond.value
	case "<=":
		return *a.Int <= cond.value
	case "=":
		return *a.Int == cond.value
	case ">=":
		return *a.Int >= cond.value
	case ">":
		return *a.Int > cond.value
	case "!=":
		return *a.Int != cond.value
	}
	return false
}

// AddTypeRewrite records a rewrite of the catalog's dialect, applied after
// the ones recorded before it.
func (c *Catalog) AddTypeRewrite(ord int, pattern, template, cond string) error {
	err := c.q.CreateTypeRewrite(context.Background(), catalogdb.CreateTypeRewriteParams{
		DialectOid: c.dialectOID,
		Ord:        int64(ord),
		Pattern:    pattern,
		Template:   template,
		Cond:       cond,
	})
	if err != nil {
		return fmt.Errorf("type rewrite %q: %w", pattern, err)
	}
	c.invalidateRules()
	return nil
}

// AddTypeAffinity records the next step of the dialect's affinity rule: a
// family whose name contains one of words, upper-cased, stands on typeOID.
// No words is the default that ends the rule.
func (c *Catalog) AddTypeAffinity(ord int, words []string, typeOID int64) error {
	err := c.q.CreateTypeAffinity(context.Background(), catalogdb.CreateTypeAffinityParams{
		DialectOid: c.dialectOID,
		Ord:        int64(ord),
		Words:      strings.ToUpper(strings.Join(words, ",")),
		TypeOid:    typeOID,
	})
	if err != nil {
		return fmt.Errorf("type affinity %d: %w", ord, err)
	}
	return nil
}

// userTypeBase is the type an unseeded family stands on by the dialect's
// affinity rule, or 0 when the dialect has none or none of it matches.
func (c *Catalog) userTypeBase(name string) (int64, error) {
	if c.dialectOID == 0 {
		return 0, nil
	}
	rows, err := c.q.ListTypeAffinities(context.Background(), c.dialectOID)
	if err != nil {
		return 0, fmt.Errorf("type affinities: %w", err)
	}
	upper := strings.ToUpper(name)
	for _, row := range rows {
		if row.Words == "" {
			return row.TypeOid, nil
		}
		for _, w := range strings.Split(row.Words, ",") {
			if w != "" && strings.Contains(upper, w) {
				return row.TypeOid, nil
			}
		}
	}
	return 0, nil
}

// canonicalize applies the dialect's identifier settings and rewrites to
// an expression: the first rewrite whose pattern matches is applied, and
// its result is not rewritten again.
func (c *Catalog) canonicalize(t *TypeExpr) (*TypeExpr, error) {
	r, err := c.loadRules()
	if err != nil {
		return nil, err
	}
	t = r.identify(t)
	name := strings.ToLower(t.Name)
	for _, rw := range r.rewrites {
		if strings.ToLower(rw.pattern.Name) != name || len(rw.pattern.Args) != len(t.Args) {
			continue
		}
		bindings, ok := matchArgs(rw.pattern.Args, t.Args)
		if !ok || !rw.cond.holds(bindings) {
			continue
		}
		out := substitute(rw.template, bindings)
		out.Nullable = t.Nullable
		return out, nil
	}
	return t, nil
}

// identify turns the bare words the dialect calls identifiers into
// identifier arguments.
func (r *rules) identify(t *TypeExpr) *TypeExpr {
	if len(t.Args) == 0 || (len(r.idents) == 0 && len(r.identArgs) == 0) {
		return t
	}
	positions := r.identArgs[strings.ToLower(t.Name)]
	var out *TypeExpr
	for i, a := range t.Args {
		if a.Type == nil || len(a.Type.Args) != 0 {
			continue
		}
		word := strings.ToLower(a.Type.Name)
		if !r.idents[word] && !containsInt(positions, i+1) {
			continue
		}
		if out == nil {
			out = t.Clone()
		}
		out.Args[i] = TypeArg{Label: a.Label, Ident: &word}
	}
	if out == nil {
		return t
	}
	return out
}

func containsInt(list []int, n int) bool {
	for _, v := range list {
		if v == n {
			return true
		}
	}
	return false
}

// matchArgs matches a pattern's arguments against an expression's, binding
// each $n to what stands in its place and requiring a literal to be equal.
func matchArgs(pattern, args []TypeArg) (map[string]TypeArg, bool) {
	bindings := map[string]TypeArg{}
	for i, p := range pattern {
		a := args[i]
		switch {
		case p.Type != nil && strings.HasPrefix(p.Type.Name, "$"):
			bindings[p.Type.Name] = a
		case p.Type != nil:
			if a.Type == nil || !strings.EqualFold(a.Type.Name, p.Type.Name) {
				return nil, false
			}
		case p.Int != nil:
			if a.Int == nil || *a.Int != *p.Int {
				return nil, false
			}
		case p.String != nil:
			if a.String == nil || *a.String != *p.String {
				return nil, false
			}
		case p.Ident != nil:
			if a.Ident == nil || *a.Ident != *p.Ident {
				return nil, false
			}
		default:
			return nil, false
		}
	}
	return bindings, true
}

// substitute fills a template's $n arguments from the bindings.
func substitute(template *TypeExpr, bindings map[string]TypeArg) *TypeExpr {
	out := template.Clone()
	for i, a := range out.Args {
		if a.Type != nil && strings.HasPrefix(a.Type.Name, "$") {
			if bound, ok := bindings[a.Type.Name]; ok {
				label := a.Label
				out.Args[i] = bound
				if label != "" {
					out.Args[i].Label = label
				}
			}
		} else if a.Type != nil {
			out.Args[i].Type = substitute(a.Type, bindings)
		}
	}
	return out
}
