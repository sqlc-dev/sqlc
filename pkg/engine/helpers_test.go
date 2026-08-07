package engine

import (
	"strings"
	"testing"
)

func TestCmdToString(t *testing.T) {
	tests := []struct {
		c    Cmd
		want string
	}{
		{Cmd_CMD_UNSPECIFIED, ""},
		{Cmd_CMD_ONE, ":one"},
		{Cmd_CMD_MANY, ":many"},
		{Cmd_CMD_EXEC, ":exec"},
		{Cmd_CMD_EXEC_RESULT, ":execresult"},
		{Cmd_CMD_EXEC_ROWS, ":execrows"},
		{Cmd_CMD_EXEC_LAST_ID, ":execlastid"},
		{Cmd_CMD_COPY_FROM, ":copyfrom"},
		{Cmd_CMD_BATCH_EXEC, ":batchexec"},
		{Cmd_CMD_BATCH_MANY, ":batchmany"},
		{Cmd_CMD_BATCH_ONE, ":batchone"},
	}
	for _, tt := range tests {
		if got := CmdToString(tt.c); got != tt.want {
			t.Errorf("CmdToString(%v) = %q, want %q", tt.c, got, tt.want)
		}
	}
}

func TestCmdFromString(t *testing.T) {
	valid := []string{":one", ":many", ":exec", ":execresult", ":execrows", ":execlastid", ":copyfrom", ":batchexec", ":batchmany", ":batchone"}
	for _, s := range valid {
		c, ok := CmdFromString(s)
		if !ok {
			t.Errorf("CmdFromString(%q) ok=false", s)
			continue
		}
		if CmdToString(c) != s {
			t.Errorf("CmdFromString(%q) -> %v -> CmdToString = %q", s, c, CmdToString(c))
		}
	}
	_, ok := CmdFromString(":invalid")
	if ok {
		t.Error("CmdFromString(:invalid) ok=true")
	}
}

func TestParseNameAndCmd(t *testing.T) {
	dash := CommentSyntax{Dash: true}
	hash := CommentSyntax{Hash: true}
	slash := CommentSyntax{SlashStar: true}

	t.Run("Parse", func(t *testing.T) {
		for _, line := range []string{
			`-- name: CreateFoo, :one`,
			`-- name: 9Foo :one`,
			`-- name: CreateFoo :two`,
			`-- name: CreateFoo`,
			`-- name: CreateFoo :one something`,
			`-- name: `,
			`--name: CreateFoo :one`,
			`--name CreateFoo :one`,
		} {
			t.Run(line, func(t *testing.T) {
				if _, _, ok := ParseNameAndCmd(line, dash); ok {
					t.Errorf("expected invalid: %q", line)
				}
			})
		}
	})

	t.Run("CommentSyntax", func(t *testing.T) {
		for line, syn := range map[string]CommentSyntax{
			`-- name: CreateFoo :one`:    dash,
			`# name: CreateFoo :one`:     hash,
			`/* name: CreateFoo :one */`: slash,
		} {
			t.Run(line, func(t *testing.T) {
				name, cmd, ok := ParseNameAndCmd(line, syn)
				if !ok {
					t.Errorf("expected valid %q", line)
					return
				}
				if name != "CreateFoo" {
					t.Errorf("ParseNameAndCmd(%q) name = %q", line, name)
				}
				if cmd != Cmd_CMD_ONE {
					t.Errorf("ParseNameAndCmd(%q) cmd = %v", line, cmd)
				}
			})
		}
	})

	t.Run("Many", func(t *testing.T) {
		name, cmd, ok := ParseNameAndCmd("-- name: ListAuthors :many", dash)
		if !ok || name != "ListAuthors" || cmd != Cmd_CMD_MANY {
			t.Errorf("-- name: ListAuthors :many -> name=%q cmd=%v ok=%v", name, cmd, ok)
		}
	})
}

func TestQueryBlocks(t *testing.T) {
	syntax := CommentSyntax{Dash: true, SlashStar: true}
	content := `-- name: GetUser :one
SELECT id, name FROM users WHERE id = $1

-- name: ListUsers :many
SELECT id, name FROM users ORDER BY id
`
	blocks, err := QueryBlocks(content, syntax)
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 2 {
		t.Fatalf("len(blocks)=%d, want 2", len(blocks))
	}
	if blocks[0].Name != "GetUser" || blocks[0].Cmd != Cmd_CMD_ONE {
		t.Errorf("block0: Name=%q Cmd=%v", blocks[0].Name, blocks[0].Cmd)
	}
	if blocks[1].Name != "ListUsers" || blocks[1].Cmd != Cmd_CMD_MANY {
		t.Errorf("block1: Name=%q Cmd=%v", blocks[1].Name, blocks[1].Cmd)
	}
	if !strings.Contains(blocks[0].SQL, "SELECT id, name FROM users WHERE id = $1") {
		t.Errorf("block0 SQL missing expected fragment: %s", blocks[0].SQL)
	}
	if !strings.Contains(blocks[1].SQL, "SELECT id, name FROM users ORDER BY id") {
		t.Errorf("block1 SQL missing expected fragment: %s", blocks[1].SQL)
	}
}

func TestStatementMeta(t *testing.T) {
	st := StatementMeta("GetUser", Cmd_CMD_ONE, "SELECT 1")
	if st.Name != "GetUser" || st.Cmd != Cmd_CMD_ONE || st.Sql != "SELECT 1" {
		t.Errorf("StatementMeta: Name=%q Cmd=%v Sql=%q", st.Name, st.Cmd, st.Sql)
	}
	if st.Parameters != nil || st.Columns != nil {
		t.Error("StatementMeta should leave Parameters/Columns nil")
	}
}
