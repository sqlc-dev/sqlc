package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime/trace"

	"github.com/spf13/cobra"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/migrations"
	"github.com/sqlc-dev/sqlc/internal/quickdb"
	pb "github.com/sqlc-dev/sqlc/internal/quickdb/v1"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
)

var createDBCmd = &cobra.Command{
	Use:   "createdb",
	Short: "Create an ephemeral database",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		defer trace.StartRegion(cmd.Context(), "createdb").End()
		stderr := cmd.ErrOrStderr()
		dir, name := getConfigPath(stderr, cmd.Flag("file"))
		qs, err := cmd.Flags().GetString("queryset")
		if err != nil {
			return err
		}
		err = CreateDB(cmd.Context(), dir, name, qs, &Options{
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

func CreateDB(ctx context.Context, dir, filename, name string, o *Options) error {
	_, conf, err := o.ReadConfig(dir, filename)
	if err != nil {
		return err
	}
	// Find the first queryset with a managed database
	var queryset *config.SQL
	var count int
	for _, sql := range conf.SQL {
		sql := sql
		if name != "" && sql.Name != name {
			continue
		}
		if sql.Database != nil && sql.Database.Managed {
			queryset = &sql
			count += 1
		}
	}
	if queryset == nil && name != "" {
		return fmt.Errorf("no queryset found with name %q", name)
	}
	if queryset == nil {
		return fmt.Errorf("no querysets configured to use a managed database")
	}
	if count > 1 {
		return fmt.Errorf("multiple querysets configured to use managed databases")
	}
	if queryset.Engine != config.EnginePostgreSQL {
		return fmt.Errorf("managed databases currently only support PostgreSQL")
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

	client, err := quickdb.NewClientFromConfig(conf.Cloud)
	if err != nil {
		return fmt.Errorf("client error: %w", err)
	}

	resp, err := client.CreateEphemeralDatabase(ctx, &pb.CreateEphemeralDatabaseRequest{
		Engine:     "postgresql",
		Region:     quickdb.GetClosestRegion(),
		Migrations: ddl,
	})
	if err != nil {
		return fmt.Errorf("managed: create database: %w", err)
	}
	fmt.Fprintln(os.Stderr, "WARNING: This database will be removed in two minutes")
	fmt.Println(resp.Uri)
	return nil
}
