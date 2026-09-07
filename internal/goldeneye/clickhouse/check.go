package clickhouse

import (
	"context"
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
	queries, err := c.Queries()
	if err != nil {
		return nil, err
	}
	out, err := analyze(ctx, local{binary: binary}, string(schema), string(fixture), queries)
	if err != nil {
		return nil, err
	}
	return endtoend.Encode(out)
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
