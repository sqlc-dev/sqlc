package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kyleconroy/sqlc/internal/dinosql"
	"github.com/kyleconroy/sqlc/internal/mysql"

	"github.com/davecgh/go-spew/spew"
	pg "github.com/lfittl/pg_query_go"
	"github.com/spf13/cobra"
)

// Do runs the command logic.
func Do(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	rootCmd := &cobra.Command{Use: "sqlc", SilenceUsage: true}
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(unstable__mysql)
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
		blob, err := json.MarshalIndent(dinosql.GenerateSettings{Version: "1"}, "", "  ")
		if err != nil {
			return err
		}
		return ioutil.WriteFile("sqlc.json", blob, 0644)
	},
}

const errMessageNoVersion = `The configuration file must have a version number.
Set the version to 1 at the top of sqlc.json:

{
  "version": "1"
  ...
}
`

const errMessageUnknownVersion = `The configuration file has an invalid version number.
The only supported version is "1".
`

const errMessageNoPackages = `No packages are configured`

var genCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Go code from SQL",
	Run: func(cmd *cobra.Command, args []string) {
		blob, err := ioutil.ReadFile("sqlc.json")
		if err != nil {
			fmt.Fprintln(os.Stderr, "error parsing sqlc.json: file does not exist")
			os.Exit(1)
		}

		settings, err := dinosql.ParseConfigFile(bytes.NewReader(blob))
		if err != nil {
			switch err {
			case dinosql.ErrMissingVersion:
				fmt.Fprintf(os.Stderr, errMessageNoVersion)
			case dinosql.ErrUnknownVersion:
				fmt.Fprintf(os.Stderr, errMessageUnknownVersion)
			case dinosql.ErrNoPackages:
				fmt.Fprintf(os.Stderr, errMessageNoPackages)
			}
			fmt.Fprintf(os.Stderr, "error parsing sqlc.json: %s\n", err)
			os.Exit(1)
		}

		var errored bool

		output := map[string]string{}

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

			q, err := dinosql.ParseQueries(c, pkg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "# package %s\n", name)
				if parserErr, ok := err.(*dinosql.ParserErr); ok {
					for _, fileErr := range parserErr.Errs {
						fmt.Fprintf(os.Stderr, "%s:%d:%d: %s\n", fileErr.Filename, fileErr.Line, fileErr.Column, fileErr.Err)
					}
				} else {
					fmt.Fprintf(os.Stderr, "error parsing queries: %s\n", err)
				}
				errored = true
				continue
			}

			files, err := dinosql.Generate(q, settings)
			if err != nil {
				fmt.Fprintf(os.Stderr, "# package %s\n", name)
				fmt.Fprintf(os.Stderr, "error generating code: %s\n", err)
				errored = true
				continue
			}

			for n, source := range files {
				filename := filepath.Join(pkg.Path, n)
				output[filename] = source
			}
		}

		if errored {
			os.Exit(1)
		}

		for filename, source := range output {
			os.MkdirAll(filepath.Dir(filename), 0755)
			if err := ioutil.WriteFile(filename, []byte(source), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", filename, err)
				os.Exit(1)
			}
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
			if _, err := dinosql.ParseQueries(c, pkg); err != nil {
				return err
			}
		}
		return nil
	},
}

var unstable__mysql = &cobra.Command{
	Use:   "unstable__mysql generate",
	Short: "Generate MySQL Queries into typesafe Go code",
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
			res, err := mysql.GeneratePkg(pkg.Name, pkg.Queries, settings)
			if err != nil {
				return err
			}
			for filename, source := range res {
				os.MkdirAll(filepath.Dir(filename), 0755)
				if err := ioutil.WriteFile(filename, []byte(source), 0644); err != nil {
					fmt.Fprintf(os.Stderr, "%s: %s\n", filename, err)
					os.Exit(1)
				}
			}
		}
		return nil
	},
}
