package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	pg "github.com/lfittl/pg_query_go"
	"github.com/spf13/cobra"

	"github.com/kyleconroy/sqlc/internal/config"
)

// Do runs the command logic.
func Do(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	rootCmd := &cobra.Command{Use: "sqlc", SilenceUsage: true}
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(genCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(parseCmd)
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

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the sqlc version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v0.0.1")
	},
}

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse and print the AST for a SQL file",
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, filename := range args {
			blob, err := ioutil.ReadFile(filename)
			if err != nil {
				return err
			}
			tree, err := pg.Parse(string(blob))
			if err != nil {
				return err
			}
			spew.Dump(tree)
		}
		return nil
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create an empty sqlc.json settings file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat("sqlc.json"); !os.IsNotExist(err) {
			return nil
		}
		blob, err := json.MarshalIndent(config.Config{Version: "1"}, "", "  ")
		if err != nil {
			return err
		}
		return ioutil.WriteFile("sqlc.json", blob, 0644)
	},
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

		output, err := Generate(dir, stderr)
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
		if _, err := Generate(dir, stderr); err != nil {
			os.Exit(1)
		}
		return nil
	},
}
