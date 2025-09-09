package authors

import (
	"context"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sqltest/local"
	_ "github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
)

func ptr(s string) *string {
	return &s
}

func TestAuthors(t *testing.T) {
	ctx := context.Background()

	db := local.YDB(t, []string{"schema.sql"})
	defer db.Close(ctx)

	q := New(db.Query())

	t.Run("InsertAuthors", func(t *testing.T) {
		authorsToInsert := []CreateOrUpdateAuthorParams{
			{P0: 1, P1: "Leo Tolstoy", P2: ptr("Russian writer, author of \"War and Peace\"")},
			{P0: 2, P1: "Alexander Pushkin", P2: ptr("Author of \"Eugene Onegin\"")},
			{P0: 3, P1: "Alexander Pushkin", P2: ptr("Russian poet, playwright, and prose writer")},
			{P0: 4, P1: "Fyodor Dostoevsky", P2: ptr("Author of \"Crime and Punishment\"")},
			{P0: 5, P1: "Nikolai Gogol", P2: ptr("Author of \"Dead Souls\"")},
			{P0: 6, P1: "Anton Chekhov", P2: nil},
			{P0: 7, P1: "Ivan Turgenev", P2: ptr("Author of \"Fathers and Sons\"")},
			{P0: 8, P1: "Mikhail Lermontov", P2: nil},
			{P0: 9, P1: "Daniil Kharms", P2: ptr("Absurdist, writer and poet")},
			{P0: 10, P1: "Maxim Gorky", P2: ptr("Author of \"At the Bottom\"")},
			{P0: 11, P1: "Vladimir Mayakovsky", P2: nil},
			{P0: 12, P1: "Sergei Yesenin", P2: ptr("Russian lyric poet")},
			{P0: 13, P1: "Boris Pasternak", P2: ptr("Author of \"Doctor Zhivago\"")},
		}

		for _, author := range authorsToInsert {
			if err := q.CreateOrUpdateAuthor(ctx, author, query.WithIdempotent()); err != nil {
				t.Fatalf("failed to insert author %q: %v", author.P1, err)
			}
		}
	})

	t.Run("ListAuthors", func(t *testing.T) {
		authors, err := q.ListAuthors(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if len(authors) == 0 {
			t.Fatal("expected at least one author, got none")
		}
		t.Log("Authors:")
		for _, a := range authors {
			bio := "Null"
			if a.Bio != nil {
				bio = *a.Bio
			}
			t.Logf("- ID: %d | Name: %s | Bio: %s", a.ID, a.Name, bio)
		}
	})

	t.Run("GetAuthor", func(t *testing.T) {
		singleAuthor, err := q.GetAuthor(ctx, 10)
		if err != nil {
			t.Fatal(err)
		}
		bio := "Null"
		if singleAuthor.Bio != nil {
			bio = *singleAuthor.Bio
		}
		t.Logf("- ID: %d | Name: %s | Bio: %s", singleAuthor.ID, singleAuthor.Name, bio)
	})

	t.Run("Delete All Authors", func(t *testing.T) {
		var i uint64
		for i = 1; i <= 13; i++ {
			if err := q.DeleteAuthor(ctx, i, query.WithIdempotent()); err != nil {
				t.Fatalf("failed to delete author: %v", err)
			}
		}
		authors, err := q.ListAuthors(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if len(authors) != 0 {
			t.Fatalf("expected no authors, got %d", len(authors))
		}
	})
	
	t.Run("Drop Table Authors", func(t *testing.T) {
		err := q.DropTable(ctx)
		if err != nil {
			t.Fatal(err)
		}
	})
}
