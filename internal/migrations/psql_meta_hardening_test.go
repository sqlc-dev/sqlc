package migrations

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestRemovePsqlMetaCommands_Hardening(t *testing.T) {
	tests := []stripCase{
		{
			name: "BlockCommentsActAsWhitespaceForSettingTracking",
			in: `SET/*comment*/standard_conforming_strings=off;
INSERT INTO foo VALUES ('x\'
\connect still_literal
');
\connect should_go
`,
			want: `SET/*comment*/standard_conforming_strings=off;
INSERT INTO foo VALUES ('x\'
\connect still_literal
');

`,
		},
		{
			name: "LineCommentsActAsWhitespaceForSettingTracking",
			in: `SET--comment
standard_conforming_strings = off;
INSERT INTO foo VALUES ('x\'
\connect still_literal
');
\connect should_go
`,
			want: `SET--comment
standard_conforming_strings = off;
INSERT INTO foo VALUES ('x\'
\connect still_literal
');

`,
		},
		{
			name: "InlineDoubleBackslashSeparatorPropagatesLexerState",
			in: `\x \\ SELECT 'open
\connect still_literal
'
\connect should_go
`,
			want: ` SELECT 'open
\connect still_literal
'

`,
		},
		{
			name: "NoWhitespaceDoubleBackslashSeparatorIsNotSpecial",
			in:   "\\\\SELECT 1;\n",
			want: "\\\\SELECT 1;\n",
		},
		{
			name: "NoWhitespaceInlineDoubleBackslashSeparatorIsNotSpecial",
			in:   "\\echo hi\\\\SELECT 1;\n",
			want: "\\echo hi\\\\SELECT 1;\n",
		},
		{
			name: "IndentedInvalidSeparatorCandidatePreservesLeadingWhitespace",
			in:   "  \\echo hi\\\\SELECT 1;\n",
			want: "  \\echo hi\\\\SELECT 1;\n",
		},
		{
			name: "LeadingDoubleBackslashLineIsNotSpecial",
			in:   "\\\\\nSELECT 1;\n",
			want: "\\\\\nSELECT 1;\n",
		},
		{
			name: "LeadingDoubleBackslashWithWhitespaceTailIsNotSpecial",
			in:   "\\\\ SELECT 1;\n",
			want: "\\\\ SELECT 1;\n",
		},
		{
			name: "CopyFromStdinTerminatorCanBeFollowedByLaterMetaAndSQL",
			in: `\copy foo from stdin
1	alpha
\.
\connect other
SELECT 1;
`,
			want: `



SELECT 1;
`,
		},
		{
			name: "CopyFromStdinWhitespacePrefixedDotIsDataNotTerminator",
			in: `\copy foo from stdin
 \.
SELECT 1;
`,
			wantErr: true,
		},
		{
			name: "CopyFromStdinTerminatorMustBeOwnLine",
			in: `\copy foo from stdin
1	alpha
\. \echo hi
SELECT 1;
`,
			wantErr: true,
		},
		{
			name: "CopyFromStdinTerminatorCannotShareLineWithSeparatorAndSQL",
			in: `\copy foo from stdin
1	alpha
\. \\ SELECT 1;
`,
			wantErr: true,
		},
		{
			name: "DoubledBackslashesInMetaCommandArgumentsDoNotDisableStripping",
			in:   "\\include \\\\server\\share\\schema.sql\nSELECT 1;\n",
			want: "\nSELECT 1;\n",
		},
		{
			name: "QuotedDoubledBackslashesInMetaCommandArgumentsDoNotDisableStripping",
			in:   "\\echo 'foo\\\\bar'\nSELECT 1;\n",
			want: "\nSELECT 1;\n",
		},
		{
			name: "LineCommentAtEOFDoesNotAffectState",
			in:   "-- just a trailing comment with $tag$ and /* and '",
			want: "-- just a trailing comment with $tag$ and /* and '",
		},
	}

	runRemovePsqlMetaCommandCases(t, tests)
}

