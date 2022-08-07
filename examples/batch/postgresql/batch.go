// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: batch.go

package batch

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
)

const booksByYear = `-- name: BooksByYear :batchmany
SELECT book_id, author_id, isbn, book_type, title, year, available, tags FROM books
WHERE year = $1
`

type BooksByYearBatchResults struct {
	br  pgx.BatchResults
	ind int
}

func (q *Queries) BooksByYear(ctx context.Context, year []int32) *BooksByYearBatchResults {
	batch := &pgx.Batch{}
	for _, a := range year {
		vals := []interface{}{
			a,
		}
		batch.Queue(booksByYear, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &BooksByYearBatchResults{br, 0}
}

func (b *BooksByYearBatchResults) Query(f func(int, []Book, error)) {
	for {
		rows, err := b.br.Query()
		if err != nil && (err.Error() == "no result" || err.Error() == "batch already closed") {
			break
		}
		defer rows.Close()
		var items []Book
		for rows.Next() {
			var i Book
			if err := rows.Scan(
				&i.BookID,
				&i.AuthorID,
				&i.Isbn,
				&i.BookType,
				&i.Title,
				&i.Year,
				&i.Available,
				&i.Tags,
			); err != nil {
				break
			}
			items = append(items, i)
		}

		if f != nil {
			f(b.ind, items, rows.Err())
		}
		b.ind++
	}
}

func (b *BooksByYearBatchResults) Close() error {
	return b.br.Close()
}

const createBook = `-- name: CreateBook :batchone
INSERT INTO books (
    author_id,
    isbn,
    book_type,
    title,
    year,
    available,
    tags
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING book_id, author_id, isbn, book_type, title, year, available, tags
`

type CreateBookBatchResults struct {
	br  pgx.BatchResults
	ind int
}

type CreateBookParams struct {
	AuthorID  int32     `json:"author_id"`
	Isbn      string    `json:"isbn"`
	BookType  BookType  `json:"book_type"`
	Title     string    `json:"title"`
	Year      int32     `json:"year"`
	Available time.Time `json:"available"`
	Tags      []string  `json:"tags"`
}

func (q *Queries) CreateBook(ctx context.Context, arg []CreateBookParams) *CreateBookBatchResults {
	batch := &pgx.Batch{}
	for _, a := range arg {
		vals := []interface{}{
			a.AuthorID,
			a.Isbn,
			a.BookType,
			a.Title,
			a.Year,
			a.Available,
			a.Tags,
		}
		batch.Queue(createBook, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &CreateBookBatchResults{br, 0}
}

func (b *CreateBookBatchResults) QueryRow(f func(int, Book, error)) {
	for {
		row := b.br.QueryRow()
		var i Book
		err := row.Scan(
			&i.BookID,
			&i.AuthorID,
			&i.Isbn,
			&i.BookType,
			&i.Title,
			&i.Year,
			&i.Available,
			&i.Tags,
		)
		if err != nil && (err.Error() == "no result" || err.Error() == "batch already closed") {
			break
		}
		if f != nil {
			f(b.ind, i, err)
		}
		b.ind++
	}
}

func (b *CreateBookBatchResults) Close() error {
	return b.br.Close()
}

const deleteBook = `-- name: DeleteBook :batchexec
DELETE FROM books
WHERE book_id = $1
`

type DeleteBookBatchResults struct {
	br  pgx.BatchResults
	ind int
}

func (q *Queries) DeleteBook(ctx context.Context, bookID []int32) *DeleteBookBatchResults {
	batch := &pgx.Batch{}
	for _, a := range bookID {
		vals := []interface{}{
			a,
		}
		batch.Queue(deleteBook, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &DeleteBookBatchResults{br, 0}
}

func (b *DeleteBookBatchResults) Exec(f func(int, error)) {
	for {
		_, err := b.br.Exec()
		if err != nil && (err.Error() == "no result" || err.Error() == "batch already closed") {
			break
		}
		if f != nil {
			f(b.ind, err)
		}
		b.ind++
	}
}

func (b *DeleteBookBatchResults) Close() error {
	return b.br.Close()
}

const updateBook = `-- name: UpdateBook :batchexec
UPDATE books
SET title = $1, tags = $2
WHERE book_id = $3
`

type UpdateBookBatchResults struct {
	br  pgx.BatchResults
	ind int
}

type UpdateBookParams struct {
	Title  string   `json:"title"`
	Tags   []string `json:"tags"`
	BookID int32    `json:"book_id"`
}

func (q *Queries) UpdateBook(ctx context.Context, arg []UpdateBookParams) *UpdateBookBatchResults {
	batch := &pgx.Batch{}
	for _, a := range arg {
		vals := []interface{}{
			a.Title,
			a.Tags,
			a.BookID,
		}
		batch.Queue(updateBook, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &UpdateBookBatchResults{br, 0}
}

func (b *UpdateBookBatchResults) Exec(f func(int, error)) {
	for {
		_, err := b.br.Exec()
		if err != nil && (err.Error() == "no result" || err.Error() == "batch already closed") {
			break
		}
		if f != nil {
			f(b.ind, err)
		}
		b.ind++
	}
}

func (b *UpdateBookBatchResults) Close() error {
	return b.br.Close()
}
