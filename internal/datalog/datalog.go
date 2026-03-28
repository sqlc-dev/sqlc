// Package datalog implements a self-contained Datalog evaluator supporting
// semi-naive bottom-up evaluation with stratified negation. It provides an
// in-memory deductive database suitable for computing fixed-point analyses
// over relational data.
//
// The evaluator works by repeatedly applying rules to derive new facts until
// no more new facts can be produced (a fixpoint). Stratified negation allows
// rules to reference the absence of facts, provided there are no recursive
// dependencies through negation.
package datalog

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// ---------------------------------------------------------------------------
// Symbol table — interned strings
// ---------------------------------------------------------------------------

// Symbol is an interned string represented as an integer index.
type Symbol int

// SymbolTable maps strings to Symbol IDs and back for efficient storage and
// comparison of string values.
type SymbolTable struct {
	forward map[string]Symbol
	reverse []string
}

// NewSymbolTable creates an empty symbol table.
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		forward: make(map[string]Symbol),
	}
}

// Intern adds a string to the table (if not already present) and returns its
// Symbol ID.
func (st *SymbolTable) Intern(s string) Symbol {
	if id, ok := st.forward[s]; ok {
		return id
	}
	id := Symbol(len(st.reverse))
	st.forward[s] = id
	st.reverse = append(st.reverse, s)
	return id
}

// Resolve returns the string for a given Symbol. It panics if the symbol is
// out of range.
func (st *SymbolTable) Resolve(sym Symbol) string {
	return st.reverse[int(sym)]
}

// ---------------------------------------------------------------------------
// Tuples and relations
// ---------------------------------------------------------------------------

// Tuple is a fixed-length slice of Symbol values representing a single ground
// fact.
type Tuple []Symbol

// tupleKey produces a string key for deduplication. The null byte separator
// cannot appear in interned strings.
func tupleKey(t Tuple) string {
	var b strings.Builder
	for i, s := range t {
		if i > 0 {
			b.WriteByte(0)
		}
		fmt.Fprintf(&b, "%d", int(s))
	}
	return b.String()
}

// relation is a named collection of tuples with deduplication.
type relation struct {
	arity  int
	tuples []Tuple
	index  map[string]struct{}
}

func newRelation(arity int) *relation {
	return &relation{
		arity: arity,
		index: make(map[string]struct{}),
	}
}

// add inserts a tuple if it is not already present. Returns true if the tuple
// was new.
func (r *relation) add(t Tuple) bool {
	k := tupleKey(t)
	if _, exists := r.index[k]; exists {
		return false
	}
	r.index[k] = struct{}{}
	cp := make(Tuple, len(t))
	copy(cp, t)
	r.tuples = append(r.tuples, cp)
	return true
}

// contains checks whether the tuple exists.
func (r *relation) contains(t Tuple) bool {
	_, exists := r.index[tupleKey(t)]
	return exists
}

// clone returns a deep copy of the relation.
func (r *relation) clone() *relation {
	nr := newRelation(r.arity)
	for _, t := range r.tuples {
		nr.add(t)
	}
	return nr
}

// ---------------------------------------------------------------------------
// Terms, atoms, and rules
// ---------------------------------------------------------------------------

// TermExpr is an opaque term expression used when constructing rules via the
// builder API. It represents either a variable reference or a constant value.
type TermExpr struct {
	isVar bool
	name  string // variable name when isVar; constant string value otherwise
}

// Var creates a variable term expression.
func Var(name string) TermExpr { return TermExpr{isVar: true, name: name} }

// Const creates a constant term expression.
func Const(value string) TermExpr { return TermExpr{isVar: false, name: value} }

// Term is a resolved term inside a rule — either a variable or a constant
// symbol.
type Term struct {
	IsVariable bool
	VarName    string
	Value      Symbol
}

// Atom is a predicate name paired with a list of terms. It appears in rule
// heads and bodies.
type Atom struct {
	Predicate string
	Terms     []Term
}

