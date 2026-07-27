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

	if err := cat.CreateCast(CastSpec{
		SourceTypeOID: intOID, TargetTypeOID: bigintOID, Context: "i",
	}); err != nil {
		t.Fatal(err)
	}
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
		{intOID, intOID, "i", true},
		{intOID, bigintOID, "i", true},
		{intOID, bigintOID, "a", true},
		{intOID, bigintOID, "e", true},
		{intOID, textOID, "i", false},
		{intOID, textOID, "a", false},
		{intOID, textOID, "e", true},
		{textOID, intOID, "e", false},
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
