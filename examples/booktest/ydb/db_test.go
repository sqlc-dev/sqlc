package booktest

import (
	"context"
	"testing"
	"time"

	"github.com/sqlc-dev/sqlc/internal/sqltest/local"
	_ "github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
)

const (
	BooksBookTypeFICTION    string = "FICTION"
	BooksBookTypeNONFICTION string = "NONFICTION"
)

func TestBooks(t *testing.T) {
	ctx := context.Background()
	db := local.YDB(t, []string{"schema.sql"})
	defer db.Close(ctx)

	dq := New(db.Query())

	// create an author
	a, err := dq.CreateAuthor(ctx, "Unknown Master")
	if err != nil {
		t.Fatal(err)
	}

	// save first book
	now := time.Now()
	_, err = dq.CreateBook(ctx, CreateBookParams{
		AuthorID:  a.AuthorID,
		Isbn:      "1",
		Title:     "my book title",
		BookType:  BooksBookTypeFICTION,
		Year:      2016,
		Available: now,
		Tag:       "",
	}, query.WithIdempotent())
	if err != nil {
		t.Fatal(err)
	}

	// save second book
	b1, err := dq.CreateBook(ctx, CreateBookParams{
		AuthorID:  a.AuthorID,
		Isbn:      "2",
		Title:     "the second book",
		BookType:  BooksBookTypeFICTION,
		Year:      2016,
		Available: now,
		Tag:       "unique",
	}, query.WithIdempotent())
	if err != nil {
		t.Fatal(err)
	}

	// update the title and tags
	err = dq.UpdateBook(ctx, UpdateBookParams{
		BookID: b1.BookID,
		Title:  "changed second title",
		Tag:    "disastor",
	}, query.WithIdempotent())
	if err != nil {
		t.Fatal(err)
	}

	// save third book
	_, err = dq.CreateBook(ctx, CreateBookParams{
		AuthorID:  a.AuthorID,
		Isbn:      "3",
		Title:     "the third book",
		BookType:  BooksBookTypeFICTION,
		Year:      2001,
		Available: now,
		Tag:       "cool",
	}, query.WithIdempotent())
	if err != nil {
		t.Fatal(err)
	}

	// save fourth book
	b3, err := dq.CreateBook(ctx, CreateBookParams{
		AuthorID:  a.AuthorID,
		Isbn:      "4",
		Title:     "4th place finisher",
		BookType:  BooksBookTypeFICTION,
		Year:      2011,
		Available: now,
		Tag:       "other",
	}, query.WithIdempotent())
	if err != nil {
		t.Fatal(err)
	}

	// upsert, changing ISBN and title
	err = dq.UpdateBookISBN(ctx, UpdateBookISBNParams{
		BookID: b3.BookID,
		Isbn:   "NEW ISBN",
		Title:  "never ever gonna finish, a quatrain",
		Tag:    "someother",
	}, query.WithIdempotent())
	if err != nil {
		t.Fatal(err)
	}

	// retrieve first book
	books0, err := dq.BooksByTitleYear(ctx, BooksByTitleYearParams{
		Title: "my book title",
		Year:  2016,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, book := range books0 {
		t.Logf("Book %d (%s): %s available: %s\n", book.BookID, book.BookType, book.Title, book.Available.Format(time.RFC822Z))
		author, err := dq.GetAuthor(ctx, book.AuthorID)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Book %d author: %s\n", book.BookID, author.Name)
	}

	// find a book with either "cool" or "other" or "someother" tag
	t.Logf("---------\nTag search results:\n")
	res, err := dq.BooksByTags(ctx, []string{"cool", "other", "someother"})
	if err != nil {
		t.Fatal(err)
	}
	for _, ab := range res {
		name := ""
		if ab.Name != nil {
			name = *ab.Name
		}
		t.Logf("Book %d: '%s', Author: '%s', ISBN: '%s' Tag: '%s'\n", ab.BookID, ab.Title, name, ab.Isbn, ab.Tag)
	}

	// get book 4 and delete
	b5, err := dq.GetBook(ctx, b3.BookID)
	if err != nil {
		t.Fatal(err)
	}
	if err := dq.DeleteBook(ctx, b5.BookID); err != nil {
		t.Fatal(err)
	}
}
