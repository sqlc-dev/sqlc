package core

import "testing"

func TestCastAllowed(t *testing.T) {
	cat, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer cat.Close()

	intOID, _ := cat.CreateType("integer", 4)
	bigintOID, _ := cat.CreateType("bigint", 8)
	textOID, _ := cat.CreateType("text", 0)

	// integer -> bigint is implicit
	if err := cat.CreateCast(CastSpec{
		SourceTypeOID: intOID, TargetTypeOID: bigintOID, Context: "i",
	}); err != nil {
		t.Fatal(err)
	}
	// integer -> text is explicit only
	if err := cat.CreateCast(CastSpec{
		SourceTypeOID: intOID, TargetTypeOID: textOID, Context: "e",
	}); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		src, tgt int64
		ctx      string
		want     bool
	}{
		{intOID, intOID, "i", true},      // identity always OK
		{intOID, bigintOID, "i", true},   // declared implicit
		{intOID, bigintOID, "a", true},   // implicit subsumes assignment
		{intOID, bigintOID, "e", true},   // implicit subsumes explicit
		{intOID, textOID, "i", false},    // explicit-only blocked from implicit
		{intOID, textOID, "a", false},    // and from assignment
		{intOID, textOID, "e", true},     // but works for explicit
		{textOID, intOID, "e", false},    // no rule defined
	}
	for _, c := range cases {
		got, err := cat.CastAllowed(c.src, c.tgt, c.ctx)
		if err != nil {
			t.Errorf("CastAllowed(%d,%d,%q): %v", c.src, c.tgt, c.ctx, err)
			continue
		}
		if got != c.want {
			t.Errorf("CastAllowed(%d,%d,%q): got %v, want %v", c.src, c.tgt, c.ctx, got, c.want)
		}
	}
}
