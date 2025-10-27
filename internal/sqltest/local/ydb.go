package local

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	migrate "github.com/sqlc-dev/sqlc/internal/migrations"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func YDB(t *testing.T, migrations []string) *ydb.Driver {
	return link_YDB(t, migrations, true, false)
}

func YDBTLS(t *testing.T, migrations []string) *ydb.Driver {
	return link_YDB(t, migrations, true, true)
}

func ReadOnlyYDB(t *testing.T, migrations []string) *ydb.Driver {
	return link_YDB(t, migrations, false, false)
}

func ReadOnlyYDBTLS(t *testing.T, migrations []string) *ydb.Driver {
	return link_YDB(t, migrations, false, true)
}

func link_YDB(t *testing.T, migrations []string, rw bool, tls bool) *ydb.Driver {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Helper()

	dbuiri := os.Getenv("YDB_SERVER_URI")
	if dbuiri == "" {
		t.Skip("YDB_SERVER_URI is empty")
	}

	baseDB := os.Getenv("YDB_DATABASE")
	if baseDB == "" {
		baseDB = "/local"
	}

	var connectionString string
	if tls {
		connectionString = fmt.Sprintf("grpcs://%s%s", dbuiri, baseDB)
	} else {
		connectionString = fmt.Sprintf("grpc://%s%s", dbuiri, baseDB)
	}

	db, err := ydb.Open(ctx, connectionString)
	if err != nil {
		t.Fatalf("failed to open YDB connection: %s", err)
	}

	files, err := sqlpath.Glob(migrations)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		blob, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		stmt := migrate.RemoveRollbackStatements(string(blob))

		statements := strings.Split(stmt, ";")
		for _, singleStmt := range statements {
			singleStmt = strings.TrimSpace(singleStmt)
			if singleStmt == "" {
				continue
			}
			err = db.Query().Exec(ctx, singleStmt, query.WithTxControl(query.EmptyTxControl()))
			if err != nil {
				t.Fatalf("failed to apply migration: %s\nSQL: %s", err, singleStmt)
			}
		}
	}

	return db
}
