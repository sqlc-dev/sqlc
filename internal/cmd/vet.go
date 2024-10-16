package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sqlc-dev/sqlc/internal/constants"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime/trace"
	"slices"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/ext"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/dbmanager"
	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/migrations"
	"github.com/sqlc-dev/sqlc/internal/opts"
	"github.com/sqlc-dev/sqlc/internal/plugin"
	"github.com/sqlc-dev/sqlc/internal/quickdb"
	"github.com/sqlc-dev/sqlc/internal/shfmt"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
	"github.com/sqlc-dev/sqlc/internal/vet"
)

var ErrFailedChecks = errors.New("failed checks")

var pjson = protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}

func NewCmdVet() *cobra.Command {
	return &cobra.Command{
		Use:   "vet",
		Short: "Vet examines queries",
		RunE: func(cmd *cobra.Command, args []string) error {
			defer trace.StartRegion(cmd.Context(), "vet").End()
			stderr := cmd.ErrOrStderr()
			opts := &Options{
				Env:    ParseEnv(cmd),
				Stderr: stderr,
			}
			dir, name := getConfigPath(stderr, cmd.Flag("file"))
			if err := Vet(cmd.Context(), dir, name, opts); err != nil {
				if !errors.Is(err, ErrFailedChecks) {
					fmt.Fprintf(stderr, "%s\n", err)
				}
				os.Exit(1)
			}
			return nil
		},
	}
}

