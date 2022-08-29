package source

import (
	"fmt"
	"testing"
)

// newEdit is a testing helper for quickly generating Edits
func newEdit(loc int, old, new string) Edit {
	return Edit{Location: loc, Old: old, New: new}
}

// TestMutateSingle tests almost every possibility of a single edit
func TestMutateSingle(t *testing.T) {
	type test struct {
		input    string
		edit     Edit
		expected string
	}

	tests := []test{
		// Simple edits that replace everything
		{"", newEdit(0, "", ""), ""},
		{"a", newEdit(0, "a", "A"), "A"},
		{"abcde", newEdit(0, "abcde", "fghij"), "fghij"},
		{"", newEdit(0, "", "fghij"), "fghij"},
		{"abcde", newEdit(0, "abcde", ""), ""},

		// Edits that start at the very beginning (But don't cover the whole range)
		{"abcde", newEdit(0, "a", "A"), "Abcde"},
		{"abcde", newEdit(0, "ab", "AB"), "ABcde"},
		{"abcde", newEdit(0, "abc", "ABC"), "ABCde"},
		{"abcde", newEdit(0, "abcd", "ABCD"), "ABCDe"},

		// The above repeated, but with different lengths
		{"abcde", newEdit(0, "a", ""), "bcde"},
		{"abcde", newEdit(0, "ab", "A"), "Acde"},
		{"abcde", newEdit(0, "abc", "AB"), "ABde"},
		{"abcde", newEdit(0, "abcd", "AB"), "ABe"},

		// Edits that touch the end (but don't cover the whole range)
		{"abcde", newEdit(4, "e", "E"), "abcdE"},
		{"abcde", newEdit(3, "de", "DE"), "abcDE"},
		{"abcde", newEdit(2, "cde", "CDE"), "abCDE"},
		{"abcde", newEdit(1, "bcde", "BCDE"), "aBCDE"},

		// The above repeated, but with different lengths
		{"abcde", newEdit(4, "e", ""), "abcd"},
		{"abcde", newEdit(3, "de", "D"), "abcD"},
		{"abcde", newEdit(2, "cde", "CD"), "abCD"},
		{"abcde", newEdit(1, "bcde", "BC"), "aBC"},

		// Raw insertions / deletions
		{"abcde", newEdit(0, "", "_"), "_abcde"},
		{"abcde", newEdit(1, "", "_"), "a_bcde"},
		{"abcde", newEdit(2, "", "_"), "ab_cde"},
		{"abcde", newEdit(3, "", "_"), "abc_de"},
		{"abcde", newEdit(4, "", "_"), "abcd_e"},
		{"abcde", newEdit(5, "", "_"), "abcde_"},
	}

	origTests := tests
	// Generate the reverse mutations, for every edit - the opposite edit that makes it "undo"
	for _, spec := range origTests {
		tests = append(tests, test{
			input:    spec.expected,
			edit:     newEdit(spec.edit.Location, spec.edit.New, spec.edit.Old),
			expected: spec.input,
		})
	}

	for _, spec := range tests {
		expected := spec.expected

		actual, err := Mutate(spec.input, []Edit{spec.edit})
		testName := fmt.Sprintf("Mutate(%s, Edit{%v, %v -> %v})", spec.input, spec.edit.Location, spec.edit.Old, spec.edit.New)
		if err != nil {
			t.Errorf("%s should not error (%v)", testName, err)
			continue
		}

		if actual != expected {
			t.Errorf("%s expected %v; got %v", testName, expected, actual)
		}
	}
}

