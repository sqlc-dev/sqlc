//go:build examples
// +build examples

package update

import (
	"context"
	"testing"

	"github.com/kyleconroy/sqlc/internal/sqltest"
)

func TestUpdate(t *testing.T) {
	sdb, cleanup := sqltest.MySQL(t, []string{"schema.sql"})
	defer cleanup()

	ctx := context.Background()
	db := New(sdb)

	_, err := db.CreateT1(ctx, CreateT1Params{
		UserID: int32(2),
		Name:   "",
	})
	if err != nil {
		t.Fatal(err)
	}

	// get the data we just inserted
	oldData, err := db.GetT1(ctx, int32(2))
	if err != nil {
		t.Fatal(err)
	}

	if oldData.Name != "" {
		t.Fatal("create fail")
	}

	_, err = db.CreateT2(ctx, CreateT2Params{
		Email: "test@gmail.com",
		Name:  "test",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.CreateT3(ctx, CreateT3Params{
		UserID: int32(2),
		Email:  "test@gmail.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = db.UpdateAll(ctx)
	if err != nil {
		t.Fatal(err)
	}

	newData, err := db.GetT1(ctx, int32(2))
	if err != nil {
		t.Fatal(err)
	}

	if newData.Name != "test" {
		t.Fatal("update fail")
	}
}
