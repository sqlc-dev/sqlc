package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"

	"github.com/kyleconroy/sqlc/internal/config"
)

// Do runs the command logic.
func Do(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	rootCmd := &cobra.Command{Use: "sqlc", SilenceUsage: true}
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(genCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(versionCmd)

	rootCmd.SetArgs(args)
	rootCmd.SetIn(stdin)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)

	err := rootCmd.Execute()
	if err == nil {
		return 0
	}
	if exitError, ok := err.(*exec.ExitError); ok {
		return exitError.ExitCode()
	}
	return 1
}

var version string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the sqlc version number",
	Run: func(cmd *cobra.Command, args []string) {
		if version == "" {
			// When no version is set, return the next bug fix version
			// after the most recent tag
			fmt.Printf("%s\n", "v1.5.1-devel")
		} else {
			fmt.Printf("%s\n", version)
		}
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create an empty sqlc.yaml settings file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat("sqlc.yaml"); !os.IsNotExist(err) {
			return nil
		}
		blob, err := yaml.Marshal(config.V1GenerateSettings{Version: "1"})
		if err != nil {
			return err
		}
		return ioutil.WriteFile("sqlc.yaml", blob, 0644)
	},
}

type Env struct {
}

func ParseEnv() Env {
	return Env{}
}

var genCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Go code from SQL",
	Run: func(cmd *cobra.Command, args []string) {
		stderr := cmd.ErrOrStderr()
		dir, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(stderr, "error parsing sqlc.json: file does not exist")
			os.Exit(1)
		}

		output, err := Generate(ParseEnv(), dir, stderr)
		if err != nil {
			os.Exit(1)
		}

		for filename, source := range output {
			os.MkdirAll(filepath.Dir(filename), 0755)
			if err := ioutil.WriteFile(filename, []byte(source), 0644); err != nil {
				fmt.Fprintf(stderr, "%s: %s\n", filename, err)
				os.Exit(1)
			}
		}
	},
}

var checkCmd = &cobra.Command{
	Use:   "compile",
	Short: "Statically check SQL for syntax and type errors",
	RunE: func(cmd *cobra.Command, args []string) error {
		stderr := cmd.ErrOrStderr()
		dir, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(stderr, "error parsing sqlc.json: file does not exist")
			os.Exit(1)
		}
		if _, err := Generate(Env{}, dir, stderr); err != nil {
			os.Exit(1)
		}
		return nil
	},
}
