// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package com.example.booktest.mysql

import java.sql.Connection
import java.sql.SQLException
import java.sql.Statement
import java.time.LocalDateTime

const val booksByTags = """-- name: booksByTags :many
SELECT
  book_id,
  title,
  name,
  isbn,
  tags
FROM books
LEFT JOIN authors ON books.author_id = authors.author_id
WHERE tags = ?
"""

data class BooksByTagsRow (
  val bookId: Int,
  val title: String,
  val name: String?,
  val isbn: String,
  val tags: String
)

const val booksByTitleYear = """-- name: booksByTitleYear :many
SELECT book_id, author_id, isbn, book_type, title, yr, available, tags FROM books
WHERE title = ? AND yr = ?
"""

const val createAuthor = """-- name: createAuthor :execresult
INSERT INTO authors (name) VALUES (?)
"""

const val createBook = """-- name: createBook :execresult
INSERT INTO books (
    author_id,
    isbn,
    book_type,
    title,
    yr,
    available,
    tags
) VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
)
"""

const val deleteAuthorBeforeYear = """-- name: deleteAuthorBeforeYear :exec
DELETE FROM books
WHERE yr < ? AND author_id = ?
"""

const val deleteBook = """-- name: deleteBook :exec
DELETE FROM books
WHERE book_id = ?
"""

const val getAuthor = """-- name: getAuthor :one
SELECT author_id, name FROM authors
WHERE author_id = ?
"""

const val getBook = """-- name: getBook :one
SELECT book_id, author_id, isbn, book_type, title, yr, available, tags FROM books
WHERE book_id = ?
"""

const val updateBook = """-- name: updateBook :exec
UPDATE books
SET title = ?, tags = ?
WHERE book_id = ?
"""

const val updateBookISBN = """-- name: updateBookISBN :exec
UPDATE books
SET title = ?, tags = ?, isbn = ?
WHERE book_id = ?
"""

class QueriesImpl(private val conn: Connection) : Queries {

  @Throws(SQLException::class)
  override fun booksByTags(tags: String): List<BooksByTagsRow> {
    return conn.prepareStatement(booksByTags).use { stmt ->
      stmt.setString(1, tags)

      val results = stmt.executeQuery()
      val ret = mutableListOf<BooksByTagsRow>()
      while (results.next()) {
          ret.add(BooksByTagsRow(
                results.getInt(1),
                results.getString(2),
                results.getString(3),
                results.getString(4),
                results.getString(5)
            ))
      }
      ret
    }
  }

  @Throws(SQLException::class)
  override fun booksByTitleYear(title: String, yr: Int): List<Book> {
    return conn.prepareStatement(booksByTitleYear).use { stmt ->
      stmt.setString(1, title)
          stmt.setInt(2, yr)

      val results = stmt.executeQuery()
      val ret = mutableListOf<Book>()
      while (results.next()) {
          ret.add(Book(
                results.getInt(1),
                results.getInt(2),
                results.getString(3),
                BooksBookType.lookup(results.getString(4))!!,
                results.getString(5),
                results.getInt(6),
                results.getObject(7, LocalDateTime::class.java),
                results.getString(8)
            ))
      }
      ret
    }
  }

  @Throws(SQLException::class)
  override fun createAuthor(name: String): Long {
    return conn.prepareStatement(createAuthor, Statement.RETURN_GENERATED_KEYS).use { stmt ->
      stmt.setString(1, name)

      stmt.execute()

      val results = stmt.generatedKeys
      if (!results.next()) {
          throw SQLException("no generated key returned")
      }
	  results.getLong(1)
    }
  }

  @Throws(SQLException::class)
  override fun createBook(
      authorId: Int,
      isbn: String,
      bookType: BooksBookType,
      title: String,
      yr: Int,
      available: LocalDateTime,
      tags: String): Long {
    return conn.prepareStatement(createBook, Statement.RETURN_GENERATED_KEYS).use { stmt ->
      stmt.setInt(1, authorId)
          stmt.setString(2, isbn)
          stmt.setString(3, bookType.value)
          stmt.setString(4, title)
          stmt.setInt(5, yr)
          stmt.setObject(6, available)
          stmt.setString(7, tags)

      stmt.execute()

      val results = stmt.generatedKeys
      if (!results.next()) {
          throw SQLException("no generated key returned")
      }
	  results.getLong(1)
    }
  }

  @Throws(SQLException::class)
  override fun deleteAuthorBeforeYear(yr: Int, authorId: Int) {
    conn.prepareStatement(deleteAuthorBeforeYear).use { stmt ->
      stmt.setInt(1, yr)
          stmt.setInt(2, authorId)

      stmt.execute()
    }
  }

  @Throws(SQLException::class)
  override fun deleteBook(bookId: Int) {
    conn.prepareStatement(deleteBook).use { stmt ->
      stmt.setInt(1, bookId)

      stmt.execute()
    }
  }

  @Throws(SQLException::class)
  override fun getAuthor(authorId: Int): Author? {
    return conn.prepareStatement(getAuthor).use { stmt ->
      stmt.setInt(1, authorId)

      val results = stmt.executeQuery()
      if (!results.next()) {
        return null
      }
      val ret = Author(
                results.getInt(1),
                results.getString(2)
            )
      if (results.next()) {
          throw SQLException("expected one row in result set, but got many")
      }
      ret
    }
  }

  @Throws(SQLException::class)
  override fun getBook(bookId: Int): Book? {
    return conn.prepareStatement(getBook).use { stmt ->
      stmt.setInt(1, bookId)

      val results = stmt.executeQuery()
      if (!results.next()) {
        return null
      }
      val ret = Book(
                results.getInt(1),
                results.getInt(2),
                results.getString(3),
                BooksBookType.lookup(results.getString(4))!!,
                results.getString(5),
                results.getInt(6),
                results.getObject(7, LocalDateTime::class.java),
                results.getString(8)
            )
      if (results.next()) {
          throw SQLException("expected one row in result set, but got many")
      }
      ret
    }
  }

  @Throws(SQLException::class)
  override fun updateBook(
      title: String,
      tags: String,
      bookId: Int) {
    conn.prepareStatement(updateBook).use { stmt ->
      stmt.setString(1, title)
          stmt.setString(2, tags)
          stmt.setInt(3, bookId)

      stmt.execute()
    }
  }

  @Throws(SQLException::class)
  override fun updateBookISBN(
      title: String,
      tags: String,
      isbn: String,
      bookId: Int) {
    conn.prepareStatement(updateBookISBN).use { stmt ->
      stmt.setString(1, title)
          stmt.setString(2, tags)
          stmt.setString(3, isbn)
          stmt.setInt(4, bookId)

      stmt.execute()
    }
  }

}

