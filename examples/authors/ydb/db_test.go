package authors

import (
	"context"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sqltest/local"
	_ "github.com/ydb-platform/ydb-go-sdk/v3"
)

func ptr(s string) *string {
	return &s
}

func TestAuthors(t *testing.T) {
	ctx := context.Background()

	test := local.YDB(t, []string{"schema.sql"})
	defer test.DB.Close()

	q := New(test.DB)

	t.Run("InsertAuthors", func(t *testing.T) {
		authorsToInsert := []CreateOrUpdateAuthorParams{
			{P0: 1, P1: "Лев Толстой", P2: ptr("Русский писатель, автор \"Война и мир\"")},
			{P0: 2, P1: "Александр Пушкин", P2: ptr("Автор \"Евгения Онегина\"")},
			{P0: 3, P1: "Александр Пушкин", P2: ptr("Русский поэт, драматург и прозаик")},
			{P0: 4, P1: "Фёдор Достоевский", P2: ptr("Автор \"Преступление и наказание\"")},
			{P0: 5, P1: "Николай Гоголь", P2: ptr("Автор \"Мёртвые души\"")},
			{P0: 6, P1: "Антон Чехов", P2: nil},
			{P0: 7, P1: "Иван Тургенев", P2: ptr("Автор \"Отцы и дети\"")},
			{P0: 8, P1: "Михаил Лермонтов", P2: nil},
			{P0: 9, P1: "Даниил Хармс", P2: ptr("Абсурдист, писатель и поэт")},
			{P0: 10, P1: "Максим Горький", P2: ptr("Автор \"На дне\"")},
			{P0: 11, P1: "Владимир Маяковский", P2: nil},
			{P0: 12, P1: "Сергей Есенин", P2: ptr("Русский лирик")},
			{P0: 13, P1: "Борис Пастернак", P2: ptr("Автор \"Доктор Живаго\"")},
		}

		for _, author := range authorsToInsert {
			if _, err := q.CreateOrUpdateAuthor(ctx, author); err != nil {
				t.Fatalf("failed to insert author %q: %v", author.P1, err)
			}
		}
	})

	t.Run("CreateOrUpdateAuthorReturningBio", func(t *testing.T) {
		newBio := "Обновленная биография автора"
		arg := CreateOrUpdateAuthorReturningBioParams{
			P0: 3,
			P1: "Тестовый Автор",
			P2: &newBio,
		}

		returnedBio, err := q.CreateOrUpdateAuthorReturningBio(ctx, arg)
		if err != nil {
			t.Fatalf("failed to create or update author: %v", err)
		}

		if returnedBio == nil {
			t.Fatal("expected non-nil bio, got nil")
		}
		if *returnedBio != newBio {
			t.Fatalf("expected bio %q, got %q", newBio, *returnedBio)
		}

		t.Logf("Author created or updated successfully with bio: %s", *returnedBio)
	})

	t.Run("Update Author", func(t *testing.T) {
		arg := UpdateAuthorByIDParams{
			P0: "Максим Горький",
			P1: ptr("Обновленная биография"),
			P2: 10,
		}

		singleAuthor, err := q.UpdateAuthorByID(ctx, arg)
		if err != nil {
			t.Fatal(err)
		}
		bio := "Null"
		if singleAuthor.Bio != nil {
			bio = *singleAuthor.Bio
		}
		t.Logf("- ID: %d | Name: %s | Bio: %s", singleAuthor.ID, singleAuthor.Name, bio)
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

	t.Run("GetAuthorByName", func(t *testing.T) {
		authors, err := q.GetAuthorsByName(ctx, "Александр Пушкин")
		if err != nil {
			t.Fatal(err)
		}
		if len(authors) == 0 {
			t.Fatal("expected at least one author with this name, got none")
		}
		t.Log("Authors with this name:")
		for _, a := range authors {
			bio := "Null"
			if a.Bio != nil {
				bio = *a.Bio
			}
			t.Logf("- ID: %d | Name: %s | Bio: %s", a.ID, a.Name, bio)
		}
	})

	t.Run("ListAuthorsWithNullBio", func(t *testing.T) {
		authors, err := q.ListAuthorsWithNullBio(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if len(authors) == 0 {
			t.Fatal("expected at least one author with NULL bio, got none")
		}
		t.Log("Authors with NULL bio:")
		for _, a := range authors {
			bio := "Null"
			if a.Bio != nil {
				bio = *a.Bio
			}
			t.Logf("- ID: %d | Name: %s | Bio: %s", a.ID, a.Name, bio)
		}
	})

	t.Run("Delete All Authors", func(t *testing.T) {
		var i uint64
		for i = 1; i <= 13; i++ {
			if err := q.DeleteAuthor(ctx, i); err != nil {
				t.Fatalf("failed to delete authors: %v", err)
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
}
