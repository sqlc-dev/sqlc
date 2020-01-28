package com.example.booktest.postgresql

import com.example.dbtest.DbTestExtension
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.extension.RegisterExtension
import java.time.OffsetDateTime
import java.time.format.DateTimeFormatter

class QueriesImplTest {
    companion object {
        @JvmField @RegisterExtension val dbtest = DbTestExtension("src/main/resources/booktest/postgresql/schema.sql")
    }

    @Test
    fun testQueries() {
        val conn = dbtest.getConnection()
        val db = QueriesImpl(conn)
        val author = db.createAuthor("Unknown Master")

        // Start a transaction
        conn.autoCommit = false
        db.createBook(
            CreateBookParams(
                authorId = author.authorId,
                isbn = "1",
                title = "my book title",
                booktype = BookType.NONFICTION,
                year = 2016,
                available = OffsetDateTime.now(),
                tags = listOf()
            )
        )

        val b1 = db.createBook(
            CreateBookParams(
                authorId = author.authorId,
                isbn = "2",
                title = "the second book",
                booktype = BookType.NONFICTION,
                year = 2016,
                available = OffsetDateTime.now(),
                tags = listOf("cool", "unique")
            )
        )

        db.updateBook(
            UpdateBookParams(
                bookId = b1.bookId,
                title = "changed second title",
                tags = listOf("cool", "disastor")
            )
        )

        val b3 = db.createBook(
            CreateBookParams(
                authorId = author.authorId,
                isbn = "3",
                title = "the third book",
                booktype = BookType.NONFICTION,
                year = 2001,
                available = OffsetDateTime.now(),
                tags = listOf("cool")
            )
        )

        db.createBook(
            CreateBookParams(
                authorId = author.authorId,
                isbn = "4",
                title = "4th place finisher",
                booktype = BookType.NONFICTION,
                year = 2011,
                available = OffsetDateTime.now(),
                tags = listOf("other")
            )
        )

        // Commit transaction
        conn.commit()
        conn.autoCommit = true

        // ISBN update fails because parameters are not in sequential order. After changing $N to ?, ordering is lost,
        // and the parameters are filled into the wrong slots.
        db.updateBookISBN(
            UpdateBookISBNParams(
                bookId = b3.bookId,
                isbn = "NEW ISBN",
                title = "never ever gonna finish, a quatrain",
                tags = listOf("someother")
            )
        )

        val books0 = db.booksByTitleYear(BooksByTitleYearParams("my book title", 2016))

        val formatter = DateTimeFormatter.ISO_DATE_TIME
        for (book in books0) {
            println("Book ${book.bookId} (${book.booktype}): ${book.title} available: ${book.available.format(formatter)}")
            val author = db.getAuthor(book.authorId)
            println("Book ${book.bookId} author: ${author.name}")
        }

        // find a book with either "cool" or "other" tag
        println("---------\\nTag search results:\\n")
        val res = db.booksByTags(listOf("cool", "other", "someother"))
        for (ab in res) {
            println("Book ${ab.bookId}: '${ab.title}', Author: '${ab.name}', ISBN: '${ab.isbn}' Tags: '${ab.tags.toList()}'")
        }
    }
}