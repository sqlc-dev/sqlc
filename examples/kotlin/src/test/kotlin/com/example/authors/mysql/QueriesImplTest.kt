package com.example.authors.mysql

import com.example.dbtest.MysqlDbTestExtension
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.extension.RegisterExtension

class QueriesImplTest() {

    companion object {
        @JvmField
        @RegisterExtension
        val dbtest = MysqlDbTestExtension("src/main/resources/authors/mysql/schema.sql")
    }

    @Test
    fun testCreateAuthor() {
        val db = QueriesImpl(dbtest.getConnection())

        val initialAuthors = db.listAuthors()
        assert(initialAuthors.isEmpty())

        val name = "Brian Kernighan"
        val bio = "Co-author of The C Programming Language and The Go Programming Language"
        val id = db.createAuthor(
            name = name,
            bio = bio
        )
        assertEquals(id, 1)
        val expectedAuthor = Author(id, name, bio)

        val fetchedAuthor = db.getAuthor(id)
        assertEquals(expectedAuthor, fetchedAuthor)

        val listedAuthors = db.listAuthors()
        assertEquals(1, listedAuthors.size)
        assertEquals(expectedAuthor, listedAuthors[0])

        val id2 = db.createAuthor(
            name = name,
            bio = bio
        )
        assertEquals(id2, 2)
    }

    @Test
    fun testNull() {
        val db = QueriesImpl(dbtest.getConnection())

        val initialAuthors = db.listAuthors()
        assert(initialAuthors.isEmpty())

        val name = "Brian Kernighan"
        val bio = null
        val id = db.createAuthor(name, bio)
        val expectedAuthor = Author(id, name, bio)

        val fetchedAuthor = db.getAuthor(id)
        assertEquals(expectedAuthor, fetchedAuthor)

        val listedAuthors = db.listAuthors()
        assertEquals(1, listedAuthors.size)
        assertEquals(expectedAuthor, listedAuthors[0])
    }
}
