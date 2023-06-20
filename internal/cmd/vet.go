package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/trace"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/spf13/cobra"

	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/debug"
	"github.com/kyleconroy/sqlc/internal/opts"
	"github.com/kyleconroy/sqlc/internal/plugin"
)

var ErrFailedChecks = errors.New("failed checks")

func NewCmdVet() *cobra.Command {
	return &cobra.Command{
		Use:   "vet",
		Short: "Vet examines queries",
		RunE: func(cmd *cobra.Command, args []string) error {
			defer trace.StartRegion(cmd.Context(), "vet").End()
			stderr := cmd.ErrOrStderr()
			dir, name := getConfigPath(stderr, cmd.Flag("file"))
			if err := Vet(cmd.Context(), ParseEnv(cmd), dir, name, stderr); err != nil {
				if !errors.Is(err, ErrFailedChecks) {
					fmt.Fprintf(stderr, "%s\n", err)
				}
				os.Exit(1)
			}
			return nil
		},
	}
}

func Vet(ctx context.Context, e Env, dir, filename string, stderr io.Writer) error {
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
		cel.Types(&plugin.VetQuery{}),
		cel.Variable("query",
			cel.ObjectType("plugin.VetQuery"),
		),
	)
	if err != nil {
		return fmt.Errorf("new env; %s", err)
	}

	checks := map[string]cel.Program{}
	msgs := map[string]string{}

	for _, c := range conf.Rules {
		if c.Name == "" {
			return fmt.Errorf("checks require a name")
		}
		if _, found := checks[c.Name]; found {
			return fmt.Errorf("type-check error: a check with the name '%s' already exists", c.Name)
		}
		if c.Rule == "" {
			return fmt.Errorf("type-check error: %s is empty", c.Name)
		}
		ast, issues := env.Compile(c.Rule)
		if issues != nil && issues.Err() != nil {
			return fmt.Errorf("type-check error: %s %s", c.Name, issues.Err())
		}
		prg, err := env.Program(ast)
		if err != nil {
			return fmt.Errorf("program construction error: %s %s", c.Name, err)
		}
		checks[c.Name] = prg
		msgs[c.Name] = c.Msg
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

		result, failed := parse(ctx, name, dir, sql, combo, parseOpts, stderr)
		if failed {
			return nil
		}
		req := codeGenRequest(result, combo)
		for _, q := range vetQueries(req) {
			for _, name := range sql.Rules {
				prg, ok := checks[name]
				if !ok {
					return fmt.Errorf("type-check error: a check with the name '%s' does not exist", name)
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
					// TODO: Get line numbers in the output
					msg := msgs[name]
					if msg == "" {
						fmt.Fprintf(stderr, q.Path+": %s: %s\n", q.Name, name, msg)
					} else {
						fmt.Fprintf(stderr, q.Path+": %s: %s: %s\n", q.Name, name, msg)
					}
					errored = true
				}
			}
		}
	}
	if errored {
		return ErrFailedChecks
	}
	return nil
}

func vetQueries(req *plugin.CodeGenRequest) []*plugin.VetQuery {
	var out []*plugin.VetQuery
	for _, q := range req.Queries {
		var params []*plugin.VetParameter
		for _, p := range q.Params {
			params = append(params, &plugin.VetParameter{
				Number: p.Number,
			})
		}
		out = append(out, &plugin.VetQuery{
			Sql:    q.Text,
			Name:   q.Name,
			Cmd:    strings.TrimPrefix(":", q.Cmd),
			Engine: req.Settings.Engine,
			Params: params,
			Path:   q.Filename,
		})
	}
	return out
}
