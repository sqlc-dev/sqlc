package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/sqlc-dev/sqlc/internal/engine/clickhouse"
	"github.com/sqlc-dev/sqlc/internal/engine/dolphin"
	"github.com/sqlc-dev/sqlc/internal/engine/postgresql"
	"github.com/sqlc-dev/sqlc/internal/engine/sqlite"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

var parseCmd = &cobra.Command{
	Use:   "parse [file]",
	Short: "Parse SQL and output the AST as JSON",
	Long: `Parse SQL from a file or stdin and output the abstract syntax tree as JSON.

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

		// Parse SQL based on dialect
		var stmts []ast.Statement
		switch dialect {
		case "postgresql", "postgres", "pg":
			parser := postgresql.NewParser()
			stmts, err = parser.Parse(input)
		case "mysql":
			parser := dolphin.NewParser()
			stmts, err = parser.Parse(input)
		case "sqlite":
			parser := sqlite.NewParser()
			stmts, err = parser.Parse(input)
		case "clickhouse":
			parser := clickhouse.NewParser()
			stmts, err = parser.Parse(input)
		default:
			return fmt.Errorf("unsupported dialect: %s (use postgresql, mysql, sqlite, or clickhouse)", dialect)
		}
		if err != nil {
			return fmt.Errorf("parse error: %w", err)
		}

		// Output AST as JSON
		stdout := cmd.OutOrStdout()
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")

		for _, stmt := range stmts {
			if err := encoder.Encode(stmt.Raw); err != nil {
				return fmt.Errorf("failed to encode AST: %w", err)
			}
		}

		return nil
	},
}
