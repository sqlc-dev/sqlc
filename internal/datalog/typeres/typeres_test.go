package typeres

import (
	"testing"
)

func TestDirectTypePropagation(t *testing.T) {
	r := New()
	LoadPostgreSQLTypes(r)

	r.SetExprType("e1", "int4")
	r.SetExprType("e2", "text")

	res, err := r.Resolve()
	if err != nil {
		t.Fatal(err)
	}

	if typ, ok := res.TypeOf("e1"); !ok || typ != "int4" {
		t.Errorf("expected e1=int4, got %q (ok=%v)", typ, ok)
	}
	if typ, ok := res.TypeOf("e2"); !ok || typ != "text" {
		t.Errorf("expected e2=text, got %q (ok=%v)", typ, ok)
	}
}

func TestImplicitCastTransitivity(t *testing.T) {
	r := New()
	LoadPostgreSQLTypes(r)

	res, err := r.Resolve()
	if err != nil {
		t.Fatal(err)
	}

	// int2 -> int4 -> int8 should be transitively castable
	if !res.CanCast("int2", "int8") {
		t.Error("expected int2 castable to int8")
	}
	if !res.CanCast("int4", "float8") {
		t.Error("expected int4 castable to float8")
	}
	// text should not be castable to int4 implicitly
	if res.CanCast("text", "int4") {
		t.Error("unexpected: text castable to int4")
	}
}

func TestBinaryExprResolution(t *testing.T) {
	r := New()
	LoadPostgreSQLTypes(r)

	// e1:int4 + e2:int4 = e3:int4
	r.SetExprType("e1", "int4")
	r.SetExprType("e2", "int4")
	r.AddBinaryExpr("e3", "+", "e1", "e2")

	res, err := r.Resolve()
	if err != nil {
		t.Fatal(err)
	}

	if typ, ok := res.TypeOf("e3"); !ok || typ != "int4" {
		t.Errorf("expected e3=int4, got %q (ok=%v)", typ, ok)
	}
}

func TestBinaryExprWithPromotion(t *testing.T) {
	r := New()
	LoadPostgreSQLTypes(r)

	// e1:int4 + e2:int8 should promote to int8
	r.SetExprType("e1", "int4")
	r.SetExprType("e2", "int8")
	r.AddBinaryExpr("e3", "+", "e1", "e2")

	res, err := r.Resolve()
	if err != nil {
		t.Fatal(err)
	}

	typ, ok := res.TypeOf("e3")
	if !ok {
		t.Fatal("e3 type not resolved")
	}
	// Should resolve to int8 (int4 promoted to int8)
	if typ != "int8" {
		t.Errorf("expected e3=int8, got %q", typ)
	}
}

func TestCastExpr(t *testing.T) {
	r := New()
	LoadPostgreSQLTypes(r)

	r.SetExprType("e1", "int4")
	r.AddCastExpr("e2", "e1", "text")

	res, err := r.Resolve()
	if err != nil {
		t.Fatal(err)
	}

	if typ, ok := res.TypeOf("e2"); !ok || typ != "text" {
		t.Errorf("expected e2=text, got %q (ok=%v)", typ, ok)
	}
}

func TestFunctionCall(t *testing.T) {
	r := New()
	LoadPostgreSQLTypes(r)

	// length('hello') -> int4
	r.SetExprType("e1", "text")
	r.AddFuncCall("e2", "length", 0, "e1")

	res, err := r.Resolve()
	if err != nil {
		t.Fatal(err)
	}

	if typ, ok := res.TypeOf("e2"); !ok || typ != "int4" {
		t.Errorf("expected e2=int4, got %q (ok=%v)", typ, ok)
	}
}

func TestCommonType(t *testing.T) {
	r := New()
	LoadPostgreSQLTypes(r)

	res, err := r.Resolve()
	if err != nil {
		t.Fatal(err)
	}

	// int4 and int8 should have a common type
	ct, ok := res.CommonType("int4", "int8")
	if !ok {
		t.Fatal("expected common type for int4, int8")
	}
	// int4 is castable to int8, so common type should be int8
	if ct != "int8" {
		t.Errorf("expected common type int8, got %q", ct)
	}
}

func TestCoalesceExpr(t *testing.T) {
	r := New()
	LoadPostgreSQLTypes(r)

	// COALESCE(e1:int4, e2:int4) -> should resolve to int4
	r.SetExprType("e1", "int4")
	r.SetExprType("e2", "int4")
	r.AddCoalesceExpr("e3", 0, "e1")
	r.AddCoalesceExpr("e3", 1, "e2")

	res, err := r.Resolve()
	if err != nil {
		t.Fatal(err)
	}

	typ, ok := res.TypeOf("e3")
	if !ok {
		t.Fatal("e3 type not resolved")
	}
	if typ != "int4" {
		t.Errorf("expected e3=int4, got %q", typ)
	}
}

func TestAllTypes(t *testing.T) {
	r := New()
	LoadPostgreSQLTypes(r)

	r.SetExprType("e1", "int4")
	r.SetExprType("e2", "text")

	res, err := r.Resolve()
	if err != nil {
		t.Fatal(err)
	}

	all := res.AllTypes()
	if len(all) < 2 {
		t.Errorf("expected at least 2 resolved types, got %d", len(all))
	}
	if all["e1"] != "int4" {
		t.Errorf("expected e1=int4 in AllTypes, got %q", all["e1"])
	}
	if all["e2"] != "text" {
		t.Errorf("expected e2=text in AllTypes, got %q", all["e2"])
	}
}
