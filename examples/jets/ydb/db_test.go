package jets

import (
	"context"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sqltest/local"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
)

func TestJets(t *testing.T) {
	ctx := context.Background()
	db := local.YDB(t, []string{"schema.sql"})
	defer db.Close(ctx)

	q := New(db.Query())

	// insert test data
	pilots := []struct {
		id   int32
		name string
	}{
		{1, "John Doe"},
		{2, "Jane Smith"},
		{3, "Bob Johnson"},
	}

	for _, p := range pilots {
		parameters := ydb.ParamsBuilder()
		parameters = parameters.Param("$id").Int32(p.id)
		parameters = parameters.Param("$name").Text(p.name)
		err := db.Query().Exec(ctx, "UPSERT INTO pilots (id, name) VALUES ($id, $name)",
			query.WithParameters(parameters.Build()),
		)
		if err != nil {
			t.Fatalf("failed to insert pilot: %v", err)
		}
	}

	// list all pilots
	pilotsList, err := q.ListPilots(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Pilots:", pilotsList)

	// count pilots
	count, err := q.CountPilots(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Total pilots: %d", count)

	if count != 3 {
		t.Errorf("expected 3 pilots, got %d", count)
	}

	// delete a pilot
	err = q.DeletePilot(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}

	// count after delete
	count, err = q.CountPilots(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Errorf("expected 2 pilots after delete, got %d", count)
	}
}
