package sdk

import (
	"testing"
)

func TestLowerTitle(t *testing.T) {

	testCases := []struct {
		name  string
		value string
		out   string
		err   string
	}{
		{
			name:  "Empty",
			value: "",
			out:   "",
			err:   "expected empty title to remain empty",
		},
		{
			name:  "All Lowercase",
			value: "lowercase",
			out:   "lowercase",
			err:   "expected no changes when input is all lowercase",
		},
		{
			name:  "All Uppercase",
			value: "UPPERCASE",
			out:   "uPPERCASE",
			err:   "expected first rune to be lower when input is all uppercase",
		},
		{
			name:  "Title Case",
			value: "Title Case",
			out:   "title Case",
			err:   "expected first rune to be lower when input is Title Case",
		},
	}

	for i := range testCases {

		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			out := LowerTitle(tc.value)
			if out != tc.out {
				t.Fatal(tc.err)
			}
		})
	}
}

func TestTitle(t *testing.T) {

	testCases := []struct {
		name  string
		value string
		out   string
		err   string
	}{
		{
			name:  "Empty",
			value: "",
			out:   "",
			err:   "expected empty title to remain empty",
		},
		{
			name:  "Lowercase",
			value: "lowercase",
			out:   "Lowercase",
			err:   "expected frist rune to be uppercase",
		},
		{
			name:  "CamelCase",
			value: "camelCase",
			out:   "CamelCase",
			err:   "expected only first rune to be converted to uppercase",
		},
		{
			name:  "Digit Prefix",
			value: "1title",
			out:   "1title",
			err:   "expected 1title to remain 1title",
		},
	}

	for i := range testCases {

		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			out := Title(tc.value)
			if out != tc.out {
				t.Fatal(tc.err)
			}
		})
	}
}
