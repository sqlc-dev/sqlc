package rewrite

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/engine/postgresql"
)

func TestNamedParametersInOrderBy(t *testing.T) {
	query := `
SELECT ID
FROM Sequence
WHERE SeriesID = sqlc.arg(series_id)
ORDER BY (Name = sqlc.arg(name)) DESC, ID
LIMIT 1;
`

	p := postgresql.NewParser()

	stmts, err := p.Parse(strings.NewReader(query))
	if err != nil {
		t.Fatal(err)
	}

raw := stmts[0].Raw

_, params, _ := NamedParameters(
	config.EngineSQLite,
	raw,
	map[int]bool{},
	false,
)

if params == nil {
	t.Fatalf("params should not be nil")
}
}
