package authors

import (
	"context"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sqltest/local"
	_ "github.com/ydb-platform/ydb-go-sdk/v3"
)

func TestAuthors(t *testing.T) {
	ctx := context.Background()

	test := local.YDB(t, []string{"schema.sql"})
	defer test.DB.Close()

	q := New(test.DB)

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
			bio := "NULL"
			if a.Bio.Valid {
				bio = a.Bio.String
			}
			t.Logf("- ID: %d | Name: %s | Bio: %s", a.ID, a.Name, bio)
		}
	})

	t.Run("GetAuthor", func(t *testing.T) {
		singleAuthor, err := q.GetAuthor(ctx, 10)
		if err != nil {
			t.Fatal(err)
		}
		bio := "NULL"
		if singleAuthor.Bio.Valid {
			bio = singleAuthor.Bio.String
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
			bio := "NULL"
			if a.Bio.Valid {
				bio = a.Bio.String
			}
			t.Logf("- ID: %d | Name: %s | Bio: %s", a.ID, a.Name, bio)
		}
	})

	t.Run("ListAuthorsWithIdModulo", func(t *testing.T) {
		authors, err := q.ListAuthorsWithIdModulo(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if len(authors) == 0 {
			t.Fatal("expected at least one author with even ID, got none")
		}
		t.Log("Authors with even IDs:")
		for _, a := range authors {
			bio := "NULL"
			if a.Bio.Valid {
				bio = a.Bio.String
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
			bio := "NULL"
			if a.Bio.Valid {
				bio = a.Bio.String
			}
			t.Logf("- ID: %d | Name: %s | Bio: %s", a.ID, a.Name, bio)
		}
	})
}
