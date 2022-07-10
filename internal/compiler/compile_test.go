package compiler

import (
	"testing"

	"github.com/kyleconroy/sqlc/internal/config"
)

func TestQuoteIdent(t *testing.T) {
	type test struct {
		engine config.Engine
		in     string
		want   string
	}
	tests := []test{
		{config.EnginePostgreSQL, "age", "age"},
		{config.EnginePostgreSQL, "Age", `"Age"`},
		{config.EnginePostgreSQL, "CamelCase", `"CamelCase"`},
		{config.EngineMySQL, "CamelCase", "CamelCase"},
		// keywords
		{config.EnginePostgreSQL, "select", `"select"`},
		{config.EngineMySQL, "select", "`select`"},
	}

	for _, spec := range tests {
		compiler := NewCompiler(config.SQL{
			Engine: spec.engine,
		}, config.CombinedSettings{})

		t.Run(spec.in, func(t *testing.T) {
			got := compiler.quoteIdent(spec.in)
			if got != spec.want {
				t.Error("quoteIdent: engine " + string(spec.engine) + " failed for " + spec.in + ", want " + spec.want + ", got " + got)
			}
		})
	}

}
