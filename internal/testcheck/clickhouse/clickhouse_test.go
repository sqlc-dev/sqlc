package clickhouse

import (
	"context"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/testcheck/endtoend"
)

// TestEndToEnd verifies every ClickHouse analyze case against ClickHouse.
// It skips unless the clickhouse binary is installed.
func TestEndToEnd(t *testing.T) {
	binary, err := Locate()
	if err != nil {
		t.Skip(err)
	}
	cases, err := endtoend.Cases(Engine)
	if err != nil {
		t.Fatal(err)
	}
	if len(cases) == 0 {
		t.Fatal("no clickhouse analyze cases found")
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			diff, err := Check(context.Background(), binary, c)
			if err != nil {
				t.Fatal(err)
			}
			if diff != "" {
				t.Errorf("%s does not match what ClickHouse reports (-committed +clickhouse):\n%s", c.Output, diff)
			}
		})
	}
}
