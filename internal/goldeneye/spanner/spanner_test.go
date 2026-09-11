package spanner

import (
	"context"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/dialect"
	"github.com/sqlc-dev/sqlc/internal/goldeneye/endtoend"
)

// TestDialect verifies the committed GoogleSQL dialect against what a
// Spanner server reports. It skips unless SPANNER_SERVER_URI names one.
func TestDialect(t *testing.T) {
	endpoint, err := Locate()
	if err != nil {
		t.Skip(err)
	}
	ctx := context.Background()
	version, err := Version(ctx, endpoint)
	if err != nil {
		t.Fatal(err)
	}
	files, err := Generate(ctx, endpoint)
	if err != nil {
		t.Fatal(err)
	}
	dir, err := dialect.Dir(Dir)
	if err != nil {
		t.Fatal(err)
	}
	report, err := dialect.Check(dir, files)
	if err != nil {
		t.Fatal(err)
	}
	if report != "" {
		t.Errorf("%s does not match what %s reports:\n%s", dir, version, report)
	}
}

// TestAnalyzeCases verifies every GoogleSQL analyze case under
// internal/endtoend/testdata against what a Spanner server reports. It
// skips unless SPANNER_SERVER_URI names one.
func TestAnalyzeCases(t *testing.T) {
	endpoint, err := Locate()
	if err != nil {
		t.Skip(err)
	}
	cases, err := endtoend.Cases(Cases)
	if err != nil {
		t.Fatal(err)
	}
	if len(cases) == 0 {
		t.Fatal("no googlesql analyze cases found")
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			diff, err := Check(context.Background(), endpoint, c)
			if err != nil {
				t.Fatal(err)
			}
			if diff != "" {
				t.Errorf("%s does not match what Spanner reports (-committed +spanner):\n%s", c.Output, diff)
			}
		})
	}
}
