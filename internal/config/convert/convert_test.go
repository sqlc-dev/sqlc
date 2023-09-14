package convert

import (
	"testing"

	"gopkg.in/yaml.v3"
)

const anchor = `
sql:
  - schema: query.sql
    queries: query.sql
    engine: postgresql
    codegen:
      - out: gateway/src/gateway/services/organization
        plugin: py
        options: &base-options
          query_parameter_limit: 1
          package: gateway
          output_models_file_name: null
          emit_module: true
          emit_generators: false
          emit_async: true

  - schema: query.sql
    queries: query.sql
    engine: postgresql
    codegen:
      - out: gateway/src/gateway/services/project
        plugin: py
        options: *base-options
`

type config struct {
	SQL yaml.Node `yaml:"sql"`
}

func TestAlias(t *testing.T) {
	var a config
	err := yaml.Unmarshal([]byte(anchor), &a)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := gen(&a.SQL); err != nil {
		t.Fatal(err)
	}
}
