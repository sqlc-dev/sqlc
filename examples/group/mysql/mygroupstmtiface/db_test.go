// +build examples

package mygroupstmtiface_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/kyleconroy/sqlc/examples/group/mysql/mygroupstmtiface"
	"github.com/kyleconroy/sqlc/internal/sqltest"
)

func TestBooks(t *testing.T) {
	db, cleanup := sqltest.MySQL(t, []string{"../sql/schema/schema.sql"})
	defer cleanup()

	ctx := context.Background()

	// create an author
	authorID, err := func() (retID int32, retErr error) {
		dbauthors := mygroupstmtiface.NewAuthors(db)
		dbcreate, err := dbauthors.PrepareCreate(ctx)
		if err != nil {
			retErr = err
			return
		}
		defer func() {
			if err := dbcreate.Close(); err != nil && retErr == nil {
				retErr = err
			}
		}()
		result, err := dbcreate.Exec(ctx, "Unknown Master")
		if err != nil {
			retErr = err
			return
		}
		authorID, err := result.LastInsertId()
		if err != nil {
			retErr = err
			return
		}
		retID = int32(authorID)
		return
	}()
	if err != nil {
		t.Fatal(err)
	}

	// create transaction
	bookThreeID, err := func() (retID int32, retErr error) {
		tx, err := db.Begin()
		if err != nil {
			retErr = err
			return
		}
		defer func() {
			if err := tx.Rollback(); err != nil && err != sql.ErrTxDone && retErr == nil {
				retErr = err
			}
		}()

		txbooks := mygroupstmtiface.NewBooks(tx)

		txcreate, err := txbooks.PrepareCreate(ctx)
		if err != nil {
			retErr = err
			return
		}
		defer func() {
			if err := txcreate.Close(); err != nil && retErr == nil {
				retErr = err
			}
		}()

		// save first book
		now := time.Now()
		_, err = txcreate.Exec(ctx, mygroupstmtiface.BooksCreateParams{
			AuthorID:  authorID,
			Isbn:      "1",
			Title:     "my book title",
			BookType:  mygroupstmtiface.BooksBookTypeFICTION,
			Yr:        2016,
			Available: now,
		})
		if err != nil {
			retErr = err
			return
		}

		// save second book
		result, err := txcreate.Exec(ctx, mygroupstmtiface.BooksCreateParams{
			AuthorID:  authorID,
			Isbn:      "2",
			Title:     "the second book",
			BookType:  mygroupstmtiface.BooksBookTypeFICTION,
			Yr:        2016,
			Available: now,
			Tags:      "cool,unique",
		})
		if err != nil {
			retErr = err
			return
		}
		bookOneID, err := result.LastInsertId()
		if err != nil {
			retErr = err
			return
		}

		// update the title and tags
		txupdate, err := txbooks.PrepareUpdate(ctx)
		if err != nil {
			retErr = err
			return
		}
		defer func() {
			if err := txupdate.Close(); err != nil && retErr == nil {
				retErr = err
			}
		}()
		err = txupdate.Exec(ctx, mygroupstmtiface.BooksUpdateParams{
			BookID: int32(bookOneID),
			Title:  "changed second title",
			Tags:   "cool,disastor",
		})
		if err != nil {
			retErr = err
			return
		}

		// save third book
		_, err = txcreate.Exec(ctx, mygroupstmtiface.BooksCreateParams{
			AuthorID:  authorID,
			Isbn:      "3",
			Title:     "the third book",
			BookType:  mygroupstmtiface.BooksBookTypeFICTION,
			Yr:        2001,
			Available: now,
			Tags:      "cool",
		})
		if err != nil {
			retErr = err
			return
		}

		// save fourth book
		result, err = txcreate.Exec(ctx, mygroupstmtiface.BooksCreateParams{
			AuthorID:  authorID,
			Isbn:      "4",
			Title:     "4th place finisher",
			BookType:  mygroupstmtiface.BooksBookTypeNONFICTION,
			Yr:        2011,
			Available: now,
			Tags:      "other",
		})
		if err != nil {
			retErr = err
			return
		}
		bookThreeID, err := result.LastInsertId()
		if err != nil {
			retErr = err
			return
		}

		// tx commit
		err = tx.Commit()
		if err != nil {
			retErr = err
			return
		}
		retID = int32(bookThreeID)
		return
	}()
	if err != nil {
		t.Fatal(err)
	}

	// upsert, changing ISBN and title
	err = func() (retErr error) {
		dbbooks := mygroupstmtiface.NewBooks(db)
		dbupdate, err := dbbooks.PrepareUpdateISBN(ctx)
		if err != nil {
			retErr = err
			return
		}
		defer func() {
			if err := dbupdate.Close(); err != nil && retErr == nil {
				retErr = err
			}
		}()
		retErr = dbupdate.Exec(ctx, mygroupstmtiface.BooksUpdateISBNParams{
			BookID: bookThreeID,
			Isbn:   "NEW ISBN",
			Title:  "never ever gonna finish, a quatrain",
			Tags:   "someother",
		})
		return
	}()
	if err != nil {
		t.Fatal(err)
	}

	// retrieve first book
	err = func() (retErr error) {
		dbbooks := mygroupstmtiface.NewBooks(db)
		dblist, err := dbbooks.PrepareListByTitleYear(ctx)
		if err != nil {
			retErr = err
			return
		}
		defer func() {
			if err := dblist.Close(); err != nil && retErr == nil {
				retErr = err
				return
			}
		}()
		books0, err := dblist.Exec(ctx, mygroupstmtiface.BooksListByTitleYearParams{
			Title: "my book title",
			Yr:    2016,
		})
		if err != nil {
			retErr = err
			return
		}
		dbauthors := mygroupstmtiface.NewAuthors(db)
		dbget, err := dbauthors.PrepareGet(ctx)
		if err != nil {
			retErr = err
			return
		}
		defer func() {
			if err := dbget.Close(); err != nil && retErr == nil {
				retErr = err
				return
			}
		}()
		for _, book := range books0 {
			t.Logf("Book %d (%s): %s available: %s\n", book.BookID, book.BookType, book.Title, book.Available.Format(time.RFC822Z))
			author, err := dbget.Exec(ctx, book.AuthorID)
			if err != nil {
				retErr = err
				return
			}
			t.Logf("Book %d author: %s\n", book.BookID, author.Name)
		}
		return
	}()
	if err != nil {
		t.Fatal(err)
	}

	// find a book with either "cool" or "other" tag
	err = func() (retErr error) {
		t.Logf("---------\nTag search results:\n")
		dbbooks := mygroupstmtiface.NewBooks(db)
		dblist, err := dbbooks.PrepareListByTags(ctx)
		if err != nil {
			retErr = err
			return
		}
		defer func() {
			if err := dblist.Close(); err != nil && retErr == nil {
				retErr = err
				return
			}
		}()
		res, err := dblist.Exec(ctx, "cool")
		if err != nil {
			retErr = err
			return
		}
		for _, ab := range res {
			t.Logf("Book %d: '%s', Author: '%s', ISBN: '%s' Tags: '%v'\n", ab.BookID, ab.Title, ab.Name, ab.Isbn, ab.Tags)
		}
		return
	}()
	if err != nil {
		t.Fatal(err)
	}

	// get book 4 and delete
	err = func() (retErr error) {
		dbbooks := mygroupstmtiface.NewBooks(db)
		dbget, err := dbbooks.PrepareGet(ctx)
		if err != nil {
			retErr = err
			return
		}
		defer func() {
			if err := dbget.Close(); err != nil && retErr == nil {
				retErr = err
			}
		}()
		b5, err := dbget.Exec(ctx, bookThreeID)
		if err != nil {
			retErr = err
			return
		}
		dbdelete, err := dbbooks.PrepareDelete(ctx)
		if err != nil {
			retErr = err
			return
		}
		defer func() {
			if err := dbdelete.Close(); err != nil && retErr == nil {
				retErr = err
			}
		}()
		if err := dbdelete.Exec(ctx, b5.BookID); err != nil {
			retErr = err
			return
		}
		return
	}()
	if err != nil {
		t.Fatal(err)
	}
}
