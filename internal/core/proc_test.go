package core

import "testing"

func TestProcCreateAndFind(t *testing.T) {
	cat, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer cat.Close()

	intOID, err := cat.CreateType("integer", 4)
	if err != nil {
		t.Fatal(err)
	}
	textOID, err := cat.CreateType("text", 0)
	if err != nil {
		t.Fatal(err)
	}

	// length(text) -> integer
	procOID, err := cat.CreateProc(ProcSpec{
		Name:           "length",
		Kind:           "f",
		ReturnTypeOID:  intOID,
		ReturnNullable: true,
		Args:           []ProcArg{{Name: "s", TypeOID: textOID}},
	})
	if err != nil {
		t.Fatalf("create length: %v", err)
	}
	if procOID == 0 {
		t.Fatal("expected non-zero proc oid")
	}

	// concat(text, text) -> text
	if _, err := cat.CreateProc(ProcSpec{
		Name:          "concat",
		ReturnTypeOID: textOID,
		Args: []ProcArg{
			{Name: "a", TypeOID: textOID},
			{Name: "b", TypeOID: textOID},
		},
	}); err != nil {
		t.Fatal(err)
	}

	got, err := cat.FindProcs("length", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("FindProcs length: want 1, got %d", len(got))
	}
	if got[0].ReturnTypeOID != intOID {
		t.Errorf("length return type: got %d, want %d", got[0].ReturnTypeOID, intOID)
	}
	if len(got[0].ArgTypes) != 1 || got[0].ArgTypes[0] != textOID {
		t.Errorf("length args: got %v, want [%d]", got[0].ArgTypes, textOID)
	}

	got, err = cat.FindProcs("concat", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || len(got[0].ArgTypes) != 2 {
		t.Fatalf("FindProcs concat: got %+v", got)
	}
}
