package seed_test

import (
	"testing"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/core/seed"
	"github.com/sqlc-dev/sqlc/internal/engine/clickhouse"
	"github.com/sqlc-dev/sqlc/internal/engine/dolphin"
	"github.com/sqlc-dev/sqlc/internal/engine/googlesql"
	"github.com/sqlc-dev/sqlc/internal/engine/postgresql"
	"github.com/sqlc-dev/sqlc/internal/engine/sqlite"
)

// TestDialects loads every engine's dialect.json into a catalog. A malformed
// document, a misspelled key or an operator over a type the document does not
// define fails here rather than on the first query analyzed.
func TestDialects(t *testing.T) {
	t.Parallel()

	dialects := map[string]struct {
		option core.Option
		// types are spellings a schema written for the dialect may use, each
		// of which has to resolve and to compare against itself.
		types []string
	}{
		"postgresql": {postgresql.Dialect(), []string{"int4", "integer", "bigserial", "text", "timestamptz"}},
		"mysql":      {dolphin.Dialect(), []string{"bigint", "varchar", "datetime", "signed"}},
		"sqlite":     {sqlite.Dialect(), []string{"integer", "text", "real", "blob"}},
		"clickhouse": {clickhouse.Dialect(), []string{"uint64", "string", "datetime"}},
		"googlesql":  {googlesql.Dialect(), []string{"int64", "string", "timestamp"}},
	}

	for name, d := range dialects {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cat, err := core.New(d.option)
			if err != nil {
				t.Fatal(err)
			}
			defer cat.Close()

			if _, err := cat.DialectOID(name); err != nil {
				t.Errorf("dialect %s not registered: %s", name, err)
			}

			// Every dialect types the four kinds of literal.
			for _, kind := range []string{core.ConstInteger, core.ConstFloat, core.ConstString, core.ConstBool} {
				if _, err := cat.ConstTypeOID(kind); err != nil {
					t.Errorf("%s literal: %s", kind, err)
				}
			}

			for _, typeName := range d.types {
				oid, err := cat.TypeOID(typeName)
				if err != nil {
					t.Errorf("type %q: %s", typeName, err)
					continue
				}
				ops, err := cat.FindOperators("=", oid, oid)
				if err != nil {
					t.Fatal(err)
				}
				if len(ops) == 0 {
					t.Errorf("type %q has no %q operator", typeName, "=")
				}
			}
		})
	}
}

func TestLoadRejectsUnknownField(t *testing.T) {
	t.Parallel()

	if _, err := seed.Load([]byte(`{"dialect": "x", "typos": []}`)); err == nil {
		t.Fatal("expected an unknown field to be rejected")
	}
}

func TestLoadRequiresName(t *testing.T) {
	t.Parallel()

	if _, err := seed.Load([]byte(`{"types": []}`)); err == nil {
		t.Fatal("expected a dialect without a name to be rejected")
	}
}
