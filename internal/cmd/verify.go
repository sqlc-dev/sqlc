package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/dbmanager"
	"github.com/sqlc-dev/sqlc/internal/migrations"
	"github.com/sqlc-dev/sqlc/internal/plugin"
	"github.com/sqlc-dev/sqlc/internal/quickdb"
	pb "github.com/sqlc-dev/sqlc/internal/quickdb/v1"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
)

func init() {
	verifyCmd.Flags().String("against", "", "compare against this tag")
}

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify schema, queries, and configuration for this project",
	RunE: func(cmd *cobra.Command, args []string) error {
		stderr := cmd.ErrOrStderr()
		dir, name := getConfigPath(stderr, cmd.Flag("file"))
		against, err := cmd.Flags().GetString("against")
		if err != nil {
			return err
		}
		opts := &Options{
			Env:     ParseEnv(cmd),
			Stderr:  stderr,
			Against: against,
		}
		if err := Verify(cmd.Context(), dir, name, opts); err != nil {
			fmt.Fprintf(stderr, "Error verifying queries: %s\n", err)
			os.Exit(1)
		}
		return nil
	},
}

func Verify(ctx context.Context, dir, filename string, opts *Options) error {
	stderr := opts.Stderr
	_, conf, err := readConfig(stderr, dir, filename)
	if err != nil {
		return err
	}

	client, err := quickdb.NewClientFromConfig(conf.Cloud)
	if err != nil {
		return fmt.Errorf("client init failed: %w", err)
	}

	manager := dbmanager.NewClient(conf.Servers)

	// Get query sets from a previous archive by tag. If no tag is provided, get
	// the latest query sets.
	previous, err := client.GetQuerySets(ctx, &pb.GetQuerySetsRequest{
		Tag: opts.Against,
	})
	if err != nil {
		return err
	}

	// Create a mapping of name to query set
	existing := map[string]config.SQL{}
	for _, qs := range conf.SQL {
		existing[qs.Name] = qs
	}

	var verr error
	for _, qs := range previous.QuerySets {
		// TODO: Create a function for this so that we can return early on errors

		check := func() error {
			if qs.Name == "" {
				return fmt.Errorf("unnamed query set")
			}

			current, found := existing[qs.Name]
			if !found {
				return fmt.Errorf("unknown query set: %s", qs.Name)
			}

			// Read the schema files into memory, removing rollback statements
			var ddl []string
			files, err := sqlpath.Glob(current.Schema)
			if err != nil {
				return err
			}
			for _, schema := range files {
				contents, err := os.ReadFile(schema)
				if err != nil {
					return fmt.Errorf("read file: %w", err)
				}
				ddl = append(ddl, migrations.RemoveRollbackStatements(string(contents)))
			}

			var codegen plugin.GenerateRequest
			if err := proto.Unmarshal(qs.CodegenRequest.Contents, &codegen); err != nil {
				return err
			}

			// Create (or re-use) a database to verify against
			resp, err := manager.CreateDatabase(ctx, &dbmanager.CreateDatabaseRequest{
				Engine:     string(current.Engine),
				Migrations: ddl,
			})
			if err != nil {
				return err
			}

			db, err := sql.Open("pgx", resp.Uri)
			if err != nil {
				return err
			}
			defer db.Close()

			var qerr error
			for _, query := range codegen.Queries {
				stmt, err := db.PrepareContext(ctx, query.Text)
				if err != nil {
					fmt.Fprintf(stderr, "Failed to prepare the following query:\n")
					fmt.Fprintf(stderr, "%s\n", query.Text)
					fmt.Fprintf(stderr, "Error was: %s\n", err)
					qerr = err
					continue
				}
				if err := stmt.Close(); err != nil {
					slog.Error("stmt.Close failed", "err", err)
				}
			}

			return qerr
		}

		if err := check(); err != nil {
			verr = errors.New("errored")
			fmt.Fprintf(stderr, "FAIL\t%s\n", qs.Name)
			fmt.Fprintf(stderr, "  ERROR\t%s\n", err)
		} else {
			fmt.Fprintf(stderr, "ok\t%s\n", qs.Name)
		}
	}

	return verr
}
