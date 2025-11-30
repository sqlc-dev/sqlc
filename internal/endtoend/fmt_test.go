package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/engine/postgresql"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func TestFormat(t *testing.T) {
	t.Parallel()
	for _, tc := range FindTests(t, "testdata", "base") {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			// Parse the config file to determine the engine
			configPath := filepath.Join(tc.Path, tc.ConfigName)
			configFile, err := os.Open(configPath)
			if err != nil {
				t.Fatal(err)
			}
			conf, err := config.ParseConfig(configFile)
			configFile.Close()
			if err != nil {
				t.Fatal(err)
			}

			// Skip if there are no SQL packages configured
			if len(conf.SQL) == 0 {
				return
			}

			// For now, only test PostgreSQL since that's the only engine with Format support
			engine := conf.SQL[0].Engine
			if engine != config.EnginePostgreSQL {
				return
			}

			// Find query files from config
			var queryFiles []string
			for _, sql := range conf.SQL {
				for _, q := range sql.Queries {
					queryPath := filepath.Join(tc.Path, q)
					info, err := os.Stat(queryPath)
					if err != nil {
						continue
					}
					if info.IsDir() {
						// If it's a directory, glob for .sql files
						matches, err := filepath.Glob(filepath.Join(queryPath, "*.sql"))
						if err != nil {
							continue
						}
						queryFiles = append(queryFiles, matches...)
					} else {
						queryFiles = append(queryFiles, queryPath)
					}
				}
			}

			if len(queryFiles) == 0 {
				return
			}

			parse := postgresql.NewParser()

			for _, queryFile := range queryFiles {
				if _, err := os.Stat(queryFile); os.IsNotExist(err) {
					continue
				}

				contents, err := os.ReadFile(queryFile)
				if err != nil {
					t.Fatal(err)
				}

				// Parse the entire file to get proper statement boundaries
				stmts, err := parse.Parse(bytes.NewReader(contents))
				if err != nil {
					// Skip files with parse errors (e.g., syntax_errors test cases)
					return
				}

				for i, stmt := range stmts {
					stmt := stmt
					t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
						// Extract the original query text using statement location and length
						start := stmt.Raw.StmtLocation
						length := stmt.Raw.StmtLen
						if length == 0 {
							// If StmtLen is 0, it means the statement goes to the end of the input
							length = len(contents) - start
						}
						query := strings.TrimSpace(string(contents[start : start+length]))

						expected, err := postgresql.Fingerprint(query)
						if err != nil {
							t.Fatal(err)
						}

						if false {
							r, err := postgresql.Parse(query)
							debug.Dump(r, err)
						}

						out := ast.Format(stmt.Raw)
						actual, err := postgresql.Fingerprint(out)
						if err != nil {
							t.Error(err)
						}
						if expected != actual {
							debug.Dump(stmt.Raw)
							t.Errorf("- %s", expected)
							t.Errorf("- %s", query)
							t.Errorf("+ %s", actual)
							t.Errorf("+ %s", out)
						}
					})
				}
			}
		})
	}
}
