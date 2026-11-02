package rewrite

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/engine/sqlite"
	"github.com/sqlc-dev/sqlc/internal/sql/validate"
)

func TestSQLiteNamedParamAfterNot(t *testing.T) {
	parser := sqlite.NewParser()
	stmts, err := parser.Parse(strings.NewReader(`SELECT 1 WHERE NOT @argname;`))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}

	raw := stmts[0].Raw
	numbers, dollar, err := validate.ParamRef(raw)
	if err != nil {
		t.Fatalf("validate.ParamRef failed: %v", err)
	}

	_, params, _ := NamedParameters(config.EngineSQLite, raw, numbers, dollar)
	if name, ok := params.NameFor(1); !ok || name != "argname" {
		t.Fatalf("expected named parameter 'argname' at position 1, got %q (ok=%v)", name, ok)
	}
}
