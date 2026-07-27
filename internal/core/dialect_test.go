package core

import "testing"

func TestDialectAndFlags(t *testing.T) {
	cat, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer cat.Close()

	pgOID, err := cat.CreateDialect("postgresql")
	if err != nil {
		t.Fatal(err)
	}

	got, err := cat.DialectOID("postgresql")
	if err != nil || got != pgOID {
		t.Fatalf("DialectOID: got %d (err=%v), want %d", got, err, pgOID)
	}

	if err := cat.SetDialectFlag(pgOID, "identifier_case", "fold_lower"); err != nil {
		t.Fatal(err)
	}
	v, err := cat.DialectFlag(pgOID, "identifier_case")
	if err != nil || v != "fold_lower" {
		t.Errorf("DialectFlag: got %q (err=%v), want fold_lower", v, err)
	}

	if err := cat.SetDialectFlag(pgOID, "identifier_case", "fold_upper"); err != nil {
		t.Fatal(err)
	}
	v, _ = cat.DialectFlag(pgOID, "identifier_case")
	if v != "fold_upper" {
		t.Errorf("upsert: got %q, want fold_upper", v)
	}

	v, err = cat.DialectFlag(pgOID, "nonexistent")
	if err != nil || v != "" {
		t.Errorf("missing flag: got %q (err=%v), want empty", v, err)
	}
}