// Rule is a Horn clause: a head atom derived from a conjunction of positive
// body atoms and negated body atoms.
type Rule struct {
	Head    Atom
	Body    []Atom
	Negated []Atom
}

// ---------------------------------------------------------------------------
// Rule builder
// ---------------------------------------------------------------------------

type atomExpr struct {
	pred  string
	terms []TermExpr
}

// RuleBuilder provides a fluent API for constructing rules.
type RuleBuilder struct {
	st        *SymbolTable
	headPred  string
	headTerms []TermExpr
	body      []atomExpr
	negBody   []atomExpr
}

// NewRule begins building a rule with the given head predicate and terms. The
// symbol table is used to intern constant values when Build is called.
func NewRule(st *SymbolTable, headPred string, headTerms ...TermExpr) *RuleBuilder {
	return &RuleBuilder{
		st:        st,
		headPred:  headPred,
		headTerms: headTerms,
	}
}

// Where adds a positive body atom to the rule.
func (rb *RuleBuilder) Where(pred string, terms ...TermExpr) *RuleBuilder {
	rb.body = append(rb.body, atomExpr{pred: pred, terms: terms})
	return rb
}

// WhereNot adds a negated body atom to the rule.
func (rb *RuleBuilder) WhereNot(pred string, terms ...TermExpr) *RuleBuilder {
	rb.negBody = append(rb.negBody, atomExpr{pred: pred, terms: terms})
	return rb
}

// Build finalises the rule, interning any constant values that appear in the
// terms.
func (rb *RuleBuilder) Build() Rule {
	resolve := func(te TermExpr) Term {
		if te.isVar {
			return Term{IsVariable: true, VarName: te.name}
		}
		return Term{IsVariable: false, Value: rb.st.Intern(te.name)}
	}
	makeAtom := func(ae atomExpr) Atom {
		a := Atom{Predicate: ae.pred, Terms: make([]Term, len(ae.terms))}
		for i, te := range ae.terms {
			a.Terms[i] = resolve(te)
		}
		return a
	}

	head := Atom{Predicate: rb.headPred, Terms: make([]Term, len(rb.headTerms))}
	for i, te := range rb.headTerms {
		head.Terms[i] = resolve(te)
	}

	r := Rule{Head: head}
	for _, ae := range rb.body {
		r.Body = append(r.Body, makeAtom(ae))
	}
	for _, ae := range rb.negBody {
		r.Negated = append(r.Negated, makeAtom(ae))
	}
	return r
}

// ---------------------------------------------------------------------------
// Program
// ---------------------------------------------------------------------------

// Program holds the initial facts (extensional database) and the rules that
// will be evaluated to compute derived facts (intensional database).
type Program struct {
	st    *SymbolTable
	facts map[string]*relation
	rules []Rule
}

// NewProgram creates an empty program bound to the given symbol table.
func NewProgram(st *SymbolTable) *Program {
	return &Program{
		st:    st,
		facts: make(map[string]*relation),
	}
}

// AddFact adds a ground fact. String values are interned automatically.
func (p *Program) AddFact(predicate string, values ...string) {
	t := make(Tuple, len(values))
	for i, v := range values {
		t[i] = p.st.Intern(v)
	}
	rel, ok := p.facts[predicate]
	if !ok {
		rel = newRelation(len(values))
		p.facts[predicate] = rel
	}
	rel.add(t)
}

// AddRule adds a derivation rule to the program.
func (p *Program) AddRule(r Rule) {
	p.rules = append(p.rules, r)
}

// ---------------------------------------------------------------------------
// Database — query interface over computed facts
// ---------------------------------------------------------------------------

// Database holds the final set of facts after evaluation.
type Database struct {
	st    *SymbolTable
	facts map[string]*relation
}

