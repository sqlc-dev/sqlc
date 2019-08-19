package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kyleconroy/sqlc/internal/dinosql"

	"github.com/davecgh/go-spew/spew"
	pg "github.com/lfittl/pg_query_go"
	"github.com/spf13/cobra"
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
	rootCmd.SetErr(stderr)
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
		blob, err := json.MarshalIndent(dinosql.GenerateSettings{}, "  ", "")
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
		blob, err := ioutil.ReadFile("sqlc.json")
		if err != nil {
			fmt.Fprintln(os.Stderr, "error parsing sqlc.json: file does not exist")
			os.Exit(1)
		}

		var settings dinosql.GenerateSettings
		if err := json.Unmarshal(blob, &settings); err != nil {
			switch err.(type) {
			// TODO: Provide better error messages for sqlc.json parsing
			// case *json.SyntaxError:
			// case *json.InvalidUnmarshalError:
			// case *json.UnmarshalFieldError:
			// case *json.UnmarshalTypeError:
			// case *json.UnsupportedTypeError:
			// case *json.UnsupportedValueError:
			default:
				fmt.Fprintf(os.Stderr, "error parsing sqlc.json: %s\n", err)
			}
			os.Exit(1)
		}

		var errored bool

		for i, pkg := range settings.Packages {
			name := pkg.Name

			if pkg.Path == "" {
				fmt.Fprintf(os.Stderr, "package[%d]: path must be set\n", i)
				errored = true
				continue
			}

			if name == "" {
				name = filepath.Base(pkg.Path)
			}

			c, err := dinosql.ParseCatalog(pkg.Schema)
			if err != nil {
				fmt.Fprintf(os.Stderr, "# package %s\n", name)
				if parserErr, ok := err.(*dinosql.ParserErr); ok {
					for _, fileErr := range parserErr.Errs {
						fmt.Fprintf(os.Stderr, "%s:%d:%d: %s\n", fileErr.Filename, fileErr.Line, fileErr.Column, fileErr.Err)
					}
				} else {
					fmt.Fprintf(os.Stderr, "error parsing schema: %s\n", err)
				}
				errored = true
				continue
			}

			q, err := dinosql.ParseQueries(c, settings, pkg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "# package %s\n", name)
				if parserErr, ok := err.(*dinosql.ParserErr); ok {
					for _, fileErr := range parserErr.Errs {
						fmt.Fprintf(os.Stderr, "%s:%d:%d: %s\n", fileErr.Filename, fileErr.Line, fileErr.Column, fileErr.Err)
					}
				} else {
					fmt.Fprintf(os.Stderr, "error parsing schema: %s\n", err)
				}
				errored = true
				continue
			}

			files, err := dinosql.Generate(q, settings, pkg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "# package %s\n", name)
				fmt.Fprintf(os.Stderr, "error generating code: %s\n", err)
				errored = true
				continue
			}

			os.MkdirAll(pkg.Path, 0755)

			for n, source := range files {
				filename := filepath.Join(pkg.Path, n)
				if err := ioutil.WriteFile(filename, []byte(source), 0644); err != nil {
					fmt.Fprintf(os.Stderr, "# package %s\n", name)
					fmt.Fprintf(os.Stderr, "%s: %s\n", filename, err)
					os.Exit(1)
				}
			}
		}

		if errored {
			os.Exit(1)
		}
	},
}

var checkCmd = &cobra.Command{
	Use:   "compile",
	Short: "Statically check SQL for syntax and type errors",
	RunE: func(cmd *cobra.Command, args []string) error {
		blob, err := ioutil.ReadFile("sqlc.json")
		if err != nil {
			return err
		}

		var settings dinosql.GenerateSettings
		if err := json.Unmarshal(blob, &settings); err != nil {
			return err
		}

		for _, pkg := range settings.Packages {
			c, err := dinosql.ParseCatalog(pkg.Schema)
			if err != nil {
				return err
			}
			if _, err := dinosql.ParseQueries(c, settings, pkg); err != nil {
				return err
			}
		}
		return nil
	},
}
