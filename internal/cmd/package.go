package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/kyleconroy/sqlc/internal/bundler"
)

var packageCmd = &cobra.Command{
	Use:   "build",
	Short: "Create a tarball containing schema, queries, and configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		stderr := cmd.ErrOrStderr()
		dir, name := getConfigPath(stderr, cmd.Flag("file"))
		if err := createPkg(ParseEnv(cmd), dir, name, stderr); err != nil {
			os.Exit(1)
		}
		return nil
	},
}

func createPkg(e Env, dir, filename string, stderr io.Writer) error {
	configPath, conf, err := readConfig(stderr, dir, filename)
	if err != nil {
		return err
	}
	tarball, err := bundler.Build(configPath, conf)
	if err != nil {
		return err
	}

	// TODO: Gzip?
	if err := os.WriteFile("sqlc-package-test.tar", tarball, 0644); err != nil {
		return err
	}

	return nil
}
