//go:build examples
// +build examples

package batch

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/sqlc-dev/sqlc/internal/sqltest/local"
)

func TestBatchBooks(t *testing.T) {
	uri := local.PostgreSQL(t, []string{"schema.sql"})

	ctx := context.Background()

	db, err := pgx.Connect(ctx, uri)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close(ctx)

	dq := New(db)

	// create an author
	a, err := dq.CreateAuthor(ctx, "Unknown Master")
	if err != nil {
		t.Fatal(err)
	}

	now := pgtype.Timestamptz{Time: time.Now(), Valid: true}

	// batch insert new books
	newBooksParams := []CreateBookParams{
		{
			AuthorID:  a.AuthorID,
			Isbn:      "1",
			Title:     "my book title",
			BookType:  BookTypeFICTION,
			Year:      2016,
			Available: now,
			Tags:      []string{},
		},
		{
			AuthorID:  a.AuthorID,
			Isbn:      "2",
			Title:     "the second book",
			BookType:  BookTypeFICTION,
			Year:      2016,
			Available: now,
			Tags:      []string{"cool", "unique"},
		},
		{
			AuthorID:  a.AuthorID,
			Isbn:      "3",
			Title:     "the third book",
			BookType:  BookTypeFICTION,
			Year:      2001,
			Available: now,
			Tags:      []string{"cool"},
		},
		{
			AuthorID:  a.AuthorID,
			Isbn:      "4",
			Title:     "4th place finisher",
			BookType:  BookTypeNONFICTION,
			Year:      2011,
			Available: now,
			Tags:      []string{"other"},
		},
	}
	newBooks := make([]Book, len(newBooksParams))
	var cnt int
	dq.CreateBook(ctx, newBooksParams).QueryRow(func(i int, b Book, err error) {
		if err != nil {
			t.Fatalf("failed inserting book (%s): %s", b.Title, err)
		}
		newBooks[i] = b
		cnt = i
	})
	// first i was 0, so add 1
	cnt++
	numBooksExpected := len(newBooks)
	if cnt != numBooksExpected {
		t.Fatalf("expected to insert %d books; got %d", numBooksExpected, cnt)
	}

	// batch update the title and tags
	updateBooksParams := []UpdateBookParams{
		{
			BookID: newBooks[1].BookID,
			Title:  "changed second title",
			Tags:   []string{"cool", "disastor"},
		},
	}
	dq.UpdateBook(ctx, updateBooksParams).Exec(func(i int, err error) {
		if err != nil {
			t.Fatalf("error updating book %d: %s", updateBooksParams[i].BookID, err)
		}
	})

	// batch many to retrieve books by year
	selectBooksByTitleYearParams := []int32{2001, 2016}
	var books0 []Book
	dq.BooksByYear(ctx, selectBooksByTitleYearParams).Query(func(i int, books []Book, err error) {
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("num books for %d: %d", selectBooksByTitleYearParams[i], len(books))
		books0 = append(books0, books...)
	})

	for _, book := range books0 {
		t.Logf("Book %d (%s): %s available: %s\n", book.BookID, book.BookType, book.Title, book.Available.Time.Format(time.RFC822Z))
		author, err := dq.GetAuthor(ctx, book.AuthorID)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Book %d author: %s\n", book.BookID, author.Name)
	}

	// batch delete books
	deleteBooksParams := make([]int32, len(newBooks))
	for i, book := range newBooks {
		deleteBooksParams[i] = book.BookID
	}
	batchDelete := dq.DeleteBook(ctx, deleteBooksParams)
	numDeletesProcessed := 0
	wantNumDeletesProcessed := 2
	batchDelete.Exec(func(i int, err error) {
		if err != nil && err.Error() != "batch already closed" {
			t.Fatalf("error deleting book %d: %s", deleteBooksParams[i], err)
		}

		if err == nil {
			numDeletesProcessed++
		}

		if i == wantNumDeletesProcessed-1 {
			// close batch operation before processing all errors from delete operation
			if err := batchDelete.Close(); err != nil {
				t.Fatalf("failed to close batch operation: %s", err)
			}
		}
	})
	if numDeletesProcessed != wantNumDeletesProcessed {
		t.Fatalf("expected Close to short-circuit record processing (expected %d; got %d)", wantNumDeletesProcessed, numDeletesProcessed)
	}
}
