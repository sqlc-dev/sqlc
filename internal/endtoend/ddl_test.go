package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/migrations"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
	"github.com/sqlc-dev/sqlc/internal/sqltest/pgtest"
)

func TestValidSchema(t *testing.T) {
	ctx := context.Background()

	dburi := os.Getenv("POSTGRESQL_SERVER_URI")
	if dburi == "" {
		t.Skip("POSTGRESQL_SERVER_URI is empty")
	}

	pool, err := pgxpool.New(ctx, dburi)
	if err != nil {
		t.Fatal(err)
	}

	for _, replay := range FindTests(t, "testdata", "managed-db") {
		replay := replay // https://golang.org/doc/faq#closures_and_goroutines

		if len(replay.Stderr) > 0 {
			continue
		}

		if replay.Exec != nil {
			if !slices.Contains(replay.Exec.Contexts, "managed-db") {
				continue
			}
		}

		file := filepath.Join(replay.Path, replay.ConfigName)
		rd, err := os.Open(file)
		if err != nil {
			t.Fatal(err)
		}

		conf, err := config.ParseConfig(rd)
		if err != nil {
			t.Fatal(err)
		}

		for j, pkg := range conf.SQL {
			j, pkg := j, pkg
			if pkg.Engine != config.EnginePostgreSQL {
				continue
			}
			t.Run(fmt.Sprintf("endtoend-%s-%d", file, j), func(t *testing.T) {
				t.Parallel()

				if strings.Contains(file, "pg_dump") {
					t.Skip("loading pg_dump not supported")
				}

				var schema []string
				for _, path := range pkg.Schema {
					schema = append(schema, filepath.Join(filepath.Dir(file), path))
				}

				files, err := sqlpath.Glob(schema)
				if err != nil {
					t.Fatal(err)
				}

				var sqls []string
				for _, f := range files {
					contents, err := os.ReadFile(f)
					if err != nil {
						t.Fatalf("%s: %s", f, err)
					}
					// Support loading pg_dump SQL files
					before := strings.ReplaceAll(string(contents), "CREATE SCHEMA public;", "CREATE SCHEMA IF NOT EXISTS public;")
					sqls = append(sqls, migrations.RemoveRollbackStatements(before))
				}

				uri, err := url.Parse(dburi)
				if err != nil {
					t.Fatal(err)
				}

				name, cleanup := pgtest.CreateDatabase(t, ctx, pool)
				t.Cleanup(cleanup)

				uri.Path = name
				source := uri.String()
				t.Log(source)

				conn, err := pgx.Connect(ctx, uri.String())
				if err != nil {
					t.Fatalf("connect %s: %s", name, err)
				}
				defer conn.Close(ctx)
			})
		}
	}
}
