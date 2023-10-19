package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime/trace"
	"strings"

	"github.com/spf13/cobra"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/opts"
	"github.com/sqlc-dev/sqlc/internal/quickdb"
	pb "github.com/sqlc-dev/sqlc/internal/quickdb/v1"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
)

var createDBCmd = &cobra.Command{
	Use:   "createdb",
	Short: "Create an ephemeral database",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		defer trace.StartRegion(cmd.Context(), "createdb").End()
		stderr := cmd.ErrOrStderr()
		dir, name := getConfigPath(stderr, cmd.Flag("file"))
		env, err := cmd.Flags().GetString("env")
		if err != nil {
			return err
		}
		code, err := CreateDB(cmd.Context(), dir, name, args, env, &Options{
			Env:    ParseEnv(cmd),
			Stderr: stderr,
		})
		if err != nil {
			fmt.Fprintln(stderr, err.Error())
			os.Exit(code)
		}
		return nil
	},
}

func CreateDB(ctx context.Context, dir, filename string, args []string, env string, o *Options) (int, error) {
	dbg := opts.DebugFromEnv()
	if !dbg.ProcessPlugins {
		return 1, fmt.Errorf("process-plugins disabled")
	}
	_, conf, err := o.ReadConfig(dir, filename)
	if err != nil {
		return 1, err
	}
	// Find the first SQL with a managed database
	var pkg *config.SQL
	for _, sql := range conf.SQL {
		sql := sql
		if sql.Database != nil && sql.Database.Managed {
			pkg = &sql
			break
		}
	}
	if pkg == nil {
		return 1, fmt.Errorf("no managed database found")
	}
	if pkg.Engine != config.EnginePostgreSQL {
		return 1, fmt.Errorf("managed: only PostgreSQL currently")
	}

	var migrations []string
	files, err := sqlpath.Glob(pkg.Schema)
	if err != nil {
		return 1, err
	}
	for _, schema := range files {
		contents, err := os.ReadFile(schema)
		if err != nil {
			return 1, fmt.Errorf("read file: %w", err)
		}
		migrations = append(migrations, string(contents))
	}
	client, err := quickdb.NewClientFromConfig(conf.Cloud)
	if err != nil {
		return 1, fmt.Errorf("client error: %w", err)
	}

	resp, err := client.CreateEphemeralDatabase(ctx, &pb.CreateEphemeralDatabaseRequest{
		Engine:     "postgresql",
		Region:     quickdb.GetClosestRegion(),
		Migrations: migrations,
	})
	if err != nil {
		return 1, fmt.Errorf("managed: create database: %w", err)
	}

	defer func() {
		client.DropEphemeralDatabase(ctx, &pb.DropEphemeralDatabaseRequest{
			DatabaseId: resp.DatabaseId,
		})
	}()

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", env, resp.Uri))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = []string{fmt.Sprintf("%s=%s", env, resp.Uri)}
	for _, val := range os.Environ() {
		if strings.HasPrefix(val, "SQLC_AUTH_TOKEN") {
			continue
		}
		cmd.Env = append(cmd.Env, val)
	}

	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode(), err
		}
		return 1, err
	}

	return 0, nil
}
