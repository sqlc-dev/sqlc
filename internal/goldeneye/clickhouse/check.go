package clickhouse

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/endtoend"
)

// Analyze runs a case's queries through the clickhouse binary and returns
// the analysis in the JSON shape sqlc analyze prints.
func Analyze(ctx context.Context, binary string, c endtoend.Case) ([]byte, error) {
	schema, err := os.ReadFile(c.Schema)
	if err != nil {
		return nil, err
	}
	var fixture []byte
	if c.Fixture != "" {
		if fixture, err = os.ReadFile(c.Fixture); err != nil {
			return nil, err
		}
	}
	src, err := os.ReadFile(c.Query)
	if err != nil {
		return nil, err
	}
	queries, err := parseQueries(string(src))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", c.Query, err)
	}
	out, err := analyze(ctx, local{binary: binary}, string(schema), string(fixture), queries)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Check compares what ClickHouse reports for a case with the output the
// case committed, returning a diff when they differ.
func Check(ctx context.Context, binary string, c endtoend.Case) (string, error) {
	got, err := Analyze(ctx, binary, c)
	if err != nil {
		return "", err
	}
	return c.Compare(got)
}
