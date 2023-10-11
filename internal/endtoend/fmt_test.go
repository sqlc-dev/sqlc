package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	pg_query "github.com/pganalyze/pg_query_go/v4"
	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/engine/postgresql"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func TestFormat(t *testing.T) {
	t.Parallel()
	var queries []string
	err := filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.Contains(path, filepath.Join("pgx/v5")) {
			return nil
		}
		if info.Name() == "query.sql" {
			queries = append(queries, path)
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	parse := postgresql.NewParser()
	for _, q := range queries {
		q := q
		t.Run(filepath.Dir(q), func(t *testing.T) {
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
					expected, err := pg_query.Fingerprint(string(query))
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
					out := ast.Format(stmts[0].Raw)
					actual, err := pg_query.Fingerprint(out)
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
