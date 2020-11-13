package com.example.booktest.mysql

import com.example.dbtest.MysqlDbTestExtension
import com.example.dbtest.PostgresDbTestExtension
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.extension.RegisterExtension
import java.time.LocalDateTime
import java.time.OffsetDateTime
import java.time.format.DateTimeFormatter

class QueriesImplTest {
    companion object {
        @JvmField @RegisterExtension val dbtest = MysqlDbTestExtension("src/main/resources/booktest/mysql/schema.sql")
    }

    @Test
    fun testQueries() {
        val conn = dbtest.getConnection()
        val db = QueriesImpl(conn)
        val authorId = db.createAuthor("Unknown Master")
        val author = db.getAuthor(authorId.toInt())!!

        // Start a transaction
        conn.autoCommit = false
        db.createBook(
                authorId = author.authorId,
                isbn = "1",
                title = "my book title",
                bookType = BooksBookType.NONFICTION,
                yr = 2016,
                available = LocalDateTime.now(),
                tags = ""
        )

        val b1Id = db.createBook(
                authorId = author.authorId,
                isbn = "2",
                title = "the second book",
                bookType = BooksBookType.NONFICTION,
                yr = 2016,
                available = LocalDateTime.now(),
                tags = listOf("cool", "unique").joinToString(",")
        )

        db.updateBook(
                bookId = b1Id.toInt(),
                title = "changed second title",
                tags = listOf("cool", "disastor").joinToString(",")
        )

        val b3Id = db.createBook(
                authorId = author.authorId,
                isbn = "3",
                title = "the third book",
                bookType = BooksBookType.NONFICTION,
                yr = 2001,
                available = LocalDateTime.now(),
                tags = listOf("cool").joinToString(",")
        )

        db.createBook(
                authorId = author.authorId,
                isbn = "4",
                title = "4th place finisher",
                bookType = BooksBookType.NONFICTION,
                yr = 2011,
                available = LocalDateTime.now(),
                tags = listOf("other").joinToString(",")
        )

        // Commit transaction
        conn.commit()
        conn.autoCommit = true

        db.updateBookISBN(
                bookId = b3Id.toInt(),
                isbn = "NEW ISBN",
                title = "never ever gonna finish, a quatrain",
                tags = listOf("someother").joinToString(",")
        )

        val books0 = db.booksByTitleYear("my book title", 2016)

        val formatter = DateTimeFormatter.ISO_DATE_TIME
        for (book in books0) {
            println("Book ${book.bookId} (${book.bookType}): ${book.title} available: ${book.available.format(formatter)}")
            val author2 = db.getAuthor(book.authorId)!!
            println("Book ${book.bookId} author: ${author2.name}")
        }

        // find a book with either "cool" or "other" tag
        println("---------\\nTag search results:\\n")
        val res = db.booksByTags(listOf("cool", "other", "someother").joinToString(","))
        for (ab in res) {
            println("Book ${ab.bookId}: '${ab.title}', Author: '${ab.name}', ISBN: '${ab.isbn}' Tags: '${ab.tags.toList()}'")
        }
    }
}