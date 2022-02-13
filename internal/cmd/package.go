package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/kyleconroy/sqlc/internal/bundler"
)

var packageCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload the schema, queries, and configuration for this project",
	RunE: func(cmd *cobra.Command, args []string) error {
		stderr := cmd.ErrOrStderr()
		dir, name := getConfigPath(stderr, cmd.Flag("file"))
		if err := createPkg(cmd.Context(), ParseEnv(cmd), dir, name, stderr); err != nil {
			fmt.Fprintf(stderr, "error uploading project: %s\n", err)
			os.Exit(1)
		}
		return nil
	},
}

func createPkg(ctx context.Context, e Env, dir, filename string, stderr io.Writer) error {
	output, err := Generate(ctx, e, dir, filename, stderr)
	if err != nil {
		os.Exit(1)
	}
	configPath, conf, err := readConfig(stderr, dir, filename)
	if err != nil {
		return err
	}
	up := bundler.NewUploader(configPath, dir, conf)
	if err != nil {
		return err
	}
	return up.Upload(ctx, output)
}
