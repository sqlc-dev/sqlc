package compiler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/multierr"
	"github.com/sqlc-dev/sqlc/internal/opts"
)

func TestParseCatalogManagedAnalyzerRejectsSemanticPsqlCommandsForApply(t *testing.T) {
	dir := t.TempDir()
	schema := filepath.Join(dir, "schema.sql")
	if err := os.WriteFile(schema, []byte("\\include extra.sql\nCREATE TABLE foo (id int);\n"), 0600); err != nil {
		t.Fatal(err)
	}

	c, err := NewCompiler(config.SQL{
		Engine: config.EnginePostgreSQL,
		Schema: []string{schema},
		Database: &config.Database{
			Managed: true,
		},
	}, config.CombinedSettings{}, opts.Parser{})
	if err != nil {
		t.Fatal(err)
	}

	err = c.ParseCatalog([]string{schema})
	if err == nil {
		t.Fatal("expected managed analyzer schema preprocessing to reject semantic psql command")
	}
	merr, ok := err.(*multierr.Error)
	if !ok || len(merr.Errs()) != 1 {
		t.Fatalf("expected one schema error, got %T: %v", err, err)
	}
	if !strings.Contains(merr.Errs()[0].Err.Error(), `psql meta-command \include is not supported`) {
		t.Fatalf("unexpected error: %v", err)
	}
}