// Query returns all tuples for the given predicate, with symbol values
// resolved back to strings. The results are sorted lexicographically.
func (db *Database) Query(predicate string) [][]string {
	rel, ok := db.facts[predicate]
	if !ok {
		return nil
	}
	result := make([][]string, len(rel.tuples))
	for i, t := range rel.tuples {
		row := make([]string, len(t))
		for j, sym := range t {
			row[j] = db.st.Resolve(sym)
		}
		result[i] = row
	}
	return result
}

// Contains checks whether a specific ground fact exists in the database.
func (db *Database) Contains(predicate string, values ...string) bool {
	rel, ok := db.facts[predicate]
	if !ok {
		return false
	}
	t := make(Tuple, len(values))
	for i, v := range values {
		sym, exists := db.st.forward[v]
		if !exists {
			return false
		}
		t[i] = sym
	}
	return rel.contains(t)
}

// ---------------------------------------------------------------------------
// Stratification
// ---------------------------------------------------------------------------

// stratum groups predicates that can be evaluated together.
type stratum struct {
	predicates map[string]bool
	rules      []Rule
}

// stratify computes strata by analyzing negation dependencies among rules.
// Predicates that appear in negated body atoms must be fully computed before
// the stratum that negates them. Returns an error if a cycle through negation
// is detected, which makes stratification impossible.
func stratify(rules []Rule) ([]stratum, error) {
	// Collect all predicates mentioned in rule heads.
	preds := make(map[string]bool)
	for _, r := range rules {
		preds[r.Head.Predicate] = true
	}

	// Build dependency graph.
	// posEdges: head depends positively on body predicate.
	// negEdges: head depends (through negation) on negated body predicate.
	posEdges := make(map[string]map[string]bool)
	negEdges := make(map[string]map[string]bool)
	ensureSet := func(m map[string]map[string]bool, k string) {
		if m[k] == nil {
			m[k] = make(map[string]bool)
		}
	}

	for _, r := range rules {
		h := r.Head.Predicate
		ensureSet(posEdges, h)
		ensureSet(negEdges, h)
		for _, b := range r.Body {
			if preds[b.Predicate] {
				posEdges[h][b.Predicate] = true
			}
		}
		for _, b := range r.Negated {
			preds[b.Predicate] = true // ensure negated predicates are tracked
			negEdges[h][b.Predicate] = true
		}
	}

	// Assign stratum numbers using iterative relaxation.
	// A predicate's stratum must be:
	//   >= stratum of any positive dependency
	//   >  stratum of any negated dependency
	stratumOf := make(map[string]int)
	for p := range preds {
		stratumOf[p] = 0
	}

	changed := true
	maxIter := len(preds) + 1
	for iter := 0; changed && iter < maxIter; iter++ {
		changed = false
		for p := range preds {
			cur := stratumOf[p]
			for dep := range posEdges[p] {
				if stratumOf[dep] > cur {
					stratumOf[p] = stratumOf[dep]
					changed = true
					cur = stratumOf[p]
				}
			}
			for dep := range negEdges[p] {
				need := stratumOf[dep] + 1
				if need > cur {
					stratumOf[p] = need
					changed = true
					cur = stratumOf[p]
				}
			}
		}
	}

	// If still changing after maxIter iterations, a negation cycle exists.
	if changed {
		return nil, errors.New("datalog: negation cycle detected; stratification is impossible")
	}

	// Group predicates by stratum number.
	maxStratum := 0
	for _, s := range stratumOf {
		if s > maxStratum {
			maxStratum = s
		}
	}

	strata := make([]stratum, maxStratum+1)
	for i := range strata {
		strata[i].predicates = make(map[string]bool)
	}
	for p, s := range stratumOf {
		strata[s].predicates[p] = true
	}

	// Assign each rule to the stratum of its head predicate.
	for _, r := range rules {
		s := stratumOf[r.Head.Predicate]
		strata[s].rules = append(strata[s].rules, r)
	}

	return strata, nil
}

// ---------------------------------------------------------------------------
// Evaluation engine
// ---------------------------------------------------------------------------

