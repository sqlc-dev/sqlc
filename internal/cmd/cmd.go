package cmd

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/trace"

	"github.com/cubicdaiya/gonp"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	"github.com/kyleconroy/sqlc/internal/codegen/golang"
	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/debug"
	"github.com/kyleconroy/sqlc/internal/info"
	"github.com/kyleconroy/sqlc/internal/tracer"
)

func init() {
	uploadCmd.Flags().BoolP("dry-run", "", false, "dump upload request (default: false)")
}

// Do runs the command logic.
func Do(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	rootCmd := &cobra.Command{Use: "sqlc", SilenceUsage: true}
	rootCmd.PersistentFlags().StringP("file", "f", "", "specify an alternate config file (default: sqlc.yaml)")
	rootCmd.PersistentFlags().BoolP("experimental", "x", false, "enable experimental features (default: false)")

	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(genCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(uploadCmd)

	rootCmd.SetArgs(args)
	rootCmd.SetIn(stdin)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SilenceErrors = true

	ctx := context.Background()
	if debug.Debug.Trace != "" {
		tracectx, cleanup, err := tracer.Start(ctx)
		if err != nil {
			fmt.Printf("failed to start trace: %v\n", err)
			return 1
		}
		ctx = tracectx
		defer cleanup()
	}
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		} else {
			return 1
		}
	}
	return 0
}

var version string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the sqlc version number",
	RunE: func(cmd *cobra.Command, args []string) error {
		defer trace.StartRegion(cmd.Context(), "version").End()
		if version == "" {
			fmt.Printf("%s\n", info.Version)
		} else {
			fmt.Printf("%s\n", version)
		}
		return nil
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create an empty sqlc.yaml settings file",
	RunE: func(cmd *cobra.Command, args []string) error {
		defer trace.StartRegion(cmd.Context(), "init").End()
		file := "sqlc.yaml"
		if f := cmd.Flag("file"); f != nil && f.Changed {
			file = f.Value.String()
			if file == "" {
				return fmt.Errorf("file argument is empty")
			}
		}
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			return nil
		}
		blob, err := yaml.Marshal(config.V1GenerateSettings{Version: "1"})
		if err != nil {
			return err
		}
		return os.WriteFile(file, blob, 0644)
	},
}

type Env struct {
	ExperimentalFeatures bool
	DryRun               bool
}

func ParseEnv(c *cobra.Command) Env {
	x := c.Flag("experimental")
	dr := c.Flag("dry-run")
	return Env{
		ExperimentalFeatures: x != nil && x.Changed,
		DryRun:               dr != nil && dr.Changed,
	}
}

func (e *Env) Validate(cfg *config.Config) error {
	for _, sql := range cfg.SQL {
		if sql.Gen.Go != nil && sql.Gen.Go.SQLPackage == golang.SQLPackagePGXV5 && !e.ExperimentalFeatures {
			return fmt.Errorf("'pgx/v5' golang sql package requires enabled '--experimental' flag")
		}
	}

	return nil
}

func getConfigPath(stderr io.Writer, f *pflag.Flag) (string, string) {
	if f != nil && f.Changed {
		file := f.Value.String()
		if file == "" {
			fmt.Fprintln(stderr, "error parsing config: file argument is empty")
			os.Exit(1)
		}
		abspath, err := filepath.Abs(file)
		if err != nil {
			fmt.Fprintf(stderr, "error parsing config: absolute file path lookup failed: %s\n", err)
			os.Exit(1)
		}
		return filepath.Dir(abspath), filepath.Base(abspath)
	} else {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(stderr, "error parsing sqlc.json: file does not exist")
			os.Exit(1)
		}
		return wd, ""
	}
}

var genCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Go code from SQL",
	RunE: func(cmd *cobra.Command, args []string) error {
		defer trace.StartRegion(cmd.Context(), "generate").End()
		stderr := cmd.ErrOrStderr()
		dir, name := getConfigPath(stderr, cmd.Flag("file"))
		output, err := Generate(cmd.Context(), ParseEnv(cmd), dir, name, stderr)
		if err != nil {
			return err
		}
		defer trace.StartRegion(cmd.Context(), "writefiles").End()
		for filename, source := range output {
			os.MkdirAll(filepath.Dir(filename), 0755)
			if err := os.WriteFile(filename, []byte(source), 0644); err != nil {
				fmt.Fprintf(stderr, "%s: %s\n", filename, err)
				return err
			}
		}
		return nil
	},
}

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload the schema, queries, and configuration for this project",
	RunE: func(cmd *cobra.Command, args []string) error {
		stderr := cmd.ErrOrStderr()
		dir, name := getConfigPath(stderr, cmd.Flag("file"))
		if err := createPkg(cmd.Context(), ParseEnv(cmd), dir, name, stderr); err != nil {
			fmt.Fprintf(stderr, "error uploading: %s\n", err)
			os.Exit(1)
		}
		return nil
	},
}

var checkCmd = &cobra.Command{
	Use:   "compile",
	Short: "Statically check SQL for syntax and type errors",
	RunE: func(cmd *cobra.Command, args []string) error {
		defer trace.StartRegion(cmd.Context(), "compile").End()
		stderr := cmd.ErrOrStderr()
		dir, name := getConfigPath(stderr, cmd.Flag("file"))
		if _, err := Generate(cmd.Context(), ParseEnv(cmd), dir, name, stderr); err != nil {
			os.Exit(1)
		}
		return nil
	},
}

func getLines(f []byte) []string {
	fp := bytes.NewReader(f)
	scanner := bufio.NewScanner(fp)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func filterHunks[T gonp.Elem](uniHunks []gonp.UniHunk[T]) []gonp.UniHunk[T] {
	var out []gonp.UniHunk[T]
	for i, uniHunk := range uniHunks {
		var changed bool
		for _, e := range uniHunk.GetChanges() {
			switch e.GetType() {
			case gonp.SesDelete:
				changed = true
			case gonp.SesAdd:
				changed = true
			}
		}
		if changed {
			out = append(out, uniHunks[i])
		}
	}
	return out
}

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Compare the generated files to the existing files",
	RunE: func(cmd *cobra.Command, args []string) error {
		defer trace.StartRegion(cmd.Context(), "diff").End()
		stderr := cmd.ErrOrStderr()
		dir, name := getConfigPath(stderr, cmd.Flag("file"))
		if err := Diff(cmd.Context(), ParseEnv(cmd), dir, name, stderr); err != nil {
			os.Exit(1)
		}
		return nil
	},
}
