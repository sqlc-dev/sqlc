package sqlite

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// The analysis shell is driven with a script on its standard input, made
// of SQL and of the shell's own dot-commands: `.print` labels the sections
// of the output so that each statement's results can be found again,
// `.mode` chooses how rows print — json for the catalog and the bytecode,
// quote for a query's own rows, since quote mode is the one that tells a
// blob from a string and a real from an integer — and `.stats stmt` makes
// the shell print what it knows about each result column of the statement
// it has just run, which with column metadata compiled in is the column's
// name, declared type, database, table and origin column.

// script is a shell script under construction, counting its lines so that
// an error the shell reports "near line N" can be placed.
type script struct {
	buf  strings.Builder
	line int // the number of the next line written
}

func newScript() *script {
	return &script{line: 1}
}

// add writes text on lines of its own, and returns the line it starts on.
func (s *script) add(text string) int {
	start := s.line
	text = strings.TrimSpace(text)
	if text == "" {
		return start
	}
	s.buf.WriteString(text)
	s.buf.WriteByte('\n')
	s.line += strings.Count(text, "\n") + 1
	return start
}

// sql writes statements, ending them with a semicolon when they lack one.
func (s *script) sql(text string) int {
	text = strings.TrimRight(strings.TrimSpace(text), "; \t\r\n")
	if text == "" {
		return s.line
	}
	return s.add(text + ";")
}

// section labels what follows, up to the next label.
func (s *script) section(name string) {
	s.add(".print @@" + name)
}

// columnMeta is what `.stats stmt` prints about one result column. Table
// and Origin name the table column a result column is read from, and are
// empty for an expression.
type columnMeta struct {
	Name     string
	DeclType string
	Database string
	Table    string
	Origin   string
}

// section is the output between two labels.
type section struct {
	// blocks holds every JSON result set printed in the section.
	blocks []json.RawMessage
	// rows holds every other line printed before the statement's
	// statistics: in quote mode, one row each.
	rows []string
	// prepared says a statement's statistics were printed, which the shell
	// does only for a statement it could prepare.
	prepared bool
	columns  []columnMeta
}

// decode reads the section's nth JSON result set.
func (s *section) decode(n int, v any) error {
	if n >= len(s.blocks) {
		return fmt.Errorf("expected at least %d result set(s), got %d", n+1, len(s.blocks))
	}
	return json.Unmarshal(s.blocks[n], v)
}

// output is what one run of the shell printed, by section.
type output struct {
	sections map[string]*section
	stderr   string
}

var errLineRe = regexp.MustCompile(`near line (\d+):`)

// errorsBefore returns the errors the shell reported for statements before
// the given line: the ones that are not the statement under analysis
// failing to run, which it may, with nothing bound to its parameters.
func (o *output) errorsBefore(line int) []string {
	var errs []string
	for _, l := range strings.Split(o.stderr, "\n") {
		if m := errLineRe.FindStringSubmatch(l); m != nil {
			if n, _ := strconv.Atoi(m[1]); n < line {
				errs = append(errs, strings.TrimSpace(l))
			}
		}
	}
	return errs
}

// run feeds a script to the shell over a fresh in-memory database and
// parses what it printed. A statement that fails does not stop the shell,
// so a failure is read from the output: a section whose statement has no
// statistics was not prepared, and stderr says why.
func run(ctx context.Context, binary string, s *script) (*output, error) {
	cmd := exec.CommandContext(ctx, binary, ":memory:")
	cmd.Stdin = strings.NewReader(s.buf.String())
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		var exit *exec.ExitError
		if !errors.As(err, &exit) {
			return nil, fmt.Errorf("sqlite3: %w", err)
		}
	}
	out := &output{sections: map[string]*section{}, stderr: stderr.String()}
	var (
		cur   *section
		block []string // lines of the JSON result set being read
		row   []string // lines of the quote-mode row being read
	)
	for _, line := range strings.Split(stdout.String(), "\n") {
		switch {
		case block != nil:
			block = append(block, line)
			if strings.HasSuffix(line, "]") {
				cur.blocks = append(cur.blocks, json.RawMessage(strings.Join(block, "\n")))
				block = nil
			}
		case row != nil:
			row = append(row, line)
			if joined := strings.Join(row, "\n"); strings.Count(joined, "'")%2 == 0 {
				cur.rows = append(cur.rows, joined)
				row = nil
			}
		case strings.HasPrefix(line, "@@"):
			cur = &section{}
			out.sections[line[2:]] = cur
		case cur == nil || line == "":
		case strings.HasPrefix(line, "["):
			if strings.HasSuffix(line, "]") {
				cur.blocks = append(cur.blocks, json.RawMessage(line))
			} else {
				block = []string{line}
			}
		case strings.HasPrefix(line, "Number of output columns:"):
			cur.prepared = true
		case cur.prepared:
			if i, field, value, ok := statLine(line); ok {
				for len(cur.columns) <= i {
					cur.columns = append(cur.columns, columnMeta{})
				}
				c := &cur.columns[i]
				switch field {
				case "name":
					c.Name = value
				case "declared type":
					c.DeclType = value
				case "database name":
					c.Database = value
				case "table name":
					c.Table = value
				case "origin name":
					c.Origin = value
				}
			}
		default:
			// A quote-mode row, complete unless a string in it spans lines.
			if strings.Count(line, "'")%2 == 0 {
				cur.rows = append(cur.rows, line)
			} else {
				row = []string{line}
			}
		}
	}
	return out, nil
}

var statLineRe = regexp.MustCompile(`^Column (\d+) ([a-z ]+):`)

// statLine reads one line of `.stats stmt` column output, which the shell
// prints as the label padded to 36 columns, a space and the value; a value
// the library has no answer for prints as the C library prints a null
// string.
func statLine(line string) (int, string, string, bool) {
	m := statLineRe.FindStringSubmatch(line)
	if m == nil {
		return 0, "", "", false
	}
	i, _ := strconv.Atoi(m[1])
	label := m[0]
	value := ""
	if start := max(len(label), 36) + 1; start < len(line) {
		value = line[start:]
	}
	if value == "(null)" {
		value = ""
	}
	return i, m[2], value, true
}

// cells splits a quote-mode row into its values as printed.
func cells(row string) []string {
	var out []string
	start := 0
	quoted := false
	for i := 0; i < len(row); i++ {
		switch {
		case row[i] == '\'':
			quoted = !quoted
		case row[i] == ',' && !quoted:
			out = append(out, row[start:i])
			start = i + 1
		}
	}
	return append(out, row[start:])
}

// storageClass reads the storage class of a value as quote mode prints
// it: a string is quoted, a blob is quoted with an X, NULL is spelled
// out, and a real has a point, an exponent or is infinite or not a number
// where an integer is digits.
func storageClass(cell string) string {
	switch {
	case cell == "NULL":
		return "null"
	case strings.HasPrefix(cell, "'"):
		return "text"
	case strings.HasPrefix(cell, "x'") || strings.HasPrefix(cell, "X'"):
		return "blob"
	case strings.ContainsAny(cell, ".eEIN"):
		return "real"
	default:
		return "integer"
	}
}