func Vet(ctx context.Context, dir, filename string, opts *Options) error {
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

	env, err := cel.NewEnv(
		cel.StdLib(),
		ext.Strings(ext.StringsVersion(1)),
		cel.Types(
			&vet.Config{},
			&vet.Query{},
			&vet.PostgreSQL{},
			&vet.MySQL{},
		),
		cel.Variable("query",
			cel.ObjectType("vet.Query"),
		),
		cel.Variable("config",
			cel.ObjectType("vet.Config"),
		),
		cel.Variable("postgresql",
			cel.ObjectType("vet.PostgreSQL"),
		),
		cel.Variable("mysql",
			cel.ObjectType("vet.MySQL"),
		),
	)
	if err != nil {
		return fmt.Errorf("new CEL env error: %s", err)
	}

	rules := map[string]rule{
		constants.QueryRuleDbPrepare: {NeedsPrepare: true},
	}

	for _, c := range conf.Rules {
		if c.Name == "" {
			return fmt.Errorf("rules require a name")
		}
		if _, found := rules[c.Name]; found {
			return fmt.Errorf("type-check error: a rule with the name '%s' already exists", c.Name)
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
		rule := rule{Program: &prg, Message: c.Msg}

		// TODO There's probably a nicer way to do this from the ast
		// https://pkg.go.dev/github.com/google/cel-go/common/ast#AllMatcher
		if strings.Contains(c.Rule, "postgresql.explain") ||
			strings.Contains(c.Rule, "mysql.explain") {
			rule.NeedsExplain = true
		}

		rules[c.Name] = rule
	}

	c := checker{
		Rules:         rules,
		Conf:          conf,
		Dir:           dir,
		Env:           env,
		Stderr:        stderr,
		OnlyManagedDB: e.Debug.OnlyManagedDatabases,
		Replacer:      shfmt.NewReplacer(nil),
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

type preparer interface {
	Prepare(context.Context, string, string) error
}

type pgxConn struct {
	c *pgx.Conn
}

func (p *pgxConn) Prepare(ctx context.Context, name, query string) error {
	_, err := p.c.Prepare(ctx, name, query)
	return err
}

// Return a default value for a PostgreSQL column based on its type. Returns nil
// if the type is unknown.
func pgDefaultValue(col *plugin.Column) any {
	if col == nil {
		return nil
	}
	if col.Type == nil {
		return nil
	}
	typname := strings.TrimPrefix(col.Type.Name, "pg_catalog.")
	switch typname {
	case "any", "void":
		if col.IsArray {
			return []any{}
		} else {
			return nil
		}
	case "anyarray":
		return []any{}
	case "bool", "boolean":
		if col.IsArray {
			return []bool{}
		} else {
			return false
		}
	case "double", "double precision", "real":
		if col.IsArray {
			return []float32{}
		} else {
			return 0.1
		}
	case "json", "jsonb":
		if col.IsArray {
			return []string{}
		} else {
			return "{}"
		}
	case "citext", "string", "text", "varchar":
		if col.IsArray {
			return []string{}
		} else {
			return ""
		}
	case "bigint", "bigserial", "integer", "int", "int2", "int4", "int8", "serial":
		if col.IsArray {
			return []int{}
		} else {
			return 1
		}
	case "date", "time", "timestamp", "timestamptz":
		if col.IsArray {
			return []time.Time{}
		} else {
			return time.Time{}
		}
	case "uuid":
		if col.IsArray {
			return []string{}
		} else {
			return "00000000-0000-0000-0000-000000000000"
		}
	case "numeric", "decimal":
		if col.IsArray {
			return []string{}
		} else {
			return "0.1"
		}
	case "inet":
		if col.IsArray {
			return []string{}
		} else {
			return "192.168.0.1/24"
		}
	case "cidr":
		if col.IsArray {
			return []string{}
		} else {
			return "192.168.1/24"
		}
	default:
		return nil
	}
}

// Return a default value for a MySQL column based on its type. Returns nil
// if the type is unknown.
func mysqlDefaultValue(col *plugin.Column) any {
	if col == nil {
		return nil
	}
	if col.Type == nil {
		return nil
	}
	switch col.Type.Name {
	case "any":
		return nil
	case "bool":
		return false
	case "int", "bigint", "mediumint", "smallint", "tinyint", "bit":
		return 1
	case "decimal": // "numeric", "dec", "fixed"
		// No perfect choice here to avoid "Impossible WHERE" but I think
		// 0.1 is decent. It works for all cases where `scale` > 0 which
		// should be the majority. For more information refer to
		// https://dev.mysql.com/doc/refman/8.1/en/fixed-point-types.html.
		return 0.1
	case "float", "double":
		return 0.1
	case "date":
		return "0000-00-00"
	case "datetime", "timestamp":
		return "0000-00-00 00:00:00"
	case "time":
		return "00:00:00"
	case "year":
		return "0000"
	case "char", "varchar", "binary", "varbinary", "tinyblob", "blob",
		"mediumblob", "longblob", "tinytext", "text", "mediumtext", "longtext":
		return ""
	case "json":
		return "{}"
	default:
		return nil
	}
}

func (p *pgxConn) Explain(ctx context.Context, query string, args ...*plugin.Parameter) (*vetEngineOutput, error) {
	eQuery := "EXPLAIN (ANALYZE false, VERBOSE, COSTS, SETTINGS, BUFFERS, FORMAT JSON) " + query
	eArgs := make([]any, len(args))
	for i, a := range args {
		eArgs[i] = pgDefaultValue(a.Column)
	}
	row := p.c.QueryRow(ctx, eQuery, eArgs...)
	var result []json.RawMessage
	if err := row.Scan(&result); err != nil {
		return nil, err
	}
	if debug.Debug.DumpExplain {
		fmt.Println(eQuery, "with args", eArgs)
		fmt.Println(string(result[0]))
	}
	var explain vet.PostgreSQLExplain
	if err := pjson.Unmarshal(result[0], &explain); err != nil {
		return nil, err
	}
	return &vetEngineOutput{PostgreSQL: &vet.PostgreSQL{Explain: &explain}}, nil
}

type dbPreparer struct {
	db *sql.DB
}

func (p *dbPreparer) Prepare(ctx context.Context, name, query string) error {
	s, err := p.db.PrepareContext(ctx, query)
	if s != nil {
		s.Close()
	}
	return err
}

type explainer interface {
	Explain(context.Context, string, ...*plugin.Parameter) (*vetEngineOutput, error)
}

type mysqlExplainer struct {
	*sql.DB
}

func (me *mysqlExplainer) Explain(ctx context.Context, query string, args ...*plugin.Parameter) (*vetEngineOutput, error) {
	eQuery := "EXPLAIN FORMAT=JSON " + query
	eArgs := make([]any, len(args))
	for i, a := range args {
		eArgs[i] = mysqlDefaultValue(a.Column)
	}
	row := me.QueryRowContext(ctx, eQuery, eArgs...)
	var result json.RawMessage
	if err := row.Scan(&result); err != nil {
		return nil, err
	}
	if debug.Debug.DumpExplain {
		fmt.Println(eQuery, "with args", eArgs)
		fmt.Println(string(result))
	}
	var explain vet.MySQLExplain
	if err := pjson.Unmarshal(result, &explain); err != nil {
		return nil, err
	}
	if explain.QueryBlock.Message != "" {
		return nil, fmt.Errorf("mysql explain: %s", explain.QueryBlock.Message)
	}
	return &vetEngineOutput{MySQL: &vet.MySQL{Explain: &explain}}, nil
}

type rule struct {
	Program      *cel.Program
	Message      string
	NeedsPrepare bool
	NeedsExplain bool
}

type checker struct {
	Rules         map[string]rule
	Conf          *config.Config
	Dir           string
	Env           *cel.Env
	Stderr        io.Writer
	OnlyManagedDB bool
	Client        dbmanager.Client
	Replacer      *shfmt.Replacer
}

func (c *checker) fetchDatabaseUri(ctx context.Context, s config.SQL) (string, func() error, error) {
	cleanup := func() error {
		return nil
	}

	if s.Database == nil {
		panic("fetch database URI called with nil database")
	}
	if !s.Database.Managed {
		uri, err := c.DSN(s.Database.URI)
		return uri, cleanup, err
	}

	if c.Client == nil {
		// FIXME: Eventual race condition
		client := dbmanager.NewClient(c.Conf.Servers)
		c.Client = client
	}

	var ddl []string
	files, err := sqlpath.Glob(s.Schema)
	if err != nil {
		return "", cleanup, err
	}
	for _, schema := range files {
		contents, err := os.ReadFile(schema)
		if err != nil {
			return "", cleanup, fmt.Errorf("read file: %w", err)
		}
		ddl = append(ddl, migrations.RemoveRollbackStatements(string(contents)))
	}

	resp, err := c.Client.CreateDatabase(ctx, &dbmanager.CreateDatabaseRequest{
		Engine:     string(s.Engine),
		Migrations: ddl,
	})
	if err != nil {
		return "", cleanup, fmt.Errorf("managed: create database: %w", err)
	}

	var uri string
	switch s.Engine {
	case config.EngineMySQL:
		dburi, err := quickdb.MySQLReformatURI(resp.Uri)
		if err != nil {
			return "", cleanup, fmt.Errorf("reformat uri: %w", err)
		}
		uri = dburi
	default:
		uri = resp.Uri
	}

	return uri, cleanup, nil
}

func (c *checker) DSN(dsn string) (string, error) {
	return c.Replacer.Replace(dsn), nil
}

func (c *checker) checkSQL(ctx context.Context, s config.SQL) error {
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

	var prep preparer
	var expl explainer
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
			pConn := &pgxConn{conn}
			prep = pConn
			expl = pConn
		case config.EngineMySQL:
			db, err := sql.Open("mysql", dburl)
			if err != nil {
				return fmt.Errorf("database: connection error: %s", err)
			}
			if err := db.PingContext(ctx); err != nil {
				return fmt.Errorf("database: connection error: %s", err)
			}
			defer db.Close()
			prep = &dbPreparer{db}
			expl = &mysqlExplainer{db}
		case config.EngineSQLite:
			db, err := sql.Open("sqlite", dburl)
			if err != nil {
				return fmt.Errorf("database: connection error: %s", err)
			}
			if err := db.PingContext(ctx); err != nil {
				return fmt.Errorf("database: connection error: %s", err)
			}
			defer db.Close()
			prep = &dbPreparer{db}
			// SQLite really doesn't want us to depend on the output of EXPLAIN
			// QUERY PLAN: https://www.sqlite.org/eqp.html
			expl = nil
		default:
			return fmt.Errorf("unsupported database uri: %s", s.Engine)
		}
	}

	errored := false
	req := codeGenRequest(result, combo)
	cfg := vetConfig(req)
	for i, query := range req.Queries {
		md := result.Queries[i].Metadata
		if md.Flags[constants.QueryFlagSqlcVetDisable] {
			// If the vet disable flag is specified without any rules listed, all rules are ignored.
			if len(md.RuleSkiplist) == 0 {
				if debug.Active {
					log.Printf("Skipping all vet rules for query: %s\n", query.Name)
				}
				continue
			}

			// Rules which are listed to be disabled but not declared in the config file are rejected.
			for r := range md.RuleSkiplist {
				if !slices.Contains(s.Rules, r) {
					fmt.Fprintf(c.Stderr, "%s: %s: rule-check error: rule %q does not exist in the config file\n", query.Filename, query.Name, r)
					errored = true
				}
			}
		}

		evalMap := map[string]any{
			"query":  vetQuery(query),
			"config": cfg,
		}

		for _, name := range s.Rules {
			if _, skip := md.RuleSkiplist[name]; skip {
				if debug.Active {
					log.Printf("Skipping vet rule %q for query: %s\n", name, query.Name)
				}
			} else {
				rule, ok := c.Rules[name]
				if !ok {
					return fmt.Errorf("type-check error: a rule with the name '%s' does not exist", name)
				}

				if rule.NeedsPrepare {
					if prep == nil {
						fmt.Fprintf(c.Stderr, "%s: %s: %s: error preparing query: database connection required\n", query.Filename, query.Name, name)
						errored = true
						continue
					}
					prepName := fmt.Sprintf("sqlc_vet_%d_%d", time.Now().Unix(), i)
					if err := prep.Prepare(ctx, prepName, query.Text); err != nil {
						fmt.Fprintf(c.Stderr, "%s: %s: %s: error preparing query: %s\n", query.Filename, query.Name, name, err)
						errored = true
						continue
					}
				}

				// short-circuit for "sqlc/db-prepare" rule which doesn't have a CEL program
				if rule.Program == nil {
					continue
				}

				// Get explain output for this query if we need it
				_, pgsqlOK := evalMap["postgresql"]
				_, mysqlOK := evalMap["mysql"]
				if rule.NeedsExplain && !(pgsqlOK || mysqlOK) {
					if expl == nil {
						fmt.Fprintf(c.Stderr, "%s: %s: %s: error explaining query: database connection required\n", query.Filename, query.Name, name)
						errored = true
						continue
					}
					engineOutput, err := expl.Explain(ctx, query.Text, query.Params...)
					if err != nil {
						fmt.Fprintf(c.Stderr, "%s: %s: %s: error explaining query: %s\n", query.Filename, query.Name, name, err)
						errored = true
						continue
					}

					evalMap["postgresql"] = engineOutput.PostgreSQL
					evalMap["mysql"] = engineOutput.MySQL
				}

				if debug.Debug.DumpVetEnv {
					fmt.Printf("vars for rule '%s' evaluating against query '%s':\n", name, query.Name)
					debug.DumpAsJSON(evalMap)
				}

				out, _, err := (*rule.Program).Eval(evalMap)
				if err != nil {
					return err
				}
				tripped, ok := out.Value().(bool)
				if !ok {
					return fmt.Errorf("expression returned non-bool value: %v", out.Value())
				}
				if tripped {
					// TODO: Get line numbers in the output
					if rule.Message == "" {
						fmt.Fprintf(c.Stderr, "%s: %s: %s\n", query.Filename, query.Name, name)
					} else {
						fmt.Fprintf(c.Stderr, "%s: %s: %s: %s\n", query.Filename, query.Name, name, rule.Message)
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

func vetConfig(req *plugin.GenerateRequest) *vet.Config {
	return &vet.Config{
		Version: req.Settings.Version,
		Engine:  req.Settings.Engine,
		Schema:  req.Settings.Schema,
		Queries: req.Settings.Queries,
	}
}

func vetQuery(q *plugin.Query) *vet.Query {
	var params []*vet.Parameter
	for _, p := range q.Params {
		params = append(params, &vet.Parameter{
			Number: p.Number,
		})
	}
	return &vet.Query{
		Sql:    q.Text,
		Name:   q.Name,
		Cmd:    strings.TrimPrefix(q.Cmd, ":"),
		Params: params,
	}
}

type vetEngineOutput struct {
	PostgreSQL *vet.PostgreSQL
	MySQL      *vet.MySQL
}
