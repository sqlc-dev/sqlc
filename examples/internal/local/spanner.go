package local

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// spannerInstance is the instance a Spanner Omni single server provides,
// which every database is created in.
const spannerInstance = "projects/default/instances/default"

// Spanner creates a database on the Spanner Omni server SPANNER_SERVER_URI
// names by its gRPC endpoint, such as localhost:15000, runs the migrations
// as its DDL and returns a data source name for the go-sql-spanner driver
// that connects to it without TLS or credentials, as Omni is reached. The
// test is skipped when no server is named.
func Spanner(t *testing.T, migrations []string) string {
	ctx := context.Background()
	t.Helper()

	endpoint := os.Getenv("SPANNER_SERVER_URI")
	if endpoint == "" {
		t.Skip("SPANNER_SERVER_URI is empty")
	}
	if strings.Contains(endpoint, "://") {
		t.Fatalf("SPANNER_SERVER_URI: %q should be a host:port, the gRPC endpoint of a Spanner Omni server", endpoint)
	}

	ddl, err := statements(migrations)
	if err != nil {
		t.Fatal(err)
	}

	admin, err := database.NewDatabaseAdminClient(ctx,
		option.WithEndpoint(endpoint),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { admin.Close() })

	name := fmt.Sprintf("sqlc_test_%s", strings.ToLower(id()))
	op, err := admin.CreateDatabase(ctx, &databasepb.CreateDatabaseRequest{
		Parent:          spannerInstance,
		CreateStatement: "CREATE DATABASE `" + name + "`",
		ExtraStatements: ddl,
	})
	if err != nil {
		t.Fatal(err)
	}
	db, err := op.Wait(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := admin.DropDatabase(ctx, &databasepb.DropDatabaseRequest{Database: db.Name}); err != nil {
			t.Fatalf("failed cleaning up: %s", err)
		}
	})

	return endpoint + "/" + db.Name + ";usePlainText=true"
}
