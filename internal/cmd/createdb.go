package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime/trace"
	"time"

	"github.com/spf13/cobra"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/dbmanager"
	"github.com/sqlc-dev/sqlc/internal/migrations"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
)

var createDBCmd = &cobra.Command{
	Use:   "createdb",
	Short: "Create an ephemeral database",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		defer trace.StartRegion(cmd.Context(), "createdb").End()
		stderr := cmd.ErrOrStderr()
		dir, filename := getConfigPath(stderr, cmd.Flag("file"))
		querySetName, err := cmd.Flags().GetString("queryset")
		if err != nil {
			return err
		}
		err = CreateDB(cmd.Context(), dir, filename, querySetName, &Options{
			Env:    ParseEnv(cmd),
			Stderr: stderr,
		})
		if err != nil {
			fmt.Fprintln(stderr, err.Error())
			os.Exit(1)
		}
		return nil
	},
}

func CreateDB(ctx context.Context, dir, filename, querySetName string, o *Options) error {
	_, conf, err := o.ReadConfig(dir, filename)
	if err != nil {
		return err
	}
	// Find the first queryset with a managed database
	var queryset *config.SQL
	var count int
	for _, sql := range conf.SQL {
		sql := sql
		if querySetName != "" && sql.Name != querySetName {
			continue
		}
		if sql.Database != nil && sql.Database.Managed {
			queryset = &sql
			count += 1
		}
	}
	if queryset == nil && querySetName != "" {
		return fmt.Errorf("no queryset found with name %q", querySetName)
	}
	if queryset == nil {
		return fmt.Errorf("no querysets configured to use a managed database")
	}
	if count > 1 {
		return fmt.Errorf("multiple querysets configured to use managed databases")
	}

	switch queryset.Engine {
	case config.EngineMySQL:
		// pass
	case config.EnginePostgreSQL:
		// pass
	default:
		return fmt.Errorf("createdb does not support the %s engine", queryset.Engine)
	}

	var ddl []string
	files, err := sqlpath.Glob(queryset.Schema)
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

	now := time.Now().UTC().UnixNano()
	client := dbmanager.NewClient(conf.Servers)
	resp, err := client.CreateDatabase(ctx, &dbmanager.CreateDatabaseRequest{
		Engine:     string(queryset.Engine),
		Migrations: ddl,
		Prefix:     fmt.Sprintf("sqlc_createdb_%d", now),
	})
	if err != nil {
		return fmt.Errorf("managed: create database: %w", err)
	}
	fmt.Println(resp.Uri)
	return nil
}
