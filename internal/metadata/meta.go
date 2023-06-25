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
	CmdExecLastId = ":execlastid"
	CmdMany       = ":many"
	CmdOne        = ":one"
	CmdCopyFrom   = ":copyfrom"
	CmdBatchExec  = ":batchexec"
	CmdBatchMany  = ":batchmany"
	CmdBatchOne   = ":batchone"
)

type CmdParams struct {
	ManyKey        string
	InsertMultiple bool
	NoInference    bool
}

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

func Parse(t string, commentStyle CommentSyntax) (string, string, CmdParams, error) {
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
			return "", "", CmdParams{}, fmt.Errorf("invalid metadata: %s", line)
		}

		part := strings.Split(strings.TrimSpace(line), " ")
		if prefix == "/*" {
			part = part[:len(part)-1] // removes the trailing "*/" element
		}
		if len(part) < 4 {
			return "", "", CmdParams{}, fmt.Errorf("missing query type [':one', ':many', ':exec', ':execrows', ':execlastid', ':execresult', ':copyfrom', 'batchexec', 'batchmany', 'batchone']: %s", line)
		}
		queryName := part[2]
		queryType := strings.TrimSpace(part[3])
		switch queryType {
		case CmdOne, CmdMany, CmdExec, CmdExecResult, CmdExecRows, CmdExecLastId, CmdCopyFrom, CmdBatchExec, CmdBatchMany, CmdBatchOne:
		default:
			return "", "", CmdParams{}, fmt.Errorf("invalid query type: %s", queryType)
		}
		if err := validateQueryName(queryName); err != nil {
			return "", "", CmdParams{}, err
		}
		cmdParams, err := parseCmdParams(part[4:], queryType)
		if err != nil {
			return "", "", CmdParams{}, err
		}
		return queryName, queryType, cmdParams, nil
	}
	return "", "", CmdParams{}, nil
}

func parseCmdParams(part []string, queryType string) (CmdParams, error) {
	var ret CmdParams
	for _, p := range part {
		if p == "multiple" {
			if queryType != CmdExec && queryType != CmdExecResult && queryType != CmdExecRows && queryType != CmdExecLastId {
				return ret, fmt.Errorf("query command parameter multiple is invalid for query type %s", queryType)
			}
			ret.InsertMultiple = true
		} else if p == "no-inference" {
			if queryType != CmdExec && queryType != CmdExecResult && queryType != CmdExecRows && queryType != CmdExecLastId {
				return ret, fmt.Errorf("query command parameter no-inference is invalid for query type %s", queryType)
			}
			ret.NoInference = true
		} else if strings.HasPrefix(p, "key=") {
			if queryType != CmdMany {
				return ret, fmt.Errorf("query command parameter %s is invalid for query type %s", p, queryType)
			}
			ret.ManyKey = strings.TrimPrefix(p, "key=")
		} else {
			return ret, fmt.Errorf("invalid query command parameter %q", p)
		}
	}
	return ret, nil
}
