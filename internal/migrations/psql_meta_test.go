package migrations

import (
	"fmt"
	"strings"
	"testing"
)

// allPsqlMetaCommands enumerates the documented psql backslash commands from
// the current PostgreSQL `app-psql` reference so the stripper test covers the
// full command vocabulary. That list includes `\restrict` and `\unrestrict`,
// which became relevant here after the CVE-2025-8714 backpatches in PostgreSQL
// 17.6 / 16.10 / 15.14 / 14.19 / 13.22 taught pg_dump/plain-text restore flows
// to emit them. This list is a reference fixture for the known documented
// command surface; the stripping policy itself is intentionally broader and
// also removes unknown top-level backslash directives. The special `\\`
// separator is covered separately because it is only valid after another
// meta-command on the same line, not as a standalone line-start token.
var allPsqlMetaCommands = []string{
	`\a`, `\bind`, `\bind_named`, `\c`, `\connect`, `\C`, `\cd`, `\close_prepared`, `\conninfo`, `\copy`,
	`\copyright`, `\crosstabview`, `\d`, `\da`, `\dA`, `\dAc`, `\dAf`, `\dAo`, `\dAp`, `\db`,
	`\dc`, `\dconfig`, `\dC`, `\dd`, `\dD`, `\ddp`, `\dE`, `\di`, `\dm`, `\ds`,
	`\dt`, `\dv`, `\des`, `\det`, `\deu`, `\dew`, `\df`, `\dF`, `\dFd`, `\dFp`,
	`\dFt`, `\dg`, `\dl`, `\dL`, `\dn`, `\do`, `\dO`, `\dp`, `\dP`, `\drds`,
	`\drg`, `\dRp`, `\dRs`, `\dT`, `\du`, `\dx`, `\dX`, `\dy`, `\e`, `\edit`,
	`\echo`, `\ef`, `\encoding`, `\ev`, `\f`, `\g`, `\gdesc`, `\getenv`, `\gexec`, `\gset`,
	`\gx`, `\h`, `\help`, `\H`, `\html`, `\i`, `\include`, `\if`, `\elif`, `\else`,
	`\endif`, `\ir`, `\include_relative`, `\l`, `\list`, `\lo_export`, `\lo_import`, `\lo_list`, `\lo_unlink`, `\o`,
	`\out`, `\p`, `\print`, `\parse`, `\password`, `\prompt`, `\pset`, `\q`, `\quit`, `\qecho`,
	`\r`, `\reset`, `\restrict`, `\s`, `\set`, `\setenv`, `\sf`, `\sv`, `\startpipeline`, `\sendpipeline`,
	`\syncpipeline`, `\endpipeline`, `\flushrequest`, `\flush`, `\getresults`, `\t`, `\T`, `\timing`, `\unrestrict`, `\unset`,
	`\w`, `\write`, `\warn`, `\watch`, `\x`, `\z`, `\!`, `\?`, `\;`,
}

type stripCase struct {
	name    string
	in      string
	want    string
	wantErr bool
}

func runRemovePsqlMetaCommandCases(t *testing.T, tests []stripCase) {
	t.Helper()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, _, err := RemovePsqlMetaCommands(tc.in)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got output %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error after stripping meta commands: %v", err)
			}
			if got != tc.want {
				t.Fatalf("unexpected output after stripping meta commands:\nwant=%q\ngot =%q", tc.want, got)
			}
		})
	}
}

