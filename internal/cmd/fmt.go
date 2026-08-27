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
	"github.com/sqlc-dev/sqlc/internal/engine/postgresql"
	"github.com/sqlc-dev/sqlc/internal/engine/sqlite"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/format"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
)

// noLineLimit renders without a maximum line width: like gofmt, fmt never
// rewraps a line on its own. Breaks come from the author's own line breaks
// at the boundaries the printer models, and from comments, which cannot
// share a line with the code after them.
const noLineLimit = 1 << 30

// queryFormatter is what the fmt command needs from an engine: a parser
// that surfaces the comments its lexer already scans (so statements and
// comments come from one pass, and interior comments format along with
// their statements) plus dialect-aware formatting. Comments cannot be
// reprinted without parser support, so ParseFile is the price of admission.
type queryFormatter interface {
	ParseFile(io.Reader) (*ast.File, error)
	format.Dialect
}

// newQueryFormatter returns the formatter for engines fmt supports —
// SQLite and PostgreSQL today. An engine joins by teaching its parser to
// surface comments (meyer's and oliphant's ParseFile are the templates) and
// adding its case here.
func newQueryFormatter(engine config.Engine) queryFormatter {
	switch engine {
	case config.EnginePostgreSQL:
		return postgresql.NewParser()
	case config.EngineSQLite:
		return sqlite.NewParser()
	default:
		return nil
	}
}

func newFmtCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fmt",
		Short: "Format SQL queries",
		Long: `Format the SQL query files referenced by the configuration file.

Each query is parsed with the engine's parser and printed back in a canonical
form, with its comments kept where they were written. A statement that cannot
be proven to survive formatting unchanged is left exactly as written. Files
are rewritten in place; pass --diff to print the changes to stdout instead.`,
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
	unsupported := map[config.Engine]bool{}
	for _, file := range files {
		f := newQueryFormatter(engines[file])
		if f == nil {
			if !unsupported[engines[file]] {
				unsupported[engines[file]] = true
				fmt.Fprintf(stderr, "sqlc fmt does not yet support the %s engine; query files left unchanged\n", engines[file])
			}
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
		formatted, err := formatQueries(f, string(contents))
		if err != nil {
			// Not every query file parses without the compiler's help — for
			// example, some engines only accept sqlc.slice() after
			// preprocessing. Leave those files as they are.
			fmt.Fprintf(stderr, "%s: skipped: %s\n", rel, err)
			continue
		}
		// File-level belt for the reprinter path: no formatting result may
		// change the file's comments. A violation is a formatter bug; keep
		// the file as written and say so.
		before, err1 := f.ParseFile(strings.NewReader(string(contents)))
		after, err2 := f.ParseFile(strings.NewReader(formatted))
		if err1 == nil && (err2 != nil || !sameComments(before.Comments, after.Comments)) {
			fmt.Fprintf(stderr, "%s: skipped: formatting would alter comments (this is a bug in sqlc fmt)\n", rel)
			continue
		}
		output[file] = formatted
	}
	return output, nil
}

