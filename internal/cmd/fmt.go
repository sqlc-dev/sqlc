package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/trace"
	"sort"
	"strings"

	"github.com/cubicdaiya/gonp"
	"github.com/spf13/cobra"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/engine/clickhouse"
	"github.com/sqlc-dev/sqlc/internal/engine/dolphin"
	"github.com/sqlc-dev/sqlc/internal/engine/postgresql"
	"github.com/sqlc-dev/sqlc/internal/engine/sqlite"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/format"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
)

// fmtLineWidth is the line width formatted queries are wrapped to. A
// statement whose flat form fits stays on one line; a longer one breaks at
// clause boundaries, and any list or parenthesized region that still does
// not fit breaks again one indentation level deeper.
const fmtLineWidth = 80

// queryFormatter is the subset of the engine parsers that the fmt command
// needs: parsing SQL into statements and dialect-aware formatting.
type queryFormatter interface {
	Parse(io.Reader) ([]ast.Statement, error)
	format.Dialect
}

func newQueryFormatter(engine config.Engine) queryFormatter {
	switch engine {
	case config.EnginePostgreSQL:
		return postgresql.NewParser()
	case config.EngineMySQL:
		return dolphin.NewParser()
	case config.EngineSQLite:
		return sqlite.NewParser()
	case config.EngineClickHouse:
		return clickhouse.NewParser()
	default:
		return nil
	}
}

func newFmtCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fmt",
		Short: "Format query files",
		Long: `Format the SQL query files referenced by the configuration file.

Each query is parsed with the engine's parser and printed back in a canonical
form, keeping the comments above it (including the "-- name:" annotation).
Comments are never deleted: a statement with comments inside it, or one that
cannot be proven to survive formatting unchanged, is left exactly as written.
Files are rewritten in place; pass --diff to print the changes to stdout
instead.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			defer trace.StartRegion(cmd.Context(), "fmt").End()
			stderr := cmd.ErrOrStderr()
			stdout := cmd.OutOrStdout()
			showDiff, err := cmd.Flags().GetBool("diff")
			if err != nil {
				return err
			}
			dir, name := getConfigPath(stderr, cmd.Flag("file"))
			output, err := Format(cmd.Context(), dir, name, &Options{
				Env:    ParseEnv(cmd),
				Stderr: stderr,
			})
			if err != nil {
				os.Exit(1)
			}

			filenames := make([]string, 0, len(output))
			for filename := range output {
				filenames = append(filenames, filename)
			}
			sort.Strings(filenames)

			for _, filename := range filenames {
				existing, err := os.ReadFile(filename)
				if err != nil {
					fmt.Fprintf(stderr, "%s: %s\n", filename, err)
					return err
				}
				source := output[filename]
				if string(existing) == source {
					continue
				}
				if showDiff {
					diff := gonp.New(getLines(existing), getLines([]byte(source)))
					diff.Compose()
					uniHunks := filterHunks(diff.UnifiedHunks())
					if len(uniHunks) == 0 {
						continue
					}
					rel := filepath.ToSlash(strings.TrimPrefix(filename, dir))
					fmt.Fprintf(stdout, "--- a%s\n", rel)
					fmt.Fprintf(stdout, "+++ b%s\n", rel)
					diff.FprintUniHunks(stdout, uniHunks)
					continue
				}
				if err := os.WriteFile(filename, []byte(source), 0644); err != nil {
					fmt.Fprintf(stderr, "%s: %s\n", filename, err)
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().Bool("diff", false, "display diffs instead of rewriting files")
	return cmd
}

// Format formats every query file referenced by the configuration with its
// engine's SQL formatter. It returns the formatted contents of each file,
// keyed by absolute path. Query files for engines without a formatter are
// skipped.
func Format(ctx context.Context, dir, filename string, o *Options) (map[string]string, error) {
	defer trace.StartRegion(ctx, "format").End()
	stderr := o.Stderr

	configPath, conf, err := o.ReadConfig(dir, filename)
	if err != nil {
		return nil, err
	}
	base := filepath.Base(configPath)
	if err := config.Validate(conf); err != nil {
		fmt.Fprintf(stderr, "error validating %s: %s\n", base, err)
		return nil, err
	}

	engines := map[string]config.Engine{}
	var files []string
	for _, sql := range conf.SQL {
		joined := make([]string, 0, len(sql.Queries))
		for _, q := range sql.Queries {
			joined = append(joined, filepath.Join(dir, q))
		}
		list, err := sqlpath.Glob(joined)
		if err != nil {
			fmt.Fprintf(stderr, "error listing queries: %s\n", err)
			return nil, err
		}
		for _, file := range list {
			if _, ok := engines[file]; ok {
				continue
			}
			engines[file] = sql.Engine
			files = append(files, file)
		}
	}

	output := map[string]string{}
	for _, file := range files {
		f := newQueryFormatter(engines[file])
		if f == nil {
			continue
		}
		rel, err := filepath.Rel(dir, file)
		if err != nil {
			rel = file
		}
		contents, err := os.ReadFile(file)
		if err != nil {
			fmt.Fprintf(stderr, "%s: skipped: %s\n", rel, err)
			continue
		}
		formatted, err := formatQueries(engines[file], f, string(contents))
		if err != nil {
			// Not every query file parses without the compiler's help — for
			// example, some engines only accept sqlc.slice() after
			// preprocessing. Leave those files as they are.
			fmt.Fprintf(stderr, "%s: skipped: %s\n", rel, err)
			continue
		}
		output[file] = formatted
	}
	return output, nil
}

// formatQueries formats each statement in a query file, preserving the
// comments above every statement (that's where the "-- name:" annotation
// lives).
func formatQueries(engine config.Engine, f queryFormatter, src string) (string, error) {
	stmts, err := f.Parse(strings.NewReader(src))
	if err != nil {
		return "", err
	}
	var blocks []string
	prevEnd := 0
	for _, stmt := range stmts {
		if stmt.Raw == nil {
			continue
		}
		start := stmt.Raw.StmtLocation
		length := stmt.Raw.StmtLen
		if length == 0 {
			// An unterminated final statement is reported with a zero length.
			length = len(src) - start
		}
		if start < 0 || start+length > len(src) {
			return "", fmt.Errorf("statement location out of bounds")
		}
		// Engines differ on whether a statement's span covers the comments
		// above it, so start each segment where the previous one ended: the
		// gap holds nothing but the previous terminator, whitespace and
		// comments.
		segStart := prevEnd
		if start < segStart {
			segStart = start
		}
		if start+length <= prevEnd {
			continue
		}
		prevEnd = start + length
		seg := src[segStart : start+length]
		// A comment on the same line as the previous statement's terminator
		// belongs to that statement, not to this one's header.
		if len(blocks) > 0 {
			if trailing, remainder := splitTrailingComment(engine, seg); trailing != "" {
				blocks[len(blocks)-1] += " " + trailing
				seg = remainder
			}
		}
		orig := strings.TrimSpace(strings.TrimLeft(seg, "; \t\r\n"))
		if orig == "" {
			continue
		}
		header, rest := splitLeadingComments(engine, orig)
		block := append([]string{}, header...)
		block = append(block, formatStmt(engine, f, stmt.Raw, rest))
		blocks = append(blocks, strings.Join(block, "\n"))
	}
	// Everything after the last statement is its terminator, whitespace and
	// trailing comments; keep the comments.
	tailSeg := src[prevEnd:]
	if len(blocks) > 0 {
		if trailing, remainder := splitTrailingComment(engine, tailSeg); trailing != "" {
			blocks[len(blocks)-1] += " " + trailing
			tailSeg = remainder
		}
	}
	if tail := strings.TrimSpace(strings.TrimLeft(tailSeg, "; \t\r\n")); tail != "" {
		blocks = append(blocks, tail)
	}
	if len(blocks) == 0 {
		return src, nil
	}
	return strings.Join(blocks, "\n\n") + "\n", nil
}

// splitTrailingComment splits a comment sitting on the same line as the
// previous statement's terminator off the front of the text between two
// statements, returning the comment (or "") and the remaining text.
func splitTrailingComment(engine config.Engine, seg string) (string, string) {
	k := 0
	for k < len(seg) && (seg[k] == ';' || seg[k] == ' ' || seg[k] == '\t' || seg[k] == '\r') {
		k++
	}
	rest := seg[k:]
	line, _, _ := strings.Cut(rest, "\n")
	trimmed := strings.TrimSpace(line)
	switch {
	case strings.HasPrefix(trimmed, "--"),
		engine == config.EngineMySQL && strings.HasPrefix(trimmed, "#"):
		return trimmed, seg[k+len(line):]
	case strings.HasPrefix(trimmed, "/*") && !strings.HasPrefix(trimmed, "/*+") &&
		strings.HasSuffix(trimmed, "*/") && strings.Count(trimmed, "*/") == 1:
		// A block comment contained on the terminator's line.
		return trimmed, seg[k+len(line):]
	}
	return "", seg
}

// splitLeadingComments splits a statement's text into the comment lines above
// it and the statement itself. Comment lines are returned with surrounding
// whitespace trimmed (block comment continuation lines keep their leading
// whitespace); blank lines between comments are kept.
func splitLeadingComments(engine config.Engine, text string) ([]string, string) {
	var header []string
	rest := text
	for rest != "" {
		line, remainder, found := strings.Cut(rest, "\n")
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "/*") && !strings.HasPrefix(trimmed, "/*+") && !strings.Contains(trimmed, "*/") {
			// A block comment spanning several lines: commit it to the
			// header only when its closing line arrives with nothing after
			// it. Otherwise leave the whole block with the statement.
			block := []string{trimmed}
			blockRest := remainder
			for closed := false; !closed; {
				if !found || blockRest == "" {
					return header, rest
				}
				var blockLine string
				blockLine, blockRest, found = strings.Cut(blockRest, "\n")
				if end := strings.TrimSpace(blockLine); strings.Contains(end, "*/") {
					if !strings.HasSuffix(end, "*/") || strings.Count(end, "*/") != 1 {
						return header, rest
					}
					closed = true
				}
				block = append(block, strings.TrimRight(blockLine, " \t\r"))
			}
			header = append(header, block...)
			rest = blockRest
			continue
		}
		if !isCommentLine(engine, trimmed) {
			break
		}
		header = append(header, trimmed)
		if !found {
			rest = ""
			break
		}
		rest = remainder
	}
	// A statement never ends with a comment-only line here (the parser found a
	// statement in this segment), but drop any blank lines left at the end of
	// the header so comments sit directly above their query.
	for len(header) > 0 && header[len(header)-1] == "" {
		header = header[:len(header)-1]
	}
	return header, rest
}

func isCommentLine(engine config.Engine, line string) bool {
	switch {
	case line == "":
		// A blank line between comments.
		return true
	case strings.HasPrefix(line, "--"):
		return true
	case engine == config.EngineMySQL && strings.HasPrefix(line, "#"):
		return true
	case strings.HasPrefix(line, "/*") && strings.HasSuffix(line, "*/") && strings.Count(line, "*/") == 1:
		// A block comment contained on a single line. MySQL optimizer hints
		// share this syntax but belong to the query, not the header.
		return !strings.HasPrefix(line, "/*+")
	default:
		return false
	}
}

// formatStmt returns the canonical form of a single statement, or the
// original text when formatting cannot be proven to preserve the query.
func formatStmt(engine config.Engine, f queryFormatter, raw *ast.RawStmt, orig string) string {
	fallback := strings.TrimSuffix(strings.TrimSpace(orig), ";") + ";"
	if hasComment(engine, orig) {
		// Comments never survive the trip through the AST — the parsers
		// discard them — so a statement carrying comments (or optimizer
		// hints, which share their syntax) is left exactly as written.
		return fallback
	}
	out := formatRaw(raw, f)
	if strings.TrimSpace(strings.TrimSuffix(out, ";")) == "" {
		return fallback
	}
	if !verifyFormatted(engine, f, orig, out) {
		return fallback
	}
	return out
}

// hasComment reports whether the statement text contains a comment outside
// of string literals and quoted identifiers. Escaping conventions inside
// strings vary with server settings (e.g. PostgreSQL's
// standard_conforming_strings), so the text is scanned under both
// conventions and a comment found by either counts: a false positive only
// leaves a statement unformatted, while a miss would delete the comment.
func hasComment(engine config.Engine, sql string) bool {
	return scanComment(engine, sql, false) || scanComment(engine, sql, true)
}

// scanComment reports whether sql contains a comment token outside of
// quoted regions. backslash controls whether a backslash escapes the next
// character inside single-quoted strings. The scan is conservative: an
// unterminated quoted region also reports true, since anything after it
// cannot be classified.
func scanComment(engine config.Engine, sql string, backslash bool) bool {
	n := len(sql)
	i := 0
	for i < n {
		c := sql[i]
		switch {
		case c == '-' && i+1 < n && sql[i+1] == '-':
			return true
		case c == '/' && i+1 < n && sql[i+1] == '*':
			return true
		case c == '#' && engine == config.EngineMySQL:
			return true
		case c == '\'' || c == '"' || c == '`':
			term := c
			i++
			closed := false
			for i < n {
				if backslash && term == '\'' && sql[i] == '\\' {
					i += 2
					continue
				}
				if sql[i] == term {
					// A doubled quote is an escaped quote, not the end.
					if i+1 < n && sql[i+1] == term {
						i += 2
						continue
					}
					closed = true
					i++
					break
				}
				i++
			}
			if !closed {
				return true
			}
		case c == '$' && engine == config.EnginePostgreSQL:
			// A dollar-quoted string ($tag$ ... $tag$) can hold anything.
			j := i + 1
			for j < n && (isTagChar(sql[j])) {
				j++
			}
			if j < n && sql[j] == '$' {
				tag := sql[i : j+1]
				end := strings.Index(sql[j+1:], tag)
				if end < 0 {
					return true
				}
				i = j + 1 + end + len(tag)
			} else {
				i++
			}
		default:
			i++
		}
	}
	return false
}

