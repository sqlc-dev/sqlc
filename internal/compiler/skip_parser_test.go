package compiler

import (
	"testing"

	"github.com/sqlc-dev/sqlc/internal/config"
)

func TestSkipParserRequiresDatabase(t *testing.T) {
	skipParser := true
	conf := config.SQL{
		Engine: config.EnginePostgreSQL,
		Analyzer: config.Analyzer{
			SkipParser: &skipParser,
		},
	}

	combo := config.CombinedSettings{
		Package: conf,
	}

	_, err := NewCompiler(conf, combo)
	if err == nil {
		t.Fatal("expected error when skip_parser is true without database config")
	}
	if err.Error() != "skip_parser requires database configuration" {
		t.Errorf("unexpected error message: %s", err)
	}
}

func TestSkipParserRequiresDatabaseAnalyzer(t *testing.T) {
	skipParser := true
	analyzerDisabled := false
	conf := config.SQL{
		Engine: config.EnginePostgreSQL,
		Database: &config.Database{
			URI: "postgresql://localhost/test",
		},
		Analyzer: config.Analyzer{
			SkipParser: &skipParser,
			Database:   &analyzerDisabled,
		},
	}

	combo := config.CombinedSettings{
		Package: conf,
	}

	_, err := NewCompiler(conf, combo)
	if err == nil {
		t.Fatal("expected error when skip_parser is true but database analyzer is disabled")
	}
	if err.Error() != "skip_parser requires database analyzer to be enabled" {
		t.Errorf("unexpected error message: %s", err)
	}
}

func TestSkipParserOnlyPostgreSQL(t *testing.T) {
	skipParser := true
	engines := []config.Engine{
		config.EngineMySQL,
		config.EngineSQLite,
	}

	for _, engine := range engines {
		conf := config.SQL{
			Engine: engine,
			Database: &config.Database{
				URI: "test://localhost/test",
			},
			Analyzer: config.Analyzer{
				SkipParser: &skipParser,
			},
		}

		combo := config.CombinedSettings{
			Package: conf,
		}

		_, err := NewCompiler(conf, combo)
		if err == nil {
			t.Fatalf("expected error for engine %s with skip_parser", engine)
		}
		if err.Error() != "skip_parser is only supported for PostgreSQL" {
			t.Errorf("unexpected error message for %s: %s", engine, err)
		}
	}
}

func TestSkipParserValidConfig(t *testing.T) {
	skipParser := true
	conf := config.SQL{
		Engine: config.EnginePostgreSQL,
		Database: &config.Database{
			URI: "postgresql://localhost/test",
		},
		Analyzer: config.Analyzer{
			SkipParser: &skipParser,
		},
	}

	combo := config.CombinedSettings{
		Package: conf,
	}

	c, err := NewCompiler(conf, combo)
	if err != nil {
		t.Fatalf("unexpected error with valid skip_parser config: %s", err)
	}

	// Verify parser and catalog are nil when skip_parser is true
	if c.parser != nil {
		t.Error("expected parser to be nil when skip_parser is true")
	}
	if c.catalog != nil {
		t.Error("expected catalog to be nil when skip_parser is true")
	}
	// Analyzer should still be set (but we can't check it without a real DB connection)
}

func TestSkipParserDisabledNormalOperation(t *testing.T) {
	skipParser := false
	conf := config.SQL{
		Engine: config.EnginePostgreSQL,
		Analyzer: config.Analyzer{
			SkipParser: &skipParser,
		},
	}

	combo := config.CombinedSettings{
		Package: conf,
	}

	c, err := NewCompiler(conf, combo)
	if err != nil {
		t.Fatalf("unexpected error with skip_parser=false: %s", err)
	}

	// Verify parser and catalog ARE set when skip_parser is false
	if c.parser == nil {
		t.Error("expected parser to be set when skip_parser is false")
	}
	if c.catalog == nil {
		t.Error("expected catalog to be set when skip_parser is false")
	}
}

func TestSkipParserDefaultNormalOperation(t *testing.T) {
	// When skip_parser is not specified (nil), should work normally
	conf := config.SQL{
		Engine: config.EnginePostgreSQL,
	}

	combo := config.CombinedSettings{
		Package: conf,
	}

	c, err := NewCompiler(conf, combo)
	if err != nil {
		t.Fatalf("unexpected error with default config: %s", err)
	}

	// Verify parser and catalog ARE set by default
	if c.parser == nil {
		t.Error("expected parser to be set by default")
	}
	if c.catalog == nil {
		t.Error("expected catalog to be set by default")
	}
}
