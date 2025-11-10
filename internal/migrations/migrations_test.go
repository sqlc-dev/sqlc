package migrations

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func containsWarningToken(warning, token string) bool {
	return strings.Contains(warning, token)
}

const inputGoose = `
-- +goose Up
ALTER TABLE archived_jobs ADD COLUMN expires_at TIMESTAMP WITH TIME ZONE;

-- +goose Down
ALTER TABLE archived_jobs DROP COLUMN expires_at;
`

const outputGoose = `
-- +goose Up
ALTER TABLE archived_jobs ADD COLUMN expires_at TIMESTAMP WITH TIME ZONE;
`

const inputMigrate = `
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE people (id int);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE people;
`

const outputMigrate = `
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE people (id int);
`

const inputTern = `
-- Write your migrate up statements here
ALTER TABLE todo RENAME COLUMN done TO is_done;
---- create above / drop below ----
ALTER TABLE todo RENAME COLUMN is_done TO done;
`

const outputTern = `
-- Write your migrate up statements here
ALTER TABLE todo RENAME COLUMN done TO is_done;`

const inputDbmate = `
-- migrate:up
CREATE TABLE foo (bar int);
-- migrate:down
DROP TABLE foo;`

const outputDbmate = `
-- migrate:up
CREATE TABLE foo (bar int);`

const inputPsqlMeta = `\restrict auwherpfqaiuwrhgp

CREATE TABLE foo (id int);

\unrestrict auwherpfqaiuwrhgp
`

const outputPsqlMeta = `

CREATE TABLE foo (id int);


`

func TestRemoveRollback(t *testing.T) {
	if diff := cmp.Diff(outputGoose, RemoveRollbackStatements(inputGoose)); diff != "" {
		t.Errorf("goose migration mismatch:\n%s", diff)
	}
	if diff := cmp.Diff(outputMigrate, RemoveRollbackStatements(inputMigrate)); diff != "" {
		t.Errorf("sql-migrate migration mismatch:\n%s", diff)
	}
	if diff := cmp.Diff(outputTern, RemoveRollbackStatements(inputTern)); diff != "" {
		t.Errorf("tern migration mismatch:\n%s", diff)
	}
	if diff := cmp.Diff(outputDbmate, RemoveRollbackStatements(inputDbmate)); diff != "" {
		t.Errorf("dbmate migration mismatch:\n%s", diff)
	}
}

func TestRemovePsqlMetaCommands(t *testing.T) {
	got, warnings, err := RemovePsqlMetaCommands(inputPsqlMeta)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	if diff := cmp.Diff(outputPsqlMeta, got); diff != "" {
		t.Errorf("psql meta-command mismatch:\n%s", diff)
	}
}

func TestPreprocessSchema(t *testing.T) {
	input := `\restrict key

CREATE TABLE foo (id int);
`
	wantPostgreSQL := `

CREATE TABLE foo (id int);`
	wantMySQL := `\restrict key

CREATE TABLE foo (id int);`

	gotPostgreSQL, warningsPostgreSQL, err := PreprocessSchema(input, "postgresql")
	if err != nil {
		t.Fatalf("unexpected postgresql error: %v", err)
	}
	if len(warningsPostgreSQL) != 0 {
		t.Fatalf("unexpected postgresql warnings: %v", warningsPostgreSQL)
	}
	if diff := cmp.Diff(wantPostgreSQL, gotPostgreSQL); diff != "" {
		t.Errorf("postgresql preprocess mismatch:\n%s", diff)
	}

	gotMySQL, warningsMySQL, err := PreprocessSchema(input, "mysql")
	if err != nil {
		t.Fatalf("unexpected mysql error: %v", err)
	}
	if len(warningsMySQL) != 0 {
		t.Fatalf("unexpected mysql warnings: %v", warningsMySQL)
	}
	if diff := cmp.Diff(wantMySQL, gotMySQL); diff != "" {
		t.Errorf("mysql preprocess mismatch:\n%s", diff)
	}
}

