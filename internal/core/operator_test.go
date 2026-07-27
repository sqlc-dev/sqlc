package core

import "testing"

func TestOperatorCreateAndFind(t *testing.T) {
	cat, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer cat.Close()

	intOID, _ := cat.CreateType("integer", 4)
	boolOID, _ := cat.CreateType("boolean", 1)

	gtOID, err := cat.CreateOperator(OperatorSpec{
		Name:          ">",
		LeftTypeOID:   intOID,
		RightTypeOID:  intOID,
		ResultTypeOID: boolOID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if gtOID == 0 {
		t.Fatal("expected non-zero operator oid")
	}

	got, err := cat.FindOperators(">", intOID, intOID)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("FindOperators: want 1, got %d", len(got))
	}
	if got[0].ResultTypeOID != boolOID {
		t.Errorf("result type: got %d, want %d", got[0].ResultTypeOID, boolOID)
	}

	all, err := cat.FindOperators(">", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 1 {
		t.Errorf("listing >: want 1, got %d", len(all))
	}
}