// binding maps variable names to symbol values during rule evaluation.
type binding map[string]Symbol

// copyBinding returns a shallow copy of a binding.
func copyBinding(b binding) binding {
	nb := make(binding, len(b))
	for k, v := range b {
		nb[k] = v
	}
	return nb
}

// matchAtom attempts to unify an atom against a tuple under the given binding.
// It returns an extended binding on success, or nil if unification fails.
func matchAtom(a Atom, t Tuple, b binding) binding {
	if len(a.Terms) != len(t) {
		return nil
	}
	nb := copyBinding(b)
	for i, term := range a.Terms {
		if term.IsVariable {
			if prev, bound := nb[term.VarName]; bound {
				if prev != t[i] {
					return nil
				}
			} else {
				nb[term.VarName] = t[i]
			}
		} else {
			if term.Value != t[i] {
				return nil
			}
		}
	}
	return nb
}

// projectHead builds a tuple from the head atom using the given binding.
func projectHead(head Atom, b binding) Tuple {
	t := make(Tuple, len(head.Terms))
	for i, term := range head.Terms {
		if term.IsVariable {
			t[i] = b[term.VarName]
		} else {
			t[i] = term.Value
		}
	}
	return t
}

// evaluateRuleBody performs a nested-loop join over the positive body atoms,
// collecting all satisfying bindings. The useDelta parameter specifies which
// body atom index should read from the delta relation instead of the full
// fact set; pass -1 to use full relations for all atoms.
func evaluateRuleBody(
	body []Atom,
	facts map[string]*relation,
	delta map[string]*relation,
	useDelta int,
) []binding {
	bindings := []binding{make(binding)}

	for i, atom := range body {
		source := facts[atom.Predicate]
		if i == useDelta {
			if d, ok := delta[atom.Predicate]; ok {
				source = d
			} else {
				return nil
			}
		}
		if source == nil {
			return nil
		}

		var next []binding
		for _, b := range bindings {
			for _, t := range source.tuples {
				nb := matchAtom(atom, t, b)
				if nb != nil {
					next = append(next, nb)
				}
			}
		}
		bindings = next
		if len(bindings) == 0 {
			return nil
		}
	}
	return bindings
}

// checkNegation returns true if none of the negated atoms are satisfied under
// the binding (i.e., all negation conditions hold).
func checkNegation(negated []Atom, facts map[string]*relation, b binding) bool {
	for _, atom := range negated {
		rel := facts[atom.Predicate]
		if rel == nil {
			continue // no facts — negation trivially satisfied
		}
		// Try to build a fully-ground tuple from the binding.
		t := make(Tuple, len(atom.Terms))
		fullyBound := true
		for i, term := range atom.Terms {
			if term.IsVariable {
				val, ok := b[term.VarName]
				if !ok {
					fullyBound = false
					break
				}
				t[i] = val
			} else {
				t[i] = term.Value
			}
		}
		if fullyBound {
			if rel.contains(t) {
				return false
			}
		} else {
			// Partially bound: scan for any matching tuple.
			if anyMatch(atom, rel, b) {
				return false
			}
		}
	}
	return true
}

// anyMatch checks whether any tuple in rel matches the partially-bound atom.
func anyMatch(atom Atom, rel *relation, b binding) bool {
	for _, t := range rel.tuples {
		if matchAtom(atom, t, b) != nil {
			return true
		}
	}
	return false
}

