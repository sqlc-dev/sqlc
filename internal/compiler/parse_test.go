package compiler

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/engine/dolphin"
	"github.com/kyleconroy/sqlc/internal/engine/postgresql"
	"github.com/kyleconroy/sqlc/internal/engine/sqlite"
	"github.com/kyleconroy/sqlc/internal/opts"
)

func Test_ParseQueryErrors(t *testing.T) {
	for _, tc := range []struct {
		name       string
		engine     config.Engine
		createStmt string
		selectStmt string
		parser     Parser
		wantErr    error
	}{
		{
			name:       "unreferenced order column postgresql",
			engine:     config.EnginePostgreSQL,
			createStmt: `CREATE TABLE authors (id INT);`,
			selectStmt: `SELECT id FROM authors ORDER BY foo;`,
			parser:     postgresql.NewParser(),
			wantErr:    fmt.Errorf(`column reference "foo" not found`),
		},
		{
			name:       "unreferenced order column mysql",
			engine:     config.EngineMySQL,
			createStmt: `CREATE TABLE authors (id INT);`,
			selectStmt: `SELECT id FROM authors ORDER BY foo;`,
			parser:     dolphin.NewParser(),
			wantErr:    fmt.Errorf(`column reference "foo" not found`),
		},
		{
			name:       "unreferenced order column sqlite",
			engine:     config.EngineSQLite,
			createStmt: `CREATE TABLE authors (id INT);`,
			selectStmt: `SELECT id FROM authors ORDER BY foo;`,
			parser:     sqlite.NewParser(),
			wantErr:    fmt.Errorf(`column reference "foo" not found`),
		},
		{
			name:       "unreferenced group column postgresql",
			engine:     config.EnginePostgreSQL,
			createStmt: `CREATE TABLE authors ( id INT );`,
			selectStmt: `SELECT id FROM authors GROUP BY foo;`,
			parser:     postgresql.NewParser(),
			wantErr:    fmt.Errorf(`column reference "foo" not found`),
		},
		{
			name:       "unreferenced group column mysql",
			engine:     config.EngineMySQL,
			createStmt: `CREATE TABLE authors (id INT);`,
			selectStmt: `SELECT id FROM authors GROUP BY foo;`,
			parser:     dolphin.NewParser(),
			wantErr:    fmt.Errorf(`column reference "foo" not found`),
		},
		{
			name:       "unreferenced group column sqlite",
			engine:     config.EngineSQLite,
			createStmt: `CREATE TABLE authors (id INT);`,
			selectStmt: `SELECT id FROM authors GROUP BY foo;`,
			parser:     sqlite.NewParser(),
			wantErr:    fmt.Errorf(`column reference "foo" not found`),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			conf := config.SQL{
				Engine: tc.engine,
			}
			combo := config.CombinedSettings{}
			comp := NewCompiler(conf, combo)
			stmts, err := tc.parser.Parse(strings.NewReader(tc.createStmt))
			if err != nil {
				t.Fatalf("cannot parse test catalog: %v", err)
			}
			err = comp.catalog.Update(stmts[0], comp)
			if err != nil {
				t.Fatalf("cannot update test catalog: %v", err)
			}
			stmts, err = tc.parser.Parse(strings.NewReader(tc.selectStmt))
			if err != nil {
				t.Errorf("Parse failed: %v", err)
			}
			if len(stmts) != 1 {
				t.Errorf("expected one statement, got %d", len(stmts))
			}

			_, err = comp.parseQuery(stmts[0].Raw, tc.selectStmt, opts.Parser{})
			if err == nil {
				t.Fatalf("expected parseQuery to return an error, got nil")
			}
			if err.Error() != tc.wantErr.Error() {
				t.Errorf("error message: want %s, got %s", tc.wantErr, err.Error())
			}
		})
	}
}
