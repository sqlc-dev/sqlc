package local

import (
	"context"
	"database/sql"
	"fmt"
	"hash/fnv"
	"math/rand"
	"net"
	"os"
	"testing"
	"time"

	migrate "github.com/sqlc-dev/sqlc/internal/migrations"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
	"github.com/ydb-platform/ydb-go-sdk/v3"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func YDB(t *testing.T, migrations []string) TestYDB {
	return link_YDB(t, migrations, true)
}

func ReadOnlyYDB(t *testing.T, migrations []string) TestYDB {
	return link_YDB(t, migrations, false)
}

type TestYDB struct {
	DB     *sql.DB
	Prefix string
}

func link_YDB(t *testing.T, migrations []string, rw bool) TestYDB {
	t.Helper()

	// 1) Контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbuiri := os.Getenv("YDB_SERVER_URI")
	if dbuiri == "" {
		t.Skip("YDB_SERVER_URI is empty")
	}
	host, _, err := net.SplitHostPort(dbuiri)
	if err != nil {
		t.Fatalf("invalid YDB_SERVER_URI: %q", dbuiri)
	}

	baseDB := os.Getenv("YDB_DATABASE")
	if baseDB == "" {
		baseDB = "/local"
	}

	// собираем миграции
	var seed []string
	files, err := sqlpath.Glob(migrations)
	if err != nil {
		t.Fatal(err)
	}
	h := fnv.New64()
	for _, f := range files {
		blob, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		h.Write(blob)
		seed = append(seed, migrate.RemoveRollbackStatements(string(blob)))
	}

	var name string
	if rw {
		// name = fmt.Sprintf("sqlc_test_%s", id())
		name = fmt.Sprintf("sqlc_test_%s", "test_new")
	} else {
		name = fmt.Sprintf("sqlc_test_%x", h.Sum(nil))
	}
	prefix := fmt.Sprintf("%s/%s", baseDB, name)

	// 2) Открываем драйвер к корню "/"
	rootDSN := fmt.Sprintf("grpc://%s?database=%s", dbuiri, baseDB)
	t.Logf("→ Opening root driver: %s", rootDSN)
	driver, err := ydb.Open(ctx, rootDSN,
		ydb.WithInsecure(),
		ydb.WithDiscoveryInterval(time.Hour),
		ydb.WithNodeAddressMutator(func(_ string) string {
			return host
		}),
	)
	if err != nil {
		t.Fatalf("failed to open root YDB connection: %s", err)
	}

	connector, err := ydb.Connector(
		driver,
		ydb.WithTablePathPrefix(prefix),
		ydb.WithAutoDeclare(),
	)
	if err != nil {
		t.Fatalf("failed to create connector: %s", err)
	}

	db := sql.OpenDB(connector)

	t.Log("→ Applying migrations to prefix: ", prefix)

	schemeCtx := ydb.WithQueryMode(ctx, ydb.SchemeQueryMode)
	for _, stmt := range seed {
		_, err := db.ExecContext(schemeCtx, stmt)
		if err != nil {
			t.Fatalf("failed to apply migration: %s\nSQL: %s", err, stmt)
		}
	}
	return TestYDB{DB: db, Prefix: prefix}
}
