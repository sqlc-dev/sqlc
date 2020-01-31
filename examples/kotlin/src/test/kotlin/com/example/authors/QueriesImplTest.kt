package com.example.authors

import com.example.dbtest.DbTestExtension
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.extension.RegisterExtension

class QueriesImplTest() {

    companion object {
        @JvmField
        @RegisterExtension
        val dbtest = DbTestExtension("src/main/resources/authors/schema.sql")
    }

    @Test
    fun testCreateAuthor() {
        val db = QueriesImpl(dbtest.getConnection())

        val initialAuthors = db.listAuthors()
        assert(initialAuthors.isEmpty())

        val name = "Brian Kernighan"
        val bio = "Co-author of The C Programming Language and The Go Programming Language"
        val insertedAuthor = db.createAuthor(
            name = name,
            bio = bio
        )
        val expectedAuthor = Author(insertedAuthor.id, name, bio)
        assertEquals(expectedAuthor, insertedAuthor)

        val fetchedAuthor = db.getAuthor(insertedAuthor.id)
        assertEquals(expectedAuthor, fetchedAuthor)

        val listedAuthors = db.listAuthors()
        assertEquals(1, listedAuthors.size)
        assertEquals(expectedAuthor, listedAuthors[0])
    }

    @Test
    fun testNull() {
        val db = QueriesImpl(dbtest.getConnection())

        val initialAuthors = db.listAuthors()
        assert(initialAuthors.isEmpty())

        val name = "Brian Kernighan"
        val bio = null
        val insertedAuthor = db.createAuthor(name, bio)
        val expectedAuthor = Author(insertedAuthor.id, name, bio)
        assertEquals(expectedAuthor, insertedAuthor)

        val fetchedAuthor = db.getAuthor(insertedAuthor.id)
        assertEquals(expectedAuthor, fetchedAuthor)

        val listedAuthors = db.listAuthors()
        assertEquals(1, listedAuthors.size)
        assertEquals(expectedAuthor, listedAuthors[0])
    }
}
