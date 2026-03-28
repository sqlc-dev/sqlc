package datalog

import (
	"testing"
)

func TestSymbolTable(t *testing.T) {
	st := NewSymbolTable()
	a := st.Intern("hello")
	b := st.Intern("world")
	c := st.Intern("hello")

	if a != c {
		t.Fatalf("expected same symbol for same string, got %d and %d", a, c)
	}
	if a == b {
		t.Fatal("expected different symbols for different strings")
	}
	if st.Resolve(a) != "hello" {
		t.Fatalf("expected 'hello', got %q", st.Resolve(a))
	}
	if st.Resolve(b) != "world" {
		t.Fatalf("expected 'world', got %q", st.Resolve(b))
	}
}

func TestTransitiveClosure(t *testing.T) {
	// Classic Datalog example: compute reachability in a graph.
	// edge(a,b), edge(b,c), edge(c,d)
	// reach(X,Y) :- edge(X,Y).
	// reach(X,Y) :- reach(X,Z), edge(Z,Y).
	st := NewSymbolTable()
	p := NewProgram(st)

	p.AddFact("edge", "a", "b")
	p.AddFact("edge", "b", "c")
	p.AddFact("edge", "c", "d")

	p.AddRule(NewRule(st, "reach", Var("X"), Var("Y")).
		Where("edge", Var("X"), Var("Y")).
		Build())

	p.AddRule(NewRule(st, "reach", Var("X"), Var("Y")).
		Where("reach", Var("X"), Var("Z")).
		Where("edge", Var("Z"), Var("Y")).
		Build())

	db, err := p.Evaluate()
	if err != nil {
		t.Fatal(err)
	}

	// Should derive: reach(a,b), reach(a,c), reach(a,d), reach(b,c), reach(b,d), reach(c,d)
	for _, tc := range []struct{ from, to string }{
		{"a", "b"}, {"a", "c"}, {"a", "d"},
		{"b", "c"}, {"b", "d"},
		{"c", "d"},
	} {
		if !db.Contains("reach", tc.from, tc.to) {
			t.Errorf("expected reach(%s, %s)", tc.from, tc.to)
		}
	}

	// Should NOT have: reach(d, a), reach(b, a), etc.
	if db.Contains("reach", "d", "a") {
		t.Error("unexpected reach(d, a)")
	}
	if db.Contains("reach", "b", "a") {
		t.Error("unexpected reach(b, a)")
	}

	results := db.Query("reach")
	if len(results) != 6 {
		t.Errorf("expected 6 reach facts, got %d", len(results))
	}
}

func TestStratifiedNegation(t *testing.T) {
	// alive(X) :- person(X), not dead(X).
	st := NewSymbolTable()
	p := NewProgram(st)

	p.AddFact("person", "alice")
	p.AddFact("person", "bob")
	p.AddFact("person", "carol")
	p.AddFact("dead", "bob")

	p.AddRule(NewRule(st, "alive", Var("X")).
		Where("person", Var("X")).
		WhereNot("dead", Var("X")).
		Build())

	db, err := p.Evaluate()
	if err != nil {
		t.Fatal(err)
	}

	if !db.Contains("alive", "alice") {
		t.Error("expected alive(alice)")
	}
	if db.Contains("alive", "bob") {
		t.Error("unexpected alive(bob)")
	}
	if !db.Contains("alive", "carol") {
		t.Error("expected alive(carol)")
	}
}

func TestNegationCycleDetected(t *testing.T) {
	// p(X) :- q(X), not p(X). — this is a negation cycle
	st := NewSymbolTable()
	p := NewProgram(st)

	p.AddFact("q", "a")

	p.AddRule(NewRule(st, "p", Var("X")).
		Where("q", Var("X")).
		WhereNot("p", Var("X")).
		Build())

	_, err := p.Evaluate()
	if err == nil {
		t.Fatal("expected error for negation cycle")
	}
}

func TestConstants(t *testing.T) {
	// type_is_numeric(X) :- base_type(X, "numeric").
	st := NewSymbolTable()
	p := NewProgram(st)

	p.AddFact("base_type", "int4", "numeric")
	p.AddFact("base_type", "int8", "numeric")
	p.AddFact("base_type", "text", "string")

	p.AddRule(NewRule(st, "numeric_type", Var("X")).
		Where("base_type", Var("X"), Const("numeric")).
		Build())

	db, err := p.Evaluate()
	if err != nil {
		t.Fatal(err)
	}

	if !db.Contains("numeric_type", "int4") {
		t.Error("expected numeric_type(int4)")
	}
	if !db.Contains("numeric_type", "int8") {
		t.Error("expected numeric_type(int8)")
	}
	if db.Contains("numeric_type", "text") {
		t.Error("unexpected numeric_type(text)")
	}
}

func TestMultipleStrata(t *testing.T) {
	// Stratum 0: castable(X,Y) :- implicit_cast(X,Y).
	// Stratum 0: castable(X,Y) :- castable(X,Z), implicit_cast(Z,Y).
	// Stratum 1: not_castable(X,Y) :- type(X), type(Y), not castable(X,Y).
	st := NewSymbolTable()
	p := NewProgram(st)

	p.AddFact("type", "int4")
	p.AddFact("type", "int8")
	p.AddFact("type", "text")
	p.AddFact("implicit_cast", "int4", "int8")

	p.AddRule(NewRule(st, "castable", Var("X"), Var("Y")).
		Where("implicit_cast", Var("X"), Var("Y")).
		Build())
	p.AddRule(NewRule(st, "castable", Var("X"), Var("Y")).
		Where("castable", Var("X"), Var("Z")).
		Where("implicit_cast", Var("Z"), Var("Y")).
		Build())
	p.AddRule(NewRule(st, "not_castable", Var("X"), Var("Y")).
		Where("type", Var("X")).
		Where("type", Var("Y")).
		WhereNot("castable", Var("X"), Var("Y")).
		Build())

	db, err := p.Evaluate()
	if err != nil {
		t.Fatal(err)
	}

	if !db.Contains("castable", "int4", "int8") {
		t.Error("expected castable(int4, int8)")
	}
	if !db.Contains("not_castable", "int4", "text") {
		t.Error("expected not_castable(int4, text)")
	}
	// int4 IS castable to int8, so it should NOT be in not_castable
	if db.Contains("not_castable", "int4", "int8") {
		t.Error("unexpected not_castable(int4, int8)")
	}
}

func TestArityMismatch(t *testing.T) {
	st := NewSymbolTable()
	p := NewProgram(st)

	p.AddFact("r", "a", "b")

	// Rule head has arity 1, but fact has arity 2
	p.AddRule(NewRule(st, "r", Var("X")).
		Where("r", Var("X"), Var("Y")).
		Build())

	_, err := p.Evaluate()
	if err == nil {
		t.Fatal("expected arity mismatch error")
	}
}

func TestEmptyProgram(t *testing.T) {
	st := NewSymbolTable()
	p := NewProgram(st)

	db, err := p.Evaluate()
	if err != nil {
		t.Fatal(err)
	}

	results := db.Query("anything")
	if results != nil {
		t.Errorf("expected nil results, got %v", results)
	}
}

func TestQueryNonexistent(t *testing.T) {
	st := NewSymbolTable()
	p := NewProgram(st)
	p.AddFact("a", "1")

	db, err := p.Evaluate()
	if err != nil {
		t.Fatal(err)
	}

	if db.Contains("b", "1") {
		t.Error("should not contain fact in nonexistent relation")
	}
}