func TestPreprocessSchema_NormalizesBareCR(t *testing.T) {
	input := "SELECT 1;\r\\restrict key\rSELECT 2;\r"
	want := "SELECT 1;\n\nSELECT 2;"

	got, warnings, err := PreprocessSchema(input, "postgresql")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("bare CR normalization mismatch:\n%s", diff)
	}
}

func TestPreprocessSchema_RejectsPsqlConditionals(t *testing.T) {
	input := `\if false
SELECT invalid ;;;
\else
SELECT 42;
\endif
`
	if _, _, err := PreprocessSchema(input, "postgresql"); err == nil {
		t.Fatalf("expected psql conditional directives to be rejected")
	}
}

func TestPreprocessSchema_WarnsForSemanticPsqlCommands(t *testing.T) {
	tests := []struct {
		input string
		token string
	}{
		{input: `\connect db`, token: `\connect`},
		{input: `\i extra.sql`, token: `\i`},
		{input: `\include extra.sql`, token: `\include`},
		{input: `\ir extra.sql`, token: `\ir`},
		{input: `\include_relative extra.sql`, token: `\include_relative`},
		{input: `\copy foo from '/tmp/data.csv'`, token: `\copy`},
		{input: `\gexec`, token: `\gexec`},
	}

	for _, tc := range tests {
		t.Run(tc.token, func(t *testing.T) {
			got, warnings, err := PreprocessSchema(tc.input+"\nSELECT 42;\n", "postgresql")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff("\nSELECT 42;", got); diff != "" {
				t.Fatalf("unexpected output:\n%s", diff)
			}
			if len(warnings) != 1 {
				t.Fatalf("expected one warning, got %v", warnings)
			}
			if want := tc.token; warnings[0] == "" || !containsWarningToken(warnings[0], want) {
				t.Fatalf("warning %q does not mention %s", warnings[0], want)
			}
		})
	}
}

func TestPreprocessSchemaForApply_RejectsSemanticPsqlCommands(t *testing.T) {
	tests := []string{
		`\connect db`,
		`\i extra.sql`,
		`\include extra.sql`,
		`\ir extra.sql`,
		`\include_relative extra.sql`,
		`\copy foo from stdin`,
		`\gexec`,
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			if _, _, err := PreprocessSchemaForApply(input+"\nSELECT 42;\n", "postgresql"); err == nil {
				t.Fatalf("expected %q to be rejected in apply mode", input)
			}
		})
	}
}

func TestPreprocessSchema_RejectsUnterminatedCopyFromStdin(t *testing.T) {
	input := `\copy foo from stdin
1	alpha
SELECT 42;
`

	if _, _, err := PreprocessSchema(input, "postgresql"); err == nil {
		t.Fatalf("expected unterminated \\copy ... from stdin block to be rejected")
	}
}

func TestPreprocessSchema_WarnsForApproximateSessionSemantics(t *testing.T) {
	input := `BEGIN;
SET LOCAL standard_conforming_strings = off;
SELECT '\still best effort';
COMMIT;
`

	got, warnings, err := PreprocessSchema(input, "postgresql")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diff := cmp.Diff(input[:len(input)-1], got); diff != "" {
		t.Fatalf("unexpected output:\n%s", diff)
	}
	if len(warnings) != 1 {
		t.Fatalf("expected one approximation warning, got %v", warnings)
	}
	if !strings.Contains(warnings[0], "approximates psql session semantics") {
		t.Fatalf("unexpected warning: %q", warnings[0])
	}
}

func TestRemoveGolangMigrateRollback(t *testing.T) {
	filenames := map[string]bool{
		// make sure we let through golang-migrate files that aren't rollbacks
		"migrations/1.up.sql": false,
		// make sure we let through other sql files
		"migrations/2.sql":      false,
		"migrations/foo.sql":    false,
		"migrations/1.down.sql": true,
	}

	for filename, want := range filenames {
		got := IsDown(filename)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("IsDown mismatch: %s\n %s", filename, diff)
		}
	}
}
