package hosted

import (
	"context"
	"os"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/quickdb"
	pb "github.com/sqlc-dev/sqlc/internal/quickdb/v1"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
)

func PostgreSQL(t *testing.T, migrations []string) string {
	ctx := context.Background()
	t.Helper()

	once.Do(func() {
		if err := initClient(); err != nil {
			t.Log(err)
		}
	})

	if client == nil {
		t.Skip("client init failed")
	}

	var seed []string
	files, err := sqlpath.Glob(migrations)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		blob, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		seed = append(seed, string(blob))
	}

	resp, err := client.CreateEphemeralDatabase(ctx, &pb.CreateEphemeralDatabaseRequest{
		Engine:     "postgresql",
		Region:     quickdb.GetClosestRegion(),
		Migrations: seed,
	})
	if err != nil {
		t.Fatalf("region %s: %s", quickdb.GetClosestRegion(), err)
	}

	t.Cleanup(func() {
		_, err = client.DropEphemeralDatabase(ctx, &pb.DropEphemeralDatabaseRequest{
			DatabaseId: resp.DatabaseId,
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	return resp.Uri
}
