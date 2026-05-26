package compiler

import "testing"

func TestCamelize(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want string
	}{
		{"name_asc", "NameAsc"},
		{"recent", "Recent"},
		{"else", "Else"},
		{"created-at-desc", "CreatedAtDesc"},
		{"two words", "TwoWords"},
		{"", ""},
	} {
		if got := camelize(tc.in); got != tc.want {
			t.Errorf("camelize(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestMatchParen(t *testing.T) {
	for _, tc := range []struct {
		name    string
		in      string
		start   int
		want    int
		wantErr bool
	}{
		{"simple", "f(x)", 0, 3, false},
		{"nested", "f(g(x), h(y))", 0, 12, false},
		{"string with parens", "f('a)b', x)", 0, 10, false},
		{"escaped quote in string", "f('a''b)', x)", 0, 12, false},
		{"fragment with call", "case('coalesce(x, 0) ASC')", 0, 25, false},
		{"unbalanced", "f(x", 0, 0, true},
		{"no open paren", "abc", 0, 0, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := matchParen(tc.in, tc.start)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("matchParen(%q) expected error, got %d", tc.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("matchParen(%q) unexpected error: %v", tc.in, err)
			}
			if got != tc.want {
				t.Errorf("matchParen(%q) = %d, want %d", tc.in, got, tc.want)
			}
			if tc.in[got] != ')' {
				t.Errorf("matchParen(%q) index %d is %q, not ')'", tc.in, got, tc.in[got])
			}
		})
	}
}
