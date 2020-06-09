package metadata

import (
	"fmt"
	"strings"
	"unicode"
)

type CommentSyntax int

const (
	CommentSyntaxDash CommentSyntax = iota
	CommentSyntaxStar               // Note: this is the only style supported by the MySQL sqlparser
	CommentSyntaxHash
)

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
		if commentStyle == CommentSyntaxDash && !strings.HasPrefix(line, "-- name:") {
			continue
		}
		if commentStyle == CommentSyntaxStar && !strings.HasPrefix(line, "/* name:") {
			continue
		}
		part := strings.Split(strings.TrimSpace(line), " ")

		if commentStyle == CommentSyntaxStar {
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