// TestMutateMulti tests combinations of edits
func TestMutateMulti(t *testing.T) {
	type test struct {
		input    string
		edit1    Edit
		edit2    Edit
		expected string
	}

	tests := []test{
		// Edits that are >1 character from each other
		{"abcde", newEdit(0, "a", "A"), newEdit(2, "c", "C"), "AbCde"},
		{"abcde", newEdit(0, "a", "A"), newEdit(2, "c", "C"), "AbCde"},

		// 2 edits bump right up next to each other
		{"abcde", newEdit(0, "abc", ""), newEdit(3, "de", "DE"), "DE"},
		{"abcde", newEdit(0, "abc", "ABC"), newEdit(3, "de", ""), "ABC"},
		{"abcde", newEdit(0, "abc", "ABC"), newEdit(3, "de", "DE"), "ABCDE"},
		{"abcde", newEdit(1, "b", "BB"), newEdit(2, "c", "CC"), "aBBCCde"},

		// 2 edits bump next to each other, but don't cover the whole string
		{"abcdef", newEdit(1, "bc", "C"), newEdit(3, "de", "D"), "aCDf"},
		{"abcde", newEdit(1, "bc", "CCCC"), newEdit(3, "d", "DDD"), "aCCCCDDDe"},

		// lengthening edits
		{"abcde", newEdit(1, "b", "BBBB"), newEdit(2, "c", "CCCC"), "aBBBBCCCCde"},
	}

	origTests := tests
	// Generate the edits in opposite order mutations, source edits should be independent of
	// the order the edits are specified
	for _, spec := range origTests {
		tests = append(tests, test{
			input:    spec.input,
			edit1:    spec.edit2,
			edit2:    spec.edit1,
			expected: spec.expected,
		})
	}

	for _, spec := range tests {
		expected := spec.expected

		actual, err := Mutate(spec.input, []Edit{spec.edit1, spec.edit2})
		testName := fmt.Sprintf("Mutate(%s, Edits{(%v, %v -> %v), (%v, %v -> %v)})", spec.input,
			spec.edit1.Location, spec.edit1.Old, spec.edit1.New,
			spec.edit2.Location, spec.edit2.Old, spec.edit2.New)

		if err != nil {
			t.Errorf("%s should not error (%v)", testName, err)
			continue
		}

		if actual != expected {
			t.Errorf("%s expected %v; got %v", testName, expected, actual)
		}
	}
}

// TestMutateErrorSingle test errors are generated for trivially incorrect single edits
func TestMutateErrorSingle(t *testing.T) {
	type test struct {
		input string
		edit  Edit
	}

	tests := []test{
		// old text is longer than input text
		{"", newEdit(0, "a", "A")},
		{"a", newEdit(0, "aa", "A")},
		{"hello", newEdit(0, "hello!", "A")},

		// negative indexes
		{"aaa", newEdit(-1, "aa", "A")},
		{"aaa", newEdit(-2, "aa", "A")},
		{"aaa", newEdit(-100, "aa", "A")},
	}

	for _, spec := range tests {
		edit := spec.edit

		_, err := Mutate(spec.input, []Edit{edit})
		testName := fmt.Sprintf("Mutate(%s, Edit{%v, %v -> %v})", spec.input, edit.Location, edit.Old, edit.New)
		if err == nil {
			t.Errorf("%s should error (%v)", testName, err)
			continue
		}
	}
}

// TestMutateErrorMulti tests error that can only happen across multiple errors
func TestMutateErrorMulti(t *testing.T) {
	type test struct {
		input string
		edit1 Edit
		edit2 Edit
	}

	tests := []test{
		// These edits overlap each other, and are therefore undefined
		{"abcdef", newEdit(0, "a", ""), newEdit(0, "a", "A")},
		{"abcdef", newEdit(0, "ab", ""), newEdit(1, "ab", "AB")},
		{"abcdef", newEdit(0, "abc", ""), newEdit(2, "abc", "ABC")},

		// the last edit is longer than the string itself
		{"abcdef", newEdit(0, "abcdefghi", ""), newEdit(2, "abc", "ABC")},

		// negative indexes
		{"abcdef", newEdit(-1, "abc", ""), newEdit(3, "abc", "ABC")},
		{"abcdef", newEdit(0, "abc", ""), newEdit(-1, "abc", "ABC")},
	}

	for _, spec := range tests {
		actual, err := Mutate(spec.input, []Edit{spec.edit1, spec.edit2})
		testName := fmt.Sprintf("Mutate(%s, Edits{(%v, %v -> %v), (%v, %v -> %v)})", spec.input,
			spec.edit1.Location, spec.edit1.Old, spec.edit1.New,
			spec.edit2.Location, spec.edit2.Old, spec.edit2.New)

		if err == nil {
			t.Errorf("%s should error, but got (%v)", testName, actual)
		}
	}
}
