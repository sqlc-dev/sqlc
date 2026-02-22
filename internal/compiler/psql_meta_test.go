package compiler

import (
	"fmt"
	"strings"
	"testing"
)

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

func TestRemovePsqlMetaCommands_TableDriven(t *testing.T) {
	inDoubleQuoted := "CREATE TABLE \"foo\\bar\" (id int);\nSELECT \"foo\\bar\"." +
		"id FROM \"foo\\bar\";\n"
	inValidSQL := "CREATE TABLE t (id int);\nINSERT INTO t VALUES (1);\n"
	inWhitespaceOnly := "   \t  "
	inNoTrailingNewline := "SELECT 1"
	inBackslashNotAtStart := "SELECT '\\not_meta' AS col;\n  SELECT '\\still_not_meta';\n"
	inDoubleSingleQuotes := "INSERT INTO t VALUES ('It''s fine');\n"

	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "RemovesTopLevelMetaCommands",
			in:   "CREATE TABLE public.authors();\n\\connect test\n  \\set ON_ERROR_STOP on\nSELECT 1;\n",
			want: "CREATE TABLE public.authors();\n\n\nSELECT 1;\n",
		},
		{
			name: "IgnoresBackslashesInStrings",
			in:   "SELECT E'\\n' || E'\\' || '\n\\restrict inside';\nSELECT E'\n\\still_string\n';\n\\connect nope\n",
			want: "SELECT E'\\n' || E'\\' || '\n\\restrict inside';\nSELECT E'\n\\still_string\n';\n\n",
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
			name: "DollarTagWithIdentifier",
			in:   "DO $foo$\n\\inside\n$foo$;\n\\set should_go\n",
			want: "DO $foo$\n\\inside\n$foo$;\n\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := removePsqlMetaCommands(tc.in)
			if got != tc.want {
				t.Fatalf("unexpected output after stripping meta commands:\nwant=%q\ngot =%q", tc.want, got)
			}
		})
	}

	t.Run("CoversDocumentedMetaCommands", func(t *testing.T) {
		for _, cmd := range allPsqlMetaCommands {
			t.Run(fmt.Sprintf("strip_%s", strings.TrimPrefix(cmd, `\`)), func(t *testing.T) {
				input := fmt.Sprintf("%s -- meta command\nSELECT 42;\n", cmd)
				got := removePsqlMetaCommands(input)

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