func isTagChar(c byte) bool {
	return c == '_' ||
		('a' <= c && c <= 'z') ||
		('A' <= c && c <= 'Z') ||
		('0' <= c && c <= '9')
}

// formatRaw pretty-prints a statement's AST, turning a formatter panic into
// an empty string so the caller falls back to the original text.
func formatRaw(raw *ast.RawStmt, f queryFormatter) (out string) {
	defer func() {
		if r := recover(); r != nil {
			out = ""
		}
	}()
	return ast.Pretty(raw, f, fmtLineWidth)
}

// verifyFormatted reports whether the formatted SQL provably still means the
// same thing as the original statement.
func verifyFormatted(engine config.Engine, f queryFormatter, orig, formatted string) bool {
	if engine == config.EnginePostgreSQL {
		// PostgreSQL has a real fingerprint, independent of our formatter.
		before, err := postgresql.Fingerprint(orig)
		if err != nil {
			return false
		}
		after, err := postgresql.Fingerprint(formatted)
		if err != nil {
			return false
		}
		return before == after
	}
	// Round-trip: the formatted SQL must parse back as a single statement
	// that formats to itself. The comparison ignores case because some
	// parsers normalize keyword and identifier case.
	stmts, err := f.Parse(strings.NewReader(formatted))
	if err != nil || len(stmts) != 1 || stmts[0].Raw == nil {
		return false
	}
	return strings.EqualFold(formatRaw(stmts[0].Raw, f), formatted)
}
