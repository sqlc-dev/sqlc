package metadata

import (
	"bufio"
	"fmt"
	"strings"
	"unicode"
)

type Metadata struct {
	Name     string
	Cmd      string
	Comments []string
	Params   map[string]string
	Flags    map[string]bool

	Filename string
}

type CommentSyntax struct {
	Dash      bool
	Hash      bool
	SlashStar bool
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

func ParseQueryMetadata(rawSql string, commentStyle CommentSyntax) (Metadata, error) {
	md := Metadata{}
	s := bufio.NewScanner(strings.NewReader(strings.TrimSpace(rawSql)))
	var comments []string
	for s.Scan() {
		line := s.Text()
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
		rest = strings.TrimSuffix(rest, "*/")
		comments = append(comments, rest)

		if !strings.HasPrefix(strings.TrimSpace(rest), "name") {
			continue
		}
		if !strings.Contains(rest, ":") {
			continue
		}
		if !strings.HasPrefix(rest, " name: ") {
			return md, fmt.Errorf("invalid metadata: %s", line)
		}

		comments = comments[:len(comments)-1] // Remove tha name line from returned comments

		parts := strings.Split(strings.TrimSpace(rest), " ")

		if len(parts) == 2 {
			return md, fmt.Errorf("missing query type [':one', ':many', ':exec', ':execrows', ':execlastid', ':execresult', ':copyfrom', 'batchexec', 'batchmany', 'batchone']: %s", line)
		}
		if len(parts) > 3 {
			return md, fmt.Errorf("invalid query comment: %s", line)
		}
		queryName := parts[1]
		queryType := parts[2]
		switch queryType {
		case CmdOne, CmdMany, CmdExec, CmdExecResult, CmdExecRows, CmdExecLastId, CmdCopyFrom, CmdBatchExec, CmdBatchMany, CmdBatchOne:
		default:
			return md, fmt.Errorf("invalid query type: %s", queryType)
		}
		if err := validateQueryName(queryName); err != nil {
			return md, err
		}
		md.Name = queryName
		md.Cmd = queryType
	}

	md.Comments = comments

	var err error
	md.Params, md.Flags, err = parseParamsAndFlags(md.Comments)
	if err != nil {
		return md, err
	}

	return md, s.Err()
}

func parseParamsAndFlags(comments []string) (map[string]string, map[string]bool, error) {
	params := make(map[string]string)
	flags := make(map[string]bool)

	for _, line := range comments {
		s := bufio.NewScanner(strings.NewReader(line))
		s.Split(bufio.ScanWords)

		s.Scan() // The first token is always the comment indicator, e.g. "--"
		s.Scan()
		token := s.Text()

		if !strings.HasPrefix(token, "@") {
			continue
		}

		switch token {
		case "@param":
			s.Scan()
			name := s.Text()
			var rest []string
			for s.Scan() {
				paramToken := s.Text()
				if paramToken == "*/" {
					break
				}
				rest = append(rest, paramToken)
			}
			params[name] = strings.Join(rest, " ")
		default:
			flags[token] = true
		}

		if s.Err() != nil {
			return params, flags, s.Err()
		}
	}

	return params, flags, nil
}
