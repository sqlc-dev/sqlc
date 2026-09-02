package clickhouse

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// local runs SQL through an ephemeral `clickhouse local` process.
type local struct {
	binary string
}

// resultSet is one JSON-format result printed by clickhouse local. Only
// statements that return rows print one; DDL and INSERT print nothing.
type resultSet struct {
	Meta []resultColumn               `json:"meta"`
	Data []map[string]json.RawMessage `json:"data"`
	Rows int                          `json:"rows"`
}

type resultColumn struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// run executes a multi-statement script in a fresh process with fresh
// storage and returns the result sets in statement order. Any statement
// failing fails the whole run with ClickHouse's own error message.
func (l local) run(ctx context.Context, script string) ([]resultSet, error) {
	dir, err := os.MkdirTemp("", "sqlc-clickhouse-testgen-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	queries := filepath.Join(dir, "queries.sql")
	if err := os.WriteFile(queries, []byte(script), 0o600); err != nil {
		return nil, err
	}

	// stdin must not be inherited: clickhouse local reads it as table data
	// and blocks until it is closed.
	cmd := exec.CommandContext(ctx, l.binary, "local",
		"--multiquery",
		"--queries-file", queries,
		"--output-format", "JSON",
		"--path", filepath.Join(dir, "data"),
	)
	cmd.Stdin = nil
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return nil, errors.New(msg)
	}

	var results []resultSet
	dec := json.NewDecoder(&stdout)
	for {
		var rs resultSet
		err := dec.Decode(&rs)
		if errors.Is(err, io.EOF) {
			return results, nil
		}
		if err != nil {
			return nil, fmt.Errorf("decoding clickhouse local output: %w", err)
		}
		results = append(results, rs)
	}
}
