package metadata

import (
	"bufio"
	"fmt"
	"github.com/sqlc-dev/sqlc/internal/constants"
	"strings"
	"unicode"

	"github.com/sqlc-dev/sqlc/internal/source"
)

type CommentSyntax source.CommentSyntax

type Metadata struct {
	Name     string
	Cmd      string
	Comments []string
	Params   map[string]string
	Flags    map[string]bool

	// RuleSkiplist contains the names of rules to disable vetting for.
	// If the map is empty, but the disable vet flag is specified, then all rules are ignored.
	RuleSkiplist map[string]struct{}

	Filename string
}

const (
	CmdExec       = ":exec"
	CmdExecResult = ":execresult"
	CmdExecRows   = ":execrows"
	CmdExecLastId = ":execlastid"
	CmdMany       = ":many"
	CmdOne        = ":one"
	CmdCopyFrom   = ":copyfrom"
	CmdBatchExec  = ":batchexec"
	CmdBatchMany  = ":batchmany"
	CmdBatchOne   = ":batchone"
)

// A query name must be a valid Go identifier
//
// https://golang.org/ref/spec#Identifiers
func validateQueryName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("invalid query name: %q", name)
	}
	for i, c := range name {
		isLetter := unicode.IsLetter(c) || c == '_'
		isDigit := unicode.IsDigit(c)
		if i == 0 && !isLetter {
			return fmt.Errorf("invalid query name %q", name)
		} else if !(isLetter || isDigit) {
			return fmt.Errorf("invalid query name %q", name)
		}
	}
	return nil
}

func ParseQueryNameAndType(t string, commentStyle CommentSyntax) (string, string, error) {
	for _, line := range strings.Split(t, "\n") {
		var prefix string
		if strings.HasPrefix(line, "--") {
			if !commentStyle.Dash {
				continue
			}
			prefix = "--"
		}
		if strings.HasPrefix(line, "/*") {
			if !commentStyle.SlashStar {
				continue
			}
			prefix = "/*"
		}
		if strings.HasPrefix(line, "#") {
			if !commentStyle.Hash {
				continue
			}
			prefix = "#"
		}
		if prefix == "" {
			continue
		}
		rest := line[len(prefix):]
		if !strings.HasPrefix(strings.TrimSpace(rest), "name") {
			continue
		}
		if !strings.Contains(rest, ":") {
			continue
		}
		if !strings.HasPrefix(rest, " name: ") {
			return "", "", fmt.Errorf("invalid metadata: %s", line)
		}

		part := strings.Split(strings.TrimSpace(line), " ")
		if prefix == "/*" {
			part = part[:len(part)-1] // removes the trailing "*/" element
		}
		if len(part) == 3 {
			return "", "", fmt.Errorf("missing query type [':one', ':many', ':exec', ':execrows', ':execlastid', ':execresult', ':copyfrom', 'batchexec', 'batchmany', 'batchone']: %s", line)
		}
		if len(part) != 4 {
			return "", "", fmt.Errorf("invalid query comment: %s", line)
		}
		queryName := part[2]
		queryType := strings.TrimSpace(part[3])
		switch queryType {
		case CmdOne, CmdMany, CmdExec, CmdExecResult, CmdExecRows, CmdExecLastId, CmdCopyFrom, CmdBatchExec, CmdBatchMany, CmdBatchOne:
		default:
			return "", "", fmt.Errorf("invalid query type: %s", queryType)
		}
		if err := validateQueryName(queryName); err != nil {
			return "", "", err
		}
		return queryName, queryType, nil
	}
	return "", "", nil
}

// QueryBlock is one named query block (from " name: X :cmd" to the next such line or EOF).
type QueryBlock struct {
	SQL  string
	Name string
	Cmd  string
}

// isBlockStartLine reports whether the line starts a named query block (e.g. "-- name: GetUser :one").
func isBlockStartLine(line string, commentStyle CommentSyntax) bool {
	line = strings.TrimSpace(line)
	var rest string
	switch {
	case strings.HasPrefix(line, "--"):
		if !commentStyle.Dash {
			return false
		}
		rest = line[2:] // keep " name: X :cmd" with leading space to match ParseQueryNameAndType
	case strings.HasPrefix(line, "#"):
		if !commentStyle.Hash {
			return false
		}
		rest = line[1:]
	case strings.HasPrefix(line, "/*"):
		if !commentStyle.SlashStar {
			return false
		}
		rest = line[2:]
		if strings.HasSuffix(rest, "*/") {
			rest = rest[:len(rest)-2]
		}
		rest = strings.TrimSpace(rest)
	default:
		return false
	}
	if !strings.HasPrefix(rest, " name: ") {
		return false
	}
	part := strings.Split(rest, " ")
	if len(part) < 4 {
		return false
	}
	queryType := strings.TrimSpace(part[3])
	switch queryType {
	case CmdOne, CmdMany, CmdExec, CmdExecResult, CmdExecRows, CmdExecLastId, CmdCopyFrom, CmdBatchExec, CmdBatchMany, CmdBatchOne:
		return true
	}
	return false
}

// QueryBlocks splits content into named query blocks. Each block runs from a " name: X :cmd" line
// to the next such line (or EOF). Returns one entry per block with non-empty name.
func QueryBlocks(content string, commentStyle CommentSyntax) ([]QueryBlock, error) {
	lines := strings.Split(content, "\n")
	var starts []int
	pos := 0
	for _, line := range lines {
		if isBlockStartLine(line, commentStyle) {
			starts = append(starts, pos)
		}
		pos += len(line) + 1
	}
	var out []QueryBlock
	for i := 0; i < len(starts); i++ {
		end := len(content)
		if i+1 < len(starts) {
			end = starts[i+1]
		}
		blockSQL := content[starts[i]:end]
		name, cmd, err := ParseQueryNameAndType(blockSQL, commentStyle)
		if err != nil {
			return nil, err
		}
		if name == "" {
			continue
		}
		out = append(out, QueryBlock{SQL: blockSQL, Name: name, Cmd: cmd})
	}
	return out, nil
}

// ParseCommentFlags processes the comments provided with queries to determine the metadata params, flags and rules to skip.
// All flags in query comments are prefixed with `@`, e.g. @param, @@sqlc-vet-disable.
func ParseCommentFlags(comments []string) (map[string]string, map[string]bool, map[string]struct{}, error) {
	params := make(map[string]string)
	flags := make(map[string]bool)
	ruleSkiplist := make(map[string]struct{})

	for _, line := range comments {
		s := bufio.NewScanner(strings.NewReader(line))
		s.Split(bufio.ScanWords)

		s.Scan()
		token := s.Text()

		if !strings.HasPrefix(token, "@") {
			continue
		}

		switch token {
		case constants.QueryFlagParam:
			s.Scan()
			name := s.Text()
			var rest []string
			for s.Scan() {
				paramToken := s.Text()
				rest = append(rest, paramToken)
			}
			params[name] = strings.Join(rest, " ")

		case constants.QueryFlagSqlcVetDisable:
			flags[token] = true

			// Vet rules can all be disabled in the same line or split across lines .i.e.
			// /* @sqlc-vet-disable sqlc/db-prepare delete-without-where */
			// is equivalent to:
			// /* @sqlc-vet-disable sqlc/db-prepare */
			// /* @sqlc-vet-disable delete-without-where */
			for s.Scan() {
				ruleSkiplist[s.Text()] = struct{}{}
			}

		default:
			flags[token] = true
		}

		if s.Err() != nil {
			return params, flags, ruleSkiplist, s.Err()
		}
	}

	return params, flags, ruleSkiplist, nil
}
