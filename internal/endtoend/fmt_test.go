package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/engine/postgresql"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func TestFormat(t *testing.T) {
	t.Parallel()
	parse := postgresql.NewParser()
	for _, tc := range FindTests(t, "testdata", "base") {
		tc := tc

		if !strings.Contains(tc.Path, filepath.Join("pgx/v5")) {
			continue
		}

		q := filepath.Join(tc.Path, "query.sql")
		if _, err := os.Stat(q); os.IsNotExist(err) {
			continue
		}

		t.Run(tc.Name, func(t *testing.T) {
			contents, err := os.ReadFile(q)
			if err != nil {
				t.Fatal(err)
			}
			for i, query := range bytes.Split(bytes.TrimSpace(contents), []byte(";")) {
				if len(query) <= 1 {
					continue
				}
				query := query
				t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
					expected, err := postgresql.Fingerprint(string(query))
					if err != nil {
						t.Fatal(err)
					}
					stmts, err := parse.Parse(bytes.NewReader(query))
					if err != nil {
						t.Fatal(err)
					}
					if len(stmts) != 1 {
						t.Fatal("expected one statement")
					}
					if false {
						r, err := postgresql.Parse(string(query))
						debug.Dump(r, err)
					}

					out := ast.Format(stmts[0].Raw)
					actual, err := postgresql.Fingerprint(out)
					if err != nil {
						t.Error(err)
					}
					if expected != actual {
						debug.Dump(stmts[0].Raw)
						t.Errorf("- %s", expected)
						t.Errorf("- %s", string(query))
						t.Errorf("+ %s", actual)
						t.Errorf("+ %s", out)
					}
				})
			}
		})
	}
}
