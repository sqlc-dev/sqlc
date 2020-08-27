package metadata

import (
	"fmt"
	"strings"
	"unicode"
)

type CommentSyntax struct {
	Dash      bool
	Hash      bool
	SlashStar bool
}

const (
	CmdExec       = ":exec"
	CmdExecResult = ":execresult"
	CmdExecRows   = ":execrows"
	CmdMany       = ":many"
	CmdOne        = ":one"
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

func Parse(t string, commentStyle CommentSyntax) (string, string, error) {
	for _, line := range strings.Split(t, "\n") {
		var prefix string
		if strings.HasPrefix(line, "--") {
			if !commentStyle.Dash {
				continue
			}
			prefix = "-- name:"
		}
		if strings.HasPrefix(line, "/*") {
			if !commentStyle.SlashStar {
				continue
			}
			prefix = "/* name:"
		}
		if strings.HasPrefix(line, "#") {
			if !commentStyle.Hash {
				continue
			}
			prefix = "# name:"
		}
		if prefix == "" {
			continue
		}
		if !strings.HasPrefix(line, prefix) {
			continue
		}

		part := strings.Split(strings.TrimSpace(line), " ")
		if strings.HasPrefix(line, "/*") {
			part = part[:len(part)-1] // removes the trailing "*/" element
		}
		if len(part) == 2 {
			return "", "", fmt.Errorf("missing query type [':one', ':many', ':exec', ':execrows', ':execresult']: %s", line)
		}
		if len(part) != 4 {
			return "", "", fmt.Errorf("invalid query comment: %s", line)
		}
		queryName := part[2]
		queryType := strings.TrimSpace(part[3])
		switch queryType {
		case CmdOne, CmdMany, CmdExec, CmdExecResult, CmdExecRows:
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
