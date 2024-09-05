package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/trace"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/cobra"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/dbmanager"
	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/opts"
	"github.com/sqlc-dev/sqlc/internal/plugin"
	"github.com/sqlc-dev/sqlc/internal/shfmt"
	"gopkg.in/yaml.v3"
)

func NewCmdExplain() *cobra.Command {
	return &cobra.Command{
		Use:   "explain",
		Short: "Explain queries",
		RunE: func(cmd *cobra.Command, args []string) error {
			defer trace.StartRegion(cmd.Context(), "vet").End()
			stderr := cmd.ErrOrStderr()
			opts := &Options{
				Env:    ParseEnv(cmd),
				Stderr: stderr,
			}
			dir, name := getConfigPath(stderr, cmd.Flag("file"))
			if err := Explain(cmd.Context(), dir, name, opts); err != nil {
				if !errors.Is(err, ErrFailedChecks) {
					fmt.Fprintf(stderr, "%s\n", err)
				}
				os.Exit(1)
			}
			return nil
		},
	}
}

func Explain(ctx context.Context, dir, filename string, opts *Options) error {
	e := opts.Env
	stderr := opts.Stderr
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

	c := rawExplainer{
		Conf:          conf,
		Dir:           dir,
		Stderr:        stderr,
		OnlyManagedDB: e.Debug.OnlyManagedDatabases,
		Replacer:      shfmt.NewReplacer(nil),
	}
	var errs error
	for _, sql := range conf.SQL {
		if err := c.explainSQL(ctx, sql); err != nil {
			fmt.Fprintf(stderr, "%s\n", err)
			errs = errors.Join(errs, err)
		}
	}

	return errs
}

type rawExplainer struct {
	Conf          *config.Config
	Dir           string
	Stderr        io.Writer
	OnlyManagedDB bool
	Client        dbmanager.Client
	Replacer      *shfmt.Replacer
}

func (c *rawExplainer) fetchDatabaseUri(ctx context.Context, s config.SQL) (string, func() error, error) {
	return (&checker{
		Conf:          c.Conf,
		Dir:           c.Dir,
		Stderr:        c.Stderr,
		OnlyManagedDB: c.OnlyManagedDB,
		Client:        c.Client,
		Replacer:      c.Replacer,
	}).fetchDatabaseUri(ctx, s)
}

func (c *rawExplainer) explainSQL(ctx context.Context, s config.SQL) error {
	// TODO: Create a separate function for this logic so we can
	combo := config.Combine(*c.Conf, s)

	// TODO: This feels like a hack that will bite us later
	joined := make([]string, 0, len(s.Schema))
	for _, s := range s.Schema {
		joined = append(joined, filepath.Join(c.Dir, s))
	}
	s.Schema = joined

	joined = make([]string, 0, len(s.Queries))
	for _, q := range s.Queries {
		joined = append(joined, filepath.Join(c.Dir, q))
	}
	s.Queries = joined

	var name string
	parseOpts := opts.Parser{
		Debug: debug.Debug,
	}

	result, failed := parse(ctx, name, c.Dir, s, combo, parseOpts, c.Stderr)
	if failed {
		return ErrFailedChecks
	}

	var expl rawDBExplainer
	if s.Database != nil { // TODO only set up a database connection if a rule evaluation requires it
		if s.Database.URI != "" && c.OnlyManagedDB {
			return fmt.Errorf("database: connections disabled via SQLCDEBUG=databases=managed")
		}
		dburl, cleanup, err := c.fetchDatabaseUri(ctx, s)
		if err != nil {
			return err
		}
		defer func() {
			if err := cleanup(); err != nil {
				fmt.Fprintf(c.Stderr, "error cleaning up: %s\n", err)
			}
		}()

		switch s.Engine {
		case config.EnginePostgreSQL:
			conn, err := pgx.Connect(ctx, dburl)
			if err != nil {
				return fmt.Errorf("database: connection error: %s", err)
			}
			if err := conn.Ping(ctx); err != nil {
				return fmt.Errorf("database: connection error: %s", err)
			}
			defer conn.Close(ctx)

			expl = &rawPostgresExplainer{c: conn}
		case config.EngineMySQL:
			db, err := sql.Open("mysql", dburl)
			if err != nil {
				return fmt.Errorf("database: connection error: %s", err)
			}
			if err := db.PingContext(ctx); err != nil {
				return fmt.Errorf("database: connection error: %s", err)
			}
			defer db.Close()
			expl = &rawMySQLExplainer{db}
		case config.EngineSQLite:
			db, err := sql.Open("sqlite", dburl)
			if err != nil {
				return fmt.Errorf("database: connection error: %s", err)
			}
			if err := db.PingContext(ctx); err != nil {
				return fmt.Errorf("database: connection error: %s", err)
			}
			defer db.Close()
			// SQLite really doesn't want us to depend on the output of EXPLAIN
			// QUERY PLAN: https://www.sqlite.org/eqp.html
			expl = nil
		default:
			return fmt.Errorf("unsupported database uri: %s", s.Engine)
		}

		req := codeGenRequest(result, combo)
		for _, query := range req.Queries {
			if expl == nil {
				fmt.Fprintf(c.Stderr, "%s: %s: %s: error explaining query: database connection required\n", query.Filename, query.Name, name)
				continue
			}
			results, err := expl.Explain(ctx, query.Text, query.Params...)
			if err != nil {
				fmt.Fprintf(c.Stderr, "%s: %s: %s: error explaining query: %s\n", query.Filename, query.Name, name, err)
			}

			err = yaml.NewEncoder(os.Stdout).Encode([]struct {
				Name      string
				Query     string
				Arguments []any
				Output    interface{}
			}{{
				Name:      query.Name,
				Query:     query.Text,
				Arguments: expl.DefaultValues(query.Params),
				Output:    string(results),
			}})
			if err != nil {
				fmt.Fprintf(c.Stderr, "%s: %s: %s: fail marshal results: %s\n", query.Filename, query.Name, name, err)
			}

		}
	}
	return nil
}

type rawDBExplainer interface {
	DefaultValues([]*plugin.Parameter) []any
	Explain(context.Context, string, ...*plugin.Parameter) ([]byte, error)
}

type rawPostgresExplainer struct {
	c *pgx.Conn
}

func (p *rawPostgresExplainer) DefaultValues(args []*plugin.Parameter) []any {
	eArgs := make([]any, len(args))
	for i, a := range args {
		eArgs[i] = pgDefaultValue(a.Column)
	}
	return eArgs
}

func (p *rawPostgresExplainer) Explain(ctx context.Context, query string, args ...*plugin.Parameter) ([]byte, error) {
	eQuery := "EXPLAIN " + query
	eArgs := p.DefaultValues(args)
	var results []byte

	rows, err := p.c.Query(ctx, eQuery, eArgs...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var result []byte
		err := rows.Scan(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, append(result, '\n')...)
	}
	return results, nil
}

type rawMySQLExplainer struct {
	*sql.DB
}

func (me *rawMySQLExplainer) DefaultValues(args []*plugin.Parameter) []any {
	eArgs := make([]any, len(args))
	for i, a := range args {
		eArgs[i] = mysqlDefaultValue(a.Column)
	}
	return eArgs
}
func (me *rawMySQLExplainer) Explain(ctx context.Context, query string, args ...*plugin.Parameter) ([]byte, error) {
	eQuery := "EXPLAIN " + query
	eArgs := me.DefaultValues(args)
	var results []byte
	rows, err := me.QueryContext(ctx, eQuery, eArgs...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var result []byte
		err := rows.Scan(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, append(result, '\n')...)
	}
	return results, nil
}
