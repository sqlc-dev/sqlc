package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/sqlc-dev/sqlc/internal/engine/dolphin"
	"github.com/sqlc-dev/sqlc/internal/engine/postgresql"
	"github.com/sqlc-dev/sqlc/internal/engine/sqlite"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func NewCmdParse() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parse [file]",
		Short: "Parse SQL and output the AST as JSON (experimental)",
		Long: `Parse SQL from a file or stdin and output the abstract syntax tree as JSON.

This command is experimental and requires the 'parsecmd' experiment to be enabled.
Enable it by setting: SQLCEXPERIMENT=parsecmd

Examples:
  # Parse a SQL file with PostgreSQL dialect
  SQLCEXPERIMENT=parsecmd sqlc parse --dialect postgresql schema.sql

  # Parse from stdin with MySQL dialect
  echo "SELECT * FROM users" | SQLCEXPERIMENT=parsecmd sqlc parse --dialect mysql

  # Parse SQLite SQL
  SQLCEXPERIMENT=parsecmd sqlc parse --dialect sqlite queries.sql`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env := ParseEnv(cmd)
			if !env.Experiment.ParseCmd {
				return fmt.Errorf("parse command requires the 'parsecmd' experiment to be enabled.\nSet SQLCEXPERIMENT=parsecmd to use this command")
			}

			dialect, err := cmd.Flags().GetString("dialect")
			if err != nil {
				return err
			}
			if dialect == "" {
				return fmt.Errorf("--dialect flag is required (postgresql, mysql, or sqlite)")
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
			default:
				return fmt.Errorf("unsupported dialect: %s (use postgresql, mysql, or sqlite)", dialect)
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

	cmd.Flags().StringP("dialect", "d", "", "SQL dialect to use (postgresql, mysql, or sqlite)")

	return cmd
}
