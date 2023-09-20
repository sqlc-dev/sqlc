//go:build examples
// +build examples

package enums

import (
	"context"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sqltest"
)

func TestUsers(t *testing.T) {
	sdb, cleanup := sqltest.PostgreSQLPgxV5(t, []string{"schema.sql", "fixtures.sql"})
	defer cleanup()

	ctx := context.Background()
	db := New(sdb)

	for _, name := range EnumNames {
		enumType, err := sdb.LoadType(ctx, name) // register enum
		if err != nil {
			t.Fatal(err)
		}
		sdb.TypeMap().RegisterType(enumType)

		enumType, err = sdb.LoadType(ctx, "_"+name) // register slice of enum
		if err != nil {
			t.Fatal(err)
		}
		sdb.TypeMap().RegisterType(enumType)
	}

	err := db.UserCreate(ctx, UserCreateParams{
		FirstName: "Alex",
		LastName:  "Brown",
		Age:       25,
		ShirtSize: SizeMedium,
	})
	if err != nil {
		t.Fatal(err)
	}

	// list all users
	users, err := db.ListUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(users)

	// list users with large shirt_size
	usersWithLarge, err := db.ListUsersByShirtSizes(ctx, []Size{SizeLarge})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(usersWithLarge)
}
