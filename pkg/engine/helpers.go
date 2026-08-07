// Package engine helpers for splitting query files and parsing " name: X :cmd" comments.
// Plugins may use these to split query.sql into blocks and build Statement metadata.
// ParseNameAndCmd and QueryBlocks delegate to internal/metadata so behavior matches
// the built-in engines (postgresql, mysql, sqlite) exactly.

package engine

import (
	"strings"

	"github.com/sqlc-dev/sqlc/internal/metadata"
	"github.com/sqlc-dev/sqlc/internal/source"
)

// CommentSyntax defines which comment styles to recognize (e.g. "--", "/* */", "#").
// Same meaning as source.CommentSyntax / metadata.CommentSyntax.
type CommentSyntax struct {
	Dash      bool // "--"
	SlashStar bool // "/* */"
	Hash      bool // "#"
}

func toMetadataCommentSyntax(s CommentSyntax) metadata.CommentSyntax {
	return metadata.CommentSyntax(source.CommentSyntax{
		Dash: s.Dash, SlashStar: s.SlashStar, Hash: s.Hash,
	})
}

// CmdToString returns the sqlc command string for c (e.g. ":one", ":many").
// Returns "" for CMD_UNSPECIFIED.
func CmdToString(c Cmd) string {
	switch c {
	case Cmd_CMD_ONE:
		return ":one"
	case Cmd_CMD_MANY:
		return ":many"
	case Cmd_CMD_EXEC:
		return ":exec"
	case Cmd_CMD_EXEC_RESULT:
		return ":execresult"
	case Cmd_CMD_EXEC_ROWS:
		return ":execrows"
	case Cmd_CMD_EXEC_LAST_ID:
		return ":execlastid"
	case Cmd_CMD_COPY_FROM:
		return ":copyfrom"
	case Cmd_CMD_BATCH_EXEC:
		return ":batchexec"
	case Cmd_CMD_BATCH_MANY:
		return ":batchmany"
	case Cmd_CMD_BATCH_ONE:
		return ":batchone"
	default:
		return ""
	}
}

// CmdFromString parses s (e.g. ":one", ":many") into Cmd. Returns (cmd, true) on success.
func CmdFromString(s string) (Cmd, bool) {
	switch strings.TrimSpace(s) {
	case ":one":
		return Cmd_CMD_ONE, true
	case ":many":
		return Cmd_CMD_MANY, true
	case ":exec":
		return Cmd_CMD_EXEC, true
	case ":execresult":
		return Cmd_CMD_EXEC_RESULT, true
	case ":execrows":
		return Cmd_CMD_EXEC_ROWS, true
	case ":execlastid":
		return Cmd_CMD_EXEC_LAST_ID, true
	case ":copyfrom":
		return Cmd_CMD_COPY_FROM, true
	case ":batchexec":
		return Cmd_CMD_BATCH_EXEC, true
	case ":batchmany":
		return Cmd_CMD_BATCH_MANY, true
	case ":batchone":
		return Cmd_CMD_BATCH_ONE, true
	default:
		return Cmd_CMD_UNSPECIFIED, false
	}
}

// ParseNameAndCmd parses a single comment line like "-- name: ListAuthors :many" or
// "# name: GetUser :one" or "/* name: CreateFoo :exec */". Returns (name, cmd, true) on success.
// Uses metadata.ParseQueryNameAndType so behavior matches built-in engines exactly.
func ParseNameAndCmd(line string, syntax CommentSyntax) (name string, cmd Cmd, ok bool) {
	name, cmdStr, err := metadata.ParseQueryNameAndType(line+"\n", toMetadataCommentSyntax(syntax))
	if err != nil || name == "" {
		return "", Cmd_CMD_UNSPECIFIED, false
	}
	c, ok := CmdFromString(cmdStr)
	if !ok {
		return "", Cmd_CMD_UNSPECIFIED, false
	}
	return name, c, true
}

// QueryBlock is one named query block (from " name: X :cmd" to the next such line or EOF).
type QueryBlock struct {
	Name string
	Cmd  Cmd
	SQL  string
}

// QueryBlocks splits content into named query blocks. Each block runs from a
// " name: X :cmd" line to the next such line (or EOF). Returns one entry per block.
// Uses metadata.QueryBlocks so rules match the built-in engines exactly.
func QueryBlocks(content string, syntax CommentSyntax) ([]QueryBlock, error) {
	blocks, err := metadata.QueryBlocks(content, toMetadataCommentSyntax(syntax))
	if err != nil {
		return nil, err
	}
	out := make([]QueryBlock, 0, len(blocks))
	for _, b := range blocks {
		c, ok := CmdFromString(b.Cmd)
		if !ok {
			continue
		}
		out = append(out, QueryBlock{Name: b.Name, Cmd: c, SQL: b.SQL})
	}
	return out, nil
}

// StatementMeta builds a Statement with name, cmd, and sql set. Parameters and columns
// are left empty for the plugin to fill. Plugins can use this after splitting query.sql
// with QueryBlocks and then attach their own parameters/columns per block.
func StatementMeta(name string, cmd Cmd, sql string) *Statement {
	return &Statement{
		Name: name,
		Cmd:  cmd,
		Sql:  sql,
	}
}
