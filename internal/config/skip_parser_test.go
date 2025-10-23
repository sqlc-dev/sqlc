package config

import (
	"strings"
	"testing"
)

func TestSkipParserConfig(t *testing.T) {
	yaml := `
version: "2"
sql:
  - name: "test"
    engine: "postgresql"
    queries: "query.sql"
    schema: "schema.sql"
    database:
      uri: "postgresql://localhost/test"
    analyzer:
      skip_parser: true
    gen:
      go:
        package: "test"
        out: "test"
`

	conf, err := ParseConfig(strings.NewReader(yaml))
	if err != nil {
		t.Fatalf("failed to parse config: %s", err)
	}

	if len(conf.SQL) != 1 {
		t.Fatalf("expected 1 SQL config, got %d", len(conf.SQL))
	}

	sql := conf.SQL[0]
	if sql.Analyzer.SkipParser == nil {
		t.Fatal("expected skip_parser to be set")
	}
	if !*sql.Analyzer.SkipParser {
		t.Error("expected skip_parser to be true")
	}
}

func TestSkipParserConfigDefault(t *testing.T) {
	yaml := `
version: "2"
sql:
  - name: "test"
    engine: "postgresql"
    queries: "query.sql"
    schema: "schema.sql"
    gen:
      go:
        package: "test"
        out: "test"
`

	conf, err := ParseConfig(strings.NewReader(yaml))
	if err != nil {
		t.Fatalf("failed to parse config: %s", err)
	}

	if len(conf.SQL) != 1 {
		t.Fatalf("expected 1 SQL config, got %d", len(conf.SQL))
	}

	sql := conf.SQL[0]
	if sql.Analyzer.SkipParser != nil {
		t.Errorf("expected skip_parser to be nil (default), got %v", *sql.Analyzer.SkipParser)
	}
}

func TestSkipParserConfigFalse(t *testing.T) {
	yaml := `
version: "2"
sql:
  - name: "test"
    engine: "postgresql"
    queries: "query.sql"
    schema: "schema.sql"
    analyzer:
      skip_parser: false
    gen:
      go:
        package: "test"
        out: "test"
`

	conf, err := ParseConfig(strings.NewReader(yaml))
	if err != nil {
		t.Fatalf("failed to parse config: %s", err)
	}

	if len(conf.SQL) != 1 {
		t.Fatalf("expected 1 SQL config, got %d", len(conf.SQL))
	}

	sql := conf.SQL[0]
	if sql.Analyzer.SkipParser == nil {
		t.Fatal("expected skip_parser to be set")
	}
	if *sql.Analyzer.SkipParser {
		t.Error("expected skip_parser to be false")
	}
}