// evaluateRule derives new tuples for a single rule using the semi-naive
// strategy: for each positive body atom position, it evaluates the rule with
// that position reading from the delta relation, ensuring at least one delta
// atom participates.
func evaluateRule(
	r Rule,
	facts map[string]*relation,
	delta map[string]*relation,
) []Tuple {
	if len(r.Body) == 0 {
		// Rules with no body are treated as ground facts. They are handled via
		// EDB seeding; returning nothing here avoids infinite derivation.
		return nil
	}

	seen := make(map[string]struct{})
	var results []Tuple

	for i, atom := range r.Body {
		if _, ok := delta[atom.Predicate]; !ok {
			continue
		}
		bindings := evaluateRuleBody(r.Body, facts, delta, i)
		for _, b := range bindings {
			if !checkNegation(r.Negated, facts, b) {
				continue
			}
			t := projectHead(r.Head, b)
			k := tupleKey(t)
			if _, dup := seen[k]; dup {
				continue
			}
			seen[k] = struct{}{}
			results = append(results, t)
		}
	}
	return results
}

// ensureRelation returns the relation for pred, creating it with the given
// arity if absent.
func ensureRelation(m map[string]*relation, pred string, arity int) *relation {
	if r, ok := m[pred]; ok {
		return r
	}
	r := newRelation(arity)
	m[pred] = r
	return r
}

// ---------------------------------------------------------------------------
// Evaluate — main entry point
// ---------------------------------------------------------------------------

// Evaluate runs the Datalog program to a fixpoint using semi-naive evaluation
// with stratified negation and returns the resulting database. It returns an
// error if stratification fails (e.g., due to a negation cycle) or if arity
// mismatches are detected between rules and facts for the same predicate.
func (p *Program) Evaluate() (*Database, error) {
	// Validate arity consistency between rules and existing facts.
	arities := make(map[string]int)
	for pred, rel := range p.facts {
		arities[pred] = rel.arity
	}
	for _, r := range p.rules {
		headArity := len(r.Head.Terms)
		if prev, ok := arities[r.Head.Predicate]; ok {
			if prev != headArity {
				return nil, fmt.Errorf(
					"datalog: arity mismatch for predicate %q: have %d, rule head has %d",
					r.Head.Predicate, prev, headArity,
				)
			}
		} else {
			arities[r.Head.Predicate] = headArity
		}
	}

	// Stratify rules.
	strata, err := stratify(p.rules)
	if err != nil {
		return nil, err
	}

	// Initialize the fact base with a deep copy of EDB facts.
	facts := make(map[string]*relation)
	for pred, rel := range p.facts {
		facts[pred] = rel.clone()
	}

	// Process each stratum in order.
	for _, st := range strata {
		if len(st.rules) == 0 {
			continue
		}

		// Seed delta with current facts for predicates relevant to this stratum.
		// This includes both IDB predicates defined in this stratum and any EDB
		// predicates referenced in rule bodies, so the first iteration can find
		// matches.
		delta := make(map[string]*relation)
		for pred := range st.predicates {
			if rel, ok := facts[pred]; ok {
				delta[pred] = rel.clone()
			}
		}
		for _, r := range st.rules {
			for _, atom := range r.Body {
				if _, ok := delta[atom.Predicate]; !ok {
					if rel, ok := facts[atom.Predicate]; ok {
						delta[atom.Predicate] = rel.clone()
					}
				}
			}
		}

		// Semi-naive fixpoint loop.
		for {
			nextDelta := make(map[string]*relation)
			for _, r := range st.rules {
				derived := evaluateRule(r, facts, delta)
				for _, t := range derived {
					pred := r.Head.Predicate
					rel := ensureRelation(facts, pred, len(t))
					if rel.add(t) {
						nd := ensureRelation(nextDelta, pred, len(t))
						nd.add(t)
					}
				}
			}
			if len(nextDelta) == 0 {
				break // fixpoint reached
			}
			delta = nextDelta
		}
	}

	// Sort tuples for deterministic query results.
	for _, rel := range facts {
		sort.Slice(rel.tuples, func(i, j int) bool {
			a, b := rel.tuples[i], rel.tuples[j]
			for k := 0; k < len(a) && k < len(b); k++ {
				if a[k] != b[k] {
					return a[k] < b[k]
				}
			}
			return len(a) < len(b)
		})
	}

	return &Database{st: p.st, facts: facts}, nil
}