func TestPsqlMetaHelperHardening(t *testing.T) {
	t.Run("isUnsupportedConditionalMetaCommand", func(t *testing.T) {
		if !isUnsupportedConditionalMetaCommand(`\if`) {
			t.Fatalf(`expected \if to be unsupported`)
		}
		if !isUnsupportedConditionalMetaCommand(`\endif`) {
			t.Fatalf(`expected \endif to be unsupported`)
		}
		if isUnsupportedConditionalMetaCommand(`\connect`) {
			t.Fatalf(`expected \connect to remain supported`)
		}
	})

	t.Run("isSemanticMetaCommand", func(t *testing.T) {
		for _, token := range []string{`\c`, `\connect`, `\i`, `\include`, `\ir`, `\include_relative`, `\copy`, `\gexec`} {
			if !isSemanticMetaCommand(token) {
				t.Fatalf("expected %s to be semantic", token)
			}
		}
		for _, token := range []string{`\restrict`, `\set`} {
			if isSemanticMetaCommand(token) {
				t.Fatalf("expected %s not to be semantic", token)
			}
		}
	})

	t.Run("readsPsqlCopyData", func(t *testing.T) {
		if readsPsqlCopyData(``, 0, 0) {
			t.Fatalf(`empty input should not be treated as \copy data`)
		}
		if !readsPsqlCopyData(`\copy foo from stdin`, 0, len(`\copy foo from stdin`)) {
			t.Fatalf(`expected \copy ... from stdin to be recognized`)
		}
		if !readsPsqlCopyData("\\copy foo FROM\tSTDIN", 0, len("\\copy foo FROM\tSTDIN")) {
			t.Fatalf(`expected whitespace/case variants of \copy ... from stdin to be recognized`)
		}
		if readsPsqlCopyData(`\copy foo from '/tmp/data.csv'`, 0, len(`\copy foo from '/tmp/data.csv'`)) {
			t.Fatalf(`did not expect file-based \copy to be treated as stdin data`)
		}
	})

	t.Run("hasInvalidSeparatorCandidate", func(t *testing.T) {
		if hasInvalidSeparatorCandidate(`\echo hi \\ SELECT 1;`, 0, len(`\echo hi \\ SELECT 1;`)) {
			t.Fatalf(`did not expect valid whitespace-delimited separator to be treated as invalid`)
		}
		if !hasInvalidSeparatorCandidate(`\echo hi\\SELECT 1;`, 0, len(`\echo hi\\SELECT 1;`)) {
			t.Fatalf(`expected glued separator candidate to be recognized`)
		}
		if hasInvalidSeparatorCandidate(`\echo "hi""there" \\ SELECT 1;`, 0, len(`\echo "hi""there" \\ SELECT 1;`)) {
			t.Fatalf(`did not expect doubled-quote escape path with valid separator to be treated as invalid`)
		}
		if !hasInvalidSeparatorCandidate(`\echo "hi""there"\\SELECT 1;`, 0, len(`\echo "hi""there"\\SELECT 1;`)) {
			t.Fatalf(`expected doubled-quote exit path to still recognize glued separator`)
		}
		if hasInvalidSeparatorCandidate(`\echo 'it''s' \\ SELECT 1;`, 0, len(`\echo 'it''s' \\ SELECT 1;`)) {
			t.Fatalf(`did not expect doubled-single-quote escape path with valid separator to be treated as invalid`)
		}
		if hasInvalidSeparatorCandidate(`\include \\server\share\schema.sql`, 0, len(`\include \\server\share\schema.sql`)) {
			t.Fatalf(`did not expect UNC-style argument to be treated as an invalid separator candidate`)
		}
		if hasInvalidSeparatorCandidate(`\echo 'foo\\bar'`, 0, len(`\echo 'foo\\bar'`)) {
			t.Fatalf(`did not expect quoted doubled backslashes to be treated as an invalid separator candidate`)
		}
	})

	t.Run("parseStandardConformingStringsSetting", func(t *testing.T) {
		type testCase struct {
			name         string
			in           string
			wantOK       bool
			wantLocal    bool
			wantScopeDef bool
			wantValue    *bool
		}
		boolPtr := func(v bool) *bool { return &v }
		tests := []testCase{
			{name: "empty", in: "", wantOK: false},
			{name: "off", in: "SET standard_conforming_strings = off;", wantOK: true, wantValue: boolPtr(true)},
			{name: "local on", in: "set local standard_conforming_strings to on;", wantOK: true, wantLocal: true, wantValue: boolPtr(false)},
			{name: "quoted off", in: "set standard_conforming_strings = 'off';", wantOK: true, wantValue: boolPtr(true)},
			{name: "quoted local on", in: "set local standard_conforming_strings to 'on';", wantOK: true, wantLocal: true, wantValue: boolPtr(false)},
			{name: "boolean false", in: "set standard_conforming_strings = false;", wantOK: true, wantValue: boolPtr(true)},
			{name: "boolean true", in: "set session standard_conforming_strings = true;", wantOK: true, wantValue: boolPtr(false)},
			{name: "session on", in: "set session standard_conforming_strings = on;", wantOK: true, wantValue: boolPtr(false)},
			{name: "local default", in: "set local standard_conforming_strings = default;", wantOK: true, wantLocal: true},
			{name: "session default", in: "set standard_conforming_strings = default;", wantOK: true, wantScopeDef: true},
			{name: "from current", in: "set standard_conforming_strings from current;", wantOK: false},
			{name: "unknown value", in: "set standard_conforming_strings = maybe;", wantOK: false},
			{name: "missing operator", in: "set standard_conforming_strings off;", wantOK: false},
			{name: "truncated", in: "set standard_conforming_strings =", wantOK: false},
			{name: "truncated local", in: "set local standard_conforming_strings =", wantOK: false},
			{name: "other setting", in: "set other_setting = off;", wantOK: false},
			{name: "non set", in: "SELECT 1;", wantOK: false},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				got, ok := parseStandardConformingStringsSetting(tc.in)
				if ok != tc.wantOK {
					t.Fatalf("expected ok=%v, got ok=%v value=%+v", tc.wantOK, ok, got)
				}
				if !ok {
					return
				}
				if got.local != tc.wantLocal {
					t.Fatalf("expected local=%v, got %+v", tc.wantLocal, got)
				}
				if got.scopeDefault != tc.wantScopeDef {
					t.Fatalf("expected scopeDefault=%v, got %+v", tc.wantScopeDef, got)
				}
				switch {
				case tc.wantValue == nil && got.value != nil:
					t.Fatalf("expected nil value, got %+v", got)
				case tc.wantValue != nil && got.value == nil:
					t.Fatalf("expected value %v, got nil", *tc.wantValue)
				case tc.wantValue != nil && got.value != nil && *tc.wantValue != *got.value:
					t.Fatalf("expected value %v, got %+v", *tc.wantValue, got)
				}
			})
		}
	})

	t.Run("consumeStatementFragment", func(t *testing.T) {
		var statement strings.Builder
		standardConformingStringsOff := false
		warnedApproximateSessionSemantics := false

		consumeStatementFragment("SET standard_conforming_strings = off;", &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, true)
		if !standardConformingStringsOff {
			t.Fatalf("expected standard_conforming_strings to switch off")
		}
		if !warnedApproximateSessionSemantics {
			t.Fatalf("expected standard_conforming_strings tracking to emit an approximation warning")
		}
		if statement.Len() != 0 {
			t.Fatalf("expected statement buffer to reset after semicolon, got %q", statement.String())
		}

		consumeStatementFragment("SELECT 1\r", &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, true)
		if statement.String() != "SELECT 1 " {
			t.Fatalf("expected line break normalization in statement buffer, got %q", statement.String())
		}

		statement.Reset()
		consumeStatementFragment("SET standard_conforming_strings = on;", &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, true)
		if standardConformingStringsOff {
			t.Fatalf("expected standard_conforming_strings to switch back on")
		}

		statement.Reset()
		standardConformingStringsOff = false
		warnedApproximateSessionSemantics = false
		consumeStatementFragment("SET standard_conforming_strings = 'off", &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, true)
		consumeStatementFragment("'", &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, false)
		consumeStatementFragment(";", &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, true)
		if !standardConformingStringsOff {
			t.Fatalf("expected quoted standard_conforming_strings to switch off")
		}
		if !warnedApproximateSessionSemantics {
			t.Fatalf("expected quoted standard_conforming_strings tracking to warn")
		}

		statement.Reset()
		consumeStatementFragment("RESET ALL;", &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, true)
		if standardConformingStringsOff {
			t.Fatalf("expected RESET ALL to restore default standard_conforming_strings")
		}

		statement.Reset()
		standardConformingStringsOff = false
		warnedApproximateSessionSemantics = false
		applyStandardConformingStringsStatement("", &standardConformingStringsOff, &warnedApproximateSessionSemantics)
		if standardConformingStringsOff || warnedApproximateSessionSemantics {
			t.Fatalf("expected empty statement to be ignored")
		}

		statement.Reset()
		standardConformingStringsOff = false
		warnedApproximateSessionSemantics = false
		consumeStatementFragment("SET", &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, true)
		consumeStatementFragment(" ", &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, false)
		consumeStatementFragment("standard_conforming_strings=off;", &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, true)
		if !standardConformingStringsOff {
			t.Fatalf("expected comment-separated statement fragments to preserve whitespace boundaries")
		}

		statement.Reset()
		standardConformingStringsOff = false
		warnedApproximateSessionSemantics = false
		consumeStatementFragment("START TRANSACTION;", &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, true)
		if !warnedApproximateSessionSemantics {
			t.Fatalf("expected START TRANSACTION to trigger approximation warning")
		}

		statement.Reset()
		standardConformingStringsOff = false
		warnedApproximateSessionSemantics = false
		applyStandardConformingStringsStatement("START WORK;", &standardConformingStringsOff, &warnedApproximateSessionSemantics)
		if warnedApproximateSessionSemantics {
			t.Fatalf("expected unrelated START variant to remain ignored")
		}

		statement.Reset()
		standardConformingStringsOff = true
		warnedApproximateSessionSemantics = false
		consumeStatementFragment("RESET standard_conforming_strings;", &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, true)
		if standardConformingStringsOff {
			t.Fatalf("expected RESET standard_conforming_strings to restore default")
		}
		if !warnedApproximateSessionSemantics {
			t.Fatalf("expected RESET standard_conforming_strings to trigger approximation warning")
		}

		statement.Reset()
		standardConformingStringsOff = true
		warnedApproximateSessionSemantics = false
		applyStandardConformingStringsStatement("RESET search_path;", &standardConformingStringsOff, &warnedApproximateSessionSemantics)
		if !standardConformingStringsOff || warnedApproximateSessionSemantics {
			t.Fatalf("expected unrelated RESET to remain ignored")
		}

		statement.Reset()
		standardConformingStringsOff = false
		warnedApproximateSessionSemantics = false
		consumeStatementFragment("\n", &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, false)
		if statement.String() != " " {
			t.Fatalf("expected newline to normalize to space, got %q", statement.String())
		}
	})

	t.Run("isEscapeStringPrefix", func(t *testing.T) {
		if isEscapeStringPrefix("", 0) {
			t.Fatalf("empty input should not report an escape-string prefix")
		}
		if !isEscapeStringPrefix("E'", 1) {
			t.Fatalf("single-rune E prefix should be recognized")
		}
		if isEscapeStringPrefix("x'", 1) {
			t.Fatalf("non-E prefix should not be recognized")
		}
	})

	t.Run("matchDollarTagStart", func(t *testing.T) {
		if got := matchDollarTagStart("plain", 0); got != "" {
			t.Fatalf("non-dollar input should not match, got %q", got)
		}
		if got := matchDollarTagStart("x$foo$", 1); got != "" {
			t.Fatalf("identifier-adjacent dollar tag should not match, got %q", got)
		}
		if got := matchDollarTagStart("$", 0); got != "" {
			t.Fatalf("truncated dollar tag should not match, got %q", got)
		}
		if got := matchDollarTagStart("$foo", 0); got != "" {
			t.Fatalf("unterminated dollar tag should not match, got %q", got)
		}
		if got := matchDollarTagStart("$1$", 0); got != "" {
			t.Fatalf("digit-start tag should not match, got %q", got)
		}
		if got := matchDollarTagStart("$a-$", 0); got != "" {
			t.Fatalf("invalid rune inside tag should not match, got %q", got)
		}
		if got := matchDollarTagStart("$$", 0); got != "$$" {
			t.Fatalf("empty tag should match $$, got %q", got)
		}
		if got := matchDollarTagStart("$ä$", 0); got != "$ä$" {
			t.Fatalf("unicode tag should match, got %q", got)
		}
		if got := matchDollarTagStart(string([]byte{'$', 0xff, '$'}), 0); got != "" {
			t.Fatalf("invalid utf-8 after $ should not match, got %q", got)
		}
		if got := matchDollarTagStart("$a\xff$", 0); got != "" {
			t.Fatalf("invalid utf-8 inside tag should not match, got %q", got)
		}
	})

	t.Run("lastRuneBefore", func(t *testing.T) {
		if got := lastRuneBefore("", 0); got != utf8.RuneError {
			t.Fatalf("expected RuneError for missing previous rune, got %q", got)
		}
		if got := lastRuneBefore("äb", len("ä")); got != 'ä' {
			t.Fatalf("expected previous rune ä, got %q", got)
		}
	})

	t.Run("writeLineEnding", func(t *testing.T) {
		var out strings.Builder
		if next := writeLineEnding(&out, "", 0); next != 0 || out.String() != "" {
			t.Fatalf("expected no-op at end of input, next=%d out=%q", next, out.String())
		}

		out.Reset()
		if next := writeLineEnding(&out, "x", 0); next != 0 || out.String() != "" {
			t.Fatalf("expected non-line-ending input to be untouched, next=%d out=%q", next, out.String())
		}

		out.Reset()
		if next := writeLineEnding(&out, "\r", 0); next != 1 || out.String() != "\n" {
			t.Fatalf("expected bare CR normalization, next=%d out=%q", next, out.String())
		}
	})
}