// formatQueries formats each statement in a query file. The comments above
// every statement (that's where the "-- name:" annotation lives) are kept
// as written; comments inside a statement are threaded through the printer
// by position.
func formatQueries(f queryFormatter, src string) (string, error) {
	// The parser hands over statements and comments from one lexer pass;
	// statement bodies slice the comment list by span.
	file, err := f.ParseFile(strings.NewReader(src))
	if err != nil {
		return "", err
	}
	stmts, fileComments := file.Stmts, file.Comments
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
		segOff := segStart
		// A comment on the same line as the previous statement's terminator
		// belongs to that statement, not to this one's header.
		if len(blocks) > 0 {
			if trailing, remainder := splitTrailingComment(seg); trailing != "" {
				blocks[len(blocks)-1] += " " + trailing
				segOff += len(seg) - len(remainder)
				seg = remainder
			}
		}
		orig := strings.TrimSpace(strings.TrimLeft(seg, "; \t\r\n"))
		if orig == "" {
			continue
		}
		header, rest := splitLeadingComments(orig)
		// The statement body's own file offset, so its interior comments
		// can be sliced out of the file's comment list in the same
		// coordinates the parser stamped on the nodes.
		restOff := segOff + strings.Index(seg, orig) + len(orig) - len(rest)
		var interior []ast.Comment
		for _, c := range fileComments {
			if c.Start >= restOff && c.End <= start+length {
				interior = append(interior, c)
			}
		}
		block := append([]string{}, header...)
		block = append(block, formatStmt(f, stmt.Raw, rest, interior, src))
		blocks = append(blocks, strings.Join(block, "\n"))
	}
	// Everything after the last statement is its terminator, whitespace and
	// trailing comments; keep the comments.
	tailSeg := src[prevEnd:]
	if len(blocks) > 0 {
		if trailing, remainder := splitTrailingComment(tailSeg); trailing != "" {
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
func splitTrailingComment(seg string) (string, string) {
	k := 0
	for k < len(seg) && (seg[k] == ';' || seg[k] == ' ' || seg[k] == '\t' || seg[k] == '\r') {
		k++
	}
	rest := seg[k:]
	line, _, _ := strings.Cut(rest, "\n")
	trimmed := strings.TrimSpace(line)
	switch {
	case strings.HasPrefix(trimmed, "--"):
		return trimmed, seg[k+len(line):]
	case strings.HasPrefix(trimmed, "/*") &&
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
func splitLeadingComments(text string) ([]string, string) {
	var header []string
	rest := text
	for rest != "" {
		line, remainder, found := strings.Cut(rest, "\n")
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "/*") && !strings.Contains(trimmed, "*/") {
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
		if !isCommentLine(trimmed) {
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

func isCommentLine(line string) bool {
	switch {
	case line == "":
		// A blank line between comments.
		return true
	case strings.HasPrefix(line, "--"):
		return true
	case strings.HasPrefix(line, "/*") && strings.HasSuffix(line, "*/") && strings.Count(line, "*/") == 1:
		// A block comment contained on a single line.
		return true
	default:
		return false
	}
}

// fingerprinter is implemented by engines that can reduce a query to a
// fingerprint that survives changes in whitespace, case and layout —
// PostgreSQL via oliphant's pg_query-compatible Fingerprint. Where it is
// available, fmt gets a proof the other checks cannot give: the formatted
// statement still parses to the same query as the original.
type fingerprinter interface {
	Fingerprint(string) (string, error)
}

// formatStmt returns the canonical form of a single statement, or the
// original text when formatting cannot be proven to preserve the query.
func formatStmt(f queryFormatter, raw *ast.RawStmt, orig string, interior []ast.Comment, src string) string {
	fallback := strings.TrimSuffix(strings.TrimSpace(orig), ";") + ";"
	if out, ok := formatWithComments(f, raw, interior, src); ok {
		if fp, ok := f.(fingerprinter); ok && !sameFingerprint(fp, orig, out) {
			return fallback
		}
		return out
	}
	return fallback
}

// sameFingerprint reports that both texts fingerprint successfully to the
// same value. Anything less is not a proof, so the caller keeps the
// statement as written.
func sameFingerprint(fp fingerprinter, orig, out string) bool {
	a, err1 := fp.Fingerprint(orig)
	b, err2 := fp.Fingerprint(out)
	return err1 == nil && err2 == nil && a == b
}

// formatWithComments pretty-prints a statement with its interior comments
// and proves the result faithful three ways before accepting it: every
// comment survives (multiset equality after reparsing the output), the SQL
// still parses as one statement, and the output is a fixed point (printing
// the reparsed statement with its reparsed comments reproduces it).
func formatWithComments(f queryFormatter, raw *ast.RawStmt, interior []ast.Comment, src string) (string, bool) {
	out := prettyCommented(raw, f, interior, src)
	if strings.TrimSpace(strings.TrimSuffix(out, ";")) == "" {
		return "", false
	}
	file, err := f.ParseFile(strings.NewReader(out))
	if err != nil || len(file.Stmts) != 1 || file.Stmts[0].Raw == nil {
		return "", false
	}
	if !sameComments(interior, file.Comments) {
		return "", false
	}
	again := prettyCommented(file.Stmts[0].Raw, f, file.Comments, out)
	if !strings.EqualFold(again, out) {
		return "", false
	}
	return out, true
}

// prettyCommented attaches a statement's comments to its nodes and renders
// it, returning "" when the printer panics or fails to place every comment.
func prettyCommented(raw *ast.RawStmt, f queryFormatter, comments []ast.Comment, src string) (out string) {
	defer func() {
		if r := recover(); r != nil {
			out = ""
		}
	}()
	ct := ast.AttachComments(raw, f, comments, src)
	out = ast.PrettyWithComments(raw, f, noLineLimit, ct)
	if !ct.Exhausted() {
		return ""
	}
	return out
}

// sameComments reports whether two comment lists carry the same comment
// texts, as multisets, ignoring surrounding whitespace.
func sameComments(a, b []ast.Comment) bool {
	if len(a) != len(b) {
		return false
	}
	counts := make(map[string]int, len(a))
	for _, c := range a {
		counts[strings.TrimSpace(c.Text)]++
	}
	for _, c := range b {
		counts[strings.TrimSpace(c.Text)]--
	}
	for _, n := range counts {
		if n != 0 {
			return false
		}
	}
	return true
}
