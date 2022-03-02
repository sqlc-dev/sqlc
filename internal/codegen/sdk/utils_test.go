package sdk

import (
	"testing"
)

func TestLowerTitle(t *testing.T) {

	// empty
	if LowerTitle("") != "" {
		t.Fatal("expected empty title to remain empty")
	}

	// all lowercase
	if LowerTitle("lowercase") != "lowercase" {
		t.Fatal("expected no changes when input is all lowercase")
	}

	// all uppercase
	if LowerTitle("UPPERCASE") != "uPPERCASE" {
		t.Fatal("expected first rune to be lower when input is all uppercase")
	}

	// Title Case
	if LowerTitle("Title Case") != "title Case" {
		t.Fatal("expected first rune to be lower when input is Title Case")
	}
}
