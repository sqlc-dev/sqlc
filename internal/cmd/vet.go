package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/trace"

	"github.com/google/cel-go/cel"
	"github.com/spf13/cobra"

	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/debug"
	"github.com/kyleconroy/sqlc/internal/opts"
	"github.com/kyleconroy/sqlc/internal/plugin"
)

func NewCmdVet() *cobra.Command {
	return &cobra.Command{
		Use:   "vet",
		Short: "Vet examines queries",
		RunE: func(cmd *cobra.Command, args []string) error {
			defer trace.StartRegion(cmd.Context(), "vet").End()
			stderr := cmd.ErrOrStderr()
			dir, name := getConfigPath(stderr, cmd.Flag("file"))
			if err := examine(cmd.Context(), ParseEnv(cmd), dir, name, stderr); err != nil {
				fmt.Fprintf(stderr, "%s\n", err)
				os.Exit(1)
			}
			return nil
		},
	}
}

func examine(ctx context.Context, e Env, dir, filename string, stderr io.Writer) error {
	configPath, conf, err := readConfig(stderr, dir, filename)
	if err != nil {
		return err
	}

	base := filepath.Base(configPath)
	if err := config.Validate(conf); err != nil {
		fmt.Fprintf(stderr, "error validating %s: %s\n", base, err)
		return err
	}

	if err := e.Validate(conf); err != nil {
		fmt.Fprintf(stderr, "error validating %s: %s\n", base, err)
		return err
	}

	env, err := cel.NewEnv(
		cel.StdLib(),
		cel.Types(&plugin.Query{}),
		cel.Variable("query",
			cel.ObjectType("plugin.Query"),
		),
	)
	if err != nil {
		return fmt.Errorf("new env; %s", err)
	}

	checks := map[string]cel.Program{}

	for _, c := range conf.Checks {
		// TODO: Verify check has a name
		// TODO: Verify that check names are unique
		ast, issues := env.Compile(c.Expr)
		if issues != nil && issues.Err() != nil {
			return fmt.Errorf("type-check error: %s %s", c.Name, issues.Err())
		}
		prg, err := env.Program(ast)
		if err != nil {
			return fmt.Errorf("program construction error: %s %s", c.Name, err)
		}
		checks[c.Name] = prg
	}

	errored := true
	for _, sql := range conf.SQL {
		combo := config.Combine(*conf, sql)

		// TODO: This feels like a hack that will bite us later
		joined := make([]string, 0, len(sql.Schema))
		for _, s := range sql.Schema {
			joined = append(joined, filepath.Join(dir, s))
		}
		sql.Schema = joined

		joined = make([]string, 0, len(sql.Queries))
		for _, q := range sql.Queries {
			joined = append(joined, filepath.Join(dir, q))
		}
		sql.Queries = joined

		var name string
		parseOpts := opts.Parser{
			Debug: debug.Debug,
		}

		var errout bytes.Buffer
		result, failed := parse(ctx, name, dir, sql, combo, parseOpts, &errout)
		if failed {
			return nil
		}
		req := codeGenRequest(result, combo)
		for _, q := range req.Queries {
			for _, name := range sql.Checks {
				prg, ok := checks[name]
				if !ok {
					// TODO: Return a helpful error message
					continue
				}
				out, _, err := prg.Eval(map[string]any{
					"query": q,
				})
				if err != nil {
					return err
				}
				tripped, ok := out.Value().(bool)
				if !ok {
					return fmt.Errorf("expression returned non-bool: %s", err)
				}
				if tripped {
					// internal/cmd/vet.go:123:13: fmt.Errorf format %s has arg false of wrong type bool
					fmt.Fprintf(stderr, q.Filename+":17:1: query uses :exec\n")
					errored = true
				}
			}
		}
	}
	if errored {
		return fmt.Errorf("errored")
	}
	return nil
}
