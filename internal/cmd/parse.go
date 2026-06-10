package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sqlc-dev/sqlc/internal/engine/clickhouse"
	"github.com/sqlc-dev/sqlc/internal/engine/dolphin"
	"github.com/sqlc-dev/sqlc/internal/engine/postgresql"
	"github.com/sqlc-dev/sqlc/internal/engine/sqlite"
	"github.com/sqlc-dev/sqlc/internal/metadata"
	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// dialectParser is the subset of the engine parsers that the parse command
// needs: parsing SQL into statements and reporting the dialect's comment syntax
// (used to extract the sqlc query name and command).
type dialectParser interface {
	Parse(io.Reader) ([]ast.Statement, error)
	CommentSyntax() source.CommentSyntax
}

// parsedStatement is the JSON representation of a single parsed statement. The
// name and cmd are extracted from the sqlc query annotation (e.g.
// "-- name: GetAuthor :one") and are omitted when the statement has none.
type parsedStatement struct {
	Name string       `json:"name,omitempty"`
	Cmd  string       `json:"cmd,omitempty"`
	AST  *ast.RawStmt `json:"ast"`
}

var parseCmd = &cobra.Command{
	Use:   "parse [file]",
	Short: "Parse SQL and output the AST as JSON",
	Long: `Parse SQL from a file or stdin and output the abstract syntax tree as JSON.

Each statement is reported with its sqlc query name and command (when the
statement carries a "-- name:" annotation) alongside the AST.

Examples:
  # Parse a SQL file with PostgreSQL dialect
  sqlc parse --dialect postgresql schema.sql

  # Parse from stdin with MySQL dialect
  echo "SELECT * FROM users" | sqlc parse --dialect mysql

  # Parse SQLite SQL
  sqlc parse --dialect sqlite queries.sql

  # Parse ClickHouse SQL
  sqlc parse --dialect clickhouse queries.sql`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dialect, err := cmd.Flags().GetString("dialect")
		if err != nil {
			return err
		}
		if dialect == "" {
			return fmt.Errorf("--dialect flag is required (postgresql, mysql, sqlite, or clickhouse)")
		}

		// Determine input source
		var input io.Reader
		if len(args) == 1 {
			file, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer file.Close()
			input = file
		} else {
			// Check if stdin has data
			stat, err := os.Stdin.Stat()
			if err != nil {
				return fmt.Errorf("failed to stat stdin: %w", err)
			}
			if (stat.Mode() & os.ModeCharDevice) != 0 {
				return fmt.Errorf("no input provided. Specify a file path or pipe SQL via stdin")
			}
			input = cmd.InOrStdin()
		}

		// Select the parser for the requested dialect
		var parser dialectParser
		switch dialect {
		case "postgresql", "postgres", "pg":
			parser = postgresql.NewParser()
		case "mysql":
			parser = dolphin.NewParser()
		case "sqlite":
			parser = sqlite.NewParser()
		case "clickhouse":
			parser = clickhouse.NewParser()
		default:
			return fmt.Errorf("unsupported dialect: %s (use postgresql, mysql, sqlite, or clickhouse)", dialect)
		}

		// Read the full source so each statement's name and command can be
		// extracted from its annotation comment.
		src, err := io.ReadAll(input)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		stmts, err := parser.Parse(strings.NewReader(string(src)))
		if err != nil {
			return fmt.Errorf("parse error: %w", err)
		}

		commentSyntax := metadata.CommentSyntax(parser.CommentSyntax())

		// Output the AST as a single JSON document
		out := make([]parsedStatement, 0, len(stmts))
		for _, stmt := range stmts {
			ps := parsedStatement{AST: stmt.Raw}
			rawSQL, err := source.Pluck(string(src), stmt.Raw.StmtLocation, stmt.Raw.StmtLen)
			if err != nil {
				return fmt.Errorf("failed to read statement source: %w", err)
			}
			name, cmd, err := metadata.ParseQueryNameAndType(rawSQL, commentSyntax)
			if err != nil {
				return fmt.Errorf("failed to parse query annotation: %w", err)
			}
			ps.Name = name
			ps.Cmd = cmd
			out = append(out, ps)
		}

		stdout := cmd.OutOrStdout()
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(out); err != nil {
			return fmt.Errorf("failed to encode AST: %w", err)
		}

		return nil
	},
}
