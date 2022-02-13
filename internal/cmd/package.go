package cmd

import (
	"crypto/sha256"
	"fmt"
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
			fmt.Fprintf(stderr, "error building package: %s\n", err)
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

	// TODO: Move this to the configuration file
	owner := "tabbed"
	project := "sqlc"

	checksum := sha256.Sum256(tarball)
	sha := fmt.Sprintf("%x", checksum)
	output := fmt.Sprintf("%s_%s_%s.tar.gz", owner, project, sha[:10])
	if err := os.WriteFile(output, tarball, 0644); err != nil {
		return err
	}

	return nil
}
