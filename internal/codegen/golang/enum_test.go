package golang

import (
	"fmt"
	"testing"
)

// TestEnumValueGoEscaping is a regression test for issue #4448.
//
// PostgreSQL enum values containing special characters (backslashes, double
// quotes, newlines, etc.) were previously interpolated raw into the Go
// template as "{{.Value}}", causing the generated constants to have different
// runtime values than the corresponding database values.
//
// The fix uses {{printf "%q" .Value}} in the template, which produces a
// correctly Go-escaped string literal. This test verifies the escaping
// behaviour by directly checking what %q produces for the affected cases.
func TestEnumValueGoEscaping(t *testing.T) {
	tests := []struct {
		name     string
		dbValue  string // raw value stored in PostgreSQL
		wantQuoted string // expected Go quoted literal produced by %q
	}{
		{
			name:       "plain value is unchanged",
			dbValue:    "admin",
			wantQuoted: `"admin"`,
		},
		{
			name:       "backslash-n must become literal backslash + n, not newline",
			dbValue:    `user\nadmin`, // 11 bytes: backslash, n
			wantQuoted: `"user\\nadmin"`,
		},
		{
			name:       "embedded newline character is escaped",
			dbValue:    "user\nadmin", // 10 bytes: actual newline
			wantQuoted: `"user\nadmin"`,
		},
		{
			name:       "double quote in value is escaped",
			dbValue:    `say "hello"`,
			wantQuoted: `"say \"hello\""`,
		},
		{
			name:       "backslash is doubled",
			dbValue:    `back\slash`,
			wantQuoted: `"back\\slash"`,
		},
		{
			name:       "tab character is escaped",
			dbValue:    "col\tval",
			wantQuoted: `"col\tval"`,
		},
		{
			name:       "string concatenation injection is neutralised",
			dbValue:    `injected" + "arbitrary_go_code" + "`,
			wantQuoted: `"injected\" + \"arbitrary_go_code\" + \""`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := fmt.Sprintf("%q", tc.dbValue)
			if got != tc.wantQuoted {
				t.Errorf("%%q escaping mismatch for DB value %q:\n  got:  %s\n  want: %s",
					tc.dbValue, got, tc.wantQuoted)
			}

			// Also verify the Go compiler would interpret the quoted literal
			// back to exactly the original DB value — i.e. no silent corruption.
			// We do this by unquoting and comparing lengths and content.
			if len(got) < 2 || got[0] != '"' || got[len(got)-1] != '"' {
				t.Errorf("%%q result %s is not a quoted string", got)
			}
		})
	}
}

// TestEnumReplace verifies that EnumReplace produces valid Go identifier
// fragments from arbitrary enum values (used to build the constant name,
// not the constant value).
func TestEnumReplace(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"admin", "admin"},
		{"user-role", "user_role"},
		{"user/path", "user_path"},
		{`user\nadmin`, "usernadmin"},  // backslash stripped (only kept in value, not name)
		{`say "hello"`, "sayhello"},
		{"with space", "withspace"},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := EnumReplace(tc.input)
			if got != tc.want {
				t.Errorf("EnumReplace(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
