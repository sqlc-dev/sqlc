package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/migrations"
	"github.com/sqlc-dev/sqlc/internal/quickdb"
	pb "github.com/sqlc-dev/sqlc/internal/quickdb/v1"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
)

func TestValidSchema(t *testing.T) {
	ctx := context.Background()

	projectID := os.Getenv("DDL_SQLC_PROJECT_ID")
	authToken := os.Getenv("DDL_SQLC_AUTH_TOKEN")

	if projectID == "" || authToken == "" {
		if os.Getenv("CI") == "" {
			t.Skip("skiping ddl tests outside of CI")
		} else {
			t.Fatal("missing project id and auth token")
		}
	}

	client, err := quickdb.NewClient(projectID, authToken)
	if err != nil {
		t.Fatal(err)
	}

	files := []string{}

	// Find all tests that do not have a stderr.txt file
	err = filepath.Walk("testdata", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Base(path) == "sqlc.json" || filepath.Base(path) == "sqlc.yaml" {
			stderr := filepath.Join(filepath.Dir(path), "stderr.txt")
			if _, err := os.Stat(stderr); !os.IsNotExist(err) {
				return nil
			}
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		file := file // https://golang.org/doc/faq#closures_and_goroutines
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
					// TODO: Split schema into separate files
					before, _, _ := strings.Cut(string(contents), "-- name:")
					before, _, _ = strings.Cut(before, "/* name:")
					// Support loading pg_dump SQL files
					before = strings.ReplaceAll(before, "CREATE SCHEMA public;", "CREATE SCHEMA IF NOT EXISTS public;")
					sqls = append(sqls, migrations.RemoveRollbackStatements(before))
				}

				start := time.Now()
				resp, err := client.CreateEphemeralDatabase(ctx, &pb.CreateEphemeralDatabaseRequest{
					Engine:     "postgresql",
					Region:     "sjc",
					Migrations: sqls,
				})
				t.Logf("%s", time.Since(start))
				if err != nil {
					t.Fatal(err)
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
