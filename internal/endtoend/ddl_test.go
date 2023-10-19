package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/migrations"
	"github.com/sqlc-dev/sqlc/internal/quickdb"
	pb "github.com/sqlc-dev/sqlc/internal/quickdb/v1"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
)

func TestValidSchema(t *testing.T) {
	ctx := context.Background()

	projectID := os.Getenv("CI_SQLC_PROJECT_ID")
	authToken := os.Getenv("CI_SQLC_AUTH_TOKEN")
	if projectID == "" || authToken == "" {
		t.Skip("missing project id or auth token")
	}

	client, err := quickdb.NewClient(projectID, authToken)
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

				resp, err := client.CreateEphemeralDatabase(ctx, &pb.CreateEphemeralDatabaseRequest{
					Engine:     "postgresql",
					Region:     quickdb.GetClosestRegion(),
					Migrations: sqls,
				})
				if err != nil {
					t.Fatalf("region %s: %s", quickdb.GetClosestRegion(), err)
				}

				t.Cleanup(func() {
					_, err = client.DropEphemeralDatabase(ctx, &pb.DropEphemeralDatabaseRequest{
						DatabaseId: resp.DatabaseId,
					})
					if err != nil {
						t.Fatal(err)
					}
				})

				conn, err := pgx.Connect(ctx, resp.Uri)
				if err != nil {
					t.Fatalf("connect %s: %s", resp.Uri, err)
				}
				defer conn.Close(ctx)
			})
		}
	}
}
