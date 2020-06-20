// Code generated by sqlc. DO NOT EDIT.

package com.example.booktest.postgresql

import java.sql.Connection
import java.sql.SQLException
import java.sql.Types
import java.time.OffsetDateTime

import sqlc.runtime.ExecuteQuery
import sqlc.runtime.ListQuery
import sqlc.runtime.RowQuery

interface Queries {
  @Throws(SQLException::class)
  fun booksByTags(dollar1: List<String>): ListQuery<BooksByTagsRow>
  
  @Throws(SQLException::class)
  fun booksByTitleYear(title: String, year: Int): ListQuery<Book>
  
  @Throws(SQLException::class)
  fun createAuthor(name: String): RowQuery<Author>
  
  @Throws(SQLException::class)
  fun createBook(
      authorId: Int,
      isbn: String,
      booktype: BookType,
      title: String,
      year: Int,
      available: OffsetDateTime,
      tags: List<String>): RowQuery<Book>
  
  @Throws(SQLException::class)
  fun deleteBook(bookId: Int): ExecuteQuery
  
  @Throws(SQLException::class)
  fun getAuthor(authorId: Int): RowQuery<Author>
  
  @Throws(SQLException::class)
  fun getBook(bookId: Int): RowQuery<Book>
  
  @Throws(SQLException::class)
  fun updateBook(
      title: String,
      tags: List<String>,
      bookId: Int): ExecuteQuery
  
  @Throws(SQLException::class)
  fun updateBookISBN(
      title: String,
      tags: List<String>,
      bookId: Int,
      isbn: String): ExecuteQuery
  
}

