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
	"time"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/ext"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/cobra"

	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/debug"
	"github.com/kyleconroy/sqlc/internal/opts"
	"github.com/kyleconroy/sqlc/internal/plugin"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
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
		ext.Strings(ext.StringsVersion(1)),
		cel.Types(
			&plugin.VetConfig{},
			&plugin.VetQuery{},
		),
		cel.Variable("query",
			cel.ObjectType("plugin.VetQuery"),
		),
		cel.Variable("config",
			cel.ObjectType("plugin.VetConfig"),
		),
	)
	if err != nil {
		return fmt.Errorf("new env: %s", err)
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

	dbenv, err := cel.NewEnv(
		cel.StdLib(),
		ext.Strings(ext.StringsVersion(1)),
		cel.Variable("env",
			cel.MapType(cel.StringType, cel.StringType),
		),
	)
	if err != nil {
		return fmt.Errorf("new dbenv; %s", err)
	}

	c := checker{
		Checks: checks,
		Conf:   conf,
		Dbenv:  dbenv,
		Dir:    dir,
		Env:    env,
		Envmap: map[string]string{},
		Msgs:   msgs,
		Stderr: stderr,
	}
	errored := false
	for _, sql := range conf.SQL {
		if err := c.checkSQL(ctx, sql); err != nil {
			if !errors.Is(err, ErrFailedChecks) {
				fmt.Fprintf(stderr, "%s\n", err)
			}
			errored = true
		}
	}
	if errored {
		return ErrFailedChecks
	}
	return nil
}

type checker struct {
	Checks map[string]cel.Program
	Conf   *config.Config
	Dbenv  *cel.Env
	Dir    string
	Env    *cel.Env
	Envmap map[string]string
	Msgs   map[string]string
	Stderr io.Writer
}

// Determine if a query can be prepared based on the engine and the statement
// type.
func prepareable(sql config.SQL, raw *ast.RawStmt) bool {
	if sql.Engine == config.EnginePostgreSQL {
		// TOOD: Add support for MERGE and VALUES stmts
		switch raw.Stmt.(type) {
		case *ast.DeleteStmt:
			return true
		case *ast.InsertStmt:
			return true
		case *ast.SelectStmt:
			return true
		case *ast.UpdateStmt:
			return true
		default:
			return false
		}
	}
	return false
}

func (c *checker) checkSQL(ctx context.Context, sql config.SQL) error {
	// TODO: Create a separate function for this logic so we can
	combo := config.Combine(*c.Conf, sql)

	// TODO: This feels like a hack that will bite us later
	joined := make([]string, 0, len(sql.Schema))
	for _, s := range sql.Schema {
		joined = append(joined, filepath.Join(c.Dir, s))
	}
	sql.Schema = joined

	joined = make([]string, 0, len(sql.Queries))
	for _, q := range sql.Queries {
		joined = append(joined, filepath.Join(c.Dir, q))
	}
	sql.Queries = joined

	var name string
	parseOpts := opts.Parser{
		Debug: debug.Debug,
	}

	result, failed := parse(ctx, name, c.Dir, sql, combo, parseOpts, c.Stderr)
	if failed {
		return ErrFailedChecks
	}

	// TODO: Add MySQL support
	var pgconn *pgx.Conn
	if sql.Engine == config.EnginePostgreSQL && sql.Database != nil {
		ast, issues := c.Dbenv.Compile(sql.Database.URL)
		if issues != nil && issues.Err() != nil {
			return fmt.Errorf("type-check error: database url %s", issues.Err())
		}
		prg, err := c.Dbenv.Program(ast)
		if err != nil {
			return fmt.Errorf("program construction error: database url %s", err)
		}
		// Populate the environment variable map if it is empty
		if len(c.Envmap) == 0 {
			for _, e := range os.Environ() {
				k, v, _ := strings.Cut(e, "=")
				c.Envmap[k] = v
			}
		}
		out, _, err := prg.Eval(map[string]any{
			"env": c.Envmap,
		})
		if err != nil {
			return fmt.Errorf("expression error: %s", err)
		}
		dburl, ok := out.Value().(string)
		if !ok {
			return fmt.Errorf("expression returned non-string value: %v", out.Value())
		}
		fmt.Println("URL", dburl)
		conn, err := pgx.Connect(ctx, dburl)
		if err != nil {
			return fmt.Errorf("database: connection error: %s", err)
		}
		defer conn.Close(ctx)
		pgconn = conn
	}

	errored := false
	req := codeGenRequest(result, combo)
	cfg := vetConfig(req)
	for i, query := range req.Queries {
		original := result.Queries[i]
		if pgconn != nil && prepareable(sql, original.RawStmt) {
			name := fmt.Sprintf("sqlc_vet_%d_%d", time.Now().Unix(), i)
			_, err := pgconn.Prepare(ctx, name, query.Text)
			if err != nil {
				fmt.Fprintf(c.Stderr, "%s: error preparing %s: %s\n", query.Filename, query.Name, err)
				errored = true
				continue
			}
		}
		q := vetQuery(query)
		for _, name := range sql.Rules {
			prg, ok := c.Checks[name]
			if !ok {
				return fmt.Errorf("type-check error: a check with the name '%s' does not exist", name)
			}
			out, _, err := prg.Eval(map[string]any{
				"query":  q,
				"config": cfg,
			})
			if err != nil {
				return err
			}
			tripped, ok := out.Value().(bool)
			if !ok {
				return fmt.Errorf("expression returned non-bool value: %v", out.Value())
			}
			if tripped {
				// TODO: Get line numbers in the output
				msg := c.Msgs[name]
				if msg == "" {
					fmt.Fprintf(c.Stderr, "%s: %s: %s\n", query.Filename, q.Name, name)
				} else {
					fmt.Fprintf(c.Stderr, "%s: %s: %s: %s\n", query.Filename, q.Name, name, msg)
				}
				errored = true
			}
		}
	}
	if errored {
		return ErrFailedChecks
	}
	return nil
}

func vetConfig(req *plugin.CodeGenRequest) *plugin.VetConfig {
	return &plugin.VetConfig{
		Version: req.Settings.Version,
		Engine:  req.Settings.Engine,
		Schema:  req.Settings.Schema,
		Queries: req.Settings.Queries,
	}
}

func vetQuery(q *plugin.Query) *plugin.VetQuery {
	var params []*plugin.VetParameter
	for _, p := range q.Params {
		params = append(params, &plugin.VetParameter{
			Number: p.Number,
		})
	}
	return &plugin.VetQuery{
		Sql:    q.Text,
		Name:   q.Name,
		Cmd:    strings.TrimPrefix(":", q.Cmd),
		Params: params,
	}
}
