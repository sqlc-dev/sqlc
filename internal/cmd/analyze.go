package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sqlc-dev/sqlc/internal/compiler"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/multierr"
	"github.com/sqlc-dev/sqlc/internal/opts"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze [query-file]",
	Short: "Analyze a query against a schema and output the result columns and parameters",
	Long: `Analyze a query file against a schema file and output the inferred result
columns and parameters as JSON.

Unlike "sqlc generate", this command does not require a configuration file and
does not connect to a database. It uses sqlc's native static analysis to infer
types from the provided schema.

Examples:
  # Analyze a PostgreSQL query
  sqlc analyze --dialect postgresql --schema schema.sql query.sql

  # Analyze a MySQL query
  sqlc analyze --dialect mysql --schema schema.sql query.sql

  # Analyze a SQLite query
  sqlc analyze --dialect sqlite --schema schema.sql query.sql

  # Analyze a query piped via stdin
  echo "-- name: GetAuthor :one
  SELECT * FROM authors WHERE id = $1;" | sqlc analyze --dialect postgresql --schema schema.sql`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dialect, err := cmd.Flags().GetString("dialect")
		if err != nil {
			return err
		}
		if dialect == "" {
			return fmt.Errorf("--dialect flag is required (postgresql, mysql, or sqlite)")
		}

		schemaPath, err := cmd.Flags().GetString("schema")
		if err != nil {
			return err
		}
		if schemaPath == "" {
			return fmt.Errorf("--schema flag is required")
		}

		// The query comes from a file argument or, when none is given, from
		// stdin. The compiler reads queries from files, so stdin is written to
		// a temporary file.
		var queryPath string
		if len(args) == 1 {
			queryPath = args[0]
		} else {
			stat, err := os.Stdin.Stat()
			if err != nil {
				return fmt.Errorf("failed to stat stdin: %w", err)
			}
			if (stat.Mode() & os.ModeCharDevice) != 0 {
				return fmt.Errorf("no query provided. Specify a query file or pipe SQL via stdin")
			}
			data, err := io.ReadAll(cmd.InOrStdin())
			if err != nil {
				return fmt.Errorf("failed to read stdin: %w", err)
			}
			tmp, err := os.CreateTemp("", "sqlc-analyze-*.sql")
			if err != nil {
				return fmt.Errorf("failed to create temp file: %w", err)
			}
			defer os.Remove(tmp.Name())
			if _, err := tmp.Write(data); err != nil {
				tmp.Close()
				return fmt.Errorf("failed to write temp file: %w", err)
			}
			if err := tmp.Close(); err != nil {
				return fmt.Errorf("failed to close temp file: %w", err)
			}
			queryPath = tmp.Name()
		}

		var engine config.Engine
		switch dialect {
		case "postgresql", "postgres", "pg":
			engine = config.EnginePostgreSQL
		case "mysql":
			engine = config.EngineMySQL
		case "sqlite":
			engine = config.EngineSQLite
		default:
			return fmt.Errorf("unsupported dialect: %s (use postgresql, mysql, or sqlite)", dialect)
		}

		sql := config.SQL{
			Engine:  engine,
			Schema:  config.Paths{schemaPath},
			Queries: config.Paths{queryPath},
		}
		combo := config.Combine(config.Config{}, sql)
		parserOpts := opts.Parser{}

		ctx := cmd.Context()
		c, err := compiler.NewCompiler(sql, combo, parserOpts)
		if err != nil {
			return fmt.Errorf("error creating compiler: %w", err)
		}
		defer c.Close(ctx)

		if err := c.ParseCatalog(sql.Schema); err != nil {
			return fmt.Errorf("error parsing schema: %w", formatParseError(err))
		}
		if err := c.ParseQueries(sql.Queries, parserOpts); err != nil {
			return fmt.Errorf("error parsing queries: %w", formatParseError(err))
		}

		result := c.Result()

		out := make([]analyzedQuery, 0, len(result.Queries))
		for _, q := range result.Queries {
			out = append(out, newAnalyzedQuery(q))
		}

		stdout := cmd.OutOrStdout()
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(out); err != nil {
			return fmt.Errorf("failed to encode analysis: %w", err)
		}

		return nil
	},
}

// formatParseError unwraps a multierr.Error into a single error containing all
// of the underlying file errors, so the analyze command can report each one with
// its file location.
func formatParseError(err error) error {
	parserErr, ok := err.(*multierr.Error)
	if !ok {
		return err
	}
	var msgs []string
	for _, fileErr := range parserErr.Errs() {
		msgs = append(msgs, fmt.Sprintf("%s:%d:%d: %s",
			fileErr.Filename, fileErr.Line, fileErr.Column, fileErr.Err))
	}
	if len(msgs) == 0 {
		return err
	}
	return fmt.Errorf("%s", strings.Join(msgs, "; "))
}

type analyzedQuery struct {
	Name    string           `json:"name"`
	Cmd     string           `json:"cmd"`
	Columns []analyzedColumn `json:"columns"`
	Params  []analyzedParam  `json:"params"`
}

type analyzedColumn struct {
	Name     string `json:"name"`
	DataType string `json:"data_type"`
	NotNull  bool   `json:"not_null"`
	IsArray  bool   `json:"is_array"`
	Table    string `json:"table,omitempty"`
}

type analyzedParam struct {
	Number int            `json:"number"`
	Column analyzedColumn `json:"column"`
}

func newAnalyzedQuery(q *compiler.Query) analyzedQuery {
	aq := analyzedQuery{
		Name:    q.Metadata.Name,
		Cmd:     q.Metadata.Cmd,
		Columns: make([]analyzedColumn, 0, len(q.Columns)),
		Params:  make([]analyzedParam, 0, len(q.Params)),
	}
	for _, col := range q.Columns {
		aq.Columns = append(aq.Columns, newAnalyzedColumn(col))
	}
	for _, p := range q.Params {
		aq.Params = append(aq.Params, analyzedParam{
			Number: p.Number,
			Column: newAnalyzedColumn(p.Column),
		})
	}
	return aq
}

func newAnalyzedColumn(col *compiler.Column) analyzedColumn {
	if col == nil {
		return analyzedColumn{}
	}
	ac := analyzedColumn{
		Name:     col.Name,
		DataType: col.DataType,
		NotNull:  col.NotNull,
		IsArray:  col.IsArray,
	}
	if col.Table != nil {
		ac.Table = col.Table.Name
	}
	return ac
}