func TestRemovePsqlMetaCommands_TableDriven(t *testing.T) {
	inDoubleQuoted := "CREATE TABLE \"foo\\bar\" (id int);\nSELECT \"foo\\bar\"." +
		"id FROM \"foo\\bar\";\n"
	inValidSQL := "CREATE TABLE t (id int);\nINSERT INTO t VALUES (1);\n"
	inWhitespaceOnly := "   \t  "
	inNoTrailingNewline := "SELECT 1"
	inBackslashNotAtStart := "SELECT '\\not_meta' AS col;\n  SELECT '\\still_not_meta';\n"
	inDoubleSingleQuotes := "INSERT INTO t VALUES ('It''s fine');\n"

	tests := []stripCase{
		{
			name: "RemovesTopLevelMetaCommands",
			in:   "CREATE TABLE public.authors();\n\\connect test\n  \\set ON_ERROR_STOP on\nSELECT 1;\n",
			want: "CREATE TABLE public.authors();\n\n\nSELECT 1;\n",
		},
		{
			name: "IgnoresBackslashesInStrings",
			in: `SELECT E'\n' || E'\\\\' || '
\restrict inside';
SELECT E'
\still_string
';
\connect nope
`,
			want: `SELECT E'\n' || E'\\\\' || '
\restrict inside';
SELECT E'
\still_string
';

`,
		},
		{
			name: "PreservesDollarQuotedBlocks",
			in:   "DO $$\n\\this_should_stay\n$$;\n\\connect other\n",
			want: "DO $$\n\\this_should_stay\n$$;\n\n",
		},
		{
			name: "IgnoresBlockComments",
			in:   "/*\n\\comment_not_meta\n*/\n\\set x 1\nSELECT 1;\n",
			want: "/*\n\\comment_not_meta\n*/\n\nSELECT 1;\n",
		},
		{
			name: "LeavesValidSqlUntouched",
			in:   inValidSQL,
			want: inValidSQL,
		},
		{
			name: "HandlesEmptyInput",
			in:   "",
			want: "",
		},
		{
			name: "PreservesWhitespaceOnlyInput",
			in:   inWhitespaceOnly,
			want: inWhitespaceOnly,
		},
		{
			name: "PreservesFinalLineWithoutNewline",
			in:   inNoTrailingNewline,
			want: inNoTrailingNewline,
		},
		{
			name: "BackslashInDoubleQuotedIdentifier",
			in:   inDoubleQuoted,
			want: inDoubleQuoted,
		},
		{
			name: "BackslashNotAtLineStart",
			in:   inBackslashNotAtStart,
			want: inBackslashNotAtStart,
		},
		{
			name: "DoubleSingleQuotesRemain",
			in:   inDoubleSingleQuotes,
			want: inDoubleSingleQuotes,
		},
		{
			name: "MetaCommandTextInsideLiteral",
			in: `INSERT INTO logs VALUES ('Remember to run \connect later');
	SELECT E'\n\connect\n' as literal;` + "\n",
			want: `INSERT INTO logs VALUES ('Remember to run \connect later');
	SELECT E'\n\connect\n' as literal;` + "\n",
		},
		{
			name: "EscapeStringsPreserveBackslashEscapedQuotes",
			in: `SELECT E'line1\'
\connect still_literal
';
\connect should_go
`,
			want: `SELECT E'line1\'
\connect still_literal
';

`,
		},
		{
			name: "StandardConformingStringsOffPreservesBackslashEscapes",
			in: `SET standard_conforming_strings = off;
INSERT INTO foo VALUES ('x\'
\connect still_literal
');
\connect should_go
`,
			want: `SET standard_conforming_strings = off;
INSERT INTO foo VALUES ('x\'
\connect still_literal
');

`,
		},
		{
			name: "QuotedStandardConformingStringsOffPreservesBackslashEscapes",
			in: `SET standard_conforming_strings = 'off';
INSERT INTO foo VALUES ('x\'
\connect still_literal
');
\connect should_go
`,
			want: `SET standard_conforming_strings = 'off';
INSERT INTO foo VALUES ('x\'
\connect still_literal
');

`,
		},
		{
			name: "BooleanStandardConformingStringsFalsePreservesBackslashEscapes",
			in: `SET standard_conforming_strings = false;
INSERT INTO foo VALUES ('x\'
\connect still_literal
');
\connect should_go
`,
			want: `SET standard_conforming_strings = false;
INSERT INTO foo VALUES ('x\'
\connect still_literal
');

`,
		},
		{
			name: "ResetStandardConformingStringsRestoresDefault",
			in: `SET standard_conforming_strings = off;
RESET standard_conforming_strings;
\connect should_go
`,
			want: `SET standard_conforming_strings = off;
RESET standard_conforming_strings;

`,
		},
		{
			name: "ResetAllRestoresDefault",
			in: `SET standard_conforming_strings = off;
RESET ALL;
\connect should_go
`,
			want: `SET standard_conforming_strings = off;
RESET ALL;

`,
		},
		{
			name: "DefaultStandardConformingStringsRestoresDefault",
			in: `SET standard_conforming_strings = off;
SET standard_conforming_strings = default;
\connect should_go
`,
			want: `SET standard_conforming_strings = off;
SET standard_conforming_strings = default;

`,
		},
		{
			name: "BlockCommentsPreserveMetaText",
			in: `/* outer block begins
/* nested: run \connect test_db for interactive work */
documenting with \connect text shouldn't strip SQL
*/
SELECT 1;
/* Change instructions:
\connect reporting

Reason: run maintenance scripts as reporting user.
*/
\connect should_go
`,
			want: `/* outer block begins
/* nested: run \connect test_db for interactive work */
documenting with \connect text shouldn't strip SQL
*/
SELECT 1;
/* Change instructions:
\connect reporting

Reason: run maintenance scripts as reporting user.
*/

`,
		},
		{
			name: "LineCommentsDoNotAffectParserState",
			in: `-- $tag$ and /* and ' are just comment text here
\connect should_go
SELECT 1;
`,
			want: `-- $tag$ and /* and ' are just comment text here

SELECT 1;
`,
		},
		{
			name: "DollarTagWithIdentifier",
			in:   "DO $foo$\n\\inside\n$foo$;\n\\set should_go\n",
			want: "DO $foo$\n\\inside\n$foo$;\n\n",
		},
		{
			name: "DollarLikeIdentifierDoesNotStartDollarQuote",
			in:   "CREATE TABLE foo$bar$baz (id int);\n\\connect should_go\n",
			want: "CREATE TABLE foo$bar$baz (id int);\n\n",
		},
		{
			name: "UnknownBackslashDirectiveWithDigitIsStripped",
			in:   "\\1 not_meta\nSELECT 1;\n",
			want: "\nSELECT 1;\n",
		},
		{
			name: "UnknownBackslashDirectiveWithUnderscoreIsStripped",
			in:   "\\_oops\nSELECT 1;\n",
			want: "\nSELECT 1;\n",
		},
		{
			name: "InlineDoubleBackslashSeparatorPreservesFollowingSQL",
			in:   "\\x \\\\ SELECT * FROM foo;\n",
			want: " SELECT * FROM foo;\n",
		},
		{
			name: "CopyFromStdinBlockIsRemovedInParseMode",
			in: `\copy foo from stdin
1	alpha
2	beta
\.
SELECT 1;
`,
			want: `



SELECT 1;
`,
		},
		{
			name: "CopyFromStdinBlockIsRemovedWithWhitespaceVariants",
			in:   "\\copy foo FROM\tSTDIN\n1\talpha\n\\.\nSELECT 1;\n",
			want: "\n\n\nSELECT 1;\n",
		},
		{
			name: "UnicodeDollarQuoteTagsArePreserved",
			in:   "DO $ä$\n\\connect still_body\n$ä$;\n\\connect should_go\n",
			want: "DO $ä$\n\\connect still_body\n$ä$;\n\n",
		},
		{
			name: "UnicodeIdentifierWithDollarDoesNotStartDollarQuote",
			in:   "CREATE TABLE ä$foo$bar (id int);\n\\connect should_go\n",
			want: "CREATE TABLE ä$foo$bar (id int);\n\n",
		},
		{
			name: "MultilineDoubleQuotedIdentifiersPreserveBackslashes",
			in:   "CREATE TABLE \"foo\n\\connect still_identifier\nbar\" (id int);\n\\connect should_go\n",
			want: "CREATE TABLE \"foo\n\\connect still_identifier\nbar\" (id int);\n\n",
		},
		{
			name: "DoubleQuotedIdentifierEscapesRemain",
			in:   "CREATE TABLE \"foo\"\"bar\" (id int);\n\\connect should_go\n",
			want: "CREATE TABLE \"foo\"\"bar\" (id int);\n\n",
		},
		{
			name: "PreservesCRLFWhenRemovingMetaCommands",
			in:   "\\connect db\r\nSELECT 1;\r\n",
			want: "\r\nSELECT 1;\r\n",
		},
		{
			name: "PreservesBareCRWhenRemovingMetaCommands",
			in:   "SELECT 1;\r\\connect db\rSELECT 2;\r",
			want: "SELECT 1;\n\nSELECT 2;\n",
		},
	}

	runRemovePsqlMetaCommandCases(t, tests)

	t.Run("CoversDocumentedMetaCommands", func(t *testing.T) {
		// Keep explicit coverage of the documented psql command set even though
		// RemovePsqlMetaCommands intentionally strips a broader class of
		// top-level backslash directives for forward compatibility.
		for _, cmd := range allPsqlMetaCommands {
			t.Run(fmt.Sprintf("strip_%s", strings.TrimPrefix(cmd, `\`)), func(t *testing.T) {
				input := fmt.Sprintf("%s -- meta command\nSELECT 42;\n", cmd)
				got, warnings, err := RemovePsqlMetaCommands(input)
				if isUnsupportedConditionalMetaCommand(cmd) {
					if err == nil {
						t.Fatalf("unsupported meta-command %q should be rejected", cmd)
					}
					return
				}
				if err != nil {
					t.Fatalf("unexpected error after stripping meta commands: %v", err)
				}
				if isSemanticMetaCommand(cmd) {
					if len(warnings) == 0 {
						t.Fatalf("expected semantic meta-command %q to emit a warning", cmd)
					}
				} else if len(warnings) != 0 {
					t.Fatalf("unexpected warnings after stripping meta commands: %v", warnings)
				}

				if strings.Contains(got, cmd+" -- meta command") {
					t.Fatalf("meta command %q line was not removed", cmd)
				}
				if !strings.Contains(got, "SELECT 42;") {
					t.Fatalf("SQL content was unexpectedly removed for %q", cmd)
				}
			})
		}
	})
}
